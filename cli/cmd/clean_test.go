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
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// executeCleanCommandWithArgs runs a fresh 'clean' command instance with the
// given args, capturing stdout/stderr. Callers must t.Chdir into the
// target directory first, since 'bal clean' operates on the cwd (unless
// --target-dir is given).
func executeCleanCommandWithArgs(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	cmd := createCleanCmd()
	var outBuf, errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetArgs(args)
	err = cmd.Execute()
	return outBuf.String(), errBuf.String(), err
}

func TestCleanCommand_DeletesExistingTarget(t *testing.T) {
	dir := t.TempDir()
	writeBallerinaToml(t, dir)
	targetDir := filepath.Join(dir, "target")
	if err := os.MkdirAll(filepath.Join(targetDir, "cache"), 0755); err != nil {
		t.Fatalf("failed to create target/cache: %v", err)
	}
	t.Chdir(dir)

	stdout, stderr, err := executeCleanCommandWithArgs(t)
	if err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Successfully deleted") {
		t.Errorf("stdout = %q, want to contain 'Successfully deleted'", stdout)
	}
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		t.Errorf("expected target/ to be deleted, stat err = %v", err)
	}
}

// TestCleanCommand_WithPathArgument covers `bal clean <package-dir>` run
// from a DIFFERENT cwd than the package itself — regression test for a bug
// where the target directory (resolved from the loaded project's own
// sourceRoot-relative TargetDir(), always "." internally) was deleted
// relative to the process cwd instead of the given package directory,
// silently no-op'ing instead of actually deleting anything.
func TestCleanCommand_WithPathArgument(t *testing.T) {
	parent := t.TempDir()
	pkgDir := filepath.Join(parent, "pkg")
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		t.Fatalf("failed to create pkg dir: %v", err)
	}
	writeBallerinaToml(t, pkgDir)
	targetDir := filepath.Join(pkgDir, "target")
	if err := os.MkdirAll(filepath.Join(targetDir, "bala"), 0755); err != nil {
		t.Fatalf("failed to create target/bala: %v", err)
	}
	t.Chdir(parent)

	stdout, stderr, err := executeCleanCommandWithArgs(t, "pkg")
	if err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Successfully deleted") {
		t.Errorf("stdout = %q, want to contain 'Successfully deleted'", stdout)
	}
	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		t.Errorf("expected pkg/target to be deleted, stat err = %v", err)
	}
}

func TestCleanCommand_NoOpWhenTargetAbsent(t *testing.T) {
	dir := t.TempDir()
	writeBallerinaToml(t, dir)
	t.Chdir(dir)

	stdout, stderr, err := executeCleanCommandWithArgs(t)
	if err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr)
	}
	if stdout != "" {
		t.Errorf("stdout = %q, want empty (silent no-op when target/ is already absent)", stdout)
	}
}

func TestCleanCommand_InvalidProjectDirectory(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	_, stderr, err := executeCleanCommandWithArgs(t)
	if err == nil {
		t.Fatal("expected an error for a directory with no Ballerina.toml")
	}
	if !strings.Contains(stderr, "invalid project directory") {
		t.Errorf("stderr = %q, want 'invalid project directory' message", stderr)
	}
}

