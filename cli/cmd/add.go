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

var addLongDesc = `Creates a new directory called '<directory_name>' inside the 'modules' directory of the package.

Any top-level directory inside the 'modules' directory becomes a Ballerina module and its name can be derived as:

    module-name := <package-name>.<directory_name>

Use the 'module-name' when importing the module.`

var addExamples = []templates.Example{
	{
		Description: "Add a new module to the current package.",
		Commands:    []string{"bal add mymodule"},
	},
}

func NewAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "add <directory_name>",
		Short:         "Add a new module to the current package",
		Long:          addLongDesc,
		Example:       templates.FormatExamples(addExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runAdd,
	}

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'add' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
