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

// TomlStringValueNode represents a string value.
// Maps to: io.ballerina.toml.semantic.ast.TomlStringValueNode
type TomlStringValueNode struct {
	TomlBasicValueNode
}

// NewTomlStringValueNode creates a new string value node.
func NewTomlStringValueNode(value string, location *Location) *TomlStringValueNode {
	return &TomlStringValueNode{
		TomlBasicValueNode: TomlBasicValueNode{
			tomlNodeBase: tomlNodeBase{
				location:    location,
				diagnostics: make([]Diagnostic, 0),
			},
			kind:  TomlTypeString,
			value: value,
		},
	}
}

// Accept implements the visitor pattern.
// Maps to: TomlStringValueNode.accept(TomlNodeVisitor)
func (n *TomlStringValueNode) Accept(visitor TomlNodeVisitor) {
	visitor.VisitTomlStringValueNode(n)
}
