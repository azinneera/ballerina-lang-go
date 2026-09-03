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

// Package testdiscovery finds @test:* annotated functions in a compiled
// package's AST and extracts their configuration.
//
// jballerina reads annotation field values at runtime via `(typeof f).@Config`
// (annotation_processor.bal). That mechanism depends on `typeof`, which this
// interpreter doesn't implement (see TODO.md). This package substitutes a
// static, compile-time read of the same field values directly from the AST,
// which the compiler already resolves as ordinary annotation attachments — no
// runtime introspection needed. Whatever this package finds is expected to be
// baked into generated glue code as literal arguments (see the P8.5 glue-code
// generation step), rather than re-discovered at runtime.
package testdiscovery

import (
	"github.com/ballerina-nutcracker/ballerina/ast"
)

const (
	orgName    = "ballerina"
	moduleName = "test"
)

// TestFunc is a discovered @test:Config-annotated function.
type TestFunc struct {
	Name            string
	Enable          bool
	Groups          []string
	DependsOn       []string // referenced function names
	Before          string   // referenced function name, "" if none
	After           string   // referenced function name, "" if none
	SerialExecution bool
	DataProvider    string // referenced function name, "" if none
}

// SuiteHook is a discovered @test:BeforeSuite/@test:BeforeEach/@test:AfterEach
// function — these are marker annotations with no configurable fields.
type SuiteHook struct {
	Name string
}

// AfterSuiteHook is a discovered @test:AfterSuite function.
type AfterSuiteHook struct {
	Name      string
	AlwaysRun bool
}

// GroupHook is a discovered @test:BeforeGroups/@test:AfterGroups function.
// AlwaysRun is only meaningful for AfterGroups.
type GroupHook struct {
	Name      string
	Groups    []string
	AlwaysRun bool
}

// Discovered holds every @test:* annotated function found in a package.
type Discovered struct {
	Tests        []TestFunc
	BeforeSuite  []SuiteHook
	AfterSuite   []AfterSuiteHook
	BeforeEach   []SuiteHook
	AfterEach    []SuiteHook
	BeforeGroups []GroupHook
	AfterGroups  []GroupHook
}

// HasAny reports whether anything was discovered at all — a package with no
// @test:* annotations anywhere (e.g. no test files, or test files with none)
// needs no ballerina/test registration glue at all.
func (d *Discovered) HasAny() bool {
	return len(d.Tests) > 0 || len(d.BeforeSuite) > 0 || len(d.AfterSuite) > 0 ||
		len(d.BeforeEach) > 0 || len(d.AfterEach) > 0 || len(d.BeforeGroups) > 0 || len(d.AfterGroups) > 0
}

// Discover walks every top-level function in pkg (pkg.Functions already
// excludes class/service member functions — the same scope jballerina's own
// TestFunctionVisitor restricts itself to) and extracts @test:* annotation
// metadata. Returns an empty Discovered if pkg is nil.
func Discover(pkg *ast.BLangPackage) *Discovered {
	if pkg == nil {
		return &Discovered{}
	}
	alias := resolveTestAlias(pkg)

	d := &Discovered{}
	for _, fn := range pkg.Functions {
		name := identifierValue(fn.Name)
		for _, att := range fn.GetAnnotationAttachments() {
			if !isTestAnnotation(att, alias) {
				continue
			}
			switch identifierValue(att.AnnotationName) {
			case "Config":
				d.Tests = append(d.Tests, parseTestConfig(name, att))
			case "BeforeSuite":
				d.BeforeSuite = append(d.BeforeSuite, SuiteHook{Name: name})
			case "AfterSuite":
				d.AfterSuite = append(d.AfterSuite, parseAfterSuiteConfig(name, att))
			case "BeforeEach":
				d.BeforeEach = append(d.BeforeEach, SuiteHook{Name: name})
			case "AfterEach":
				d.AfterEach = append(d.AfterEach, SuiteHook{Name: name})
			case "BeforeGroups":
				d.BeforeGroups = append(d.BeforeGroups, parseGroupsHook(name, att, false))
			case "AfterGroups":
				d.AfterGroups = append(d.AfterGroups, parseGroupsHook(name, att, true))
			}
		}
	}
	return d
}

// resolveTestAlias returns the alias ballerina/test was imported under in pkg
// (e.g. "test" for a plain `import ballerina/test;`, or a custom alias for
// `import ballerina/test as t;`). Falls back to the default alias "test" when
// pkg.Imports doesn't say — the compiled AST this package actually receives
// (via projects.PackageCompilation.ModuleAST) has Imports cleared after
// symbol/type resolution (see projects/module_context.go's
// resolveTypesAndSymbols), so this fallback is what handles every real
// invocation; testImportAlias's import-based lookup only matters for the
// (still useful, custom-alias-precise) hand-built-AST unit tests.
func resolveTestAlias(pkg *ast.BLangPackage) string {
	if alias := testImportAlias(pkg); alias != "" {
		return alias
	}
	return moduleName
}

