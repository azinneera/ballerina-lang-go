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
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ballerina-nutcracker/ballerina/projects"
)

// executeAddCommandWithArgs runs a fresh 'add' command instance with the
// given args, capturing stdout/stderr. Callers must t.Chdir into the
// target package directory first, since 'bal add' operates on the cwd.
func executeAddCommandWithArgs(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	cmd := createAddCmd()
	var outBuf, errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetArgs(args)
	err = cmd.Execute()
	return outBuf.String(), errBuf.String(), err
}

func writeBallerinaToml(t *testing.T, dir string) {
	t.Helper()
	content := "[package]\norg = \"testorg\"\nname = \"testpkg\"\nversion = \"0.1.0\"\n"
	if err := os.WriteFile(filepath.Join(dir, projects.BallerinaTomlFile), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write Ballerina.toml: %v", err)
	}
}

func TestAddCommand_Success(t *testing.T) {
	tests := []struct {
		template   string
		sourceFile string
		contains   string
	}{
		{"lib", "util.bal", "public function hello"},
		{"service", "svc.bal", "service / on new http:Listener"},
	}

	for _, tc := range tests {
		t.Run(tc.template, func(t *testing.T) {
			dir := t.TempDir()
			writeBallerinaToml(t, dir)
			t.Chdir(dir)

			moduleName := strings.TrimSuffix(tc.sourceFile, ".bal")
			stdout, stderr, err := executeAddCommandWithArgs(t, moduleName, "-t", tc.template)
			if err != nil {
				t.Fatalf("command failed: %v\nstderr: %s", err, stderr)
			}

			wantMsg := "Added new Ballerina module at " + filepath.Join("modules", moduleName)
			if !strings.Contains(stdout, wantMsg) {
				t.Errorf("stdout = %q, want to contain %q", stdout, wantMsg)
			}

			modulePath := filepath.Join(dir, "modules", moduleName)
			sourceContent, err := os.ReadFile(filepath.Join(modulePath, tc.sourceFile))
			if err != nil {
				t.Fatalf("expected module source file to exist: %v", err)
			}
			if !strings.Contains(string(sourceContent), tc.contains) {
				t.Errorf("source file missing expected content %q:\n%s", tc.contains, sourceContent)
			}
		})
	}
}

func TestAddCommand_DefaultTemplateIsLib(t *testing.T) {
	dir := t.TempDir()
	writeBallerinaToml(t, dir)
	t.Chdir(dir)

	if _, stderr, err := executeAddCommandWithArgs(t, "util"); err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr)
	}

	if _, err := os.Stat(filepath.Join(dir, "modules", "util", "util.bal")); err != nil {
		t.Errorf("expected default (lib) template to produce util.bal: %v", err)
	}
}

func TestAddCommand_NotAPackage(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	_, stderr, err := executeAddCommandWithArgs(t, "util")
	if err == nil {
		t.Fatal("expected an error when not run inside a Ballerina package")
	}
	if !strings.Contains(stderr, "not a Ballerina project") {
		t.Errorf("stderr = %q, want to contain 'not a Ballerina project'", stderr)
	}
}

func TestAddCommand_DuplicateModule(t *testing.T) {
	dir := t.TempDir()
	writeBallerinaToml(t, dir)
	t.Chdir(dir)

	if _, stderr, err := executeAddCommandWithArgs(t, "util"); err != nil {
		t.Fatalf("first add failed: %v\nstderr: %s", err, stderr)
	}
	modulePath := filepath.Join(dir, "modules", "util")
	original, err := os.ReadFile(filepath.Join(modulePath, "util.bal"))
	if err != nil {
		t.Fatalf("failed to read the module created by the first add: %v", err)
	}

	_, stderr, err := executeAddCommandWithArgs(t, "util")
	if err == nil {
		t.Fatal("expected an error when re-adding an existing module")
	}
	if !strings.Contains(stderr, "a module already exists with the given name") {
		t.Errorf("stderr = %q, want 'a module already exists' message", stderr)
	}

	// The rejected second add must leave the existing module untouched —
	// createModule's atomic os.Mkdir claim fails before ever writing or
	// removing anything under modulePath.
	after, err := os.ReadFile(filepath.Join(modulePath, "util.bal"))
	if err != nil {
		t.Fatalf("expected the existing module to remain intact: %v", err)
	}
	if string(after) != string(original) {
		t.Errorf("existing module content changed:\nbefore:\n%s\nafter:\n%s", original, after)
	}
}

