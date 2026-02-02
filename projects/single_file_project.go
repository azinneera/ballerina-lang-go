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

// SingleFileProject represents a standalone .bal file.
type SingleFileProject struct {
	sourceRoot   string // directory containing the file
	filePath     string // full path to the .bal file
	pkg          *Package
	buildOptions BuildOptions
}

// Compile-time interface check.
var _ Project = (*SingleFileProject)(nil)

// Kind returns SingleFileProjectKind.
func (p *SingleFileProject) Kind() ProjectKind {
	return SingleFileProjectKind
}

// SourceRoot returns the directory containing the .bal file.
func (p *SingleFileProject) SourceRoot() string {
	return p.sourceRoot
}

// CurrentPackage returns the main package.
func (p *SingleFileProject) CurrentPackage() *Package {
	return p.pkg
}

// BuildOptions returns the build configuration.
func (p *SingleFileProject) BuildOptions() BuildOptions {
	return p.buildOptions
}

// TargetDir returns the build output directory.
func (p *SingleFileProject) TargetDir() string {
	if p.buildOptions.TargetDir() != "" {
		return p.buildOptions.TargetDir()
	}
	// Single file projects use the source directory
	return p.sourceRoot
}

// FilePath returns the full path to the .bal file.
func (p *SingleFileProject) FilePath() string {
	return p.filePath
}

// DocumentId returns the DocumentId for the given file path.
func (p *SingleFileProject) DocumentId(path string) (DocumentId, error) {
	baseName := filepath.Base(path)
	for _, doc := range p.pkg.DefaultModule().Documents() {
		if doc.name == baseName {
			return doc.id, nil
		}
	}
	return DocumentId{}, fmt.Errorf("document not found: %s", path)
}

// DocumentPath returns the file path for the given DocumentId.
func (p *SingleFileProject) DocumentPath(id DocumentId) (string, error) {
	doc := p.pkg.DefaultModule().Document(id)
	if doc != nil {
		return p.filePath, nil
	}
	return "", fmt.Errorf("document not found: %s", id)
}

// ClearCaches clears compilation caches.
func (p *SingleFileProject) ClearCaches() {
	// Single file projects have minimal caching
}

// Duplicate creates an independent copy (not yet implemented).
func (p *SingleFileProject) Duplicate() (Project, error) {
	return nil, ErrUnsupported
}

// Save persists changes to the filesystem (not yet implemented).
func (p *SingleFileProject) Save() error {
	return ErrUnsupported
}
