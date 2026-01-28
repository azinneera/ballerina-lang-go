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

var newLongDesc = `Creates the given path if it does not exist and initializes a Ballerina package in it. It
generates the Ballerina.toml, main.bal, .gitignore and .devcontainer.json files inside the package
directory.

This command also provides an option to generate files based on predefined templates.`

var newExamples = []templates.Example{
	{
		Description: "Create a new Ballerina package.",
		Commands:    []string{"bal new mypackage"},
	},
	{
		Description: "Create a new package using a template.",
		Commands:    []string{"bal new mypackage -t lib"},
	},
}

func NewNewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "new <package-path> [--workspace] [-t|--template <template-name>]",
		Short:         "Create a new Ballerina package",
		Long:          newLongDesc,
		Example:       templates.FormatExamples(newExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runNew,
	}

	return cmd
}

func runNew(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'new' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
