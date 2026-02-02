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

// DocumentConfig holds immutable document configuration from filesystem scan.
type DocumentConfig struct {
	id        DocumentId
	name      string
	contentFn func() string // lazy content loading
}

// NewDocumentConfig creates a new DocumentConfig with lazy content loading.
func NewDocumentConfig(id DocumentId, name string, contentFn func() string) DocumentConfig {
	return DocumentConfig{id: id, name: name, contentFn: contentFn}
}

// DocumentId returns the document's unique identifier.
func (c DocumentConfig) DocumentId() DocumentId {
	return c.id
}

// Name returns the document file name.
func (c DocumentConfig) Name() string {
	return c.name
}

// Content returns the document content (loaded lazily).
func (c DocumentConfig) Content() string {
	if c.contentFn == nil {
		return ""
	}
	return c.contentFn()
}

// ModuleConfig holds immutable module configuration from filesystem scan.
type ModuleConfig struct {
	id         ModuleId
	descriptor ModuleDescriptor
	srcDocs    []DocumentConfig
	testDocs   []DocumentConfig
}

// NewModuleConfig creates a new ModuleConfig.
func NewModuleConfig(id ModuleId, desc ModuleDescriptor, srcDocs []DocumentConfig) ModuleConfig {
	return ModuleConfig{id: id, descriptor: desc, srcDocs: srcDocs}
}

// NewModuleConfigWithTests creates a new ModuleConfig with test documents.
func NewModuleConfigWithTests(id ModuleId, desc ModuleDescriptor, srcDocs, testDocs []DocumentConfig) ModuleConfig {
	return ModuleConfig{id: id, descriptor: desc, srcDocs: srcDocs, testDocs: testDocs}
}

// ModuleId returns the module's unique identifier.
func (c ModuleConfig) ModuleId() ModuleId {
	return c.id
}

// Descriptor returns the module descriptor.
func (c ModuleConfig) Descriptor() ModuleDescriptor {
	return c.descriptor
}

// SrcDocs returns the source document configurations.
func (c ModuleConfig) SrcDocs() []DocumentConfig {
	return c.srcDocs
}

// TestDocs returns the test document configurations.
func (c ModuleConfig) TestDocs() []DocumentConfig {
	return c.testDocs
}

// PackageConfig holds immutable package configuration from filesystem scan.
type PackageConfig struct {
	id            PackageId
	path          string
	descriptor    PackageDescriptor
	defaultModule ModuleConfig
	otherModules  []ModuleConfig
}

// NewPackageConfig creates a new PackageConfig.
func NewPackageConfig(id PackageId, path string, desc PackageDescriptor, defaultModule ModuleConfig) PackageConfig {
	return PackageConfig{id: id, path: path, descriptor: desc, defaultModule: defaultModule}
}

// NewPackageConfigWithModules creates a new PackageConfig with additional modules.
func NewPackageConfigWithModules(id PackageId, path string, desc PackageDescriptor,
	defaultModule ModuleConfig, otherModules []ModuleConfig) PackageConfig {
	return PackageConfig{id: id, path: path, descriptor: desc, defaultModule: defaultModule, otherModules: otherModules}
}

// PackageId returns the package's unique identifier.
func (c PackageConfig) PackageId() PackageId {
	return c.id
}

// Path returns the package root path.
func (c PackageConfig) Path() string {
	return c.path
}

// Descriptor returns the package descriptor.
func (c PackageConfig) Descriptor() PackageDescriptor {
	return c.descriptor
}

// DefaultModule returns the default module configuration.
func (c PackageConfig) DefaultModule() ModuleConfig {
	return c.defaultModule
}

// OtherModules returns configurations for non-default modules.
func (c PackageConfig) OtherModules() []ModuleConfig {
	return c.otherModules
}

// AllModules returns all module configurations (default + others).
func (c PackageConfig) AllModules() []ModuleConfig {
	modules := make([]ModuleConfig, 0, 1+len(c.otherModules))
	modules = append(modules, c.defaultModule)
	modules = append(modules, c.otherModules...)
	return modules
}
