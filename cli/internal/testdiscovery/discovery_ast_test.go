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

// This file exercises Discover() against hand-built AST fixtures rather than
// a fully compiled project — the compiled-project integration test
// (discovery_test.go) is currently blocked on ballerina/test's pending
// unused-symbol state (see TODO.md), so this is the only currently-passing
// verification of Discover()'s actual field-extraction logic.
package testdiscovery_test

import (
	"slices"
	"testing"

	"github.com/ballerina-nutcracker/ballerina/ast"
	"github.com/ballerina-nutcracker/ballerina/cli/internal/testdiscovery"
	"github.com/ballerina-nutcracker/ballerina/cli/internal/testplan"
	"github.com/ballerina-nutcracker/ballerina/tools/diagnostics"
)

func ident(value string) *ast.BLangIdentifier {
	id := ast.NewBLangIdentifier(diagnostics.NewBuiltinLocation(), value, value)
	return &id
}

func fn(name string) *ast.BLangFunction {
	return ast.NewBLangFunction(ast.InvokableData{
		Position: diagnostics.NewBuiltinLocation(),
		Name:     ident(name),
	})
}

func testImport(alias string) *ast.BLangImportPackage {
	imp := &ast.BLangImportPackage{
		OrgName:      ident("ballerina"),
		PkgNameComps: []ast.BLangIdentifier{*ident("test")},
	}
	if alias != "" && alias != "test" {
		imp.Alias = ident(alias)
	}
	return imp
}

func attach(alias, annotationName string, hasValue bool, expr ast.BLangExpression) ast.BLangAnnotationAttachment {
	return ast.BLangAnnotationAttachment{
		PkgAlias:       ident(alias),
		AnnotationName: ident(annotationName),
		HasValue:       hasValue,
		Expr:           expr,
	}
}

func boolLit(v bool) *ast.BLangLiteral {
	return ast.NewBLangLiteral(diagnostics.NewBuiltinLocation(), ast.LiteralKindBoolean, v, "", false)
}

func strLit(v string) *ast.BLangLiteral {
	return ast.NewBLangLiteral(diagnostics.NewBuiltinLocation(), ast.LiteralKindString, v, v, false)
}

func stringList(values ...string) *ast.BLangListConstructorExpr {
	exprs := make([]ast.BLangExpression, len(values))
	for i, v := range values {
		exprs[i] = strLit(v)
	}
	return ast.NewBLangListConstructorExpr(diagnostics.NewBuiltinLocation(), exprs, nil)
}

func funcRefList(fns ...*ast.BLangFunction) *ast.BLangListConstructorExpr {
	exprs := make([]ast.BLangExpression, len(fns))
	for i, f := range fns {
		exprs[i] = &ast.BLangVarRef{VariableName: f.Name}
	}
	return ast.NewBLangListConstructorExpr(diagnostics.NewBuiltinLocation(), exprs, nil)
}

func mapping(fields map[string]ast.BLangExpression) *ast.BLangMappingConstructorExpr {
	m := &ast.BLangMappingConstructorExpr{}
	for key, valueExpr := range fields {
		m.Fields = append(m.Fields, &ast.BLangMappingKeyValueField{
			Key:       &ast.BLangMappingKey{Expr: &ast.BLangVarRef{VariableName: ident(key)}},
			ValueExpr: valueExpr,
		})
	}
	return m
}

func TestDiscover_NoTestImport(t *testing.T) {
	pkg := &ast.BLangPackage{Functions: []*ast.BLangFunction{fn("notATest")}}
	d := testdiscovery.Discover(pkg)
	if len(d.Tests) != 0 || len(d.BeforeSuite) != 0 {
		t.Fatalf("expected nothing discovered without a ballerina/test import, got %+v", d)
	}
}

