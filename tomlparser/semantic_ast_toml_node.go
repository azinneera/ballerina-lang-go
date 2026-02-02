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

// TomlNode is the base interface for all semantic AST nodes.
// Maps to: io.ballerina.toml.semantic.ast.TomlNode
type TomlNode interface {
	Kind() TomlType
	Location() *Location
	Diagnostics() []Diagnostic
	AddDiagnostic(diag Diagnostic)
	Accept(visitor TomlNodeVisitor)
}

// tomlNodeBase provides common functionality for semantic nodes.
// Maps to: io.ballerina.toml.semantic.ast.TomlNode (base fields)
type tomlNodeBase struct {
	location    *Location
	diagnostics []Diagnostic
}

func (n *tomlNodeBase) Location() *Location {
	return n.location
}

func (n *tomlNodeBase) Diagnostics() []Diagnostic {
	return n.diagnostics
}

func (n *tomlNodeBase) AddDiagnostic(diag Diagnostic) {
	n.diagnostics = append(n.diagnostics, diag)
}
