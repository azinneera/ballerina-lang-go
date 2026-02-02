/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/bir"
	debugcommon "ballerina-lang-go/common"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/runtime"

	"ballerina-lang-go/cli/pkg/templates"

	"github.com/spf13/cobra"
)

var runLongDesc = `Compile the current package and run it.

The 'run' command compiles and executes the given Ballerina source file.

A Ballerina program consists of one or more modules; one of these modules is distinguished as the root
module, which is the default module of current package.

Ballerina program execution consists of two consecutive phases. The initialization phase initializes all
modules of a program one after another. If a module defines a function named 'init()', it will be invoked
during this phase. If the root module of the program defines a public function named 'main()', then it
will be invoked.

If the initialization phase of program execution completes successfully, then execution proceeds to the
listening phase. If there are no module listeners, then the listening phase immediately terminates
successfully. Otherwise, the listening phase initializes the module listeners.

A service declaration is the syntactic sugar for creating a service object and attaching it to the module
listener specified in the service declaration.`

var runExamples = []templates.Example{
	{
		Description: "Run the current package.",
		Commands:    []string{"bal run"},
	},
	{
		Description: "Run a single Ballerina source file.",
		Commands:    []string{"bal run main.bal"},
	},
	{
		Description: "Run with program arguments.",
		Commands:    []string{"bal run main.bal -- arg1 arg2"},
	},
}

// runOpts holds the command-line options for the run command
var runOpts struct {
	dumpTokens    bool
	dumpST        bool
	dumpAST       bool
	dumpBIR       bool
	traceRecovery bool
	logFile       string
}

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "run [OPTIONS] [<package>|<source-file>] [-- <args...> <(-Ckey=value)...>]",
		Short:         "Compile and run the current package",
		Long:          runLongDesc,
		Example:       templates.FormatExamples(runExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          validateSourceFile,
		RunE:          runRun,
	}

	addRunFlags(cmd)

	return cmd
}

func addRunFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&runOpts.dumpTokens, "dump-tokens", false, "Dump lexer tokens")
	cmd.Flags().BoolVar(&runOpts.dumpST, "dump-st", false, "Dump syntax tree")
	cmd.Flags().BoolVar(&runOpts.dumpAST, "dump-ast", false, "Dump abstract syntax tree")
	cmd.Flags().BoolVar(&runOpts.dumpBIR, "dump-bir", false, "Dump Ballerina Intermediate Representation")
	cmd.Flags().BoolVar(&runOpts.traceRecovery, "trace-recovery", false, "Enable error recovery tracing")
	cmd.Flags().StringVar(&runOpts.logFile, "log-file", "", "Write debug output to specified file")
}

func runRun(cmd *cobra.Command, args []string) error {
	path := args[0]

	var debugCtx *debugcommon.DebugContext
	var wg sync.WaitGroup
	flags := uint16(0)

	if runOpts.dumpTokens {
		flags |= debugcommon.DUMP_TOKENS
	}
	if runOpts.dumpST {
		flags |= debugcommon.DUMP_ST
	}
	if runOpts.traceRecovery {
		flags |= debugcommon.DEBUG_ERROR_RECOVERY
	}

	if flags != 0 {
		debugcommon.Init(flags)
		debugCtx = &debugcommon.DebugCtx

		var logWriter *os.File
		var err error
		if runOpts.logFile != "" {
			logWriter, err = os.Create(runOpts.logFile)
			if err != nil {
				cmdErr := fmt.Errorf("error creating log file %s: %w", runOpts.logFile, err)
				printError(cmdErr, "", false, cmd.Name())
				return cmdErr
			}
		} else {
			logWriter = os.Stderr
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if runOpts.logFile != "" {
				defer logWriter.Close()
			}
			for msg := range debugCtx.Channel {
				fmt.Fprintf(logWriter, "%s\n", msg)
			}
		}()

		// Ensure debug context cleanup on any exit path
		defer func() {
			if debugCtx != nil {
				close(debugCtx.Channel)
				wg.Wait()
			}
		}()
	}

	// Load the project using ProjectLoader
	project, err := projects.Load(path)
	if err != nil {
		printError(fmt.Errorf("failed to load project: %w", err), "", false, cmd.Name())
		return fmt.Errorf("failed to load project: %w", err)
	}

	// Handle different project types
	switch project.Kind() {
	case projects.BuildProjectKind:
		return runBuildProject(cmd, project.(*projects.BuildProject), debugCtx)
	case projects.SingleFileProjectKind:
		return runSingleFileProject(cmd, project.(*projects.SingleFileProject), debugCtx)
	default:
		return fmt.Errorf("project type not yet supported: %v", project.Kind())
	}
}

