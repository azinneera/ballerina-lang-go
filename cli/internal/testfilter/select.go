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

// Package testfilter computes, in Go, which of a testplan.TestPlan's tests
// should actually be registered for a given `bal test` invocation — the
// --tests/--groups/--disable-groups/--rerun-failed inclusion decision that
// lib/stdlibs/ballerina/test/0.0.1/go1.27/filter.bal has historically made
// entirely in Ballerina, at registration time.
//
// This is a step toward upstream's compiler-plugin-based test framework
// design (github.com/ballerina-nutcracker/ballerina PRs #921/#922), whose
// Go-native Schedule makes this exact decision with no Ballerina involved at
// all. This package doesn't go that far yet: cli/cmd/test.go still passes
// the original, unmodified CLI flags to ballerina/test's own setTestOptions
// unchanged, so filter.bal's matching logic keeps running as a correctness
// backstop over the (now smaller) set this package decides to register.
// filter.bal itself is intentionally untouched — see TODO.md for why, and
// for the cleanup this leaves for later once this path has proven itself.
package testfilter

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/ballerina-nutcracker/ballerina/cli/internal/testplan"
	"github.com/ballerina-nutcracker/ballerina/lib/stdlibs/ballerina/test/0.0.1/go1.27/native"
)

// SelectOptions mirrors the raw CLI flags/context filter.bal's setTestOptions
// already receives — see cli/cmd/test.go's existing call.
type SelectOptions struct {
	Tests, Groups, DisableGroups string // raw, comma-separated CLI flag values
	RerunFailed                  bool
	RerunJSONPath                string // <target>/rerun_test.json; only read if RerunFailed
	ModuleName, PackageName      string // ModuleName is already fully-qualified (e.g. "pkg" or "pkg.submod")
}

// Select returns the set of plan.Tests' names that should be registered for
// execution. A name present and true in the result should be registered; any
// other discovered test name should not. Suite/each/group hooks are never
// filtered by this decision — see cli/internal/testglue.
//
// The returned map is never nil, distinguishing "restrict to this set"
// (including an empty set, meaning "register nothing") from the nil sentinel
// testglue.GenerateSource treats as "no filtering at all."
func Select(plan *testplan.TestPlan, opts SelectOptions) (map[string]bool, error) {
	var included map[string]bool
	if opts.RerunFailed {
		included = selectRerunFailed(opts)
	} else {
		included = selectByFilters(plan, opts)
	}
	closeDependsOn(plan, included)
	return included, nil
}

// selectByFilters ports filter.bal#registerTestConfig's three independent,
// ANDed checks (groups / disable-groups / hasTest), applied here before
// registration instead of inside it. If Tests/Groups/DisableGroups are all
// empty, every test matches — the same as filter.bal's own "no filter set"
// default.
func selectByFilters(plan *testplan.TestPlan, opts SelectOptions) map[string]bool {
	groups := splitList(opts.Groups)
	disableGroups := splitList(opts.DisableGroups)
	filters := parseTestFilters(opts.Tests)

	included := make(map[string]bool, len(plan.Tests))
	for _, t := range plan.Tests {
		if len(groups) > 0 && !hasGroup(t.Groups, groups) {
			continue
		}
		if len(disableGroups) > 0 && hasGroup(t.Groups, disableGroups) {
			continue
		}
		if len(filters) > 0 && !hasTest(t.Name, opts.ModuleName, opts.PackageName, filters) {
			continue
		}
		included[t.Name] = true
	}
	return included
}

// hasGroup reports whether any of testGroups appears in filterGroups.
func hasGroup(testGroups, filterGroups []string) bool {
	for _, g := range testGroups {
		if slicesContain(filterGroups, g) {
			return true
		}
	}
	return false
}

// testFilter is one parsed --tests entry: a base name or wildcard pattern,
// with any module qualifier and data-provider sub-key already stripped. The
// sub-key is deliberately discarded — it names a data-driven sub-case that
// only exists once the provider function runs, a runtime concern this
// package doesn't resolve (see the package doc comment). The original,
// unmodified --tests string still reaches ballerina/test's own filter.bal,
// which continues to narrow data-driven sub-cases by key exactly as today.
type testFilter struct {
	Name   string // base name or wildcard pattern
	Module string // "" if this entry wasn't module-qualified
}

