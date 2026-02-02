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
	"fmt"
	"strconv"
	"strings"
)

// PackageName represents a package name.
type PackageName string

// String returns the package name as a string.
func (n PackageName) String() string {
	return string(n)
}

// PackageOrg represents an organization.
type PackageOrg string

// Well-known organization constants.
const (
	OrgBallerina  PackageOrg = "ballerina"
	OrgBallerinai PackageOrg = "ballerinai"
	OrgBallerinax PackageOrg = "ballerinax"
)

// String returns the organization as a string.
func (o PackageOrg) String() string {
	return string(o)
}

// IsBallerinaOrg returns true if this is a Ballerina built-in organization.
func (o PackageOrg) IsBallerinaOrg() bool {
	return o == OrgBallerina || o == OrgBallerinai
}

// IsBallerinaxOrg returns true if this is the ballerinax organization.
func (o PackageOrg) IsBallerinaxOrg() bool {
	return o == OrgBallerinax
}

// IsEmpty returns true if the organization is empty (anonymous).
func (o PackageOrg) IsEmpty() bool {
	return o == ""
}

// PackageVersion represents a semantic version.
type PackageVersion struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
}

// NewPackageVersion creates a new PackageVersion with the given components.
func NewPackageVersion(major, minor, patch int) PackageVersion {
	return PackageVersion{Major: major, Minor: minor, Patch: patch}
}

// ParsePackageVersion parses a version string in SemVer format (e.g., "1.2.3" or "1.2.3-beta").
func ParsePackageVersion(version string) (PackageVersion, error) {
	if version == "" {
		return PackageVersion{}, fmt.Errorf("empty version string")
	}

	parts := strings.SplitN(version, "-", 2)
	versionParts := strings.Split(parts[0], ".")
	if len(versionParts) != 3 {
		return PackageVersion{}, fmt.Errorf("invalid version format: %s (expected major.minor.patch)", version)
	}

	major, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return PackageVersion{}, fmt.Errorf("invalid major version: %s", versionParts[0])
	}

	minor, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return PackageVersion{}, fmt.Errorf("invalid minor version: %s", versionParts[1])
	}

	patch, err := strconv.Atoi(versionParts[2])
	if err != nil {
		return PackageVersion{}, fmt.Errorf("invalid patch version: %s", versionParts[2])
	}

	v := PackageVersion{Major: major, Minor: minor, Patch: patch}
	if len(parts) > 1 {
		v.PreRelease = parts[1]
	}
	return v, nil
}

// String returns the version as a string.
func (v PackageVersion) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.PreRelease != "" {
		s += "-" + v.PreRelease
	}
	return s
}

// IsZero returns true if this is a zero-value version.
func (v PackageVersion) IsZero() bool {
	return v.Major == 0 && v.Minor == 0 && v.Patch == 0 && v.PreRelease == ""
}

// IsStable returns true if this is a stable version (major > 0, no pre-release).
func (v PackageVersion) IsStable() bool {
	return v.Major > 0 && v.PreRelease == ""
}

// Compare compares two versions. Returns -1 if v < other, 0 if v == other, 1 if v > other.
func (v PackageVersion) Compare(other PackageVersion) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}
	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}
	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}
	// Pre-release versions are considered less than release versions
	if v.PreRelease == "" && other.PreRelease != "" {
		return 1
	}
	if v.PreRelease != "" && other.PreRelease == "" {
		return -1
	}
	// Both have pre-release, compare lexicographically
	if v.PreRelease < other.PreRelease {
		return -1
	}
	if v.PreRelease > other.PreRelease {
		return 1
	}
	return 0
}

// ModuleName represents a module name (package.module or just package for default).
type ModuleName struct {
	PackageName PackageName
	ModulePart  string // empty for default module
}

// NewModuleName creates a new ModuleName.
func NewModuleName(pkgName PackageName, modulePart string) ModuleName {
	return ModuleName{PackageName: pkgName, ModulePart: modulePart}
}

// IsDefault returns true if this is the default module (no module part).
func (m ModuleName) IsDefault() bool {
	return m.ModulePart == ""
}

// String returns the full module name.
func (m ModuleName) String() string {
	if m.ModulePart == "" {
		return string(m.PackageName)
	}
	return string(m.PackageName) + "." + m.ModulePart
}