func runBuildProject(cmd *cobra.Command, project *projects.BuildProject, debugCtx *debugcommon.DebugContext) error {
	pkg := project.CurrentPackage()

	// Compile the source
	fmt.Println("Compiling source")
	fmt.Printf("\t%s/%s:%s\n", pkg.PackageOrg(), pkg.PackageName(), pkg.PackageVersion())

	// Compile through Project API
	compilation, err := pkg.Compilation()
	if err != nil {
		printError(fmt.Errorf("compilation failed: %w", err), "", false, cmd.Name())
		return fmt.Errorf("compilation failed: %w", err)
	}

	// Check for errors via DiagnosticResult
	diagResult := compilation.DiagnosticResult()
	if diagResult.HasErrors() {
		for _, diag := range diagResult.Errors() {
			fmt.Fprintf(os.Stderr, "%s\n", diag.Error())
		}
		return fmt.Errorf("compilation failed with %d error(s)", diagResult.ErrorCount())
	}

	// Get compiled artifacts from module
	mod := pkg.DefaultModule()

	// Dump AST if requested
	if runOpts.dumpAST {
		astPkg := mod.BLangPackage()
		if astPkg != nil {
			prettyPrinter := ast.PrettyPrinter{}
			fmt.Println(prettyPrinter.Print(astPkg))
		}
	}

	// Get BIR for execution
	birPkg := mod.BIRPackage()
	if birPkg == nil {
		return fmt.Errorf("compilation did not produce BIR")
	}

	// Dump BIR if requested
	if runOpts.dumpBIR {
		prettyPrinter := bir.PrettyPrinter{}
		fmt.Println("==================BEGIN BIR==================")
		fmt.Println(strings.TrimSpace(prettyPrinter.Print(*birPkg)))
		fmt.Println("===================END BIR===================")
	}

	// Run the executable
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Running executable")
	fmt.Fprintln(os.Stderr)

	rt := runtime.NewRuntime()
	if err := rt.Interpret(*birPkg); err != nil {
		return err
	}
	return nil
}

func runSingleFileProject(cmd *cobra.Command, project *projects.SingleFileProject, debugCtx *debugcommon.DebugContext) error {
	pkg := project.CurrentPackage()

	// Compile the source
	fmt.Println("Compiling source")
	fmt.Printf("\t%s\n", filepath.Base(project.FilePath()))

	// Compile through Project API
	compilation, err := pkg.Compilation()
	if err != nil {
		printError(fmt.Errorf("compilation failed: %w", err), "", false, cmd.Name())
		return fmt.Errorf("compilation failed: %w", err)
	}

	// Check for errors via DiagnosticResult
	diagResult := compilation.DiagnosticResult()
	if diagResult.HasErrors() {
		for _, diag := range diagResult.Errors() {
			fmt.Fprintf(os.Stderr, "%s\n", diag.Error())
		}
		return fmt.Errorf("compilation failed with %d error(s)", diagResult.ErrorCount())
	}

	// Get compiled artifacts from module
	mod := pkg.DefaultModule()

	// Dump AST if requested
	if runOpts.dumpAST {
		astPkg := mod.BLangPackage()
		if astPkg != nil {
			prettyPrinter := ast.PrettyPrinter{}
			fmt.Println(prettyPrinter.Print(astPkg))
		}
	}

	// Get BIR for execution
	birPkg := mod.BIRPackage()
	if birPkg == nil {
		return fmt.Errorf("compilation did not produce BIR")
	}

	// Dump BIR if requested
	if runOpts.dumpBIR {
		prettyPrinter := bir.PrettyPrinter{}
		fmt.Println("==================BEGIN BIR==================")
		fmt.Println(strings.TrimSpace(prettyPrinter.Print(*birPkg)))
		fmt.Println("===================END BIR===================")
	}

	// Run the executable
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Running executable")
	fmt.Fprintln(os.Stderr)

	rt := runtime.NewRuntime()
	if err := rt.Interpret(*birPkg); err != nil {
		return err
	}
	return nil
}
