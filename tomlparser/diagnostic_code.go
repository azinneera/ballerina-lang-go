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

import (
	"fmt"
	"strings"

	"ballerina-lang-go/tools/diagnostics"
)

// TomlDiagnosticCode represents a diagnostic error code for TOML parsing.
// Each code uniquely identifies a specific type of error or warning.
// The naming and IDs are aligned with the Java implementation for consistency.
type TomlDiagnosticCode struct {
	id            string
	messageFormat string
	severity      diagnostics.DiagnosticSeverity
}

// Severity returns the severity level of this diagnostic code.
func (c TomlDiagnosticCode) Severity() diagnostics.DiagnosticSeverity {
	return c.severity
}

// DiagnosticId returns the unique identifier for this diagnostic code.
func (c TomlDiagnosticCode) DiagnosticId() string {
	return c.id
}

// MessageKey returns the message key for localization lookup.
func (c TomlDiagnosticCode) MessageKey() string {
	return c.messageFormat
}

// MessageFormat returns the format string for error messages.
func (c TomlDiagnosticCode) MessageFormat() string {
	return c.messageFormat
}

// Format creates a formatted error message using the provided arguments.
func (c TomlDiagnosticCode) Format(args ...any) string {
	if len(args) == 0 {
		return c.messageFormat
	}
	return fmt.Sprintf(c.messageFormat, args...)
}

// String returns a string representation of this diagnostic code.
func (c TomlDiagnosticCode) String() string {
	return c.id
}

