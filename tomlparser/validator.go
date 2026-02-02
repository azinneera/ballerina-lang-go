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
	"errors"
	"fmt"
	"strings"

	"ballerina-lang-go/tools/diagnostics"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Validator interface {
	Validate(toml *Toml) error
}

type validatorImpl struct {
	schema       Schema
	schemaSource map[string]any // Original schema JSON for custom message lookup
}

func NewValidator(schema Schema) Validator {
	return &validatorImpl{
		schema: schema,
	}
}

// NewValidatorWithSource creates a validator with access to schema source for custom messages.
func NewValidatorWithSource(schema Schema, schemaSource map[string]any) Validator {
	return &validatorImpl{
		schema:       schema,
		schemaSource: schemaSource,
	}
}

func (v *validatorImpl) Validate(toml *Toml) error {
	if toml == nil {
		return fmt.Errorf("toml document is nil")
	}

	data := toml.ToMap()

	if err := v.schema.Validate(data); err != nil {
		// Try to extract custom error messages from validation errors
		diags := v.extractDiagnostics(err, toml)
		toml.diagnostics = append(toml.diagnostics, diags...)
		return err
	}

	return nil
}

// extractDiagnostics converts jsonschema validation errors to diagnostics with custom messages.
func (v *validatorImpl) extractDiagnostics(err error, toml *Toml) []Diagnostic {
	var validationErr *jsonschema.ValidationError
	if !errors.As(err, &validationErr) {
		// Not a validation error, return generic diagnostic
		return []Diagnostic{{
			Message:  err.Error(),
			Severity: diagnostics.Error,
		}}
	}

	// Extract individual validation failures
	var diags []Diagnostic
	v.collectValidationErrors(validationErr, toml, &diags, make(map[string]bool))
	return diags
}

// collectValidationErrors recursively collects validation errors with custom messages.
func (v *validatorImpl) collectValidationErrors(err *jsonschema.ValidationError, toml *Toml, diags *[]Diagnostic, seen map[string]bool) {
	// Process leaf errors (those without causes)
	if len(err.Causes) == 0 {
		// Get the keyword path to determine error type
		keywordPath := err.ErrorKind.KeywordPath()
		if len(keywordPath) == 0 {
			return
		}
		keyword := keywordPath[len(keywordPath)-1]

		// Skip allOf wrapper errors - they just indicate child failures
		if keyword == "allOf" {
			return
		}

		// Create unique key to avoid duplicates
		instancePath := strings.Join(err.InstanceLocation, "/")
		key := fmt.Sprintf("%s:%s", instancePath, keyword)
		if seen[key] {
			return
		}
		seen[key] = true

		// Try to get custom message from schema using SchemaURL
		customMsg := v.getCustomMessageFromURL(err.SchemaURL, keyword)
		if customMsg == "" {
			customMsg = err.Error()
		}

		// Get location from toml metadata if available
		loc := v.getLocationForPath(err.InstanceLocation, toml)

		*diags = append(*diags, Diagnostic{
			Message:  customMsg,
			Severity: diagnostics.Error,
			Location: loc,
		})
		return
	}

	// Recurse into causes
	for _, cause := range err.Causes {
		v.collectValidationErrors(cause, toml, diags, seen)
	}
}

// getCustomMessageFromURL attempts to get a custom error message using the SchemaURL.
// SchemaURL format: file:///path/to/schema.json#/properties/package/properties/org
func (v *validatorImpl) getCustomMessageFromURL(schemaURL string, keyword string) string {
	if v.schemaSource == nil {
		return ""
	}

	// Extract the fragment (JSON Pointer path) from the URL
	hashIdx := strings.Index(schemaURL, "#")
	if hashIdx < 0 {
		return ""
	}
	fragment := schemaURL[hashIdx+1:]
	if fragment == "" || fragment == "/" {
		return ""
	}

	// Parse the JSON Pointer path
	parts := strings.Split(strings.TrimPrefix(fragment, "/"), "/")

	// Navigate the schema using the path
	current := v.schemaSource
	for _, part := range parts {
		if part == "" {
			continue
		}
		next, ok := current[part]
		if !ok {
			return ""
		}
		nextMap, ok := next.(map[string]any)
		if !ok {
			// Try to handle array indices
			if arr, ok := next.([]any); ok {
				// Parse index
				var idx int
				if _, err := fmt.Sscanf(part, "%d", &idx); err == nil && idx < len(arr) {
					if m, ok := arr[idx].(map[string]any); ok {
						current = m
						continue
					}
				}
			}
			return ""
		}
		current = nextMap
	}

	// Look for "message" object with the keyword as key
	messageObj, ok := current["message"]
	if !ok {
		return ""
	}

	msgMap, ok := messageObj.(map[string]any)
	if !ok {
		return ""
	}

	customMsg, ok := msgMap[keyword]
	if !ok {
		return ""
	}

	if msg, ok := customMsg.(string); ok {
		return msg
	}
	return ""
}

// getLocationForPath converts a JSON path to a TOML location using metadata.
func (v *validatorImpl) getLocationForPath(path []string, toml *Toml) *Location {
	if len(path) == 0 {
		return nil
	}

	// Convert path to TOML key path
	tomlKey := strings.Join(path, ".")

	// Try to find the key in TOML metadata
	if toml.metadata.IsDefined(path...) {
		// Find location in content by parsing
		return v.findKeyLocationInContent(toml.content, path)
	}

	// Try to find by key name
	_ = tomlKey
	return nil
}

// findKeyLocationInContent finds the location of a key in TOML content.
func (v *validatorImpl) findKeyLocationInContent(content string, key []string) *Location {
	lines := strings.Split(content, "\n")
	targetKey := key[len(key)-1] // The last part of the key path

	for lineNum, line := range lines {
		// Look for key = value pattern
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, targetKey) {
			// Check if this is a key assignment
			rest := strings.TrimPrefix(trimmed, targetKey)
			rest = strings.TrimSpace(rest)
			if strings.HasPrefix(rest, "=") {
				// Found the key, now find the value location
				eqIdx := strings.Index(line, "=")
				if eqIdx >= 0 {
					valueStart := eqIdx + 1
					// Skip whitespace after =
					for valueStart < len(line) && (line[valueStart] == ' ' || line[valueStart] == '\t') {
						valueStart++
					}

					// Handle quoted strings - include the quotes in the span
					if valueStart < len(line) && line[valueStart] == '"' {
						startCol := valueStart // Opening quote position
						// Find closing quote
						endCol := strings.LastIndex(line, "\"")
						if endCol > startCol {
							return &Location{
								StartLine:   lineNum,
								StartColumn: startCol,
								EndLine:     lineNum,
								EndColumn:   endCol + 1, // One past the closing quote (exclusive)
							}
						}
					}
				}
			}
		}
	}

	return nil
}
