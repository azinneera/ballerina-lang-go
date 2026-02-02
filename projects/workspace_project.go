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

// WorkspaceProject represents a multi-project workspace.
// This is a stub implementation - full support is planned for Phase 6.
type WorkspaceProject struct {
	sourceRoot   string
	projects     []*BuildProject
	buildOptions BuildOptions
}

// Compile-time interface check.
var _ Project = (*WorkspaceProject)(nil)

// Kind returns WorkspaceProjectKind.
func (p *WorkspaceProject) Kind() ProjectKind {
	return WorkspaceProjectKind
}

// SourceRoot returns the workspace root directory.
func (p *WorkspaceProject) SourceRoot() string {
	return p.sourceRoot
}

// CurrentPackage returns nil (workspace has multiple packages).
func (p *WorkspaceProject) CurrentPackage() *Package {
	return nil
}

// BuildOptions returns the build configuration.
func (p *WorkspaceProject) BuildOptions() BuildOptions {
	return p.buildOptions
}

// TargetDir returns an empty string (not yet implemented).
func (p *WorkspaceProject) TargetDir() string {
	return ""
}

// DocumentId is not yet implemented for workspace projects.
func (p *WorkspaceProject) DocumentId(path string) (DocumentId, error) {
	return DocumentId{}, ErrUnsupported
}

// DocumentPath is not yet implemented for workspace projects.
func (p *WorkspaceProject) DocumentPath(id DocumentId) (string, error) {
	return "", ErrUnsupported
}

// ClearCaches clears all project caches.
func (p *WorkspaceProject) ClearCaches() {
	// Phase 6: Implement cache clearing for all projects
}

// Duplicate is not yet implemented.
func (p *WorkspaceProject) Duplicate() (Project, error) {
	return nil, ErrUnsupported
}

// Save is not yet implemented.
func (p *WorkspaceProject) Save() error {
	return ErrUnsupported
}

// Projects returns all build projects in the workspace (Phase 6).
func (p *WorkspaceProject) Projects() []*BuildProject {
	return nil
}

// Reload reloads the workspace configuration (Phase 6).
func (p *WorkspaceProject) Reload() (Project, error) {
	return nil, ErrUnsupported
}
