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
	"path/filepath"
)

// BuildProject represents a standard Ballerina project with Ballerina.toml.
type BuildProject struct {
	sourceRoot   string
	pkg          *Package
	buildOptions BuildOptions
}

// Compile-time interface check.
var _ Project = (*BuildProject)(nil)

// Kind returns BuildProjectKind.
func (p *BuildProject) Kind() ProjectKind {
	return BuildProjectKind
}

// SourceRoot returns the project root directory.
func (p *BuildProject) SourceRoot() string {
	return p.sourceRoot
}

// CurrentPackage returns the main package.
func (p *BuildProject) CurrentPackage() *Package {
	return p.pkg
}

// BuildOptions returns the build configuration.
func (p *BuildProject) BuildOptions() BuildOptions {
	return p.buildOptions
}

// TargetDir returns the build output directory.
func (p *BuildProject) TargetDir() string {
	if p.buildOptions.TargetDir() != "" {
		return p.buildOptions.TargetDir()
	}
	return filepath.Join(p.sourceRoot, "target")
}

// DocumentId returns the DocumentId for the given file path.
func (p *BuildProject) DocumentId(path string) (DocumentId, error) {
	baseName := filepath.Base(path)

	// Search in default module
	for _, doc := range p.pkg.DefaultModule().Documents() {
		if doc.name == baseName {
			return doc.id, nil
		}
	}

	// Search in other modules
	for _, mod := range p.pkg.Modules() {
		for _, doc := range mod.Documents() {
			if doc.name == baseName {
				return doc.id, nil
			}
		}
	}

	return DocumentId{}, fmt.Errorf("document not found: %s", path)
}

// DocumentPath returns the file path for the given DocumentId.
func (p *BuildProject) DocumentPath(id DocumentId) (string, error) {
	// Check default module
	doc := p.pkg.DefaultModule().Document(id)
	if doc != nil {
		return filepath.Join(p.sourceRoot, doc.name), nil
	}

	// Check other modules
	for _, mod := range p.pkg.Modules() {
		doc := mod.Document(id)
		if doc != nil {
			modPath := filepath.Join(p.sourceRoot, "modules", mod.ModuleName().ModulePart)
			return filepath.Join(modPath, doc.name), nil
		}
	}

	return "", fmt.Errorf("document not found: %s", id)
}

// ClearCaches clears compilation caches (Phase 3).
func (p *BuildProject) ClearCaches() {
	// Phase 3: Implement cache clearing
}

// Duplicate creates an independent copy (not yet implemented).
func (p *BuildProject) Duplicate() (Project, error) {
	return nil, ErrUnsupported
}

// Save persists changes to the filesystem (not yet implemented).
func (p *BuildProject) Save() error {
	return ErrUnsupported
}