func TestCleanCommand_Workspace(t *testing.T) {
	dir := t.TempDir()
	// Workspace member packages must be listed explicitly — there's no
	// auto-discovery at load time (only bal new --workspace does that).
	workspaceToml := "[workspace]\npackages = [\"pkg-a\", \"pkg-b\"]\n"
	if err := os.WriteFile(filepath.Join(dir, "Ballerina.toml"), []byte(workspaceToml), 0644); err != nil {
		t.Fatalf("failed to write workspace Ballerina.toml: %v", err)
	}

	for _, name := range []string{"pkg-a", "pkg-b"} {
		pkgDir := filepath.Join(dir, name)
		if err := os.MkdirAll(pkgDir, 0755); err != nil {
			t.Fatalf("failed to create %s: %v", name, err)
		}
		content := "[package]\norg = \"testorg\"\nname = \"" + name + "\"\nversion = \"0.1.0\"\n"
		if err := os.WriteFile(filepath.Join(pkgDir, "Ballerina.toml"), []byte(content), 0644); err != nil {
			t.Fatalf("failed to write %s/Ballerina.toml: %v", name, err)
		}
		if err := os.WriteFile(filepath.Join(pkgDir, "main.bal"), []byte("public function main() {\n}\n"), 0644); err != nil {
			t.Fatalf("failed to write %s/main.bal: %v", name, err)
		}
		if err := os.MkdirAll(filepath.Join(pkgDir, "target", "cache"), 0755); err != nil {
			t.Fatalf("failed to create %s/target/cache: %v", name, err)
		}
	}
	if err := os.MkdirAll(filepath.Join(dir, "target", "cache"), 0755); err != nil {
		t.Fatalf("failed to create workspace target/cache: %v", err)
	}
	t.Chdir(dir)

	_, stderr, err := executeCleanCommandWithArgs(t)
	if err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr)
	}

	for _, name := range []string{"pkg-a", "pkg-b"} {
		if _, err := os.Stat(filepath.Join(dir, name, "target")); !os.IsNotExist(err) {
			t.Errorf("expected %s/target to be deleted, stat err = %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "target")); !os.IsNotExist(err) {
		t.Errorf("expected workspace target to be deleted, stat err = %v", err)
	}
}

func TestCleanCommand_TargetDirFlag(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, "cache"), 0755); err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}
		if err := os.MkdirAll(filepath.Join(dir, "bala"), 0755); err != nil {
			t.Fatalf("failed to create bala: %v", err)
		}

		stdout, stderr, err := executeCleanCommandWithArgs(t, "--target-dir", dir)
		if err != nil {
			t.Fatalf("command failed: %v\nstderr: %s", err, stderr)
		}
		if !strings.Contains(stdout, "Successfully deleted") {
			t.Errorf("stdout = %q, want 'Successfully deleted'", stdout)
		}
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Errorf("expected target dir to be deleted, stat err = %v", err)
		}
	})

	t.Run("valid: bala only (bal pack never creates a cache/ dir)", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, "bala"), 0755); err != nil {
			t.Fatalf("failed to create bala: %v", err)
		}

		stdout, stderr, err := executeCleanCommandWithArgs(t, "--target-dir", dir)
		if err != nil {
			t.Fatalf("command failed: %v\nstderr: %s", err, stderr)
		}
		if !strings.Contains(stdout, "Successfully deleted") {
			t.Errorf("stdout = %q, want 'Successfully deleted'", stdout)
		}
	})

	t.Run("invalid: missing recognizable subdir", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, "unrelated"), 0755); err != nil {
			t.Fatalf("failed to create unrelated: %v", err)
		}
		// No cache/bala/bin/apidocs -> not recognizable as a target dir.

		_, stderr, err := executeCleanCommandWithArgs(t, "--target-dir", dir)
		if err == nil {
			t.Fatal("expected an error for an unrecognizable target dir")
		}
		if !strings.Contains(stderr, "is not a valid target directory") {
			t.Errorf("stderr = %q, want 'is not a valid target directory'", stderr)
		}
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("expected dir to survive a rejected clean, stat err = %v", err)
		}
	})

	t.Run("does not exist", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "nope")
		_, stderr, err := executeCleanCommandWithArgs(t, "--target-dir", missing)
		if err == nil {
			t.Fatal("expected an error for a non-existent target dir")
		}
		if !strings.Contains(stderr, "does not exist") {
			t.Errorf("stderr = %q, want 'does not exist'", stderr)
		}
	})

	t.Run("not a directory", func(t *testing.T) {
		dir := t.TempDir()
		file := filepath.Join(dir, "afile")
		if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
		_, stderr, err := executeCleanCommandWithArgs(t, "--target-dir", file)
		if err == nil {
			t.Fatal("expected an error for a non-directory target path")
		}
		if !strings.Contains(stderr, "is not a directory") {
			t.Errorf("stderr = %q, want 'is not a directory'", stderr)
		}
	})
}