// TOML Diagnostic Error Codes
// Names and IDs are aligned with Java implementation: DiagnosticErrorCode.java
var (
	// Generic syntax error - use when specific error cannot be determined
	ErrorSyntaxError = TomlDiagnosticCode{
		id:            "BCE0000",
		messageFormat: "syntax error",
		severity:      diagnostics.Error,
	}

	// Missing token errors
	ErrorMissingToken = TomlDiagnosticCode{
		id:            "BCE0001",
		messageFormat: "missing token",
		severity:      diagnostics.Error,
	}

	ErrorMissingOpenBracketToken = TomlDiagnosticCode{
		id:            "BCE0008",
		messageFormat: "missing open bracket token",
		severity:      diagnostics.Error,
	}

	ErrorMissingCloseBracketToken = TomlDiagnosticCode{
		id:            "BCE0009",
		messageFormat: "missing close bracket token",
		severity:      diagnostics.Error,
	}

	ErrorMissingEqualToken = TomlDiagnosticCode{
		id:            "BCE0010",
		messageFormat: "missing equal token",
		severity:      diagnostics.Error,
	}

	ErrorMissingCommaToken = TomlDiagnosticCode{
		id:            "BCE0011",
		messageFormat: "missing comma token",
		severity:      diagnostics.Error,
	}

	ErrorMissingPlusToken = TomlDiagnosticCode{
		id:            "BCE0012",
		messageFormat: "missing plus token",
		severity:      diagnostics.Error,
	}

	ErrorMissingNewLine = TomlDiagnosticCode{
		id:            "BCE0013",
		messageFormat: "missing new line",
		severity:      diagnostics.Error,
	}

	ErrorMissingOpenBraceToken = TomlDiagnosticCode{
		id:            "BCE0014",
		messageFormat: "missing open brace token",
		severity:      diagnostics.Error,
	}

	ErrorMissingCloseBraceToken = TomlDiagnosticCode{
		id:            "BCE0015",
		messageFormat: "missing close brace token",
		severity:      diagnostics.Error,
	}

	ErrorMissingDoubleQuoteToken = TomlDiagnosticCode{
		id:            "BCE0023",
		messageFormat: "missing double quote token",
		severity:      diagnostics.Error,
	}

	ErrorMissingTripleDoubleQuoteToken = TomlDiagnosticCode{
		id:            "BCE0024",
		messageFormat: "missing triple double quote token",
		severity:      diagnostics.Error,
	}

	ErrorMissingSingleQuoteToken = TomlDiagnosticCode{
		id:            "BCE0025",
		messageFormat: "missing single quote token",
		severity:      diagnostics.Error,
	}

	ErrorMissingTripleSingleQuoteToken = TomlDiagnosticCode{
		id:            "BCE0026",
		messageFormat: "missing triple single quote token",
		severity:      diagnostics.Error,
	}

	ErrorMissingDotToken = TomlDiagnosticCode{
		id:            "BCE0029",
		messageFormat: "missing dot token",
		severity:      diagnostics.Error,
	}

	// Keywords
	ErrorMissingTrueKeyword = TomlDiagnosticCode{
		id:            "BCE2014",
		messageFormat: "missing true keyword",
		severity:      diagnostics.Error,
	}

	ErrorMissingFalseKeyword = TomlDiagnosticCode{
		id:            "BCE2015",
		messageFormat: "missing false keyword",
		severity:      diagnostics.Error,
	}

	// Separators
	ErrorMissingHashToken = TomlDiagnosticCode{
		id:            "BCE2202",
		messageFormat: "missing hash token",
		severity:      diagnostics.Error,
	}

	// Operators
	ErrorMissingMinusToken = TomlDiagnosticCode{
		id:            "BCE2303",
		messageFormat: "missing minus token",
		severity:      diagnostics.Error,
	}

	// Literals
	ErrorMissingIdentifier = TomlDiagnosticCode{
		id:            "BCE2500",
		messageFormat: "missing identifier",
		severity:      diagnostics.Error,
	}

	ErrorMissingStringLiteral = TomlDiagnosticCode{
		id:            "BCE2501",
		messageFormat: "missing string literal",
		severity:      diagnostics.Error,
	}

	ErrorMissingDecimalIntegerLiteral = TomlDiagnosticCode{
		id:            "BCE2502",
		messageFormat: "missing decimal integer literal",
		severity:      diagnostics.Error,
	}

	ErrorMissingHexIntegerLiteral = TomlDiagnosticCode{
		id:            "BCE2503",
		messageFormat: "missing hex integer literal",
		severity:      diagnostics.Error,
	}

	ErrorMissingDecimalFloatingPointLiteral = TomlDiagnosticCode{
		id:            "BCE2504",
		messageFormat: "missing decimal floating point literal",
		severity:      diagnostics.Error,
	}

	ErrorMissingHexFloatingPointLiteral = TomlDiagnosticCode{
		id:            "BCE2505",
		messageFormat: "missing hex floating point literal",
		severity:      diagnostics.Error,
	}

	// Miscellaneous
	ErrorInvalidMetadata = TomlDiagnosticCode{
		id:            "BCE218",
		messageFormat: "invalid metadata",
		severity:      diagnostics.Error,
	}

	ErrorInvalidToken = TomlDiagnosticCode{
		id:            "BCE404",
		messageFormat: "invalid token '%s'",
		severity:      diagnostics.Error,
	}

	// Lexer errors
	ErrorLeadingZerosInNumericLiterals = TomlDiagnosticCode{
		id:            "BCE1000",
		messageFormat: "leading zeros in numeric literals",
		severity:      diagnostics.Error,
	}

	ErrorMissingDigitAfterExponentIndicator = TomlDiagnosticCode{
		id:            "BCE1001",
		messageFormat: "missing digit after exponent indicator",
		severity:      diagnostics.Error,
	}

	ErrorInvalidStringNumericEscapeSequence = TomlDiagnosticCode{
		id:            "BCE1002",
		messageFormat: "invalid string numeric escape sequence",
		severity:      diagnostics.Error,
	}

	ErrorInvalidEscapeSequence = TomlDiagnosticCode{
		id:            "BCE1003",
		messageFormat: "invalid escape sequence",
		severity:      diagnostics.Error,
	}

	ErrorMissingDoubleQuote = TomlDiagnosticCode{
		id:            "BCE1004",
		messageFormat: "missing double quote",
		severity:      diagnostics.Error,
	}

	ErrorMissingHexDigitAfterDot = TomlDiagnosticCode{
		id:            "BCE1005",
		messageFormat: "missing hex digit after dot",
		severity:      diagnostics.Error,
	}

	ErrorInvalidWhitespaceBefore = TomlDiagnosticCode{
		id:            "BCE1006",
		messageFormat: "invalid whitespace before",
		severity:      diagnostics.Error,
	}

	ErrorInvalidWhitespaceAfter = TomlDiagnosticCode{
		id:            "BCE1007",
		messageFormat: "invalid whitespace after",
		severity:      diagnostics.Error,
	}

	// Semantic errors
	ErrorMissingKey = TomlDiagnosticCode{
		id:            "BCE1500",
		messageFormat: "missing key",
		severity:      diagnostics.Error,
	}

	ErrorMissingValue = TomlDiagnosticCode{
		id:            "BCE1501",
		messageFormat: "missing value",
		severity:      diagnostics.Error,
	}

	ErrorMissingOpenDoubleBracket = TomlDiagnosticCode{
		id:            "BCE1502",
		messageFormat: "missing open double brackets",
		severity:      diagnostics.Error,
	}

	ErrorMissingCloseDoubleBracket = TomlDiagnosticCode{
		id:            "BCE1503",
		messageFormat: "missing close double brackets",
		severity:      diagnostics.Error,
	}

	ErrorExistingNode = TomlDiagnosticCode{
		id:            "BCE1504",
		messageFormat: "existing node '%s'",
		severity:      diagnostics.Error,
	}

	ErrorEmptyQuotedString = TomlDiagnosticCode{
		id:            "BCE1505",
		messageFormat: "empty quoted string",
		severity:      diagnostics.Error,
	}

	ErrorUnexpectedTopLevelNode = TomlDiagnosticCode{
		id:            "BCE1506",
		messageFormat: "unexpected top level node",
		severity:      diagnostics.Error,
	}

	// Warning codes
	WarningDeprecatedKey = TomlDiagnosticCode{
		id:            "BCW1000",
		messageFormat: "key '%s' is deprecated",
		severity:      diagnostics.Warning,
	}

	WarningUnusedKey = TomlDiagnosticCode{
		id:            "BCW1001",
		messageFormat: "key '%s' is not recognized",
		severity:      diagnostics.Warning,
	}
)

