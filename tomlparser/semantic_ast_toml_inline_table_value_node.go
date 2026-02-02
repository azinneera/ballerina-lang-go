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

// TomlInlineTableValueNode represents an inline table {...}.
// Maps to: io.ballerina.toml.semantic.ast.TomlInlineTableValueNode
type TomlInlineTableValueNode struct {
	tomlNodeBase
	elements []TopLevelNode
}

func (n *TomlInlineTableValueNode) tomlValueNode() {}

// NewTomlInlineTableValueNode creates a new inline table value node.
func NewTomlInlineTableValueNode(elements []TopLevelNode, location *Location) *TomlInlineTableValueNode {
	return &TomlInlineTableValueNode{
		tomlNodeBase: tomlNodeBase{
			location:    location,
			diagnostics: make([]Diagnostic, 0),
		},
		elements: elements,
	}
}

func (n *TomlInlineTableValueNode) Kind() TomlType {
	return TomlTypeInlineTable
}

func (n *TomlInlineTableValueNode) Elements() []TopLevelNode {
	return n.elements
}

// Diagnostics aggregates diagnostics from this node and all elements.
func (n *TomlInlineTableValueNode) Diagnostics() []Diagnostic {
	allDiags := make([]Diagnostic, len(n.diagnostics))
	copy(allDiags, n.diagnostics)

	for _, elem := range n.elements {
		allDiags = append(allDiags, elem.Diagnostics()...)
	}

	return allDiags
}

// Accept implements the visitor pattern.
// Maps to: TomlInlineTableValueNode.accept(TomlNodeVisitor)
func (n *TomlInlineTableValueNode) Accept(visitor TomlNodeVisitor) {
	visitor.VisitTomlInlineTableValueNode(n)
}

// ToTable converts an inline table to a TomlTableNode.
// Maps to: TomlInlineTableValueNode.toTable()
func (n *TomlInlineTableValueNode) ToTable() *TomlTableNode {
	table := NewTomlTableNode("", n.location, false)
	for _, elem := range n.elements {
		table.entries[elem.Key()] = elem
	}
	return table
}
