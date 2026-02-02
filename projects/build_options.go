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

// LockingMode determines version locking behavior for dependencies.
type LockingMode int

const (
	// LockingModeSoft updates to the latest minor version.
	LockingModeSoft LockingMode = iota
	// LockingModeMedium updates to the latest patch version only.
	LockingModeMedium
	// LockingModeHard uses exact versions as much as possible.
	LockingModeHard
	// LockingModeLocked strictly uses versions from Dependencies.toml only.
	LockingModeLocked
)

func (m LockingMode) String() string {
	switch m {
	case LockingModeSoft:
		return "SOFT"
	case LockingModeMedium:
		return "MEDIUM"
	case LockingModeHard:
		return "HARD"
	case LockingModeLocked:
		return "LOCKED"
	default:
		return "UNKNOWN"
	}
}

// BuildOptions contains build configuration options.
type BuildOptions struct {
	skipTests        bool
	testReport       bool
	codeCoverage     bool
	offline          bool
	experimental     bool
	targetDir        string
	nativeImage      bool
	graalVMBuildOpts string
	lockingMode      LockingMode
	cloud            string
	dumpBuildTime    bool
}

// SkipTests returns whether tests should be skipped.
func (o BuildOptions) SkipTests() bool {
	return o.skipTests
}

// TestReport returns whether test reports should be generated.
func (o BuildOptions) TestReport() bool {
	return o.testReport
}

// CodeCoverage returns whether code coverage should be collected.
func (o BuildOptions) CodeCoverage() bool {
	return o.codeCoverage
}

// Offline returns whether the build should run in offline mode.
func (o BuildOptions) Offline() bool {
	return o.offline
}

// Experimental returns whether experimental features are enabled.
func (o BuildOptions) Experimental() bool {
	return o.experimental
}

// TargetDir returns the custom target directory (empty for default).
func (o BuildOptions) TargetDir() string {
	return o.targetDir
}

// NativeImage returns whether to build a native image.
func (o BuildOptions) NativeImage() bool {
	return o.nativeImage
}

// GraalVMBuildOpts returns GraalVM build options.
func (o BuildOptions) GraalVMBuildOpts() string {
	return o.graalVMBuildOpts
}

// LockingMode returns the dependency locking mode.
func (o BuildOptions) LockingMode() LockingMode {
	return o.lockingMode
}

// Cloud returns the cloud platform target.
func (o BuildOptions) Cloud() string {
	return o.cloud
}

// DumpBuildTime returns whether to dump build time information.
func (o BuildOptions) DumpBuildTime() bool {
	return o.dumpBuildTime
}

// BuildOptionsBuilder builds BuildOptions using the builder pattern.
type BuildOptionsBuilder struct {
	opts BuildOptions
}

// NewBuildOptionsBuilder creates a new builder with default values.
func NewBuildOptionsBuilder() *BuildOptionsBuilder {
	return &BuildOptionsBuilder{
		opts: BuildOptions{
			skipTests:   true,
			lockingMode: LockingModeMedium,
		},
	}
}

// SkipTests sets whether to skip tests.
func (b *BuildOptionsBuilder) SkipTests(v bool) *BuildOptionsBuilder {
	b.opts.skipTests = v
	return b
}

// TestReport sets whether to generate test reports.
func (b *BuildOptionsBuilder) TestReport(v bool) *BuildOptionsBuilder {
	b.opts.testReport = v
	return b
}

// CodeCoverage sets whether to collect code coverage.
func (b *BuildOptionsBuilder) CodeCoverage(v bool) *BuildOptionsBuilder {
	b.opts.codeCoverage = v
	return b
}

// Offline sets whether to run in offline mode.
func (b *BuildOptionsBuilder) Offline(v bool) *BuildOptionsBuilder {
	b.opts.offline = v
	return b
}

// Experimental sets whether experimental features are enabled.
func (b *BuildOptionsBuilder) Experimental(v bool) *BuildOptionsBuilder {
	b.opts.experimental = v
	return b
}

// TargetDir sets a custom target directory.
func (b *BuildOptionsBuilder) TargetDir(v string) *BuildOptionsBuilder {
	b.opts.targetDir = v
	return b
}

// NativeImage sets whether to build a native image.
func (b *BuildOptionsBuilder) NativeImage(v bool) *BuildOptionsBuilder {
	b.opts.nativeImage = v
	return b
}

// GraalVMBuildOpts sets GraalVM build options.
func (b *BuildOptionsBuilder) GraalVMBuildOpts(v string) *BuildOptionsBuilder {
	b.opts.graalVMBuildOpts = v
	return b
}

// LockingMode sets the dependency locking mode.
func (b *BuildOptionsBuilder) LockingMode(v LockingMode) *BuildOptionsBuilder {
	b.opts.lockingMode = v
	return b
}

// Cloud sets the cloud platform target.
func (b *BuildOptionsBuilder) Cloud(v string) *BuildOptionsBuilder {
	b.opts.cloud = v
	return b
}

// DumpBuildTime sets whether to dump build time information.
func (b *BuildOptionsBuilder) DumpBuildTime(v bool) *BuildOptionsBuilder {
	b.opts.dumpBuildTime = v
	return b
}

// Build creates the BuildOptions from the builder.
func (b *BuildOptionsBuilder) Build() BuildOptions {
	return b.opts
}
