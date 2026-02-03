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

// Package tree provides syntax tree types for TOML parsing.
// Java equivalent: io.ballerina.toml.syntax.tree
package tree

// SyntaxKind defines various kinds of syntax tree nodes, tokens and minutiae.
// The value represents the tag used for ordering and identification.
// Java equivalent: io.ballerina.toml.syntax.tree.SyntaxKind
type SyntaxKind int

const (
	// Basic tokens
	NONE        SyntaxKind = 0
	LIST        SyntaxKind = 1
	EOF_TOKEN   SyntaxKind = 2
	MODULE_PART SyntaxKind = 3
	INVALID     SyntaxKind = 4

	// Newline and comments (100s)
	NEWLINE    SyntaxKind = 100
	HASH_TOKEN SyntaxKind = 101

	// Keywords (200s)
	TRUE_KEYWORD                 SyntaxKind = 200
	FALSE_KEYWORD                SyntaxKind = 201
	STRING_LITERAL_TOKEN         SyntaxKind = 203
	DECIMAL_INT_TOKEN            SyntaxKind = 204
	DECIMAL_FLOAT_TOKEN          SyntaxKind = 205
	HEX_INTEGER_LITERAL_TOKEN    SyntaxKind = 206
	OCTAL_INTEGER_LITERAL_TOKEN  SyntaxKind = 207
	BINARY_INTEGER_LITERAL_TOKEN SyntaxKind = 208

	// Separators (500s)
	OPEN_BRACKET_TOKEN        SyntaxKind = 500
	CLOSE_BRACKET_TOKEN       SyntaxKind = 501
	DOUBLE_QUOTE_TOKEN        SyntaxKind = 502
	SINGLE_QUOTE_TOKEN        SyntaxKind = 503
	TRIPLE_DOUBLE_QUOTE_TOKEN SyntaxKind = 504
	TRIPLE_SINGLE_QUOTE_TOKEN SyntaxKind = 505
	OPEN_BRACE_TOKEN          SyntaxKind = 506
	CLOSE_BRACE_TOKEN         SyntaxKind = 507

	// Operators (520s)
	DOT_TOKEN   SyntaxKind = 520
	COMMA_TOKEN SyntaxKind = 521
	EQUAL_TOKEN SyntaxKind = 522
	PLUS_TOKEN  SyntaxKind = 523
	MINUS_TOKEN SyntaxKind = 524

	// Literals (1000s)
	IDENTIFIER_LITERAL SyntaxKind = 1000
	STRING_LITERAL     SyntaxKind = 1001
	LITERAL_STRING     SyntaxKind = 1002

	// Minutiae kinds (1500s)
	WHITESPACE_MINUTIAE         SyntaxKind = 1500
	END_OF_LINE_MINUTIAE        SyntaxKind = 1501
	COMMENT_MINUTIAE            SyntaxKind = 1502
	INVALID_NODE_MINUTIAE       SyntaxKind = 1503
	MARKDOWN_DOCUMENTATION_LINE SyntaxKind = 1504

	// Invalid nodes (1600s)
	INVALID_TOKEN               SyntaxKind = 1600
	MISSING_VALUE               SyntaxKind = 1601
	INVALID_TOKEN_MINUTIAE_NODE SyntaxKind = 1602

	// Structural nodes (2000s)
	KEY          SyntaxKind = 2000
	KEY_VALUE    SyntaxKind = 2001
	TABLE        SyntaxKind = 2002
	TABLE_ARRAY  SyntaxKind = 2003
	INLINE_TABLE SyntaxKind = 2004

	// Integer types (2010s)
	DEC_INT    SyntaxKind = 2010
	HEX_INT    SyntaxKind = 2011
	OCT_INT    SyntaxKind = 2012
	BINARY_INT SyntaxKind = 2013

	// Float types (2020s)
	FLOAT     SyntaxKind = 2020
	INF_TOKEN SyntaxKind = 2021
	NAN_TOKEN SyntaxKind = 2022

	// String types (2030s)
	ML_STRING_LITERAL SyntaxKind = 2030

	// Boolean type (2040s)
	BOOLEAN SyntaxKind = 2040

	// Date and Time types (2050s)
	OFFSET_DATE_TIME SyntaxKind = 2050
	LOCAL_DATE_TIME  SyntaxKind = 2051
	LOCAL_DATE       SyntaxKind = 2052
	LOCAL_TIME       SyntaxKind = 2053

	// Array type (2060s)
	ARRAY SyntaxKind = 2060
)

// syntaxKindText maps SyntaxKind to its string representation.
// Only kinds with non-empty string values are included.
// Java equivalent: SyntaxKind.strValue field
var syntaxKindText = map[SyntaxKind]string{
	NEWLINE:                   "\n",
	HASH_TOKEN:                "#",
	TRUE_KEYWORD:              "true",
	FALSE_KEYWORD:             "false",
	OPEN_BRACKET_TOKEN:        "[",
	CLOSE_BRACKET_TOKEN:       "]",
	DOUBLE_QUOTE_TOKEN:        "\"",
	SINGLE_QUOTE_TOKEN:        "'",
	TRIPLE_DOUBLE_QUOTE_TOKEN: "\"\"\"",
	TRIPLE_SINGLE_QUOTE_TOKEN: "'''",
	OPEN_BRACE_TOKEN:          "{",
	CLOSE_BRACE_TOKEN:         "}",
	DOT_TOKEN:                 ".",
	COMMA_TOKEN:               ",",
	EQUAL_TOKEN:               "=",
	PLUS_TOKEN:                "+",
	MINUS_TOKEN:               "-",
	INF_TOKEN:                 "inf",
	NAN_TOKEN:                 "nan",
}

// StringValue returns the string representation of the SyntaxKind if it has one.
// Returns the value and true if found, empty string and false otherwise.
// Java equivalent: SyntaxKind.stringValue()
func (s SyntaxKind) StringValue() (string, bool) {
	text, ok := syntaxKindText[s]
	return text, ok
}

// Tag returns the numeric tag of the SyntaxKind.
// Java equivalent: SyntaxKind.tag field
func (s SyntaxKind) Tag() int {
	return int(s)
}
