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

import "errors"

// ErrUnsupported indicates an operation is not yet implemented.
var ErrUnsupported = errors.New("operation not supported")

// ProjectKind identifies the type of project.
type ProjectKind int

const (
	// BuildProjectKind represents a standard Ballerina project with Ballerina.toml.
	BuildProjectKind ProjectKind = iota
	// SingleFileProjectKind represents a standalone .bal file.
	SingleFileProjectKind
	// BalaProjectKind represents a compiled package archive.
	BalaProjectKind
	// WorkspaceProjectKind represents a multi-project workspace.
	WorkspaceProjectKind
)

// String returns the string representation of the ProjectKind.
func (k ProjectKind) String() string {
	switch k {
	case BuildProjectKind:
		return "BUILD_PROJECT"
	case SingleFileProjectKind:
		return "SINGLE_FILE_PROJECT"
	case BalaProjectKind:
		return "BALA_PROJECT"
	case WorkspaceProjectKind:
		return "WORKSPACE_PROJECT"
	default:
		return "UNKNOWN"
	}
}

// Project is the interface for all project types.
type Project interface {
	// Kind returns the type of this project.
	Kind() ProjectKind

	// SourceRoot returns the absolute path to the project root.
	SourceRoot() string

	// CurrentPackage returns the main package of this project.
	CurrentPackage() *Package

	// TargetDir returns the directory for build outputs.
	TargetDir() string

	// BuildOptions returns the build configuration for this project.
	BuildOptions() BuildOptions

	// DocumentId returns the DocumentId for the given file path.
	DocumentId(path string) (DocumentId, error)

	// DocumentPath returns the file path for the given DocumentId.
	DocumentPath(id DocumentId) (string, error)

	// ClearCaches clears all compilation caches.
	ClearCaches()

	// Duplicate creates an independent copy of this project.
	Duplicate() (Project, error)

	// Save persists any changes to the filesystem.
	Save() error
}
