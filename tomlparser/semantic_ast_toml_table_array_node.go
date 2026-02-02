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

// TomlTableArrayNode represents a TOML array of tables [[tableArray]].
// Maps to: io.ballerina.toml.semantic.ast.TomlTableArrayNode
type TomlTableArrayNode struct {
	tomlNodeBase
	key      string
	children []*TomlTableNode
}

// NewTomlTableArrayNode creates a new table array node.
// Maps to: TomlTableArrayNode constructor
func NewTomlTableArrayNode(key string, location *Location) *TomlTableArrayNode {
	return &TomlTableArrayNode{
		tomlNodeBase: tomlNodeBase{
			location:    location,
			diagnostics: make([]Diagnostic, 0),
		},
		key:      key,
		children: make([]*TomlTableNode, 0),
	}
}

func (n *TomlTableArrayNode) Kind() TomlType {
	return TomlTypeTableArray
}

func (n *TomlTableArrayNode) Key() string {
	return n.key
}

func (n *TomlTableArrayNode) Children() []*TomlTableNode {
	return n.children
}

// AddChild appends a table to the array.
// Maps to: TomlTableArrayNode.addChild()
func (n *TomlTableArrayNode) AddChild(child *TomlTableNode) {
	n.children = append(n.children, child)
}

// Diagnostics aggregates diagnostics from this node and all children.
func (n *TomlTableArrayNode) Diagnostics() []Diagnostic {
	allDiags := make([]Diagnostic, len(n.diagnostics))
	copy(allDiags, n.diagnostics)

	for _, child := range n.children {
		allDiags = append(allDiags, child.Diagnostics()...)
	}

	return allDiags
}

// Accept implements the visitor pattern.
// Maps to: TomlTableArrayNode.accept(TomlNodeVisitor)
func (n *TomlTableArrayNode) Accept(visitor TomlNodeVisitor) {
	visitor.VisitTomlTableArrayNode(n)
}