func parseTestFilters(raw string) []testFilter {
	var filters []testFilter
	for _, entry := range splitList(raw) {
		module := ""
		name := entry
		if colon := strings.Index(name, ":"); colon >= 0 {
			// A '#' (data-key separator) before the colon means the colon
			// belongs to something after the data key, not a module
			// qualifier — matches filter.bal#containsAPrefix.
			if hash := strings.Index(name, "#"); hash < 0 || hash > colon {
				module, name = name[:colon], name[colon+1:]
			}
		}
		if hash := strings.Index(name, "#"); hash >= 0 {
			name = name[:hash]
		}
		filters = append(filters, testFilter{Name: name, Module: module})
	}
	return filters
}

// hasTest ports filter.bal#hasTest/matchModuleName. Module qualification
// matches either the test's exact fully-qualified module name or the bare
// package name (meaning "any module in this package") — a simpler, exact
// comparison than filter.bal's own substring-containment check, safe to
// differ from since filter.bal remains the authoritative backstop: any test
// this function over-includes is still subject to filter.bal's own check
// once registered.
func hasTest(name, moduleName, packageName string, filters []testFilter) bool {
	moduleOK := func(f testFilter) bool {
		return f.Module == "" || f.Module == moduleName || f.Module == packageName
	}
	for _, f := range filters {
		if f.Name == name && !strings.Contains(f.Name, "*") && moduleOK(f) {
			return true
		}
	}
	for _, f := range filters {
		if !strings.Contains(f.Name, "*") {
			continue
		}
		if matched, err := native.MatchWildcard(name, f.Name); err == nil && matched && moduleOK(f) {
			return true
		}
	}
	return false
}

// closeDependsOn adds every function transitively reachable via an included
// test's DependsOn, regardless of whether that dependency itself matches any
// filter, belongs to an included group, or is enabled — reproducing both
// "--tests dependentFn auto-includes its dependency" and "a disabled
// dependency pulled in this way still gets registered, so the
// depends-on-a-disabled-function error still fires" (see the package doc
// comment and TODO.md).
func closeDependsOn(plan *testplan.TestPlan, included map[string]bool) {
	byName := make(map[string]testplan.TestFunction, len(plan.Tests))
	for _, t := range plan.Tests {
		byName[t.Name] = t
	}
	queue := make([]string, 0, len(included))
	for name := range included {
		queue = append(queue, name)
	}
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		t, ok := byName[name]
		if !ok {
			continue
		}
		for _, dep := range t.DependsOn {
			if !included[dep] {
				included[dep] = true
				queue = append(queue, dep)
			}
		}
	}
}

// rerunEntryJSON mirrors rerun_test.json's per-module entry shape, matching
// native/test_io.go's rerunEntryJSON and cli/cmd/testreport.go's reader.
type rerunEntryJSON struct {
	TestNames []string `json:"testNames"`
}

// selectRerunFailed reads the module's rerun_test.json entry and uses its
// test names directly, with no group/wildcard/module matching (matching
// filter.bal#parseRerunJson's own behavior of replacing --tests-based
// selection entirely). Any read/parse problem (missing file, malformed
// JSON, no entry for this module) yields an empty set rather than a Go
// error: ballerina/test's own setTestOptions call — unchanged, still
// receiving the raw --rerun-failed flag — already detects and reports these
// exact conditions with its own tested "Invalid failed test data" message,
// before registration ever runs. An empty inclusion set here doesn't change
// that outcome; it just avoids this package needing to duplicate that
// reporting.
func selectRerunFailed(opts SelectOptions) map[string]bool {
	included := map[string]bool{}
	data, err := os.ReadFile(opts.RerunJSONPath)
	if err != nil {
		return included
	}
	var entries map[string]rerunEntryJSON
	if err := json.Unmarshal(data, &entries); err != nil {
		return included
	}
	entry, ok := entries[opts.ModuleName]
	if !ok {
		return included
	}
	for _, name := range entry.TestNames {
		included[name] = true
	}
	return included
}

func splitList(raw string) []string {
	if raw == "" {
		return nil
	}
	return strings.Split(raw, ",")
}

func slicesContain(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
