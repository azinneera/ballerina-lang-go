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

// TomlKeyEntryNode represents a single key entry in a dotted key.
// Maps to: io.ballerina.toml.semantic.ast.TomlKeyEntryNode
type TomlKeyEntryNode struct {
	tomlNodeBase
	name string
	kind TomlType
}

// NewTomlKeyEntryNode creates a new key entry node.
// Maps to: TomlKeyEntryNode constructor
func NewTomlKeyEntryNode(name string, quoted bool, location *Location) TomlKeyEntryNode {
	kind := TomlTypeKeyValue
	if quoted {
		kind = TomlTypeString
	}
	return TomlKeyEntryNode{
		tomlNodeBase: tomlNodeBase{
			location:    location,
			diagnostics: make([]Diagnostic, 0),
		},
		name: name,
		kind: kind,
	}
}

func (n TomlKeyEntryNode) Kind() TomlType {
	return n.kind
}

func (n TomlKeyEntryNode) Name() string {
	return n.name
}

// Accept implements the visitor pattern.
// Maps to: TomlKeyEntryNode.accept(TomlNodeVisitor)
func (n TomlKeyEntryNode) Accept(visitor TomlNodeVisitor) {
	visitor.VisitTomlKeyEntryNode(&n)
}
