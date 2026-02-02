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

package projects

import (
	"path/filepath"
	"testing"
)

func TestCompileBuildProject(t *testing.T) {
	projectPath := filepath.Join("testdata", "valid_package")

	project, err := Load(projectPath)
	if err != nil {
		t.Fatalf("failed to load build project: %v", err)
	}

	pkg := project.CurrentPackage()

	// Test compilation
	compilation, err := pkg.Compilation()
	if err != nil {
		t.Fatalf("compilation failed: %v", err)
	}

	if compilation == nil {
		t.Fatal("Compilation() returned nil")
	}

	// Test DiagnosticResult
	diagResult := compilation.DiagnosticResult()
	if diagResult == nil {
		t.Fatal("DiagnosticResult() returned nil")
	}

	// Should have no errors for valid package
	if diagResult.HasErrors() {
		t.Errorf("expected no errors, got %d errors", diagResult.ErrorCount())
		for _, diag := range diagResult.Errors() {
			t.Logf("  error: %s", diag.Message)
		}
	}

	// Test BIR generation
	mod := pkg.DefaultModule()
	birPkg := mod.BIRPackage()
	if birPkg == nil {
		t.Fatal("BIRPackage() returned nil after compilation")
	}

	// Test BLangPackage
	astPkg := mod.BLangPackage()
	if astPkg == nil {
		t.Fatal("BLangPackage() returned nil after compilation")
	}
}

func TestCompileSingleFileProject(t *testing.T) {
	filePath := filepath.Join("testdata", "single_file.bal")

	project, err := Load(filePath)
	if err != nil {
		t.Fatalf("failed to load single file project: %v", err)
	}

	pkg := project.CurrentPackage()

	// Test compilation
	compilation, err := pkg.Compilation()
	if err != nil {
		t.Fatalf("compilation failed: %v", err)
	}

	if compilation == nil {
		t.Fatal("Compilation() returned nil")
	}

	// Test DiagnosticResult
	diagResult := compilation.DiagnosticResult()
	if diagResult.HasErrors() {
		t.Errorf("expected no errors, got %d errors", diagResult.ErrorCount())
	}

	// Test BIR generation
	mod := pkg.DefaultModule()
	birPkg := mod.BIRPackage()
	if birPkg == nil {
		t.Fatal("BIRPackage() returned nil after compilation")
	}
}

func TestCompilationCaching(t *testing.T) {
	projectPath := filepath.Join("testdata", "valid_package")

	project, err := Load(projectPath)
	if err != nil {
		t.Fatalf("failed to load build project: %v", err)
	}

	pkg := project.CurrentPackage()

	// First compilation
	compilation1, err := pkg.Compilation()
	if err != nil {
		t.Fatalf("first compilation failed: %v", err)
	}

	// Second compilation - should return same instance
	compilation2, err := pkg.Compilation()
	if err != nil {
		t.Fatalf("second compilation failed: %v", err)
	}

	// Verify same instance (caching)
	if compilation1 != compilation2 {
		t.Error("expected Compilation() to return cached instance")
	}
}

func TestDiagnosticResult(t *testing.T) {
	// Create diagnostics with different severities
	diagnostics := []Diagnostic{
		{Severity: SeverityError, Message: "error 1"},
		{Severity: SeverityError, Message: "error 2"},
		{Severity: SeverityWarning, Message: "warning 1"},
		{Severity: SeverityHint, Message: "hint 1"},
		{Severity: SeverityHint, Message: "hint 2"},
		{Severity: SeverityHint, Message: "hint 3"},
	}

	result := NewDiagnosticResult(diagnostics)

	// Test counts
	if result.DiagnosticCount() != 6 {
		t.Errorf("expected DiagnosticCount() = 6, got %d", result.DiagnosticCount())
	}

	if result.ErrorCount() != 2 {
		t.Errorf("expected ErrorCount() = 2, got %d", result.ErrorCount())
	}

	if result.WarningCount() != 1 {
		t.Errorf("expected WarningCount() = 1, got %d", result.WarningCount())
	}

	if result.HintCount() != 3 {
		t.Errorf("expected HintCount() = 3, got %d", result.HintCount())
	}

	// Test boolean checks
	if !result.HasErrors() {
		t.Error("expected HasErrors() = true")
	}

	if !result.HasWarnings() {
		t.Error("expected HasWarnings() = true")
	}

	// Test filtered collections
	errors := result.Errors()
	if len(errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errors))
	}

	warnings := result.Warnings()
	if len(warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(warnings))
	}

	hints := result.Hints()
	if len(hints) != 3 {
		t.Errorf("expected 3 hints, got %d", len(hints))
	}
}

func TestDiagnosticResultEmpty(t *testing.T) {
	result := NewDiagnosticResult([]Diagnostic{})

	if result.DiagnosticCount() != 0 {
		t.Errorf("expected DiagnosticCount() = 0, got %d", result.DiagnosticCount())
	}

	if result.HasErrors() {
		t.Error("expected HasErrors() = false for empty result")
	}

	if result.HasWarnings() {
		t.Error("expected HasWarnings() = false for empty result")
	}
}

func TestDiagnosticResultCaching(t *testing.T) {
	diagnostics := []Diagnostic{
		{Severity: SeverityError, Message: "error 1"},
		{Severity: SeverityWarning, Message: "warning 1"},
	}

	result := NewDiagnosticResult(diagnostics)

	// First call computes and caches
	errors1 := result.Errors()

	// Second call should return cached
	errors2 := result.Errors()

	// Should be same slice (cached)
	if &errors1[0] != &errors2[0] {
		t.Error("expected Errors() to return cached slice")
	}
}
