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

//go:build !js && !wasm

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeTestFixture creates a minimal package with the given Ballerina.toml
// package name and test source content, returning its directory.
func writeTestFixture(t *testing.T, pkgName, testSource string) string {
	t.Helper()
	dir := t.TempDir()
	toml := "[package]\norg = \"testorg\"\nname = \"" + pkgName + "\"\nversion = \"0.1.0\"\n"
	if err := os.WriteFile(filepath.Join(dir, "Ballerina.toml"), []byte(toml), 0o644); err != nil {
		t.Fatalf("write Ballerina.toml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.bal"), []byte("public function main() {\n}\n"), 0o644); err != nil {
		t.Fatalf("write main.bal: %v", err)
	}
	if testSource != "" {
		testsDir := filepath.Join(dir, "tests")
		if err := os.MkdirAll(testsDir, 0o755); err != nil {
			t.Fatalf("mkdir tests: %v", err)
		}
		if err := os.WriteFile(filepath.Join(testsDir, "main_test.bal"), []byte(testSource), 0o644); err != nil {
			t.Fatalf("write test file: %v", err)
		}
	}
	return dir
}

// executeTestCommand runs the test command and returns cobra's own
// stdout/stderr buffers plus the real OS-level stdout the test suite's
// console report is written to. ballerina/test's println (native/test_io.go)
// writes straight to the process's real os.Stdout — the same execution-path
// separation bal run relies on for the programs it executes — so it never
// lands in cobra's SetOut/SetErr buffers; capturing real stdout is the only
// way to observe pass/fail/skip counts and messages from these tests.
//
// Not safe to run in parallel with anything else that reads/writes os.Stdout
// in this process — none of these tests call t.Parallel().
func executeTestCommand(t *testing.T, args ...string) (stdout, cobraStdout, cobraStderr string, err error) {
	t.Helper()

	realStdout := os.Stdout
	r, w, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatalf("create stdout pipe: %v", pipeErr)
	}
	os.Stdout = w

	cmd := createTestCmd()
	var outBuf, errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetArgs(args)
	err = cmd.Execute()

	os.Stdout = realStdout
	_ = w.Close()
	captured, _ := io.ReadAll(r)

	return string(captured), outBuf.String(), errBuf.String(), err
}

func TestTestCommand_AllPassing(t *testing.T) {
	dir := writeTestFixture(t, "allpassing", `import ballerina/test;

@test:Config {}
function testOne() {
    test:assertTrue(true);
}
`)
	stdout, _, stderr, err := executeTestCommand(t, dir)
	if err != nil {
		t.Fatalf("expected success, got error: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "1 passing") {
		t.Errorf("expected '1 passing' in output, got: %s", stdout)
	}
	// The package's own main() must never run during `bal test`.
	if strings.Contains(stdout, "MAIN_RAN") {
		t.Errorf("package main() ran unexpectedly during bal test: %s", stdout)
	}
}

func TestTestCommand_FailureExitsNonZeroAndWritesRerunFile(t *testing.T) {
	dir := writeTestFixture(t, "hasfailure", `import ballerina/test;

@test:Config {}
function testGood() {
    test:assertTrue(true);
}

@test:Config {}
function testBad() {
    test:assertTrue(false, "boom");
}
`)
	stdout, _, stderr, err := executeTestCommand(t, dir)
	if err == nil {
		t.Fatalf("expected a non-nil error for a failing suite, got success. stdout: %s", stdout)
	}
	if !strings.Contains(stdout, "1 passing") || !strings.Contains(stdout, "1 failing") {
		t.Errorf("expected one pass and one failure in output, got: %s\nstderr: %s", stdout, stderr)
	}

	rerunPath := filepath.Join(dir, "target", "rerun_test.json")
	data, readErr := os.ReadFile(rerunPath)
	if readErr != nil {
		t.Fatalf("expected rerun_test.json to be written: %v", readErr)
	}
	var rerun map[string]struct {
		TestNames []string `json:"testNames"`
	}
	if err := json.Unmarshal(data, &rerun); err != nil {
		t.Fatalf("unmarshal rerun_test.json: %v", err)
	}
	entry, ok := rerun["hasfailure"]
	if !ok || len(entry.TestNames) != 1 || entry.TestNames[0] != "testBad" {
		t.Errorf("expected rerun_test.json to list only testBad, got: %+v", rerun)
	}
}

func TestTestCommand_RerunFailedOnlyRunsPreviousFailures(t *testing.T) {
	dir := writeTestFixture(t, "rerunpkg", `import ballerina/test;

@test:Config {}
function testGood() {
    test:assertTrue(true);
}

@test:Config {}
function testBad() {
    test:assertTrue(false, "boom");
}
`)
	if _, _, _, err := executeTestCommand(t, dir); err == nil {
		t.Fatal("expected the first run to fail")
	}

	stdout, _, _, err := executeTestCommand(t, dir, "--rerun-failed")
	if err == nil {
		t.Fatalf("expected the rerun to still report the persistent failure. stdout: %s", stdout)
	}
	if !strings.Contains(stdout, "0 passing") || !strings.Contains(stdout, "1 failing") {
		t.Errorf("expected only the previously-failed test to run, got: %s", stdout)
	}
}

func TestTestCommand_NoTestsFound(t *testing.T) {
	dir := writeTestFixture(t, "notests", "")
	stdout, _, stderr, err := executeTestCommand(t, dir)
	if err != nil {
		t.Fatalf("expected success for a package with no tests, got: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "No tests found") {
		t.Errorf("expected 'No tests found', got: %s", stdout)
	}
}

func TestTestCommand_ListGroups(t *testing.T) {
	dir := writeTestFixture(t, "listgroups", `import ballerina/test;

@test:Config {
    groups: ["unit", "fast"]
}
function testOne() {
}
`)
	stdout, _, stderr, err := executeTestCommand(t, dir, "--list-groups")
	if err != nil {
		t.Fatalf("unexpected error: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "unit") || !strings.Contains(stdout, "fast") {
		t.Errorf("expected both groups listed, got: %s", stdout)
	}
}

func TestTestCommand_GroupFiltering(t *testing.T) {
	dir := writeTestFixture(t, "groupfilter", `import ballerina/test;

@test:Config {
    groups: ["unit"]
}
function testUnit() {
    test:assertTrue(true);
}

@test:Config {
    groups: ["integration"]
}
function testIntegration() {
    test:assertTrue(true);
}
`)
	stdout, _, stderr, err := executeTestCommand(t, dir, "--groups", "unit")
	if err != nil {
		t.Fatalf("unexpected error: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "1 passing") {
		t.Errorf("expected exactly one test (testUnit) to run, got: %s", stdout)
	}
}

func TestTestCommand_RejectsUnsupportedFlags(t *testing.T) {
	dir := writeTestFixture(t, "rejectflags", "")
	for _, args := range [][]string{
		{dir, "--graalvm"},
		{dir, "--cloud", "docker"},
	} {
		_, _, _, err := executeTestCommand(t, args...)
		if err == nil {
			t.Errorf("expected %v to be rejected with an error", args)
		}
	}
}

func TestTestCommand_StandaloneFileNotYetSupported(t *testing.T) {
	dir := t.TempDir()
	balFile := filepath.Join(dir, "standalone.bal")
	if err := os.WriteFile(balFile, []byte("public function main() {\n}\n"), 0o644); err != nil {
		t.Fatalf("write standalone.bal: %v", err)
	}
	_, _, _, err := executeTestCommand(t, balFile)
	if err == nil {
		t.Error("expected an error for a standalone .bal file (P8.10 not yet implemented)")
	}
}
