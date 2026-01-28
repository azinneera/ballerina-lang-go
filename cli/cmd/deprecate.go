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

var deprecateLongDesc = `Deprecate a package published in Ballerina central.

If a specific version of a package is provided, that version will be deprecated. If no version is
specified, all versions of the package will be deprecated.

This command does not delete the package from Ballerina central, and the package deprecation can be
undone using the '--undo' option.`

var deprecateExamples = []templates.Example{
	{
		Description: "Deprecate all versions of a package.",
		Commands:    []string{"bal deprecate myorg/mypackage"},
	},
	{
		Description: "Deprecate a specific version.",
		Commands:    []string{"bal deprecate myorg/mypackage:1.0.0"},
	},
	{
		Description: "Undo deprecation.",
		Commands:    []string{"bal deprecate myorg/mypackage --undo"},
	},
}

func NewDeprecateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "deprecate [OPTIONS] <org-name>/<package-name>[:<version>]",
		Short:         "Deprecates a published package",
		Long:          deprecateLongDesc,
		Example:       templates.FormatExamples(deprecateExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runDeprecate,
	}

	return cmd
}

func runDeprecate(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'deprecate' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