func TestDiscover_ConfigWithAllFields(t *testing.T) {
	setup := fn("setupFn")
	before := fn("beforeFn")
	after := fn("afterFn")
	main := fn("mainTest")
	main.AddAnnotationAttachment(attach("test", "Config", true, mapping(map[string]ast.BLangExpression{
		"groups":          stringList("unit", "fast"),
		"enable":          boolLit(false),
		"serialExecution": boolLit(true),
		"dependsOn":       funcRefList(setup),
		"before":          &ast.BLangVarRef{VariableName: before.Name},
		"after":           &ast.BLangVarRef{VariableName: after.Name},
	})))

	pkg := &ast.BLangPackage{
		Imports:   []*ast.BLangImportPackage{testImport("test")},
		Functions: []*ast.BLangFunction{setup, before, after, main},
	}

	d := testdiscovery.Discover(pkg)
	if len(d.Tests) != 1 {
		t.Fatalf("expected 1 discovered test, got %d: %+v", len(d.Tests), d.Tests)
	}
	got := d.Tests[0]
	want := testplan.TestFunction{
		Name:            "mainTest",
		Enable:          false,
		Groups:          []string{"unit", "fast"},
		DependsOn:       []string{"setupFn"},
		Before:          "beforeFn",
		After:           "afterFn",
		SerialExecution: true,
		DataProvider:    "",
	}
	if got.Name != want.Name || got.Enable != want.Enable || !slices.Equal(got.Groups, want.Groups) ||
		!slices.Equal(got.DependsOn, want.DependsOn) || got.Before != want.Before || got.After != want.After ||
		got.SerialExecution != want.SerialExecution || got.DataProvider != want.DataProvider {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestDiscover_ConfigDefaults(t *testing.T) {
	// A bare @test:Config with no `{...}` body at all must fall back to the
	// annotation record's declared defaults (enable=true, groups=[], etc.).
	main := fn("bareTest")
	main.AddAnnotationAttachment(attach("test", "Config", false, nil))

	pkg := &ast.BLangPackage{
		Imports:   []*ast.BLangImportPackage{testImport("test")},
		Functions: []*ast.BLangFunction{main},
	}

	d := testdiscovery.Discover(pkg)
	if len(d.Tests) != 1 {
		t.Fatalf("expected 1 discovered test, got %d", len(d.Tests))
	}
	got := d.Tests[0]
	if !got.Enable {
		t.Error("expected Enable=true (default)")
	}
	if len(got.Groups) != 0 || len(got.DependsOn) != 0 {
		t.Errorf("expected empty Groups/DependsOn by default, got %+v", got)
	}
	if got.SerialExecution {
		t.Error("expected SerialExecution=false (default)")
	}
}

func TestDiscover_CustomImportAlias(t *testing.T) {
	main := fn("aliasedTest")
	main.AddAnnotationAttachment(attach("t", "Config", true, mapping(map[string]ast.BLangExpression{
		"groups": stringList("g1"),
	})))

	pkg := &ast.BLangPackage{
		Imports:   []*ast.BLangImportPackage{testImport("t")},
		Functions: []*ast.BLangFunction{main},
	}

	d := testdiscovery.Discover(pkg)
	if len(d.Tests) != 1 {
		t.Fatalf("expected discovery to work with a custom import alias, got %+v", d)
	}
	if !slices.Equal(d.Tests[0].Groups, []string{"g1"}) {
		t.Errorf("expected Groups=[g1], got %v", d.Tests[0].Groups)
	}

	// An annotation attached under the WRONG alias must not match.
	wrongAliasFn := fn("wrongAlias")
	wrongAliasFn.AddAnnotationAttachment(attach("test", "Config", true, mapping(nil)))
	pkg2 := &ast.BLangPackage{
		Imports:   []*ast.BLangImportPackage{testImport("t")},
		Functions: []*ast.BLangFunction{wrongAliasFn},
	}
	d2 := testdiscovery.Discover(pkg2)
	if len(d2.Tests) != 0 {
		t.Errorf("expected no match for an annotation under the wrong alias, got %+v", d2.Tests)
	}
}

func TestDiscover_AllHookKinds(t *testing.T) {
	beforeSuite := fn("bs")
	beforeSuite.AddAnnotationAttachment(attach("test", "BeforeSuite", false, nil))

	afterSuite := fn("as")
	afterSuite.AddAnnotationAttachment(attach("test", "AfterSuite", true, mapping(map[string]ast.BLangExpression{
		"alwaysRun": boolLit(true),
	})))

	beforeEach := fn("be")
	beforeEach.AddAnnotationAttachment(attach("test", "BeforeEach", false, nil))

	afterEach := fn("ae")
	afterEach.AddAnnotationAttachment(attach("test", "AfterEach", false, nil))

	beforeGroups := fn("bg")
	beforeGroups.AddAnnotationAttachment(attach("test", "BeforeGroups", true, mapping(map[string]ast.BLangExpression{
		"value": stringList("g1", "g2"),
	})))

	afterGroups := fn("ag")
	afterGroups.AddAnnotationAttachment(attach("test", "AfterGroups", true, mapping(map[string]ast.BLangExpression{
		"value":     stringList("g1"),
		"alwaysRun": boolLit(true),
	})))

	pkg := &ast.BLangPackage{
		Imports:   []*ast.BLangImportPackage{testImport("test")},
		Functions: []*ast.BLangFunction{beforeSuite, afterSuite, beforeEach, afterEach, beforeGroups, afterGroups},
	}

	d := testdiscovery.Discover(pkg)

	if len(d.BeforeSuite) != 1 || d.BeforeSuite[0].Name != "bs" {
		t.Errorf("BeforeSuite: got %+v", d.BeforeSuite)
	}
	if len(d.AfterSuite) != 1 || d.AfterSuite[0].Name != "as" || !d.AfterSuite[0].AlwaysRun {
		t.Errorf("AfterSuite: got %+v", d.AfterSuite)
	}
	if len(d.BeforeEach) != 1 || d.BeforeEach[0].Name != "be" {
		t.Errorf("BeforeEach: got %+v", d.BeforeEach)
	}
	if len(d.AfterEach) != 1 || d.AfterEach[0].Name != "ae" {
		t.Errorf("AfterEach: got %+v", d.AfterEach)
	}
	if len(d.BeforeGroups) != 1 || d.BeforeGroups[0].Name != "bg" || !slices.Equal(d.BeforeGroups[0].Groups, []string{"g1", "g2"}) {
		t.Errorf("BeforeGroups: got %+v", d.BeforeGroups)
	}
	if len(d.AfterGroups) != 1 || d.AfterGroups[0].Name != "ag" || !slices.Equal(d.AfterGroups[0].Groups, []string{"g1"}) || !d.AfterGroups[0].AlwaysRun {
		t.Errorf("AfterGroups: got %+v", d.AfterGroups)
	}
}
