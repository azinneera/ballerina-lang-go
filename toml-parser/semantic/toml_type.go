// Copyright (c) 2020, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
//
// WSO2 Inc. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Package semantic provides semantic AST types for TOML parsing.
// Java equivalent: io.ballerina.toml.semantic
package semantic

// TomlType defines various kinds of semantic TOML nodes.
// Java equivalent: io.ballerina.toml.semantic.TomlType
type TomlType int

const (
	// NONE represents no type (default zero value).
	// Java equivalent: TomlType.NONE
	NONE TomlType = iota

	// UNQUOTED_KEY represents an unquoted key.
	// Java equivalent: TomlType.UNQUOTED_KEY
	UNQUOTED_KEY

	// TABLE represents a table.
	// Java equivalent: TomlType.TABLE
	TABLE

	// KEY_VALUE represents a key-value pair.
	// Java equivalent: TomlType.KEY_VALUE
	KEY_VALUE

	// TABLE_ARRAY represents an array of tables.
	// Java equivalent: TomlType.TABLE_ARRAY
	TABLE_ARRAY

	// STRING represents a string value.
	// Java equivalent: TomlType.STRING
	STRING

	// INTEGER represents an integer value.
	// Java equivalent: TomlType.INTEGER
	INTEGER

	// DOUBLE represents a double/float value.
	// Java equivalent: TomlType.DOUBLE
	DOUBLE

	// BOOLEAN represents a boolean value.
	// Java equivalent: TomlType.BOOLEAN
	BOOLEAN

	// ARRAY represents an array value.
	// Java equivalent: TomlType.ARRAY
	ARRAY

	// INLINE_TABLE represents an inline table value.
	// Java equivalent: TomlType.INLINE_TABLE
	INLINE_TABLE
)

// String returns the string representation of the TomlType.
func (t TomlType) String() string {
	switch t {
	case UNQUOTED_KEY:
		return "UNQUOTED_KEY"
	case TABLE:
		return "TABLE"
	case KEY_VALUE:
		return "KEY_VALUE"
	case TABLE_ARRAY:
		return "TABLE_ARRAY"
	case STRING:
		return "STRING"
	case INTEGER:
		return "INTEGER"
	case DOUBLE:
		return "DOUBLE"
	case BOOLEAN:
		return "BOOLEAN"
	case ARRAY:
		return "ARRAY"
	case INLINE_TABLE:
		return "INLINE_TABLE"
	default:
		return "NONE"
	}
}
