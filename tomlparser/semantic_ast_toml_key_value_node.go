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

// TomlKeyValueNode represents a key = value pair.
// Maps to: io.ballerina.toml.semantic.ast.TomlKeyValueNode
type TomlKeyValueNode struct {
	tomlNodeBase
	key   string
	value TomlValueNode
}

// NewTomlKeyValueNode creates a new key-value node.
// Maps to: TomlKeyValueNode constructor
func NewTomlKeyValueNode(key string, value TomlValueNode, location *Location) *TomlKeyValueNode {
	return &TomlKeyValueNode{
		tomlNodeBase: tomlNodeBase{
			location:    location,
			diagnostics: make([]Diagnostic, 0),
		},
		key:   key,
		value: value,
	}
}

func (n *TomlKeyValueNode) Kind() TomlType {
	return TomlTypeKeyValue
}

func (n *TomlKeyValueNode) Key() string {
	return n.key
}

func (n *TomlKeyValueNode) Value() TomlValueNode {
	return n.value
}

// Diagnostics aggregates diagnostics from this node and its value.
func (n *TomlKeyValueNode) Diagnostics() []Diagnostic {
	allDiags := make([]Diagnostic, len(n.diagnostics))
	copy(allDiags, n.diagnostics)

	if n.value != nil {
		allDiags = append(allDiags, n.value.Diagnostics()...)
	}

	return allDiags
}

// Accept implements the visitor pattern.
// Maps to: TomlKeyValueNode.accept(TomlNodeVisitor)
func (n *TomlKeyValueNode) Accept(visitor TomlNodeVisitor) {
	visitor.VisitTomlKeyValueNode(n)
}
