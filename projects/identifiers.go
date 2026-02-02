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

import "github.com/google/uuid"

// PackageId uniquely identifies a Package instance.
// Equality is based on UUID only (name is for debugging).
type PackageId struct {
	id   uuid.UUID
	name string
}

// NewPackageId creates a new PackageId with a random UUID.
func NewPackageId(name string) PackageId {
	return PackageId{id: uuid.New(), name: name}
}

// String returns the package name for display purposes.
func (p PackageId) String() string {
	return p.name
}

// IsZero returns true if this is a zero-value PackageId.
func (p PackageId) IsZero() bool {
	return p.id == uuid.Nil
}

// Equal checks if two PackageIds are equal (based on UUID).
func (p PackageId) Equal(other PackageId) bool {
	return p.id == other.id
}

// ModuleId uniquely identifies a Module within a Package.
type ModuleId struct {
	id        uuid.UUID
	name      string
	packageId PackageId
}

// NewModuleId creates a new ModuleId with a random UUID.
func NewModuleId(name string, pkgId PackageId) ModuleId {
	return ModuleId{id: uuid.New(), name: name, packageId: pkgId}
}

// String returns the module name for display purposes.
func (m ModuleId) String() string {
	return m.name
}

// PackageId returns the parent package's ID.
func (m ModuleId) PackageId() PackageId {
	return m.packageId
}

// IsZero returns true if this is a zero-value ModuleId.
func (m ModuleId) IsZero() bool {
	return m.id == uuid.Nil
}

// Equal checks if two ModuleIds are equal (based on UUID and PackageId).
func (m ModuleId) Equal(other ModuleId) bool {
	return m.id == other.id && m.packageId.Equal(other.packageId)
}

// DocumentId uniquely identifies a Document within a Module.
type DocumentId struct {
	id       uuid.UUID
	path     string
	moduleId ModuleId
}

// NewDocumentId creates a new DocumentId with a random UUID.
func NewDocumentId(path string, modId ModuleId) DocumentId {
	return DocumentId{id: uuid.New(), path: path, moduleId: modId}
}

// String returns the document path for display purposes.
func (d DocumentId) String() string {
	return d.path
}

// ModuleId returns the parent module's ID.
func (d DocumentId) ModuleId() ModuleId {
	return d.moduleId
}

// IsZero returns true if this is a zero-value DocumentId.
func (d DocumentId) IsZero() bool {
	return d.id == uuid.Nil
}

// Equal checks if two DocumentIds are equal (based on UUID and ModuleId).
func (d DocumentId) Equal(other DocumentId) bool {
	return d.id == other.id && d.moduleId.Equal(other.moduleId)
}
