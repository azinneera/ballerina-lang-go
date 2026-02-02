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

// TomlType represents the type of a semantic AST node.
// Maps to: io.ballerina.toml.semantic.TomlType
type TomlType int

const (
	// TomlTypeNone maps to: TomlType.NONE
	TomlTypeNone TomlType = iota
	// TomlTypeUnquotedKey maps to: TomlType.UNQUOTED_KEY (for unquoted key identifiers)
	TomlTypeUnquotedKey
	// TomlTypeTable maps to: TomlType.TABLE
	TomlTypeTable
	// TomlTypeKeyValue maps to: TomlType.KEY_VALUE
	TomlTypeKeyValue
	// TomlTypeTableArray maps to: TomlType.TABLE_ARRAY
	TomlTypeTableArray
	// TomlTypeString maps to: TomlType.STRING
	TomlTypeString
	// TomlTypeInteger maps to: TomlType.INTEGER
	TomlTypeInteger
	// TomlTypeDouble maps to: TomlType.DOUBLE
	TomlTypeDouble
	// TomlTypeBoolean maps to: TomlType.BOOLEAN
	TomlTypeBoolean
	// TomlTypeArray maps to: TomlType.ARRAY
	TomlTypeArray
	// TomlTypeInlineTable maps to: TomlType.INLINE_TABLE
	TomlTypeInlineTable
)
