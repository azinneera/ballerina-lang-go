// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Package symbolspaces provides the store used to register and look up
// symbol spaces in env, keyed by a stable index. Package-level spaces are
// additionally grouped by package for external symbol lookup.
package symbolspaces

import (
	"sync"
	"sync/atomic"

	"github.com/ballerina-nutcracker/ballerina/model"
)

// perPackage holds the package-level spaces used for external symbol lookup.
type perPackage struct {
	mu            sync.RWMutex
	packageSpaces []*model.SymbolSpace
}

const spaceBlockSize = 1024

type spaceBlock struct {
	spaces [spaceBlockSize]atomic.Pointer[model.SymbolSpace]
}

type spaceDirectory struct {
	blocks []*spaceBlock
}

// Store provides lock-free reads by stable space index and package-scoped
// lookup for external symbol names.
type Store struct {
	mu          sync.RWMutex
	pkg         map[model.PackageIdentifier]*perPackage
	nextSpaceID atomic.Int64
	blockMu     sync.Mutex
	directory   atomic.Pointer[spaceDirectory]
}

func NewStore() Store {
	store := Store{pkg: make(map[model.PackageIdentifier]*perPackage)}
	store.directory.Store(&spaceDirectory{})
	return store
}

func (s *Store) bucket(pkgID model.PackageIdentifier) *perPackage {
	s.mu.RLock()
	b, ok := s.pkg[pkgID]
	s.mu.RUnlock()
	if ok {
		return b
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.pkg[pkgID]; ok {
		return existing
	}
	b = &perPackage{}
	s.pkg[pkgID] = b
	return b
}

// Register assigns a stable, non-zero index in the environment-wide table.
// Safe for concurrent use.
func (s *Store) Register(space *model.SymbolSpace) int {
	index := int(s.nextSpaceID.Add(1))
	zeroBasedIndex := index - 1
	blockIndex := zeroBasedIndex / spaceBlockSize
	directory := s.directory.Load()
	if blockIndex >= len(directory.blocks) {
		s.blockMu.Lock()
		directory = s.directory.Load()
		if blockIndex >= len(directory.blocks) {
			blocks := append([]*spaceBlock(nil), directory.blocks...)
			for len(blocks) <= blockIndex {
				blocks = append(blocks, &spaceBlock{})
			}
			directory = &spaceDirectory{blocks: blocks}
			s.directory.Store(directory)
		}
		s.blockMu.Unlock()
	}
	space.SpaceIndex = index
	directory.blocks[blockIndex].spaces[zeroBasedIndex%spaceBlockSize].Store(space)
	return index
}

// RegisterPackageSpace marks a space as searchable by Find. Function and
// block scopes are indexed for SymbolRef resolution but are not package symbol
// spaces and must not participate in name lookup during deserialization.
func (s *Store) RegisterPackageSpace(space *model.SymbolSpace) {
	b := s.bucket(space.Pkg)
	b.mu.Lock()
	b.packageSpaces = append(b.packageSpaces, space)
	b.mu.Unlock()
}

// SpaceAt returns the symbol space registered at index.
func (s *Store) SpaceAt(index int) *model.SymbolSpace {
	directory := s.directory.Load()
	zeroBasedIndex := index - 1
	return directory.blocks[zeroBasedIndex/spaceBlockSize].spaces[zeroBasedIndex%spaceBlockSize].Load()
}

// Find looks up name across package-level spaces registered for pkgID. Used
// for symbol deserialization only — potentially slow, not meant for use
// anywhere a SymbolRef can be carried in the AST instead.
func (s *Store) Find(pkgID model.PackageIdentifier, name string) (model.SymbolRef, bool) {
	b := s.bucket(pkgID)
	b.mu.RLock()
	spaces := append([]*model.SymbolSpace(nil), b.packageSpaces...)
	b.mu.RUnlock()
	for _, space := range spaces {
		if ref, ok := space.GetSymbol(name); ok {
			return ref, true
		}
	}
	return model.SymbolRef{}, false
}
