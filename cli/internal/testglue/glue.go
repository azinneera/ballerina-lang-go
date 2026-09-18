// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
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

// Package testglue generates Ballerina source that registers a testplan.TestPlan's
// functions with ballerina/test.
//
// TEMPORARY, see TODO.md's `typeof` entry: this package exists only because
// `typeof` isn't implemented, so annotation field values can't be read at
// Ballerina runtime the way jballerina's annotation_processor.bal does.
// Once typeof lands, delete this package along with cli/internal/testdiscovery
// and ballerina/test's registration_api.bal — the real compiler-generated glue
// only needs `test:registerTest(name, fn)` + `test:setTestOptions(...)` +
// `test:startSuite()`, with annotation_processor.bal doing the field-value
// reading at runtime instead of this package doing it ahead of time.
package testglue

import (
	"fmt"
	"strings"

	"github.com/ballerina-nutcracker/ballerina/cli/internal/testplan"
)

// EntryFunctionName is the generated function that registers every included
// @test:* function. cli/cmd/test.go invokes this by name in-process after
// compiling the module with test sources included.
const EntryFunctionName = "__test_register__"

// GenerateSource produces Ballerina source text that, compiled alongside the
// rest of a module (as a new test-only document for packages/workspaces, or
// appended in-memory for standalone files), registers plan's functions with
// ballerina/test.
//
// included restricts which of plan.Tests actually get a registerTestConfig
// call emitted — see cli/internal/testfilter, which computes this set from
// --tests/--groups/--disable-groups/--rerun-failed before this function ever
// runs. A test absent from included (or a nil included, matching "select
// everything") is skipped entirely; it never reaches ballerina/test's own
// registry at all. Suite/each/group hooks are never filtered by included —
// they're global infrastructure that's harmless to register even when no
// test in a given group actually runs.
//
// `test:setTestOptions(...)` and `test:startSuite()` are NOT emitted here —
// cli/cmd/test.go invokes those directly via the runtime after calling this
// generated entry point, since only the CLI has the resolved flag values
// setTestOptions needs.
func GenerateSource(plan *testplan.TestPlan, included map[string]bool) string {
	var b strings.Builder
	b.WriteString("import ballerina/test;\n\n")
	fmt.Fprintf(&b, "public function %s() {\n", EntryFunctionName)

	anyEmitted := false
	for _, t := range plan.Tests {
		if included != nil && !included[t.Name] {
			continue
		}
		writeTestConfigCall(&b, t)
		anyEmitted = true
	}
	for _, h := range plan.BeforeSuite {
		fmt.Fprintf(&b, "    test:registerBeforeSuite(%s, %s);\n", balString(h.Name), h.Name)
		anyEmitted = true
	}
	for _, h := range plan.AfterSuite {
		fmt.Fprintf(&b, "    test:registerAfterSuite(%s, %s, %s);\n", balString(h.Name), h.Name, balBool(h.AlwaysRun))
		anyEmitted = true
	}
	for _, h := range plan.BeforeEach {
		fmt.Fprintf(&b, "    test:registerBeforeEach(%s, %s);\n", balString(h.Name), h.Name)
		anyEmitted = true
	}
	for _, h := range plan.AfterEach {
		fmt.Fprintf(&b, "    test:registerAfterEach(%s, %s);\n", balString(h.Name), h.Name)
		anyEmitted = true
	}
	for _, h := range plan.BeforeGroups {
		fmt.Fprintf(&b, "    test:registerBeforeGroups(%s, %s, %s);\n", balString(h.Name), h.Name, balStringArray(h.Groups))
		anyEmitted = true
	}
	for _, h := range plan.AfterGroups {
		fmt.Fprintf(&b, "    test:registerAfterGroups(%s, %s, %s, %s);\n",
			balString(h.Name), h.Name, balStringArray(h.Groups), balBool(h.AlwaysRun))
		anyEmitted = true
	}
	if !anyEmitted {
		// `import ballerina/test;` above would otherwise be a compile error
		// (unused import) on a module with no @test:* functions at all, or
		// none surviving `included` — confirmed via a real "no tests in this
		// package" run of `bal test`. ballerina/test must still end up in the
		// compiled dependency graph regardless, since cli/cmd/test.go always
		// looks up and calls test:setTestOptions/test:startSuite directly via
		// the runtime, so a harmless reference to a real public function
		// value anchors the import without registering or invoking anything.
		b.WriteString("    var _anchor = test:startSuite;\n    _ = _anchor;\n")
	}

	b.WriteString("}\n")
	return b.String()
}

func writeTestConfigCall(b *strings.Builder, t testplan.TestFunction) {
	before := "()"
	if t.Before != "" {
		before = t.Before
	}
	after := "()"
	if t.After != "" {
		after = t.After
	}
	dataProvider := "()"
	if t.DataProvider != "" {
		dataProvider = t.DataProvider
	}
	fmt.Fprintf(b, "    test:registerTestConfig(%s, %s, %s, %s, %s, %s, %s, %s, %s);\n",
		balString(t.Name), t.Name, balBool(t.Enable), balStringArray(t.Groups),
		balFunctionArray(t.DependsOn), before, after, balBool(t.SerialExecution), dataProvider)
}

func balBool(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

// balString renders a Go string as a Ballerina string literal.
func balString(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\t':
			b.WriteString(`\t`)
		case '\r':
			b.WriteString(`\r`)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}

func balStringArray(values []string) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = balString(v)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// balFunctionArray renders function names as a Ballerina array of bare
// identifier references (e.g. "[testSetup, testTeardown]") — function names
// are always valid Ballerina identifiers by construction, so no escaping
// is needed the way balString needs it for arbitrary group/test names.
func balFunctionArray(names []string) string {
	return "[" + strings.Join(names, ", ") + "]"
}
