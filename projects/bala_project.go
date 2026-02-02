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

// BalaProject represents a compiled package archive (.bala file).
// This is a stub implementation - full support is planned for Phase 5.
type BalaProject struct {
	sourceRoot string
	pkg        *Package
}

// Compile-time interface check.
var _ Project = (*BalaProject)(nil)

// Kind returns BalaProjectKind.
func (p *BalaProject) Kind() ProjectKind {
	return BalaProjectKind
}

// SourceRoot returns the path to the .bala file or extracted directory.
func (p *BalaProject) SourceRoot() string {
	return p.sourceRoot
}

// CurrentPackage returns nil (not yet implemented).
func (p *BalaProject) CurrentPackage() *Package {
	return nil
}

// BuildOptions returns empty options.
func (p *BalaProject) BuildOptions() BuildOptions {
	return BuildOptions{}
}

// TargetDir returns an empty string (not applicable for BALA projects).
func (p *BalaProject) TargetDir() string {
	return ""
}

// DocumentId is not yet implemented for BALA projects.
func (p *BalaProject) DocumentId(path string) (DocumentId, error) {
	return DocumentId{}, ErrUnsupported
}

// DocumentPath is not yet implemented for BALA projects.
func (p *BalaProject) DocumentPath(id DocumentId) (string, error) {
	return "", ErrUnsupported
}

// ClearCaches clears compilation caches.
func (p *BalaProject) ClearCaches() {
	// Phase 5: Implement cache clearing
}

// Duplicate is not yet implemented.
func (p *BalaProject) Duplicate() (Project, error) {
	return nil, ErrUnsupported
}

// Save is not applicable for BALA projects (read-only).
func (p *BalaProject) Save() error {
	return ErrUnsupported
}
