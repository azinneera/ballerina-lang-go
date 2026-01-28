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

var packLongDesc = `Compiles and packages the current package into a '.bala' file and writes it to the 'target/bala' directory.`

var packExamples = []templates.Example{
	{
		Description: "Create a .bala archive for the current package.",
		Commands:    []string{"bal pack"},
	},
}

func NewPackCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "pack [OPTIONS] [<package>]",
		Short:         "Create distribution format of the current package",
		Long:          packLongDesc,
		Example:       templates.FormatExamples(packExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runPack,
	}

	return cmd
}

func runPack(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'pack' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