// testImportAlias returns the alias ballerina/test was imported under in pkg,
// or "" if pkg.Imports doesn't contain it (either genuinely absent, or simply
// unavailable at this compilation stage — see resolveTestAlias).
func testImportAlias(pkg *ast.BLangPackage) string {
	for _, imp := range pkg.Imports {
		if imp.OrgName == nil || imp.OrgName.GetValue() != orgName {
			continue
		}
		if len(imp.PkgNameComps) != 1 || imp.PkgNameComps[0].GetValue() != moduleName {
			continue
		}
		if imp.Alias != nil {
			return imp.Alias.GetValue()
		}
		return moduleName
	}
	return ""
}

func isTestAnnotation(att ast.BLangAnnotationAttachment, alias string) bool {
	return att.PkgAlias != nil && att.PkgAlias.GetValue() == alias && att.AnnotationName != nil
}

func identifierValue(id ast.IdentifierNode) string {
	if id == nil {
		return ""
	}
	return id.GetValue()
}

func parseTestConfig(fnName string, att ast.BLangAnnotationAttachment) TestFunc {
	fields := mappingFields(att.Expr)
	return TestFunc{
		Name:            fnName,
		Enable:          boolField(fields, "enable", true),
		Groups:          stringArrayField(fields, "groups"),
		DependsOn:       functionRefArrayField(fields, "dependsOn"),
		Before:          functionRefField(fields, "before"),
		After:           functionRefField(fields, "after"),
		SerialExecution: boolField(fields, "serialExecution", false),
		DataProvider:    functionRefField(fields, "dataProvider"),
	}
}

func parseAfterSuiteConfig(fnName string, att ast.BLangAnnotationAttachment) AfterSuiteHook {
	fields := mappingFields(att.Expr)
	return AfterSuiteHook{
		Name:      fnName,
		AlwaysRun: boolField(fields, "alwaysRun", false),
	}
}

func parseGroupsHook(fnName string, att ast.BLangAnnotationAttachment, isAfter bool) GroupHook {
	fields := mappingFields(att.Expr)
	hook := GroupHook{
		Name:   fnName,
		Groups: stringArrayField(fields, "value"),
	}
	if isAfter {
		hook.AlwaysRun = boolField(fields, "alwaysRun", false)
	}
	return hook
}

// mappingFields extracts field-name -> value-expression pairs from a mapping
// constructor (e.g. `{ groups: ["a"], enable: false }`). Returns an empty map
// for a bare annotation attachment (no `{...}` at all, i.e. expr is nil) or
// any field whose key isn't a plain identifier/string literal.
func mappingFields(expr ast.BLangExpression) map[string]ast.BLangExpression {
	fields := make(map[string]ast.BLangExpression)
	mapping, ok := expr.(*ast.BLangMappingConstructorExpr)
	if !ok {
		return fields
	}
	for _, field := range mapping.Fields {
		kv, ok := field.(*ast.BLangMappingKeyValueField)
		if !ok {
			continue
		}
		key, ok := mappingKeyName(kv.Key)
		if !ok {
			continue
		}
		fields[key] = kv.ValueExpr
	}
	return fields
}

func mappingKeyName(key *ast.BLangMappingKey) (string, bool) {
	if key == nil || key.Expr == nil {
		return "", false
	}
	switch expr := key.Expr.(type) {
	case *ast.BLangLiteral:
		value, ok := expr.Value.(string)
		return value, ok
	case *ast.BLangVarRef:
		return identifierValue(expr.VariableName), true
	default:
		return "", false
	}
}

func boolField(fields map[string]ast.BLangExpression, name string, defaultValue bool) bool {
	expr, ok := fields[name]
	if !ok {
		return defaultValue
	}
	literal, ok := expr.(*ast.BLangLiteral)
	if !ok {
		return defaultValue
	}
	value, ok := literal.Value.(bool)
	if !ok {
		return defaultValue
	}
	return value
}

func stringArrayField(fields map[string]ast.BLangExpression, name string) []string {
	expr, ok := fields[name]
	if !ok {
		return nil
	}
	list, ok := expr.(*ast.BLangListConstructorExpr)
	if !ok {
		return nil
	}
	var result []string
	for _, member := range list.Exprs {
		literal, ok := member.(*ast.BLangLiteral)
		if !ok {
			continue
		}
		value, ok := literal.Value.(string)
		if !ok {
			continue
		}
		result = append(result, value)
	}
	return result
}

// functionRefField reads a single function-reference field (e.g. `before:
// myHook`) and returns the referenced function's name, or "" if the field is
// absent or isn't a plain reference.
func functionRefField(fields map[string]ast.BLangExpression, name string) string {
	expr, ok := fields[name]
	if !ok {
		return ""
	}
	return functionRefName(expr)
}

// functionRefArrayField reads a field holding an array of function references
// (e.g. `dependsOn: [fnA, fnB]`).
func functionRefArrayField(fields map[string]ast.BLangExpression, name string) []string {
	expr, ok := fields[name]
	if !ok {
		return nil
	}
	list, ok := expr.(*ast.BLangListConstructorExpr)
	if !ok {
		return nil
	}
	var result []string
	for _, member := range list.Exprs {
		if refName := functionRefName(member); refName != "" {
			result = append(result, refName)
		}
	}
	return result
}

func functionRefName(expr ast.BLangExpression) string {
	varRef, ok := expr.(*ast.BLangVarRef)
	if !ok {
		return ""
	}
	return identifierValue(varRef.VariableName)
}
