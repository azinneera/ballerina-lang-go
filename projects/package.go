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

import "sync"

// Package represents a Ballerina package.
type Package struct {
	project         Project
	id              PackageId
	descriptor      PackageDescriptor
	defaultMod      *Module
	modules         map[ModuleId]*Module
	compilation     *PackageCompilation
	compilationOnce sync.Once
	compilationErr  error
}

// NewPackageFromConfig creates a Package from a PackageConfig.
func NewPackageFromConfig(project Project, cfg PackageConfig) *Package {
	pkg := &Package{
		project:    project,
		id:         cfg.PackageId(),
		descriptor: cfg.Descriptor(),
		modules:    make(map[ModuleId]*Module),
	}

	// Create default module
	pkg.defaultMod = NewModuleFromConfig(pkg, cfg.DefaultModule())
	pkg.modules[cfg.DefaultModule().ModuleId()] = pkg.defaultMod

	// Create other modules
	for _, modCfg := range cfg.OtherModules() {
		mod := NewModuleFromConfig(pkg, modCfg)
		pkg.modules[modCfg.ModuleId()] = mod
	}

	return pkg
}

// PackageId returns the package's unique identifier.
func (p *Package) PackageId() PackageId {
	return p.id
}

// Descriptor returns the package descriptor.
func (p *Package) Descriptor() PackageDescriptor {
	return p.descriptor
}

// PackageName returns the package name.
func (p *Package) PackageName() PackageName {
	return p.descriptor.Name()
}

// PackageOrg returns the package organization.
func (p *Package) PackageOrg() PackageOrg {
	return p.descriptor.Org()
}

// PackageVersion returns the package version.
func (p *Package) PackageVersion() PackageVersion {
	return p.descriptor.Version()
}

// Project returns the parent project.
func (p *Package) Project() Project {
	return p.project
}

// DefaultModule returns the default module.
func (p *Package) DefaultModule() *Module {
	return p.defaultMod
}

// Module returns a module by ID.
func (p *Package) Module(id ModuleId) *Module {
	return p.modules[id]
}

// ModuleIds returns all module IDs in this package.
func (p *Package) ModuleIds() []ModuleId {
	ids := make([]ModuleId, 0, len(p.modules))
	for id := range p.modules {
		ids = append(ids, id)
	}
	return ids
}

// Modules returns all modules in this package.
func (p *Package) Modules() []*Module {
	mods := make([]*Module, 0, len(p.modules))
	for _, mod := range p.modules {
		mods = append(mods, mod)
	}
	return mods
}

// Manifest returns the parsed package manifest (Phase 3).
func (p *Package) Manifest() (*PackageManifest, error) {
	return nil, ErrUnsupported
}

// Compilation compiles the package and returns the compilation result.
// Results are cached; subsequent calls return the same compilation.
// Thread-safe via sync.Once.
func (p *Package) Compilation() (*PackageCompilation, error) {
	p.compilationOnce.Do(func() {
		p.compilation, p.compilationErr = newPackageCompilation(p)
	})
	return p.compilation, p.compilationErr
}

// GetResolution returns the package resolution (Phase 3).
func (p *Package) GetResolution() error {
	return ErrUnsupported
}