func TestAddCommand_InvalidModuleName(t *testing.T) {
	tests := []struct {
		name     string
		wantMsg  string
		moduleID string
	}{
		{"invalid characters", "can only contain alphanumerics", "bad-name!"},
		{"initial underscore", "cannot have initial underscore", "_bad"},
		{"trailing underscore", "cannot have trailing underscore", "bad_"},
		{"consecutive underscore", "cannot have consecutive underscore", "ba__d"},
		{"initial digit", "cannot have initial numeric", "1bad"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			writeBallerinaToml(t, dir)
			t.Chdir(dir)

			_, stderr, err := executeAddCommandWithArgs(t, tc.moduleID)
			if err == nil {
				t.Fatalf("expected an error for invalid module name %q", tc.moduleID)
			}
			if !strings.Contains(stderr, tc.wantMsg) {
				t.Errorf("stderr = %q, want to contain %q", stderr, tc.wantMsg)
			}
		})
	}
}

func TestAddCommand_InvalidTemplate(t *testing.T) {
	dir := t.TempDir()
	writeBallerinaToml(t, dir)
	t.Chdir(dir)

	_, stderr, err := executeAddCommandWithArgs(t, "util", "-t", "bogus")
	if err == nil {
		t.Fatal("expected an error for an unsupported template")
	}
	if !strings.Contains(stderr, "unsupported template provided") {
		t.Errorf("stderr = %q, want 'unsupported template provided' message", stderr)
	}
}

// TestCreateModule_CleansUpOnWriteFailure covers createModule's cleanup
// path: runAdd's own pre-check always guarantees modulePath doesn't exist
// yet, so this calls createModule directly (white-box) with a
// pre-existing, write-protected module directory to force the
// os.WriteFile failure after MkdirAll's no-op success on an already-existing
// dir, and asserts the directory is removed rather than left behind
// half-written.
func TestCreateModule_CleansUpOnWriteFailure(t *testing.T) {
	dir := t.TempDir()
	modulePath := filepath.Join(dir, "util")

	// os.Mkdir(modulePath, ...) succeeds since modulePath doesn't exist yet;
	// passing a moduleName containing a path separator makes the later
	// os.WriteFile target a path whose intermediate directory doesn't
	// exist, forcing the write itself (not the atomic directory claim) to
	// fail — portably, with no permission tricks needed.
	err := createModule(modulePath, "nonexistent/util", "public function hello() returns string {\n}\n")
	if err == nil {
		t.Fatal("expected an error writing to a path with a missing intermediate directory")
	}
	if !strings.Contains(err.Error(), "failed to create") {
		t.Errorf("err = %q, want 'failed to create' message", err)
	}
	if _, statErr := os.Stat(modulePath); !os.IsNotExist(statErr) {
		t.Errorf("expected module directory to be cleaned up, stat err = %v", statErr)
	}
}

// TestCreateModule_DuplicateModule covers createModule's atomic-claim
// branch directly: os.Mkdir on a pre-existing modulePath must fail with
// fs.ErrExist (mapped by runAdd into the duplicate-module diagnostic)
// without touching anything already inside it.
func TestCreateModule_DuplicateModule(t *testing.T) {
	dir := t.TempDir()
	modulePath := filepath.Join(dir, "util")
	if err := os.MkdirAll(modulePath, 0755); err != nil {
		t.Fatalf("failed to pre-create module dir: %v", err)
	}
	existing := filepath.Join(modulePath, "util.bal")
	if err := os.WriteFile(existing, []byte("original"), 0644); err != nil {
		t.Fatalf("failed to write pre-existing util.bal: %v", err)
	}

	err := createModule(modulePath, "util", "public function hello() returns string {\n}\n")
	if !errors.Is(err, fs.ErrExist) {
		t.Fatalf("err = %v, want errors.Is(err, fs.ErrExist)", err)
	}

	content, readErr := os.ReadFile(existing)
	if readErr != nil {
		t.Fatalf("expected the pre-existing module to remain intact: %v", readErr)
	}
	if string(content) != "original" {
		t.Errorf("pre-existing util.bal content changed: %q", content)
	}
}

func TestAddCommand_ArgCount(t *testing.T) {
	dir := t.TempDir()
	writeBallerinaToml(t, dir)
	t.Chdir(dir)

	if _, stderr, err := executeAddCommandWithArgs(t); err == nil {
		t.Error("expected an error with no module name provided")
	} else if !strings.Contains(stderr, "module name is not provided") {
		t.Errorf("stderr = %q, want 'module name is not provided' message", stderr)
	}

	if _, stderr, err := executeAddCommandWithArgs(t, "a", "b"); err == nil {
		t.Error("expected an error with too many arguments")
	} else if !strings.Contains(stderr, "too many arguments") {
		t.Errorf("stderr = %q, want 'too many arguments' message", stderr)
	}
}
