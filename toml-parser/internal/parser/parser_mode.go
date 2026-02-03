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

// Package parser provides the TOML lexer and parser implementation.
// Java equivalent: io.ballerina.toml.internal.parser
package parser

// ParserMode represents the modes of parsing.
// Java equivalent: io.ballerina.toml.internal.parser.ParserMode
type ParserMode int

const (
	// DEFAULT is the default parsing mode.
	// Java equivalent: ParserMode.DEFAULT
	DEFAULT ParserMode = iota

	// STRING is the mode for parsing double-quoted strings.
	// Java equivalent: ParserMode.STRING
	STRING

	// LITERAL_STRING is the mode for parsing single-quoted literal strings.
	// Java equivalent: ParserMode.LITERAL_STRING
	LITERAL_STRING

	// MULTILINE_STRING is the mode for parsing triple double-quoted multiline strings.
	// Java equivalent: ParserMode.MULTILINE_STRING
	MULTILINE_STRING

	// MULTILINE_LITERAL_STRING is the mode for parsing triple single-quoted multiline literal strings.
	// Java equivalent: ParserMode.MULTILINE_LITERAL_STRING
	MULTILINE_LITERAL_STRING

	// NEW_LINE is the mode for handling newlines.
	// Java equivalent: ParserMode.NEW_LINE
	NEW_LINE
)

// String returns the string representation of the ParserMode.
func (m ParserMode) String() string {
	switch m {
	case STRING:
		return "STRING"
	case LITERAL_STRING:
		return "LITERAL_STRING"
	case MULTILINE_STRING:
		return "MULTILINE_STRING"
	case MULTILINE_LITERAL_STRING:
		return "MULTILINE_LITERAL_STRING"
	case NEW_LINE:
		return "NEW_LINE"
	default:
		return "DEFAULT"
	}
}
