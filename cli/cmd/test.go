// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ballerina-nutcracker/ballerina/cli/internal/testdiscovery"
	"github.com/ballerina-nutcracker/ballerina/cli/internal/testglue"
	debugcommon "github.com/ballerina-nutcracker/ballerina/common"
	_ "github.com/ballerina-nutcracker/ballerina/lib/rt"
	"github.com/ballerina-nutcracker/ballerina/platform/palnative"
	"github.com/ballerina-nutcracker/ballerina/projects"
	"github.com/ballerina-nutcracker/ballerina/runtime"
	"github.com/ballerina-nutcracker/ballerina/tools/diagnostics"
	"github.com/ballerina-nutcracker/ballerina/values"

	"github.com/spf13/cobra"
)

// testGlueDocumentName names the synthesized test-registration document
// injected into the module's test compilation (see runTestsForProject). Not
// a real file on disk.
const testGlueDocumentName = "__test_register__.bal"

type testCmdOptions struct {
	groups        string
	disableGroups string
	tests         string
	listGroups    bool
	testReport    bool
	rerunFailed   bool
	codeCoverage  bool
	graalvm       bool
	graalvmOpts   string
	cloud         string
	parallel      bool
	offline       bool
	sticky        bool
	targetDir     string

	dumpTokens    bool
	dumpST        bool
	dumpAST       bool
	dumpCFG       bool
	dumpBIR       bool
	traceRecovery bool
	stats         bool
	statsOneline  bool
	logFile       string
	format        string
}

var testCmd = createTestCmd()

