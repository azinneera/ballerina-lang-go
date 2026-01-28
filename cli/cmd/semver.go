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
	"fmt"

	"ballerina-lang-go/cli/pkg/templates"

	"github.com/spf13/cobra"
)

var semverLongDesc = `Compare the local changes in the package with any previously released package version
available in Ballerina Central.

Provide suggestions for the next version based on the source code compatibility between the local
changes and any specified previous release version.

Note: This feature is experimental and has limited support for some advanced language constructs.`

var semverExamples = []templates.Example{
	{
		Description: "Validate SemVer compliance for the current package.",
		Commands:    []string{"bal semver"},
	},
}

func NewSemverCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "semver [OPTIONS] [<package-path>]",
		Short:         "Validate SemVer compliance of the package changes",
		Long:          semverLongDesc,
		Example:       templates.FormatExamples(semverExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runSemver,
	}

	return cmd
}

func runSemver(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'semver' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
