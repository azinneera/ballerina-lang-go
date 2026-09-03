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

package testdiscovery_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/ballerina-nutcracker/ballerina/cli/internal/testdiscovery"
	"github.com/ballerina-nutcracker/ballerina/projects"
)

func loadSampleProject(t *testing.T) projects.Project {
	t.Helper()

	absPath, err := filepath.Abs(filepath.Join("testdata", "sample"))
	if err != nil {
		t.Fatalf("resolve testdata path: %v", err)
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("resolve user home: %v", err)
	}
	ballerinaEnvPath := os.Getenv(projects.BallerinaEnvVar)
	if ballerinaEnvPath == "" {
		ballerinaEnvPath = filepath.Join(userHome, projects.UserHomeDirName)
	}

	// SkipTests must be false — the whole point of this test is to confirm
	// test-source functions and their @test:* annotations reach the compiled
	// AST at all, which only happens when tests are included in compilation
	// (see projects/module_context.go's resolveTypesAndSymbols).
	buildOpts := projects.NewBuildOptionsBuilder().
		WithOffline(true).
		WithSkipTests(false).
		Build()

	result, err := projects.Load(os.DirFS(absPath), ".", projects.ProjectLoadConfig{
		BallerinaEnvFs: os.DirFS(ballerinaEnvPath),
		BuildOptions:   &buildOpts,
	})
	if err != nil {
		t.Fatalf("load project: %v", err)
	}
	if diag := result.Diagnostics(); diag.HasErrors() {
		t.Fatalf("project has diagnostics errors: %+v", diag)
	}
	return result.Project()
}

func TestDiscoverFindsAllAnnotatedFunctions(t *testing.T) {
	project := loadSampleProject(t)
	pkg := project.CurrentPackage()
	compilation := pkg.Compilation()
	defaultModule := pkg.DefaultModule()

	if diag := compilation.DiagnosticResult(); diag.HasErrors() {
		var messages []string
		for _, d := range diag.Diagnostics() {
			messages = append(messages, d.Message())
		}
		t.Fatalf("unexpected compilation errors: %v", messages)
	}

	astPkg := compilation.ModuleAST(defaultModule.ModuleID())
	if astPkg == nil {
		t.Fatal("ModuleAST returned nil — module didn't compile")
	}

	discovered := testdiscovery.Discover(astPkg)

	if len(discovered.Tests) != 2 {
		t.Fatalf("expected 2 @test:Config functions, got %d: %+v", len(discovered.Tests), discovered.Tests)
	}

	var testMain, testDisabled *testdiscovery.TestFunc
	for i := range discovered.Tests {
		switch discovered.Tests[i].Name {
		case "testMain":
			testMain = &discovered.Tests[i]
		case "testDisabled":
			testDisabled = &discovered.Tests[i]
		}
	}
	if testMain == nil {
		t.Fatal("testMain not discovered")
	}
	if !testMain.Enable {
		t.Error("testMain: expected Enable=true (default)")
	}
	if !slices.Equal(testMain.Groups, []string{"unit", "fast"}) {
		t.Errorf("testMain: expected Groups=[unit fast], got %v", testMain.Groups)
	}
	if !slices.Equal(testMain.DependsOn, []string{"testSetup"}) {
		t.Errorf("testMain: expected DependsOn=[testSetup], got %v", testMain.DependsOn)
	}

	if testDisabled == nil {
		t.Fatal("testDisabled not discovered")
	}
	if testDisabled.Enable {
		t.Error("testDisabled: expected Enable=false")
	}

	if len(discovered.BeforeSuite) != 1 || discovered.BeforeSuite[0].Name != "beforeAll" {
		t.Errorf("expected BeforeSuite=[beforeAll], got %+v", discovered.BeforeSuite)
	}
	if len(discovered.AfterSuite) != 1 || discovered.AfterSuite[0].Name != "afterAll" || !discovered.AfterSuite[0].AlwaysRun {
		t.Errorf("expected AfterSuite=[{afterAll true}], got %+v", discovered.AfterSuite)
	}
	if len(discovered.BeforeEach) != 1 || discovered.BeforeEach[0].Name != "beforeEachTest" {
		t.Errorf("expected BeforeEach=[beforeEachTest], got %+v", discovered.BeforeEach)
	}
	if len(discovered.AfterEach) != 1 || discovered.AfterEach[0].Name != "afterEachTest" {
		t.Errorf("expected AfterEach=[afterEachTest], got %+v", discovered.AfterEach)
	}

	if len(discovered.BeforeGroups) != 1 {
		t.Fatalf("expected 1 BeforeGroups hook, got %d", len(discovered.BeforeGroups))
	}
	bg := discovered.BeforeGroups[0]
	if bg.Name != "beforeUnitGroup" || !slices.Equal(bg.Groups, []string{"unit"}) {
		t.Errorf("expected BeforeGroups=[{beforeUnitGroup [unit]}], got %+v", bg)
	}

	if len(discovered.AfterGroups) != 1 {
		t.Fatalf("expected 1 AfterGroups hook, got %d", len(discovered.AfterGroups))
	}
	ag := discovered.AfterGroups[0]
	if ag.Name != "afterUnitGroup" || !slices.Equal(ag.Groups, []string{"unit"}) || !ag.AlwaysRun {
		t.Errorf("expected AfterGroups=[{afterUnitGroup [unit] true}], got %+v", ag)
	}
}

func TestDiscoverWithoutTestImportReturnsEmpty(t *testing.T) {
	// A package with no `import ballerina/test` at all (e.g. one that hasn't
	// been compiled with SkipTests=false) must not panic — Discover should
	// just find nothing.
	discovered := testdiscovery.Discover(nil)
	if discovered == nil {
		t.Fatal("Discover(nil) returned nil, expected an empty *Discovered")
	}
	if len(discovered.Tests) != 0 {
		t.Errorf("expected no tests, got %+v", discovered.Tests)
	}
}
