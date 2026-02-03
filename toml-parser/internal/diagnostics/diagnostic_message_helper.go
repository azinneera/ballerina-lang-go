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

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed resources/toml_diagnostic_message.properties
var diagnosticMessagesFile string

// diagnosticMessages contains the parsed diagnostic message templates.
// Java equivalent: ResourceBundle loaded from toml_diagnostic_message.properties
var diagnosticMessages map[string]string

func init() {
	diagnosticMessages = parseProperties(diagnosticMessagesFile)
}

// parseProperties parses a Java-style .properties file content into a map.
func parseProperties(content string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Find the first = sign
		idx := strings.Index(line, "=")
		if idx == -1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		// Convert Java MessageFormat placeholders {0}, {1} to Go fmt placeholders %s
		value = convertPlaceholders(value)
		result[key] = value
	}
	return result
}

// convertPlaceholders converts Java MessageFormat placeholders to Go fmt placeholders.
// e.g., "invalid token ''{0}''" -> "invalid token '%s'"
func convertPlaceholders(s string) string {
	// Replace '' with a placeholder to preserve escaped single quotes
	s = strings.ReplaceAll(s, "''", "\x00")
	// Replace {0}, {1}, etc. with %s
	for i := 0; i < 10; i++ {
		s = strings.ReplaceAll(s, fmt.Sprintf("{%d}", i), "%s")
	}
	// Restore escaped single quotes
	s = strings.ReplaceAll(s, "\x00", "'")
	return s
}

// GetDiagnosticMessage returns the formatted diagnostic message for a given code.
// Java equivalent: DiagnosticMessageHelper.getDiagnosticMessage()
func GetDiagnosticMessage(code DiagnosticCode, args ...any) string {
	msgKey := code.MessageKey()
	msgTemplate, ok := diagnosticMessages[msgKey]
	if !ok {
		return msgKey
	}
	if len(args) == 0 {
		return msgTemplate
	}
	return fmt.Sprintf(msgTemplate, args...)
}