func createTestCmd() *cobra.Command {
	opts := &testCmdOptions{}
	cmd := &cobra.Command{
		Use:   "test [<source-file.bal> | <package-dir>]",
		Short: "Run package tests",
		Long: `	Run the tests of the current package, a standalone '.bal' file, or every
	member package of a workspace.

	Discovers '@test:*'-annotated functions in the package's test sources
	(or, for a standalone file, in the file itself), registers them with
	ballerina/test, and executes them in-process, printing a pass/fail/skip
	summary. When run at a workspace root, every member package's tests run
	in turn; when run inside one member's own directory, only that member's
	tests run.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTest(cmd, args, opts)
		},
	}
	cmd.Flags().StringVar(&opts.groups, "groups", "", "Test groups to be executed")
	cmd.Flags().StringVar(&opts.disableGroups, "disable-groups", "", "Test groups to be disabled")
	cmd.Flags().StringVar(&opts.tests, "tests", "", "Test functions to be executed")
	cmd.Flags().BoolVar(&opts.listGroups, "list-groups", false, "List the groups available in the tests")
	cmd.Flags().BoolVar(&opts.testReport, "test-report", false, "Enable test report generation")
	cmd.Flags().BoolVar(&opts.rerunFailed, "rerun-failed", false, "Rerun failed tests")
	cmd.Flags().BoolVar(&opts.codeCoverage, "code-coverage", false, "Enable code coverage (not supported)")
	cmd.Flags().BoolVar(&opts.graalvm, "graalvm", false, "Run test suite against native image (not supported)")
	cmd.Flags().StringVar(&opts.graalvmOpts, "graalvm-build-options", "", "Native image build options (not supported)")
	cmd.Flags().StringVar(&opts.cloud, "cloud", "", "Enable cloud artifact generation (not supported)")
	cmd.Flags().BoolVar(&opts.parallel, "parallel", false, "Enable parallel execution of tests (not supported, falls back to serial)")
	cmd.Flags().BoolVar(&opts.offline, "offline", false, "Resolve dependencies offline")
	cmd.Flags().BoolVar(&opts.sticky, "sticky", false, "Stick to exact versions locked (if exists)")
	cmd.Flags().StringVar(&opts.targetDir, "target-dir", "", "Target directory path")
	cmd.Flags().BoolVar(&opts.dumpTokens, "dump-tokens", false, "Dump lexer tokens")
	cmd.Flags().BoolVar(&opts.dumpST, "dump-st", false, "Dump syntax tree")
	cmd.Flags().BoolVar(&opts.dumpAST, "dump-ast", false, "Dump abstract syntax tree")
	cmd.Flags().BoolVar(&opts.dumpCFG, "dump-cfg", false, "Dump control flow graph")
	cmd.Flags().BoolVar(&opts.dumpBIR, "dump-bir", false, "Dump Ballerina Intermediate Representation")
	cmd.Flags().BoolVar(&opts.traceRecovery, "trace-recovery", false, "Enable error recovery tracing")
	cmd.Flags().BoolVar(&opts.stats, "stats", false, "Print per-stage compilation timing statistics")
	cmd.Flags().BoolVar(&opts.statsOneline, "stats-oneline", false, "Print per-stage compilation timing totals only")
	cmd.Flags().StringVar(&opts.logFile, "log-file", "", "Write debug output to specified file")
	cmd.Flags().StringVar(&opts.format, "format", "", "Output format for dump operations (dot)")
	return cmd
}

func testError(format string, args ...any) error {
	return usageError("test [<source-file.bal> | <package-dir>]", format, args...)
}

func runTest(cmd *cobra.Command, args []string, opts *testCmdOptions) error {
	stderr := cmd.ErrOrStderr()

	// Unsupported-scope flags: fail clearly rather than silently misbehaving
	// or crashing (see TODO.md's "Deferred bal test features").
	if opts.graalvm || opts.graalvmOpts != "" {
		return testError("--graalvm is not supported: this interpreter has no separate native-image step")
	}
	if opts.cloud != "" {
		return testError("--cloud is not supported")
	}
	if opts.codeCoverage {
		fmt.Fprintln(stderr, "warning: --code-coverage is not supported yet; running without coverage")
	}
	if opts.parallel {
		fmt.Fprintln(stderr, "warning: --parallel is not supported yet; running tests serially")
	}

	buildOpts := projects.NewBuildOptionsBuilder().
		WithSkipTests(false).
		WithOffline(opts.offline).
		WithSticky(opts.sticky).
		WithTargetDir(opts.targetDir).
		WithTestReport(opts.testReport).
		WithDumpAST(opts.dumpAST).
		WithDumpBIR(opts.dumpBIR).
		WithDumpCFG(opts.dumpCFG).
		WithDumpCFGFormat(projects.ParseCFGFormat(opts.format)).
		WithDumpTokens(opts.dumpTokens).
		WithDumpST(opts.dumpST).
		WithTraceRecovery(opts.traceRecovery).
		WithStats(opts.stats || opts.statsOneline).
		Build()

	debugFlags := uint16(0)
	if buildOpts.DumpTokens() {
		debugFlags |= debugcommon.DUMP_TOKENS
	}
	if buildOpts.DumpST() {
		debugFlags |= debugcommon.DUMP_ST
	}
	if buildOpts.TraceRecovery() {
		debugFlags |= debugcommon.DEBUG_ERROR_RECOVERY
	}
	if debugFlags != 0 {
		if opts.logFile != "" {
			logWriter, err := os.Create(opts.logFile)
			if err != nil {
				return testError("error creating log file %s: %w", opts.logFile, err)
			}
			defer func() { _ = logWriter.Close() }()
			debugcommon.InitDebug(debugFlags, logWriter)
		} else {
			debugcommon.InitDebug(debugFlags, stderr)
		}
	}

	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	info, err := os.Stat(path)
	if err != nil {
		return testError("invalid project path %q: %w", path, err)
	}

	// A single .bal file loads like bal run/bal build load one: fsys rooted
	// at the parent dir, loadPath is the filename within it. A standalone
	// file can't be a workspace member, so workspace detection only applies
	// to directories (matches build.go).
	baseDir := path
	loadPath := "."
	if !info.IsDir() {
		if filepath.Ext(path) != ".bal" {
			return testError("%q is not a package directory or a .bal file", path)
		}
		baseDir = filepath.Dir(path)
		loadPath = filepath.Base(path)
	}

	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return testError("resolve absolute path: %w", err)
	}

	// Detect whether absBaseDir sits inside a workspace without being its
	// root — e.g. cwd is a workspace member's own directory. If so, load
	// from the workspace root instead so sibling member-to-member
	// dependencies resolve, matching bal build's findWorkspaceRoot handling.
	// Only applies to directories; a standalone .bal file can't be a
	// workspace member.
	workspaceRoot := ""
	if info.IsDir() {
		workspaceRoot = findWorkspaceRoot(absBaseDir)
	}
	effectiveBaseDir, effectiveLoadPath := absBaseDir, loadPath
	if workspaceRoot != "" && workspaceRoot != absBaseDir {
		effectiveBaseDir, effectiveLoadPath = workspaceRoot, "."
	}

	ballerinaEnvPath, err := getBallerinaEnvPath()
	if err != nil {
		return testError("resolve ballerina env path: %w", err)
	}

	fsys := os.DirFS(effectiveBaseDir)
	result, err := projects.Load(fsys, effectiveLoadPath, projects.ProjectLoadConfig{
		BallerinaEnvFs: os.DirFS(ballerinaEnvPath),
		BuildOptions:   &buildOpts,
	})
	if err != nil {
		return testError("failed to load package: %w", err)
	}
	if diagResult := result.Diagnostics(); diagResult.HasErrors() || diagResult.HasWarnings() {
		printDiagnostics(fsys, stderr, diagResult, !isTerminal(), diagnostics.NewDiagnosticEnv())
		if diagResult.HasErrors() {
			return testError("package loading reported errors")
		}
	}

	project := result.Project()
	if project.Kind() == projects.ProjectKindWorkspace {
		workspace := project.(*projects.WorkspaceProject)

		if workspaceRoot != "" && workspaceRoot != absBaseDir {
			// absBaseDir names one specific member (we walked up to
			// workspaceRoot to load it) — test just that member, matching
			// bal run's disambiguation.
			memberProject := findBuildProjectByPath(workspace, workspaceRoot, absBaseDir)
			if memberProject == nil {
				return testError("no package found at path %s within workspace %s", absBaseDir, workspaceRoot)
			}
			exitCode, pkgResult, err := runTestsForProject(cmd, opts, stderr, fsys, memberProject, absBaseDir)
			return finishStandaloneTestRun(cmd, absBaseDir, exitCode, pkgResult, err)
		}

		return runTestsForWorkspace(cmd, opts, stderr, fsys, effectiveLoadPath, workspace, absBaseDir, ballerinaEnvPath, &buildOpts)
	}

	exitCode, pkgResult, err := runTestsForProject(cmd, opts, stderr, fsys, project, absBaseDir)
	return finishStandaloneTestRun(cmd, absBaseDir, exitCode, pkgResult, err)
}

// finishStandaloneTestRun applies runTestsForProject's result the same way
// for both call sites that treat a run as one self-contained package (a
// plain non-workspace invocation, and "bal test" invoked from inside one
// workspace member's own directory): write that package's own
// test_results.json if requested, then translate a non-zero suite exit
// code into the usual "there are test failures" command error.
func finishStandaloneTestRun(cmd *cobra.Command, projectDir string, exitCode int, pkgResult *packageTestResult, err error) error {
	if err != nil {
		return err
	}
	if pkgResult != nil {
		if reportErr := writeTestResultsReport(cmd, projectDir, *pkgResult); reportErr != nil {
			return reportErr
		}
	}
	if exitCode != 0 {
		// The test suite's own failure, not a test-usage mistake — no USAGE block.
		cmd.SilenceUsage = true
		return fmt.Errorf("there are test failures")
	}
	return nil
}

// runTestsForWorkspace runs every member package's tests in turn, matching
// jballerina's behavior: unlike a genuine load/compile error (which aborts
// immediately, since something is actually broken), one member's test
// failures don't stop the rest from running — every member gets a chance to
// report its own results, and the workspace's overall exit status reflects
// whether *any* member failed.
//
// Each member is tested via a freshly-reloaded workspace (like build.go's
// buildNativeStub/buildOneProject loop) so each gets its own independent
// CompilerEnvironment rather than sharing one across members.
func runTestsForWorkspace(cmd *cobra.Command, opts *testCmdOptions, stderr io.Writer, fsys fs.FS,
	loadPath string, workspace *projects.WorkspaceProject, workspaceAbsRoot, ballerinaEnvPath string,
	buildOpts *projects.BuildOptions) error {
	anyFailed := false
	var pkgResults []packageTestResult
	for _, bp := range workspace.Projects() {
		memberDir := filepath.Join(workspaceAbsRoot, bp.SourceRoot())

		memberResult, err := projects.Load(fsys, loadPath, projects.ProjectLoadConfig{
			BallerinaEnvFs: os.DirFS(ballerinaEnvPath),
			BuildOptions:   buildOpts,
		})
		if err != nil {
			return testError("failed to load package: %w", err)
		}

		memberWorkspace := memberResult.Project().(*projects.WorkspaceProject)
		memberProject := findBuildProjectByPath(memberWorkspace, workspaceAbsRoot, memberDir)
		if memberProject == nil {
			return testError("no package found at path %s within workspace %s", memberDir, workspaceAbsRoot)
		}

		pkgName := memberProject.CurrentPackage().PackageName().Value()
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nRunning tests for package '%s'\n", pkgName)

		exitCode, pkgResult, err := runTestsForProject(cmd, opts, stderr, fsys, memberProject, memberDir)
		if err != nil {
			// A genuine load/compile/internal error, not just "some tests
			// failed" — matches jballerina's re-throw for non-test-failure
			// exceptions instead of continuing to the next member.
			return err
		}
		if exitCode != 0 {
			anyFailed = true
		}
		if pkgResult != nil {
			pkgResults = append(pkgResults, *pkgResult)
		}
	}

	if len(pkgResults) > 0 {
		report := workspaceTestReport{WorkspaceName: filepath.Base(workspaceAbsRoot), Packages: pkgResults}
		for _, pkg := range pkgResults {
			report.TotalTests += pkg.TotalTests
			report.Passed += pkg.Passed
			report.Failed += pkg.Failed
			report.Skipped += pkg.Skipped
		}
		if err := writeTestResultsReport(cmd, workspaceAbsRoot, report); err != nil {
			return err
		}
	}

	if anyFailed {
		cmd.SilenceUsage = true
		return fmt.Errorf("there are test failures")
	}
	return nil
}

// runTestsForProject compiles project's package once (test sources
// included), then runs each of its modules' tests in turn — almost always
// just the default module, but a package can have named submodules under
// modules/ too. Each module gets its own independent glue-injection +
// recompile + fresh runtime + setTestOptions/register/startSuite sequence
// (mirroring runTestsForWorkspace's per-member isolation): ballerina/test's
// registries are shared module-level state within its own module, so
// running two user-modules' tests through one shared runtime would
// accumulate both modules' registrations into a single combined report
// instead of jballerina's actual per-module reports (and per-module
// `--tests mod:name` filtering couldn't work at all). Returns the combined
// exit code (0 only if every module with tests passed) separately from any
// error.
//
// Confirmed via a real multi-module package fixture (2026-09-09, see
// CHECKPOINT.md) that before this, only the default module's tests ever ran
// — a submodule's tests were silently never discovered at all, reporting
// "No tests found" / exit 0 instead of actually running them.
// runTestsForProject's third return value is the package's aggregated test
// report data (nil unless opts.testReport is set and at least one module
// actually ran tests) — the caller decides whether to write it standalone
// (a single package) or fold it into a workspace-wide report (see
// writeTestResultsReport's callers in runTest/runTestsForWorkspace).
func runTestsForProject(cmd *cobra.Command, opts *testCmdOptions, stderr io.Writer, fsys fs.FS,
	project projects.Project, projectDir string) (int, *packageTestResult, error) {
	pkg := project.CurrentPackage()
	compilation := pkg.Compilation()
	if cd := compilation.DiagnosticResult(); cd.HasErrors() || cd.HasWarnings() {
		printDiagnostics(fsys, stderr, cd, !isTerminal(), compilation.DiagnosticEnv())
		if cd.HasErrors() {
			return 1, nil, testError("compilation contains errors")
		}
	}

	if opts.statsOneline {
		_, _ = fmt.Fprint(stderr, compilation.StatsReportOneline())
	} else if opts.stats {
		_, _ = fmt.Fprint(stderr, compilation.StatsReport())
	}

	modules := pkg.Modules()
	var exitCode int
	var err error
	if len(modules) == 1 {
		// Fast path, and identical behavior to before multi-module support:
		// always run the (only) module, even with zero discovered tests —
		// startSuite() itself reports "No tests found" in that case.
		exitCode, err = runTestsForModule(cmd, opts, stderr, fsys, project, compilation, modules[0], projectDir)
	} else {
		exitCode, err = runTestsForModules(cmd, opts, stderr, fsys, project, compilation, modules, projectDir)
	}
	if err != nil {
		return 1, nil, err
	}

	if !opts.testReport {
		return exitCode, nil, nil
	}
	result, hadAny := buildPackageTestResult(pkg, modules, projectDir)
	if !hadAny {
		return exitCode, nil, nil
	}
	return exitCode, &result, nil
}

// runTestsForModules runs every module that has discovered tests, matching
// jballerina's own per-module report semantics (see runTestsForProject's
// doc comment for why each module needs its own runTestsForModule call
// rather than sharing one runtime). A module with zero discovered tests is
// skipped entirely — no header, no report — rather than running it through
// startSuite() just to print "No tests found" for each empty module.
func runTestsForModules(cmd *cobra.Command, opts *testCmdOptions, stderr io.Writer, fsys fs.FS,
	project projects.Project, compilation *projects.PackageCompilation,
	modules []*projects.Module, projectDir string) (int, error) {
	anyHadTests := false
	anyFailed := false
	for _, module := range modules {
		astPkg := compilation.ModuleAST(module.ModuleID())
		if !testdiscovery.Discover(astPkg).HasAny() {
			continue
		}
		anyHadTests = true
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nRunning tests for module '%s'\n", module.ModuleName().String())
		exitCode, err := runTestsForModule(cmd, opts, stderr, fsys, project, compilation, module, projectDir)
		if err != nil {
			return 1, err
		}
		if exitCode != 0 {
			anyFailed = true
		}
	}
	if !anyHadTests {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "\tNo tests found")
		return 0, nil
	}
	if anyFailed {
		return 1, nil
	}
	return 0, nil
}

// runTestsForModule discovers @test:* functions in one module, injects
// generated registration glue as a new test document, recompiles, and runs
// that module's suite in-process via its own fresh runtime. Returns the
// suite's own exit code (0 = all passed) separately from any error.
func runTestsForModule(cmd *cobra.Command, opts *testCmdOptions, stderr io.Writer, fsys fs.FS,
	project projects.Project, compilation *projects.PackageCompilation,
	module *projects.Module, projectDir string) (int, error) {
	astPkg := compilation.ModuleAST(module.ModuleID())
	discovered := testdiscovery.Discover(astPkg)
	glueSource := testglue.GenerateSource(discovered)

	docID := projects.NewDocumentID(testGlueDocumentName, module.ModuleID())
	docConfig := projects.NewDocumentConfig(docID, testGlueDocumentName, glueSource)
	newModule := module.Modify().AddTestDocument(docConfig).Apply()
	newPkg := newModule.PackageInstance()

	testCompilation := newPkg.Compilation()
	if cd := testCompilation.DiagnosticResult(); cd.HasErrors() || cd.HasWarnings() {
		printDiagnostics(fsys, stderr, cd, !isTerminal(), testCompilation.DiagnosticEnv())
		if cd.HasErrors() {
			return 1, testError("internal error: generated test registration glue failed to compile")
		}
	}

	backend := projects.NewBallerinaBackend(testCompilation)
	birPkgs := backend.BIRPackages()
	if len(birPkgs) == 0 {
		return 1, testError("BIR generation failed: no BIR package produced")
	}

	rootOrg := newPkg.PackageOrg().Value()
	rootName := newPkg.PackageName().Value()
	for _, birPkg := range birPkgs {
		if birPkg.PackageID.OrgName.Value() == rootOrg && birPkg.PackageID.PkgName.Value() == rootName {
			// bal test never runs the package's own main() — only the
			// generated registration entry point and ballerina/test's
			// startSuite() drive execution.
			birPkg.MainFunction = nil
		}
	}

	tyEnv := project.Environment().TypeEnv()
	pal, cleanupSignals := palnative.NewPlatform()
	defer cleanupSignals()
	rt := runtime.NewRuntime(pal, tyEnv)
	for _, birPkg := range birPkgs {
		if err := rt.Init(*birPkg); err != nil {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err)
			cmd.SilenceErrors = true
			return 1, err
		}
	}
	rt.Listen()

	targetPath := filepath.Join(projectDir, projects.TargetDir)
	moduleName := module.ModuleName().String()
	pkgName := newPkg.PackageName().Value()

	setTestOptionsFn, ok := runtime.LookupFunction(rt, "ballerina", "test", "setTestOptions")
	if !ok {
		return 1, testError("internal error: ballerina/test:setTestOptions not found")
	}
	if _, err := runtime.InvokeFunction(rt, setTestOptionsFn, []values.BalValue{
		targetPath, pkgName, moduleName,
		strconv.FormatBool(opts.testReport), strconv.FormatBool(false),
		opts.groups, opts.disableGroups, opts.tests,
		strconv.FormatBool(opts.rerunFailed), strconv.FormatBool(opts.listGroups),
	}); err != nil {
		return 1, err
	}

	registerFn, ok := runtime.LookupFunction(rt, rootOrg, moduleName, testglue.EntryFunctionName)
	if !ok {
		return 1, testError("internal error: %s not found", testglue.EntryFunctionName)
	}
	if _, err := runtime.InvokeFunction(rt, registerFn, nil); err != nil {
		return 1, err
	}

	startSuiteFn, ok := runtime.LookupFunction(rt, "ballerina", "test", "startSuite")
	if !ok {
		return 1, testError("internal error: ballerina/test:startSuite not found")
	}
	suiteResult, err := runtime.InvokeFunction(rt, startSuiteFn, nil)
	if err != nil {
		return 1, err
	}
	exitCode, _ := suiteResult.(int64)
	return int(exitCode), nil
}
