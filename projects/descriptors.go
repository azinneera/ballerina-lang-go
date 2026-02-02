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

import "strings"

// PackageDescriptor contains immutable metadata for a package.
type PackageDescriptor struct {
	name       PackageName
	org        PackageOrg
	version    PackageVersion
	repository string
}

// NewPackageDescriptor creates a new PackageDescriptor.
func NewPackageDescriptor(org PackageOrg, name PackageName, version PackageVersion) PackageDescriptor {
	return PackageDescriptor{name: name, org: org, version: version}
}

// NewPackageDescriptorWithRepo creates a new PackageDescriptor with a repository.
func NewPackageDescriptorWithRepo(org PackageOrg, name PackageName, version PackageVersion,
	repository string) PackageDescriptor {
	return PackageDescriptor{name: name, org: org, version: version, repository: repository}
}

// Name returns the package name.
func (d PackageDescriptor) Name() PackageName {
	return d.name
}

// Org returns the package organization.
func (d PackageDescriptor) Org() PackageOrg {
	return d.org
}

// Version returns the package version.
func (d PackageDescriptor) Version() PackageVersion {
	return d.version
}

// Repository returns the repository name (empty for Central).
func (d PackageDescriptor) Repository() string {
	return d.repository
}

// IsLangLibPackage returns true if this is a lang lib package (ballerina/lang.*).
func (d PackageDescriptor) IsLangLibPackage() bool {
	return d.org.IsBallerinaOrg() && strings.HasPrefix(string(d.name), "lang.")
}

// IsBuiltInPackage returns true if this is a Ballerina built-in package.
func (d PackageDescriptor) IsBuiltInPackage() bool {
	return d.org.IsBallerinaOrg()
}

// String returns a string representation of the descriptor.
func (d PackageDescriptor) String() string {
	if d.org.IsEmpty() {
		return string(d.name)
	}
	return string(d.org) + "/" + string(d.name) + ":" + d.version.String()
}

// ModuleDescriptor contains immutable metadata for a module.
type ModuleDescriptor struct {
	name        ModuleName
	packageDesc PackageDescriptor
}

// NewModuleDescriptor creates a new ModuleDescriptor.
func NewModuleDescriptor(name ModuleName, pkgDesc PackageDescriptor) ModuleDescriptor {
	return ModuleDescriptor{name: name, packageDesc: pkgDesc}
}

// Name returns the module name.
func (d ModuleDescriptor) Name() ModuleName {
	return d.name
}

// PackageDescriptor returns the parent package descriptor.
func (d ModuleDescriptor) PackageDescriptor() PackageDescriptor {
	return d.packageDesc
}

// PackageName returns the package name.
func (d ModuleDescriptor) PackageName() PackageName {
	return d.packageDesc.Name()
}

// Org returns the package organization.
func (d ModuleDescriptor) Org() PackageOrg {
	return d.packageDesc.Org()
}

// Version returns the package version.
func (d ModuleDescriptor) Version() PackageVersion {
	return d.packageDesc.Version()
}

// IsDefaultModule returns true if this is the default module.
func (d ModuleDescriptor) IsDefaultModule() bool {
	return d.name.IsDefault()
}

// String returns a string representation of the module descriptor.
func (d ModuleDescriptor) String() string {
	return d.name.String()
}
