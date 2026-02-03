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

// Package diagnostics provides diagnostic error codes and messages for TOML parsing.
// Java equivalent: io.ballerina.toml.internal.diagnostics
package diagnostics

// DiagnosticSeverity represents the severity level of a diagnostic.
// Java equivalent: io.ballerina.tools.diagnostics.DiagnosticSeverity
type DiagnosticSeverity int

const (
	ERROR DiagnosticSeverity = iota
	WARNING
	INFO
	HINT
)

// DiagnosticCode defines the interface for diagnostic error codes.
// Java equivalent: io.ballerina.tools.diagnostics.DiagnosticCode
type DiagnosticCode interface {
	Severity() DiagnosticSeverity
	DiagnosticID() string
	MessageKey() string
}

// DiagnosticErrorCode represents a diagnostic error code.
// Java equivalent: io.ballerina.toml.internal.diagnostics.DiagnosticErrorCode
type DiagnosticErrorCode struct {
	diagnosticID string
	messageKey   string
}

// Severity returns the severity of the diagnostic (always ERROR).
// Java equivalent: DiagnosticErrorCode.severity()
func (d DiagnosticErrorCode) Severity() DiagnosticSeverity {
	return ERROR
}

// DiagnosticID returns the diagnostic identifier (e.g., "BCE0000").
// Java equivalent: DiagnosticErrorCode.diagnosticId()
func (d DiagnosticErrorCode) DiagnosticID() string {
	return d.diagnosticID
}

// MessageKey returns the message key for i18n lookup.
// Java equivalent: DiagnosticErrorCode.messageKey()
func (d DiagnosticErrorCode) MessageKey() string {
	return d.messageKey
}

