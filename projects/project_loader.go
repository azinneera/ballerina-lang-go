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
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadOption configures project loading.
type LoadOption func(*loadConfig)

// loadConfig holds configuration for project loading.
type loadConfig struct {
	buildOptions       BuildOptions
	environmentBuilder *ProjectEnvironmentBuilder
}

// WithBuildOptions sets custom build options for the project.
func WithBuildOptions(opts BuildOptions) LoadOption {
	return func(c *loadConfig) {
		c.buildOptions = opts
	}
}

// WithEnvironmentBuilder sets a custom environment builder for the project.
func WithEnvironmentBuilder(builder *ProjectEnvironmentBuilder) LoadOption {
	return func(c *loadConfig) {
		c.environmentBuilder = builder
	}
}

// Load loads a project from the given path (auto-detects type).
// Returns (Project, error) following Go conventions.
func Load(path string, opts ...LoadOption) (Project, error) {
	cfg := loadConfig{
		buildOptions: NewBuildOptionsBuilder().Build(),
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("path not found: %w", err)
	}

	// Check for workspace project (Phase 6)
	if info.IsDir() && IsWorkspaceProject(absPath) {
		return nil, fmt.Errorf("workspace projects not yet supported: %w", ErrUnsupported)
	}

	// Check for build project (has Ballerina.toml)
	if info.IsDir() && IsBuildProject(absPath) {
		return LoadBuildProject(absPath, opts...)
	}

	// Check for BALA project (Phase 5)
	if IsBalaProject(absPath) {
		return nil, fmt.Errorf("BALA projects not yet supported: %w", ErrUnsupported)
	}

	// Single file project
	if !info.IsDir() && filepath.Ext(absPath) == ".bal" {
		return LoadSingleFileProject(absPath, opts...)
	}

	return nil, fmt.Errorf("cannot determine project type for: %s", path)
}

// IsBuildProject checks if the path contains a Ballerina.toml.
func IsBuildProject(path string) bool {
	_, err := os.Stat(filepath.Join(path, "Ballerina.toml"))
	return err == nil
}

// IsWorkspaceProject checks if this is a workspace project.
// Phase 6: Check for [workspace] section in Ballerina.toml.
func IsWorkspaceProject(path string) bool {
	return false
}

// IsBalaProject checks if this is a BALA project.
func IsBalaProject(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		// Check for package.json in extracted BALA
		_, err := os.Stat(filepath.Join(path, "package.json"))
		return err == nil
	}
	// Check for .bala extension
	return filepath.Ext(path) == ".bala"
}

// LoadBuildProject loads a build project from the given directory.
func LoadBuildProject(path string, opts ...LoadOption) (*BuildProject, error) {
	cfg := loadConfig{
		buildOptions: NewBuildOptionsBuilder().Build(),
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	if !IsBuildProject(absPath) {
		return nil, fmt.Errorf("not a build project (missing Ballerina.toml): %s", path)
	}

	// 1. Parse Ballerina.toml
	manifest, err := ParsePackageManifest(filepath.Join(absPath, "Ballerina.toml"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse Ballerina.toml: %w", err)
	}

	// 2. Create package descriptor
	version, err := ParsePackageVersion(manifest.Version)
	if err != nil {
		return nil, fmt.Errorf("invalid version: %w", err)
	}
	pkgDesc := NewPackageDescriptor(
		PackageOrg(manifest.Org),
		PackageName(manifest.Name),
		version,
	)

	// 3. Create package ID
	pkgId := NewPackageId(manifest.Name)

	// 4. Scan source files in root directory
	srcDocs, err := ScanSourceFiles(absPath, pkgId)
	if err != nil {
		return nil, fmt.Errorf("failed to scan source files: %w", err)
	}

	// 5. Create module config for default module
	modId := NewModuleId(manifest.Name, pkgId)
	modName := NewModuleName(PackageName(manifest.Name), "")
	modDesc := NewModuleDescriptor(modName, pkgDesc)
	modConfig := NewModuleConfig(modId, modDesc, srcDocs)

	// 6. Create package config
	pkgConfig := NewPackageConfig(pkgId, absPath, pkgDesc, modConfig)

	// 7. Create project
	project := &BuildProject{
		sourceRoot:   absPath,
		buildOptions: cfg.buildOptions,
	}

	// 8. Create package from config
	project.pkg = NewPackageFromConfig(project, pkgConfig)

	return project, nil
}

// LoadSingleFileProject loads a single file project from the given .bal file.
func LoadSingleFileProject(path string, opts ...LoadOption) (*SingleFileProject, error) {
	cfg := loadConfig{
		buildOptions: NewBuildOptionsBuilder().Build(),
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("path not found: %w", err)
	}

	if info.IsDir() || filepath.Ext(absPath) != ".bal" {
		return nil, fmt.Errorf("not a single file project (expected .bal file): %s", path)
	}

	// 1. Get file name (used as package/module name)
	fileName := filepath.Base(absPath)
	name := strings.TrimSuffix(fileName, ".bal")
	dir := filepath.Dir(absPath)

	// 2. Create package descriptor (anonymous org)
	pkgDesc := NewPackageDescriptor(
		PackageOrg(""),
		PackageName(name),
		NewPackageVersion(0, 0, 0),
	)

	// 3. Create package ID
	pkgId := NewPackageId(name)

	// 4. Create single document config
	modId := NewModuleId(name, pkgId)
	docId := NewDocumentId(fileName, modId)
	contentFn := func() string {
		data, _ := os.ReadFile(absPath)
		return string(data)
	}
	docConfig := NewDocumentConfig(docId, fileName, contentFn)

	// 5. Create module config for default module
	modName := NewModuleName(PackageName(name), "")
	modDesc := NewModuleDescriptor(modName, pkgDesc)
	modConfig := NewModuleConfig(modId, modDesc, []DocumentConfig{docConfig})

	// 6. Create package config
	pkgConfig := NewPackageConfig(pkgId, dir, pkgDesc, modConfig)

	// 7. Create project
	project := &SingleFileProject{
		sourceRoot:   dir,
		filePath:     absPath,
		buildOptions: cfg.buildOptions,
	}

	// 8. Create package from config
	project.pkg = NewPackageFromConfig(project, pkgConfig)

	return project, nil
}

// ProjectEnvironmentBuilder builds project environments (Phase 2+).
type ProjectEnvironmentBuilder struct {
	// TODO: Add environment configuration fields in later phases
}

// NewProjectEnvironmentBuilder creates a new ProjectEnvironmentBuilder.
func NewProjectEnvironmentBuilder() *ProjectEnvironmentBuilder {
	return &ProjectEnvironmentBuilder{}
}
