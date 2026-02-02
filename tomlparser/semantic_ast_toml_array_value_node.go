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

// TomlArrayValueNode represents an array value [...].
// Maps to: io.ballerina.toml.semantic.ast.TomlArrayValueNode
type TomlArrayValueNode struct {
	tomlNodeBase
	elements []TomlValueNode
}

func (n *TomlArrayValueNode) tomlValueNode() {}

// NewTomlArrayValueNode creates a new array value node.
func NewTomlArrayValueNode(elements []TomlValueNode, location *Location) *TomlArrayValueNode {
	return &TomlArrayValueNode{
		tomlNodeBase: tomlNodeBase{
			location:    location,
			diagnostics: make([]Diagnostic, 0),
		},
		elements: elements,
	}
}

func (n *TomlArrayValueNode) Kind() TomlType {
	return TomlTypeArray
}

func (n *TomlArrayValueNode) Elements() []TomlValueNode {
	return n.elements
}

// Diagnostics aggregates diagnostics from this node and all elements.
func (n *TomlArrayValueNode) Diagnostics() []Diagnostic {
	allDiags := make([]Diagnostic, len(n.diagnostics))
	copy(allDiags, n.diagnostics)

	for _, elem := range n.elements {
		allDiags = append(allDiags, elem.Diagnostics()...)
	}

	return allDiags
}

// Accept implements the visitor pattern.
// Maps to: TomlArrayValueNode.accept(TomlNodeVisitor)
func (n *TomlArrayValueNode) Accept(visitor TomlNodeVisitor) {
	visitor.VisitTomlArrayValueNode(n)
}
