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
	"ballerina-lang-go/tools/diagnostics"
)

// Validator Diagnostic Codes (TVE prefix)
// These codes map to the inline diagnostic codes used in Java validator classes:
// - io.ballerina.toml.validator.schema.Schema.java
// - io.ballerina.toml.validator.SchemaValidator.java
// - io.ballerina.toml.validator.schema.StringSchema.java
// - io.ballerina.toml.validator.schema.NumericSchema.java
// - io.ballerina.toml.validator.schema.CompositionSchema.java
//
// Note: Unlike DiagnosticErrorCode.java (BCE prefix), these codes are not
// defined as a formal enum in Java but are used inline in validator classes.

var (
	// ErrorSchemaValidation - generic schema validation error (Go-specific)
	// Used when no specific TVE code applies
	ErrorSchemaValidation = TomlDiagnosticCode{
		id:            "TVE0000",
		messageFormat: "schema validation failed: %s",
		severity:      diagnostics.Error,
	}

	// ErrorUnexpectedProperty maps to TVE0001 in Schema.java:137,166
	// Java: error.unexpected.property
	// Message: "key '%s' not supported in schema '%s'"
	ErrorUnexpectedProperty = TomlDiagnosticCode{
		id:            "TVE0001",
		messageFormat: "key '%s' not supported in schema '%s'",
		severity:      diagnostics.Error,
	}

	// ErrorInvalidType maps to TVE0002 in SchemaValidator.java:70,94,148,162,175,188,208,222
	// Also used in StringSchema.java:79, NumericSchema.java:79, BooleanSchema.java:55, etc.
	// Java: error.invalid.type
	// Message: "incompatible type for key '%s': expected '%s', found '%s'"
	ErrorInvalidType = TomlDiagnosticCode{
		id:            "TVE0002",
		messageFormat: "incompatible type for key '%s': expected '%s', found '%s'",
		severity:      diagnostics.Error,
	}

	// ErrorRegexMismatch maps to TVE0003 in StringSchema.java:91
	// Java: error.regex.mismatch
	// Message: "value for key '%s' expected to match the regex: %s"
	ErrorRegexMismatch = TomlDiagnosticCode{
		id:            "TVE0003",
		messageFormat: "value for key '%s' expected to match the regex: %s",
		severity:      diagnostics.Error,
	}

	// ErrorMinimumValueDeceed maps to TVE0004 in NumericSchema.java:98
	// Java: error.minimum.value.deceed
	// Message: "value for key '%s' can't be lower than %f"
	ErrorMinimumValueDeceed = TomlDiagnosticCode{
		id:            "TVE0004",
		messageFormat: "value for key '%s' can't be lower than %v",
		severity:      diagnostics.Error,
	}

	// ErrorMaximumValueExceed maps to TVE0005 in NumericSchema.java:90
	// Java: error.maximum.value.exceed
	// Message: "value for key '%s' can't be higher than %f"
	ErrorMaximumValueExceed = TomlDiagnosticCode{
		id:            "TVE0005",
		messageFormat: "value for key '%s' can't be higher than %v",
		severity:      diagnostics.Error,
	}

	// ErrorRequiredFieldMissing maps to TVE0006 in Schema.java:146,174
	// Java: error.required.field.missing
	// Message: "missing required field '%s'"
	ErrorRequiredFieldMissing = TomlDiagnosticCode{
		id:            "TVE0006",
		messageFormat: "missing required field '%s'",
		severity:      diagnostics.Error,
	}

	// ErrorMaxLengthExceeded maps to TVE0007 in StringSchema.java:101
	// Java: error.maxlen.exceeded
	// Message: "length of the value for key '%s' is greater than defined max length %s"
	// Note: TVE0007 is also used for anyOf in CompositionSchema.java:74 (code overlap in Java)
	ErrorMaxLengthExceeded = TomlDiagnosticCode{
		id:            "TVE0007",
		messageFormat: "length of the value for key '%s' is greater than defined max length %d",
		severity:      diagnostics.Error,
	}

	// ErrorMinLengthDeceed maps to TVE0008 in StringSchema.java:110
	// Java: error.minlen.deceed
	// Message: "length of the value for key '%s' is lower than defined min length %s"
	// Note: TVE0008 is also used for oneOf in CompositionSchema.java:83 (code overlap in Java)
	ErrorMinLengthDeceed = TomlDiagnosticCode{
		id:            "TVE0008",
		messageFormat: "length of the value for key '%s' is lower than defined min length %d",
		severity:      diagnostics.Error,
	}

	// ErrorSchemaRuleMustNotValid maps to TVE0009 in CompositionSchema.java:107
	// Java: error.schema.rule.must.not.valid
	// Message: "schema rules must `NOT` be valid"
	ErrorSchemaRuleMustNotValid = TomlDiagnosticCode{
		id:            "TVE0009",
		messageFormat: "schema rules must NOT be valid",
		severity:      diagnostics.Error,
	}

	// Go-specific validator codes (no Java equivalent)
	// These are needed because Go uses a different JSON Schema library (santhosh-tekuri/jsonschema)

	// ErrorEnumValue - value must match enum
	ErrorEnumValue = TomlDiagnosticCode{
		id:            "TVE0010",
		messageFormat: "value must be one of: %s",
		severity:      diagnostics.Error,
	}

	// ErrorMinItems - array minimum items
	ErrorMinItems = TomlDiagnosticCode{
		id:            "TVE0011",
		messageFormat: "array must have at least %d items",
		severity:      diagnostics.Error,
	}

	// ErrorMaxItems - array maximum items
	ErrorMaxItems = TomlDiagnosticCode{
		id:            "TVE0012",
		messageFormat: "array must have at most %d items",
		severity:      diagnostics.Error,
	}

	// ErrorUniqueItems - array unique items
	ErrorUniqueItems = TomlDiagnosticCode{
		id:            "TVE0013",
		messageFormat: "array items must be unique",
		severity:      diagnostics.Error,
	}
)

