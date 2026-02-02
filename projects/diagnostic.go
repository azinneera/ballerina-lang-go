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

package projects

import "fmt"

// DiagnosticSeverity indicates the severity of a diagnostic.
type DiagnosticSeverity int

const (
	SeverityError DiagnosticSeverity = iota
	SeverityWarning
	SeverityHint
)

func (s DiagnosticSeverity) String() string {
	switch s {
	case SeverityError:
		return "ERROR"
	case SeverityWarning:
		return "WARNING"
	case SeverityHint:
		return "HINT"
	default:
		return "UNKNOWN"
	}
}

// DiagnosticCode identifies the diagnostic type.
type DiagnosticCode string

// Location identifies where a diagnostic occurred.
type Location struct {
	FilePath    string
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

// Diagnostic represents a single diagnostic message.
type Diagnostic struct {
	Code     DiagnosticCode
	Message  string
	Severity DiagnosticSeverity
	Location Location
}

// Error implements the error interface for Diagnostic.
func (d Diagnostic) Error() string {
	if d.Location.FilePath != "" {
		return fmt.Sprintf("%s:%d:%d: %s: %s", d.Location.FilePath, d.Location.StartLine, d.Location.StartColumn,
			d.Severity, d.Message)
	}
	return fmt.Sprintf("%s: %s", d.Severity, d.Message)
}

// Diagnostics collects multiple diagnostics.
type Diagnostics struct {
	items []Diagnostic
}

// NewDiagnostics creates a new empty Diagnostics collection.
func NewDiagnostics() *Diagnostics {
	return &Diagnostics{}
}

// Add appends a diagnostic to the collection.
func (d *Diagnostics) Add(diag Diagnostic) {
	d.items = append(d.items, diag)
}

// AddError adds an error diagnostic with the given message.
func (d *Diagnostics) AddError(message string) {
	d.items = append(d.items, Diagnostic{
		Severity: SeverityError,
		Message:  message,
	})
}

// AddErrorAt adds an error diagnostic at the specified location.
func (d *Diagnostics) AddErrorAt(message string, loc Location) {
	d.items = append(d.items, Diagnostic{
		Severity: SeverityError,
		Message:  message,
		Location: loc,
	})
}

// All returns all diagnostics.
func (d *Diagnostics) All() []Diagnostic {
	return d.items
}

// HasErrors returns true if there are any error-level diagnostics.
func (d *Diagnostics) HasErrors() bool {
	for _, item := range d.items {
		if item.Severity == SeverityError {
			return true
		}
	}
	return false
}

// Errors returns only error-level diagnostics.
func (d *Diagnostics) Errors() []Diagnostic {
	var errs []Diagnostic
	for _, item := range d.items {
		if item.Severity == SeverityError {
			errs = append(errs, item)
		}
	}
	return errs
}

// Warnings returns only warning-level diagnostics.
func (d *Diagnostics) Warnings() []Diagnostic {
	var warnings []Diagnostic
	for _, item := range d.items {
		if item.Severity == SeverityWarning {
			warnings = append(warnings, item)
		}
	}
	return warnings
}

// Count returns the total number of diagnostics.
func (d *Diagnostics) Count() int {
	return len(d.items)
}

// ErrorCount returns the number of error diagnostics.
func (d *Diagnostics) ErrorCount() int {
	count := 0
	for _, item := range d.items {
		if item.Severity == SeverityError {
			count++
		}
	}
	return count
}
