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

import (
	"ballerina-lang-go/context"
)

// PackageCompilation represents the result of compiling a package.
// Matches Java's PackageCompilation API.
type PackageCompilation struct {
	pkg              *Package
	diagnosticResult *DiagnosticResult
	compiled         bool
	compilerCtx      *context.CompilerContext
}

// newPackageCompilation creates and runs compilation for a package.
// This is an internal factory function matching Java's PackageCompilation.from().
func newPackageCompilation(pkg *Package) (*PackageCompilation, error) {
	c := &PackageCompilation{
		pkg:         pkg,
		compilerCtx: context.NewCompilerContext(),
	}
	if err := c.compile(); err != nil {
		return nil, err
	}
	return c, nil
}

// DiagnosticResult returns the compilation diagnostics.
func (c *PackageCompilation) DiagnosticResult() *DiagnosticResult {
	return c.diagnosticResult
}

// compile runs the compilation process for all modules in the package.
func (c *PackageCompilation) compile() error {
	var allDiagnostics []Diagnostic

	// Compile each module
	for _, mod := range c.pkg.Modules() {
		moduleDiags, err := mod.compile(c.compilerCtx)
		if err != nil {
			return err
		}
		allDiagnostics = append(allDiagnostics, moduleDiags...)
	}

	c.diagnosticResult = NewDiagnosticResult(allDiagnostics)
	c.compiled = true
	return nil
}
