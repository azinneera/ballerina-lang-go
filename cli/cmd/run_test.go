/*
 * Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"bytes"
	"strings"
	"testing"

	"ballerina-lang-go/cli/pkg/templates"
)

func TestRunCommandReturnsError(t *testing.T) {
	cmd := NewRunCommand()
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Fatal("Args validation should return error for empty args")
	}
}

func TestRunCommandHelpText(t *testing.T) {
	cmd := NewRunCommand()

	buf := new(bytes.Buffer)
	templates.PrintCommandHelpToWriter(buf, cmd)

	output := buf.String()
	expectedSubstrings := []string{
		"NAME",
		"ballerina-run",
		"SYNOPSIS",
		"bal run",
		"DESCRIPTION",
		"Compile and run the current package",
		"OPTIONS",
		"EXAMPLES",
		"$ bal run",
		"$ bal run main.bal",
	}
	for _, sub := range expectedSubstrings {
		if !strings.Contains(output, sub) {
			t.Errorf("run help output missing %q", sub)
		}
	}
}

func TestRunCommandHelpLineWidth(t *testing.T) {
	cmd := NewRunCommand()

	buf := new(bytes.Buffer)
	templates.PrintCommandHelpToWriter(buf, cmd)

	output := buf.String()
	lines := strings.Split(output, "\n")

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip command examples, NAME line, and SYNOPSIS line (matches original Java help)
		if strings.HasPrefix(trimmed, "$") || strings.HasPrefix(trimmed, "ballerina-") ||
			strings.HasPrefix(trimmed, "bal run") {
			continue
		}
		if len(line) > 80 {
			t.Errorf("line %d exceeds 80 characters (len=%d): %q", i+1, len(line), line)
		}
	}
}

func TestRunCommandHasExamples(t *testing.T) {
	cmd := NewRunCommand()
	if cmd.Example == "" {
		t.Error("run command should have examples")
	}
}