// Diagnostic error code constants.
// Java equivalent: DiagnosticErrorCode enum members
var (
	// Generic syntax error
	ERROR_SYNTAX_ERROR = DiagnosticErrorCode{"BCE0000", "error.syntax.error"}

	// Missing tokens
	ERROR_MISSING_TOKEN                = DiagnosticErrorCode{"BCE0001", "error.missing.token"}
	ERROR_MISSING_OPEN_BRACKET_TOKEN   = DiagnosticErrorCode{"BCE0008", "error.missing.open.bracket.token"}
	ERROR_MISSING_CLOSE_BRACKET_TOKEN  = DiagnosticErrorCode{"BCE0009", "error.missing.close.bracket.token"}
	ERROR_MISSING_EQUAL_TOKEN          = DiagnosticErrorCode{"BCE0010", "error.missing.equal.token"}
	ERROR_MISSING_COMMA_TOKEN          = DiagnosticErrorCode{"BCE00011", "error.missing.comma.token"}
	ERROR_MISSING_PLUS_TOKEN           = DiagnosticErrorCode{"BCE00012", "error.missing.plus.token"}
	ERROR_MISSING_NEW_LINE             = DiagnosticErrorCode{"BCE00013", "error.missing.new.line"}
	ERROR_MISSING_OPEN_BRACE_TOKEN     = DiagnosticErrorCode{"BCE0014", "error.missing.open.brace.token"}
	ERROR_MISSING_CLOSE_BRACE_TOKEN    = DiagnosticErrorCode{"BCE0015", "error.missing.close.brace.token"}

	// Missing quote tokens
	ERROR_MISSING_DOUBLE_QUOTE_TOKEN        = DiagnosticErrorCode{"BCE00023", "error.missing.double.quote.token"}
	ERROR_MISSING_TRIPLE_DOUBLE_QUOTE_TOKEN = DiagnosticErrorCode{"BCE00024", "error.missing.triple.double.quote.token"}
	ERROR_MISSING_SINGLE_QUOTE_TOKEN        = DiagnosticErrorCode{"BCE00025", "error.missing.single.quote.token"}
	ERROR_MISSING_TRIPLE_SINGLE_QUOTE_TOKEN = DiagnosticErrorCode{"BCE00026", "error.missing.triple.single.quote.token"}
	ERROR_MISSING_DOT_TOKEN                 = DiagnosticErrorCode{"BCE00029", "error.missing.dot.token"}

	// Missing keywords
	ERROR_MISSING_TRUE_KEYWORD  = DiagnosticErrorCode{"BCE02014", "error.missing.true.keyword"}
	ERROR_MISSING_FALSE_KEYWORD = DiagnosticErrorCode{"BCE02015", "error.missing.false.keyword"}

	// Missing separators
	ERROR_MISSING_HASH_TOKEN = DiagnosticErrorCode{"BCE02202", "error.missing.hash.token"}

	// Missing operators
	ERROR_MISSING_MINUS_TOKEN = DiagnosticErrorCode{"BCE02303", "error.missing.minus.token"}

	// Missing literals
	ERROR_MISSING_IDENTIFIER                     = DiagnosticErrorCode{"BCE02500", "error.missing.identifier"}
	ERROR_MISSING_STRING_LITERAL                 = DiagnosticErrorCode{"BCE02501", "error.missing.string.literal"}
	ERROR_MISSING_DECIMAL_INTEGER_LITERAL        = DiagnosticErrorCode{"BCE02502", "error.missing.decimal.integer.literal"}
	ERROR_MISSING_HEX_INTEGER_LITERAL            = DiagnosticErrorCode{"BCE02503", "error.missing.hex.integer.literal"}
	ERROR_MISSING_DECIMAL_FLOATING_POINT_LITERAL = DiagnosticErrorCode{"BCE02504", "error.missing.decimal.floating.point.literal"}
	ERROR_MISSING_HEX_FLOATING_POINT_LITERAL     = DiagnosticErrorCode{"BCE02505", "error.missing.hex.floating.point.literal"}

	// Miscellaneous errors
	ERROR_INVALID_METADATA = DiagnosticErrorCode{"BCE218", "error.invalid.metadata"}
	ERROR_INVALID_TOKEN    = DiagnosticErrorCode{"BCE404", "error.invalid.token"}

	// Lexer errors
	ERROR_LEADING_ZEROS_IN_NUMERIC_LITERALS    = DiagnosticErrorCode{"BCE1000", "error.leading.zeros.in.numeric.literals"}
	ERROR_MISSING_DIGIT_AFTER_EXPONENT_INDICATOR = DiagnosticErrorCode{"BCE1001", "error.missing.digit.after.exponent.indicator"}
	ERROR_INVALID_STRING_NUMERIC_ESCAPE_SEQUENCE = DiagnosticErrorCode{"BCE1002", "error.invalid.string.numeric.escape.sequence"}
	ERROR_INVALID_ESCAPE_SEQUENCE              = DiagnosticErrorCode{"BCE1003", "error.invalid.escape.sequence"}
	ERROR_MISSING_DOUBLE_QUOTE                 = DiagnosticErrorCode{"BCE1004", "error.missing.double.quote"}
	ERROR_MISSING_HEX_DIGIT_AFTER_DOT          = DiagnosticErrorCode{"BCE1005", "error.missing.hex.digit.after.dot"}
	ERROR_INVALID_WHITESPACE_BEFORE            = DiagnosticErrorCode{"BCE1006", "error.invalid.whitespace.before"}
	ERROR_INVALID_WHITESPACE_AFTER             = DiagnosticErrorCode{"BCE1007", "error.invalid.whitespace.after"}

	// Semantic errors
	ERROR_MISSING_KEY                  = DiagnosticErrorCode{"BCE1500", "error.missing.key"}
	ERROR_MISSING_VALUE                = DiagnosticErrorCode{"BCE1501", "error.missing.value"}
	ERROR_MISSING_OPEN_DOUBLE_BRACKET  = DiagnosticErrorCode{"BCE1502", "error.missing.open.double.bracket"}
	ERROR_MISSING_CLOSE_DOUBLE_BRACKET = DiagnosticErrorCode{"BCE1503", "error.missing.close.double.bracket"}
	ERROR_EXISTING_NODE                = DiagnosticErrorCode{"BCE1504", "error.existing.node"}
	ERROR_EMPTY_QUOTED_STRING          = DiagnosticErrorCode{"BCE1505", "error.empty.quoted.string"}
	ERROR_UNEXPECTED_TOP_LEVEL_NODE    = DiagnosticErrorCode{"BCE1506", "error.unexpected.top.level.node"}
)
