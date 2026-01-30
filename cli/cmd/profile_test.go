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

func TestProfileCommandReturnsNotImplemented(t *testing.T) {
	cmd := NewProfileCommand()
	err := cmd.RunE(cmd, []string{})
	if err == nil {
		t.Fatal("RunE() should return an error")
	}
	if !strings.Contains(err.Error(), "not yet implemented") {
		t.Errorf("expected 'not yet implemented' error, got: %v", err)
	}
}

func TestProfileCommandHelpText(t *testing.T) {
	cmd := NewProfileCommand()

	buf := new(bytes.Buffer)
	templates.PrintCommandHelpToWriter(buf, cmd)

	output := buf.String()
	expectedSubstrings := []string{
		"NAME",
		"ballerina-profile",
		"Run Ballerina Profiler on the source and generate flame graph",
		"SYNOPSIS",
		"DESCRIPTION",
		"flame graph",
		"EXAMPLES",
		"$ bal profile",
	}
	for _, sub := range expectedSubstrings {
		if !strings.Contains(output, sub) {
			t.Errorf("profile help output missing %q", sub)
		}
	}
}

func TestProfileCommandHelpLineWidth(t *testing.T) {
	cmd := NewProfileCommand()

	buf := new(bytes.Buffer)
	templates.PrintCommandHelpToWriter(buf, cmd)

	output := buf.String()
	lines := strings.Split(output, "\n")

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip command examples and NAME line (matches original Java help)
		if strings.HasPrefix(trimmed, "$") || strings.HasPrefix(trimmed, "ballerina-") {
			continue
		}
		if len(line) > 80 {
			t.Errorf("line %d exceeds 80 characters (len=%d): %q", i+1, len(line), line)
		}
	}
}

func TestProfileCommandHasExamples(t *testing.T) {
	cmd := NewProfileCommand()
	if cmd.Example == "" {
		t.Error("profile command should have examples")
	}
}