// classifyValidationError maps JSON Schema validation keywords to diagnostic codes.
// This function is used by validator.go to determine the appropriate TVE code
// based on the validation keyword from the jsonschema library.
func classifyValidationError(keyword string, instancePath []string) TomlDiagnosticCode {
	switch keyword {
	case "required":
		return ErrorRequiredFieldMissing
	case "type":
		return ErrorInvalidType
	case "pattern":
		return ErrorRegexMismatch
	case "minLength":
		return ErrorMinLengthDeceed
	case "maxLength":
		return ErrorMaxLengthExceeded
	case "minimum", "exclusiveMinimum":
		return ErrorMinimumValueDeceed
	case "maximum", "exclusiveMaximum":
		return ErrorMaximumValueExceed
	case "enum":
		return ErrorEnumValue
	case "additionalProperties":
		return ErrorUnexpectedProperty
	case "minItems":
		return ErrorMinItems
	case "maxItems":
		return ErrorMaxItems
	case "uniqueItems":
		return ErrorUniqueItems
	case "not":
		return ErrorSchemaRuleMustNotValid
	default:
		return ErrorSchemaValidation
	}
}

// AllValidatorDiagnosticCodes returns a list of all validator diagnostic codes.
// This is useful for documentation and testing purposes.
func AllValidatorDiagnosticCodes() []TomlDiagnosticCode {
	return []TomlDiagnosticCode{
		ErrorSchemaValidation,        // TVE0000
		ErrorUnexpectedProperty,      // TVE0001
		ErrorInvalidType,             // TVE0002
		ErrorRegexMismatch,           // TVE0003
		ErrorMinimumValueDeceed,      // TVE0004
		ErrorMaximumValueExceed,      // TVE0005
		ErrorRequiredFieldMissing,    // TVE0006
		ErrorMaxLengthExceeded,       // TVE0007
		ErrorMinLengthDeceed,         // TVE0008
		ErrorSchemaRuleMustNotValid,  // TVE0009
		ErrorEnumValue,               // TVE0010
		ErrorMinItems,                // TVE0011
		ErrorMaxItems,                // TVE0012
		ErrorUniqueItems,             // TVE0013
	}
}
