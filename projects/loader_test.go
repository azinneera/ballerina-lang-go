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

func TestLoadBuildProject(t *testing.T) {
	projectPath := filepath.Join("testdata", "valid_package")

	project, err := Load(projectPath)
	if err != nil {
		t.Fatalf("failed to load build project: %v", err)
	}

	if project.Kind() != BuildProjectKind {
		t.Errorf("expected BuildProjectKind, got %v", project.Kind())
	}

	pkg := project.CurrentPackage()
	if pkg == nil {
		t.Fatal("CurrentPackage() returned nil")
	}

	if pkg.PackageName() != "myapp" {
		t.Errorf("expected package name 'myapp', got '%s'", pkg.PackageName())
	}

	if pkg.PackageOrg() != "testorg" {
		t.Errorf("expected org 'testorg', got '%s'", pkg.PackageOrg())
	}

	if pkg.PackageVersion().String() != "0.1.0" {
		t.Errorf("expected version '0.1.0', got '%s'", pkg.PackageVersion())
	}

	mod := pkg.DefaultModule()
	if mod == nil {
		t.Fatal("DefaultModule() returned nil")
	}

	if !mod.IsDefault() {
		t.Error("expected default module")
	}

	docs := mod.Documents()
	if len(docs) == 0 {
		t.Error("expected at least one document")
	}
}

func TestLoadSingleFileProject(t *testing.T) {
	filePath := filepath.Join("testdata", "single_file.bal")

	project, err := Load(filePath)
	if err != nil {
		t.Fatalf("failed to load single file project: %v", err)
	}

	if project.Kind() != SingleFileProjectKind {
		t.Errorf("expected SingleFileProjectKind, got %v", project.Kind())
	}

	pkg := project.CurrentPackage()
	if pkg == nil {
		t.Fatal("CurrentPackage() returned nil")
	}

	if pkg.PackageName() != "single_file" {
		t.Errorf("expected package name 'single_file', got '%s'", pkg.PackageName())
	}

	// SingleFileProject has anonymous org
	if pkg.PackageOrg() != "" {
		t.Errorf("expected empty org for single file, got '%s'", pkg.PackageOrg())
	}

	mod := pkg.DefaultModule()
	if mod == nil {
		t.Fatal("DefaultModule() returned nil")
	}

	docs := mod.Documents()
	if len(docs) != 1 {
		t.Errorf("expected exactly one document, got %d", len(docs))
	}
}

func TestLoadMissingBallerinaToml(t *testing.T) {
	// Try to load a directory without Ballerina.toml
	_, err := Load("testdata")
	if err == nil {
		t.Fatal("expected error for directory without Ballerina.toml")
	}
}

func TestLoadNonExistentPath(t *testing.T) {
	_, err := Load("nonexistent_path")
	if err == nil {
		t.Fatal("expected error for non-existent path")
	}
}

func TestProjectTypeDetection(t *testing.T) {
	tests := []struct {
		path     string
		expected ProjectKind
	}{
		{filepath.Join("testdata", "valid_package"), BuildProjectKind},
		{filepath.Join("testdata", "single_file.bal"), SingleFileProjectKind},
	}

	for _, tt := range tests {
		project, err := Load(tt.path)
		if err != nil {
			t.Fatalf("failed to load %s: %v", tt.path, err)
		}
		if project.Kind() != tt.expected {
			t.Errorf("path %s: expected %v, got %v", tt.path, tt.expected, project.Kind())
		}
	}
}
