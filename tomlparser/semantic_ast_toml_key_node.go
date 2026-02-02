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

// TomlKeyNode represents the key part of a key-value pair.
// Maps to: io.ballerina.toml.semantic.ast.TomlKeyNode
type TomlKeyNode struct {
	tomlNodeBase
	keys []TomlKeyEntryNode
}

// NewTomlKeyNode creates a new key node.
// Maps to: TomlKeyNode constructor
func NewTomlKeyNode(keys []TomlKeyEntryNode, location *Location) *TomlKeyNode {
	return &TomlKeyNode{
		tomlNodeBase: tomlNodeBase{
			location:    location,
			diagnostics: make([]Diagnostic, 0),
		},
		keys: keys,
	}
}

// Kind returns TomlTypeKeyValue.
// Maps to: TomlKeyNode constructor which passes TomlType.KEY_VALUE to super()
func (n *TomlKeyNode) Kind() TomlType {
	return TomlTypeKeyValue
}

func (n *TomlKeyNode) Keys() []TomlKeyEntryNode {
	return n.keys
}

// Name returns the full key name as a dotted string.
// Maps to: TomlKeyNode.name()
func (n *TomlKeyNode) Name() string {
	if len(n.keys) == 0 {
		return ""
	}
	names := make([]string, len(n.keys))
	for i, key := range n.keys {
		names[i] = key.Name()
	}
	result := names[0]
	for i := 1; i < len(names); i++ {
		result += "." + names[i]
	}
	return result
}

// Accept implements the visitor pattern.
// Maps to: TomlKeyNode.accept(TomlNodeVisitor)
func (n *TomlKeyNode) Accept(visitor TomlNodeVisitor) {
	visitor.VisitTomlKeyNode(n)
}
