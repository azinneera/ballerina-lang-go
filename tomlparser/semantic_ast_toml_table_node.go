// Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
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

package tomlparser

// TomlTableNode represents a TOML table [table] or the root document.
// Maps to: io.ballerina.toml.semantic.ast.TomlTableNode
type TomlTableNode struct {
	tomlNodeBase
	key       string
	entries   map[string]TopLevelNode
	generated bool
}

// NewTomlTableNode creates a new table node.
// Maps to: TomlTableNode constructor
func NewTomlTableNode(key string, location *Location, generated bool) *TomlTableNode {
	return &TomlTableNode{
		tomlNodeBase: tomlNodeBase{
			location:    location,
			diagnostics: make([]Diagnostic, 0),
		},
		key:       key,
		entries:   make(map[string]TopLevelNode),
		generated: generated,
	}
}

func (n *TomlTableNode) Kind() TomlType {
	return TomlTypeTable
}

func (n *TomlTableNode) Key() string {
	return n.key
}

func (n *TomlTableNode) Entries() map[string]TopLevelNode {
	return n.entries
}

func (n *TomlTableNode) Generated() bool {
	return n.generated
}

func (n *TomlTableNode) SetGenerated(generated bool) {
	n.generated = generated
}

// Diagnostics aggregates diagnostics from this node and all children.
// Maps to: TomlTableNode.diagnostics()
func (n *TomlTableNode) Diagnostics() []Diagnostic {
	allDiags := make([]Diagnostic, len(n.diagnostics))
	copy(allDiags, n.diagnostics)

	for _, child := range n.entries {
		allDiags = append(allDiags, child.Diagnostics()...)
	}

	return allDiags
}

// Accept implements the visitor pattern.
// Maps to: TomlTableNode.accept(TomlNodeVisitor)
func (n *TomlTableNode) Accept(visitor TomlNodeVisitor) {
	visitor.VisitTomlTableNode(n)
}

// ReplaceGeneratedTable replaces an implicitly generated table with an explicit one.
// Maps to: TomlTableNode.replaceGeneratedTable()
func (n *TomlTableNode) ReplaceGeneratedTable(newTable *TomlTableNode) {
	childNode, exists := n.entries[newTable.Key()]
	if !exists {
		return
	}

	childTable, ok := childNode.(*TomlTableNode)
	if !ok {
		return
	}

	if childTable.Generated() {
		for k, v := range childTable.Entries() {
			newTable.entries[k] = v
		}
		n.entries[newTable.Key()] = newTable
	}
}
