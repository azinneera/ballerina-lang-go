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

// Module represents a Ballerina module within a package.
type Module struct {
	pkg       *Package
	id        ModuleId
	desc      ModuleDescriptor
	documents map[DocumentId]*Document
	testDocs  map[DocumentId]*Document
}

// NewModuleFromConfig creates a Module from a ModuleConfig.
func NewModuleFromConfig(pkg *Package, cfg ModuleConfig) *Module {
	mod := &Module{
		pkg:       pkg,
		id:        cfg.ModuleId(),
		desc:      cfg.Descriptor(),
		documents: make(map[DocumentId]*Document),
		testDocs:  make(map[DocumentId]*Document),
	}

	// Create source documents
	for _, docCfg := range cfg.SrcDocs() {
		doc := NewDocumentFromConfig(mod, docCfg)
		mod.documents[docCfg.DocumentId()] = doc
	}

	// Create test documents
	for _, docCfg := range cfg.TestDocs() {
		doc := NewDocumentFromConfig(mod, docCfg)
		mod.testDocs[docCfg.DocumentId()] = doc
	}

	return mod
}

// ModuleId returns the module's unique identifier.
func (m *Module) ModuleId() ModuleId {
	return m.id
}

// Descriptor returns the module descriptor.
func (m *Module) Descriptor() ModuleDescriptor {
	return m.desc
}

// ModuleName returns the module name.
func (m *Module) ModuleName() ModuleName {
	return m.desc.Name()
}

// Package returns the parent package.
func (m *Module) Package() *Package {
	return m.pkg
}

// IsDefault returns true if this is the default module.
func (m *Module) IsDefault() bool {
	return m.desc.Name().IsDefault()
}

// DocumentIds returns all source document IDs.
func (m *Module) DocumentIds() []DocumentId {
	ids := make([]DocumentId, 0, len(m.documents))
	for id := range m.documents {
		ids = append(ids, id)
	}
	return ids
}

// Documents returns all source documents.
func (m *Module) Documents() []*Document {
	docs := make([]*Document, 0, len(m.documents))
	for _, doc := range m.documents {
		docs = append(docs, doc)
	}
	return docs
}

// Document returns a source document by ID.
func (m *Module) Document(id DocumentId) *Document {
	return m.documents[id]
}

// TestDocumentIds returns all test document IDs (Phase 2).
func (m *Module) TestDocumentIds() []DocumentId {
	ids := make([]DocumentId, 0, len(m.testDocs))
	for id := range m.testDocs {
		ids = append(ids, id)
	}
	return ids
}

// TestDocuments returns all test documents (Phase 2).
func (m *Module) TestDocuments() []*Document {
	docs := make([]*Document, 0, len(m.testDocs))
	for _, doc := range m.testDocs {
		docs = append(docs, doc)
	}
	return docs
}

// TestDocument returns a test document by ID (Phase 2).
func (m *Module) TestDocument(id DocumentId) *Document {
	return m.testDocs[id]
}

// GetCompilation returns the module compilation (Phase 2).
func (m *Module) GetCompilation() error {
	return ErrUnsupported
}
