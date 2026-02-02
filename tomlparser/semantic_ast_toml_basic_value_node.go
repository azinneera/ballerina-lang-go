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

// TomlBasicValueNode represents a primitive value (string, int, float, bool).
// Maps to: io.ballerina.toml.semantic.ast.TomlBasicValueNode<T>
type TomlBasicValueNode struct {
	tomlNodeBase
	kind  TomlType
	value any
}

func (n *TomlBasicValueNode) tomlValueNode() {}

// NewTomlBasicValueNode creates a new basic value node.
func NewTomlBasicValueNode(kind TomlType, value any, location *Location) *TomlBasicValueNode {
	return &TomlBasicValueNode{
		tomlNodeBase: tomlNodeBase{
			location:    location,
			diagnostics: make([]Diagnostic, 0),
		},
		kind:  kind,
		value: value,
	}
}

func (n *TomlBasicValueNode) Kind() TomlType {
	return n.kind
}

func (n *TomlBasicValueNode) Value() any {
	return n.value
}
