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

package testglue_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ballerina-nutcracker/ballerina/cli/internal/testglue"
	"github.com/ballerina-nutcracker/ballerina/cli/internal/testplan"
)

func samplePlan() *testplan.TestPlan {
	return &testplan.TestPlan{
		Tests: []testplan.TestFunction{
			{Name: "testMain", Enable: true, Groups: []string{"unit", "fast"}, DependsOn: []string{"testSetup"}},
			{Name: "testDisabled", Enable: false},
		},
		BeforeSuite:  []testplan.SuiteHook{{Name: "beforeAll"}},
		AfterSuite:   []testplan.AfterSuiteHook{{Name: "afterAll", AlwaysRun: true}},
		BeforeEach:   []testplan.SuiteHook{{Name: "beforeEachTest"}},
		AfterEach:    []testplan.SuiteHook{{Name: "afterEachTest"}},
		BeforeGroups: []testplan.GroupHook{{Name: "beforeUnitGroup", Groups: []string{"unit"}}},
		AfterGroups:  []testplan.GroupHook{{Name: "afterUnitGroup", Groups: []string{"unit"}, AlwaysRun: true}},
	}
}

func TestGenerateSource_ContainsExpectedCalls(t *testing.T) {
	src := testglue.GenerateSource(samplePlan(), nil)

	for _, want := range []string{
		`import ballerina/test;`,
		`public function __test_register__() {`,
		`test:registerTestConfig("testMain", testMain, true, ["unit", "fast"], [testSetup], (), (), false, ());`,
		`test:registerTestConfig("testDisabled", testDisabled, false, [], [], (), (), false, ());`,
		`test:registerBeforeSuite("beforeAll", beforeAll);`,
		`test:registerAfterSuite("afterAll", afterAll, true);`,
		`test:registerBeforeEach("beforeEachTest", beforeEachTest);`,
		`test:registerAfterEach("afterEachTest", afterEachTest);`,
		`test:registerBeforeGroups("beforeUnitGroup", beforeUnitGroup, ["unit"]);`,
		`test:registerAfterGroups("afterUnitGroup", afterUnitGroup, ["unit"], true);`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated source missing expected call:\n  %s\nfull source:\n%s", want, src)
		}
	}
}

// TestGenerateSource_IncludedFiltersTests confirms the new (2026-09-18)
// included parameter, computed by cli/internal/testfilter, controls which
// tests actually get registered: an excluded test's registerTestConfig call
// is never emitted at all, so ballerina/test's own registry never even sees
// it — a stronger guarantee than filtering at execution time. Hooks are
// never filtered, since they're global infrastructure.
func TestGenerateSource_IncludedFiltersTests(t *testing.T) {
	src := testglue.GenerateSource(samplePlan(), map[string]bool{"testMain": true})

	if !strings.Contains(src, `test:registerTestConfig("testMain", testMain,`) {
		t.Errorf("expected testMain (included) to be registered:\n%s", src)
	}
	if strings.Contains(src, `testDisabled`) {
		t.Errorf("expected testDisabled (not in included) to be entirely absent:\n%s", src)
	}
	if !strings.Contains(src, `test:registerBeforeSuite("beforeAll", beforeAll);`) {
		t.Errorf("expected hooks to still be registered regardless of included:\n%s", src)
	}
}

// TestGenerateSource_EmptyIncludedStillCompiles covers every test being
// filtered out (an empty, non-nil included map) — the anchor-import fallback
// must still trigger, the same as a plan with no tests at all, since an
// empty map means "nothing survived filtering," not "no filtering."
func TestGenerateSource_EmptyIncludedStillCompiles(t *testing.T) {
	plan := &testplan.TestPlan{Tests: []testplan.TestFunction{{Name: "testMain", Enable: true}}}
	src := testglue.GenerateSource(plan, map[string]bool{})

	if strings.Contains(src, "registerTestConfig") {
		t.Errorf("expected no registerTestConfig calls when included is empty:\n%s", src)
	}
	if !strings.Contains(src, "var _anchor = test:startSuite;") {
		t.Errorf("expected the anchor fallback when nothing survived filtering:\n%s", src)
	}
}

func TestGenerateSource_StringEscaping(t *testing.T) {
	plan := &testplan.TestPlan{
		Tests: []testplan.TestFunction{
			{Name: `weird"name`, Groups: []string{`group\with"quotes`}},
		},
	}
	src := testglue.GenerateSource(plan, nil)
	if !strings.Contains(src, `"weird\"name"`) {
		t.Errorf("expected escaped test name in source:\n%s", src)
	}
	if !strings.Contains(src, `"group\\with\"quotes"`) {
		t.Errorf("expected escaped group name in source:\n%s", src)
	}
}

// TestGenerateSource_EmptyDiscoveryStillCompiles covers a package with no
// @test:* functions at all (e.g. no test files) — confirmed via a real `bal
// test` run that an unconditional `import ballerina/test;` over an otherwise
// empty function body is a compile error (unused import).
func TestGenerateSource_EmptyDiscoveryStillCompiles(t *testing.T) {
	src := testglue.GenerateSource(&testplan.TestPlan{}, nil)
	if !strings.Contains(src, `import ballerina/test;`) {
		t.Fatalf("expected the import to still be present so ballerina/test stays in the compiled dependency graph:\n%s", src)
	}

	dir := t.TempDir()
	balFile := filepath.Join(dir, "generated_glue_empty.bal")
	src += "\npublic function main() {\n    __test_register__();\n}\n"
	if err := os.WriteFile(balFile, []byte(src), 0o644); err != nil {
		t.Fatalf("write generated source: %v", err)
	}

	cliCmdDir, err := filepath.Abs(filepath.Join("..", "..", "cmd"))
	if err != nil {
		t.Fatalf("resolve cli/cmd path: %v", err)
	}
	cmd := exec.Command("go", "run", cliCmdDir, "run", balFile)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("empty-discovery generated source failed to compile/run: %v\noutput:\n%s\nsource:\n%s", err, out, src)
	}
}

// TestGenerateSource_CompilesAndRuns exercises the generated source through
// the real compiler/interpreter: appends it to sample function definitions,
// runs it via `bal run`, and confirms the whole thing compiles and executes
// without error. This is the strongest available check that registration_api.bal's
// call shapes actually match what GenerateSource emits.
func TestGenerateSource_CompilesAndRuns(t *testing.T) {
	src := testglue.GenerateSource(samplePlan(), nil)

	src += `
function testMain() {
}

function testDisabled() {
}

function testSetup() {
}

function beforeAll() {
}

function afterAll() {
}

function beforeEachTest() {
}

function afterEachTest() {
}

function beforeUnitGroup() {
}

function afterUnitGroup() {
}

public function main() {
    __test_register__();
}
`

	dir := t.TempDir()
	balFile := filepath.Join(dir, "generated_glue.bal")
	if err := os.WriteFile(balFile, []byte(src), 0o644); err != nil {
		t.Fatalf("write generated source: %v", err)
	}

	cliCmdDir, err := filepath.Abs(filepath.Join("..", "..", "cmd"))
	if err != nil {
		t.Fatalf("resolve cli/cmd path: %v", err)
	}

	cmd := exec.Command("go", "run", cliCmdDir, "run", balFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generated source failed to compile/run: %v\noutput:\n%s\nsource:\n%s", err, out, src)
	}
}
