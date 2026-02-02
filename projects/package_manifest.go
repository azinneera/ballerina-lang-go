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

	"github.com/BurntSushi/toml"
)

// PackageManifest represents parsed Ballerina.toml (minimal for Phase 1).
type PackageManifest struct {
	Org     string
	Name    string
	Version string
}

// ballerinaToml represents the structure of Ballerina.toml for parsing.
type ballerinaToml struct {
	Package struct {
		Org     string `toml:"org"`
		Name    string `toml:"name"`
		Version string `toml:"version"`
	} `toml:"package"`
}

// ParsePackageManifest parses a Ballerina.toml file and returns the package manifest.
func ParsePackageManifest(path string) (*PackageManifest, error) {
	var config ballerinaToml

	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, fmt.Errorf("cannot parse Ballerina.toml: %w", err)
	}

	manifest := &PackageManifest{
		Org:     config.Package.Org,
		Name:    config.Package.Name,
		Version: config.Package.Version,
	}

	// Validate required fields
	if manifest.Org == "" {
		return nil, fmt.Errorf("missing required field 'org' in [package]")
	}
	if manifest.Name == "" {
		return nil, fmt.Errorf("missing required field 'name' in [package]")
	}
	if manifest.Version == "" {
		return nil, fmt.Errorf("missing required field 'version' in [package]")
	}

	return manifest, nil
}
