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

var shellLongDesc = `Run a REPL instance of Ballerina to enable users to execute small snippets of code.

Note: This is an experimental feature, which supports only a limited set of functionality.

Debug messages can be enabled using the '-d' option.`

var shellExamples = []templates.Example{
	{
		Description: "Start the Ballerina REPL.",
		Commands:    []string{"bal shell"},
	},
	{
		Description: "Start the REPL with debug messages.",
		Commands:    []string{"bal shell -d"},
	},
}

func NewShellCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "shell [OPTIONS]",
		Short:         "Run Ballerina interactive REPL",
		Long:          shellLongDesc,
		Example:       templates.FormatExamples(shellExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runShell,
	}

	return cmd
}

func runShell(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'shell' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