// errorPatterns maps common error message patterns to specific diagnostic codes.
// This is used to classify errors from the underlying BurntSushi parser.
var errorPatterns = []struct {
	pattern string
	code    TomlDiagnosticCode
}{
	// Duplicate key patterns
	{"duplicate key", ErrorExistingNode},
	{"key already exists", ErrorExistingNode},
	{"already defined", ErrorExistingNode},

	// Missing token patterns
	{"expected '='", ErrorMissingEqualToken},
	{"expected newline", ErrorMissingNewLine},
	{"expected ']'", ErrorMissingCloseBracketToken},
	{"expected '['", ErrorMissingOpenBracketToken},
	{"expected '}'", ErrorMissingCloseBraceToken},
	{"expected '{'", ErrorMissingOpenBraceToken},
	{"expected ','", ErrorMissingCommaToken},
	{"expected '\"'", ErrorMissingDoubleQuoteToken},
	{"expected \"'\"", ErrorMissingSingleQuoteToken},
	{"unclosed string", ErrorMissingDoubleQuote},
	{"unterminated string", ErrorMissingDoubleQuote},

	// Invalid patterns
	{"invalid escape", ErrorInvalidEscapeSequence},
	{"invalid character", ErrorInvalidToken},
	{"unexpected character", ErrorInvalidToken},
	{"leading zeros", ErrorLeadingZerosInNumericLiterals},
	{"leading zero", ErrorLeadingZerosInNumericLiterals},

	// Value errors
	{"expected value", ErrorMissingValue},
	{"expected key", ErrorMissingKey},
	{"expected identifier", ErrorMissingIdentifier},
}

