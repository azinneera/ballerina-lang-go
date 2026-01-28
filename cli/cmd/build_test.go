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

	"github.com/spf13/cobra"
)

func TestBuildCommandExecuteReturnsNotImplemented(t *testing.T) {
	bc := NewBuildCommand()
	err := bc.Execute([]string{})
	if err == nil {
		t.Fatal("Execute() should return an error")
	}
	expected := "command 'build' is not yet implemented"
	if err.Error() != expected {
		t.Errorf("Execute() error = %q, want %q", err.Error(), expected)
	}
}

func TestBuildCommandHelpText(t *testing.T) {
	bc := NewBuildCommand()
	cmd := bc.CobraCommand()

	buf := new(bytes.Buffer)
	templates.PrintCommandHelpToWriter(buf, cmd)

	output := buf.String()
	if output == "" {
		t.Fatal("expected help output, got empty string")
	}

	// Check for required sections matching the expected format
	expectedSubstrings := []string{
		"NAME",
		"ballerina-build - Compiles the current package",
		"SYNOPSIS",
		"bal build [OPTIONS] [<package>|<source-file>]",
		"DESCRIPTION",
		"Compile a package and its dependencies",
		"OPTIONS",
		"--offline",
		"--graalvm",
		"--target-dir <path>",  // Verify placeholder works
		"--cloud <provider>",   // Verify placeholder works
		"EXAMPLES",
		"$ bal build",
		"$ bal build app.bal",
	}
	for _, sub := range expectedSubstrings {
		if !strings.Contains(output, sub) {
			t.Errorf("help output missing %q\ngot:\n%s", sub, output)
		}
	}
}

func TestBuildCommandHelpLineWidth(t *testing.T) {
	bc := NewBuildCommand()
	cmd := bc.CobraCommand()

	buf := new(bytes.Buffer)
	templates.PrintCommandHelpToWriter(buf, cmd)

	output := buf.String()
	lines := strings.Split(output, "\n")

	for i, line := range lines {
		// Skip lines starting with "$ " as these are example commands that shouldn't be wrapped
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "$") {
			continue
		}
		if len(line) > 80 {
			t.Errorf("line %d exceeds 80 characters (len=%d): %q", i+1, len(line), line)
		}
	}
}

func TestBuildCommandHelpSanitizesInput(t *testing.T) {
	// Test that the help system handles poorly formatted input from cobra command
	cmd := &cobra.Command{
		Use:     "test [OPTIONS]",
		Short:   "  Short description  ",
		Long:    "Description with\r\n  Windows line endings\r  and weird spacing",
		Example: "",
	}

	buf := new(bytes.Buffer)
	templates.PrintCommandHelpToWriter(buf, cmd)

	output := buf.String()

	// Should not contain Windows line endings after sanitization
	if strings.Contains(output, "\r") {
		t.Error("output should not contain carriage returns")
	}

	// Short description should be trimmed (leading/trailing whitespace removed)
	if strings.Contains(output, "ballerina-test -   Short") {
		t.Error("output should have trimmed short description")
	}
}