// ClassifyError attempts to classify an error message into a specific diagnostic code.
// It returns the most appropriate TomlDiagnosticCode for the given error message.
func ClassifyError(errMsg string) TomlDiagnosticCode {
	lowerMsg := strings.ToLower(errMsg)

	for _, ep := range errorPatterns {
		if strings.Contains(lowerMsg, ep.pattern) {
			return ep.code
		}
	}

	return ErrorSyntaxError
}

// ExtractKeyFromError attempts to extract a key name from an error message.
// This is useful for formatting error messages with specific key information.
func ExtractKeyFromError(errMsg string) string {
	// Look for patterns like "key 'name'" or "key `name`" or "key "name""
	patterns := []struct {
		start string
		end   string
	}{
		{"'", "'"},
		{"`", "`"},
		{"\"", "\""},
		{"key ", " "},
		{"Key ", " "},
	}

	for _, p := range patterns {
		startIdx := strings.Index(errMsg, p.start)
		if startIdx >= 0 {
			startIdx += len(p.start)
			endIdx := strings.Index(errMsg[startIdx:], p.end)
			if endIdx > 0 {
				return errMsg[startIdx : startIdx+endIdx]
			}
		}
	}

	return ""
}

// AllDiagnosticCodes returns a list of all defined diagnostic codes.
// This is useful for documentation and testing purposes.
func AllDiagnosticCodes() []TomlDiagnosticCode {
	return []TomlDiagnosticCode{
		// Syntax errors
		ErrorSyntaxError,
		ErrorMissingToken,
		ErrorMissingOpenBracketToken,
		ErrorMissingCloseBracketToken,
		ErrorMissingEqualToken,
		ErrorMissingCommaToken,
		ErrorMissingPlusToken,
		ErrorMissingNewLine,
		ErrorMissingOpenBraceToken,
		ErrorMissingCloseBraceToken,
		ErrorMissingDoubleQuoteToken,
		ErrorMissingTripleDoubleQuoteToken,
		ErrorMissingSingleQuoteToken,
		ErrorMissingTripleSingleQuoteToken,
		ErrorMissingDotToken,
		// Keywords
		ErrorMissingTrueKeyword,
		ErrorMissingFalseKeyword,
		// Separators
		ErrorMissingHashToken,
		// Operators
		ErrorMissingMinusToken,
		// Literals
		ErrorMissingIdentifier,
		ErrorMissingStringLiteral,
		ErrorMissingDecimalIntegerLiteral,
		ErrorMissingHexIntegerLiteral,
		ErrorMissingDecimalFloatingPointLiteral,
		ErrorMissingHexFloatingPointLiteral,
		// Miscellaneous
		ErrorInvalidMetadata,
		ErrorInvalidToken,
		// Lexer errors
		ErrorLeadingZerosInNumericLiterals,
		ErrorMissingDigitAfterExponentIndicator,
		ErrorInvalidStringNumericEscapeSequence,
		ErrorInvalidEscapeSequence,
		ErrorMissingDoubleQuote,
		ErrorMissingHexDigitAfterDot,
		ErrorInvalidWhitespaceBefore,
		ErrorInvalidWhitespaceAfter,
		// Semantic errors
		ErrorMissingKey,
		ErrorMissingValue,
		ErrorMissingOpenDoubleBracket,
		ErrorMissingCloseDoubleBracket,
		ErrorExistingNode,
		ErrorEmptyQuotedString,
		ErrorUnexpectedTopLevelNode,
		// Warnings
		WarningDeprecatedKey,
		WarningUnusedKey,
	}
}
