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

var toolLongDesc = `Register and manage custom commands for the Ballerina CLI.

This command facilitates searching, pulling, and updating tools from the Ballerina Central, switching
between installed versions, and listing and removing installed tools.`

var toolExamples = []templates.Example{
	{
		Description: "List installed tools.",
		Commands:    []string{"bal tool list"},
	},
	{
		Description: "Pull a tool from Ballerina Central.",
		Commands:    []string{"bal tool pull openapi"},
	},
	{
		Description: "Remove an installed tool.",
		Commands:    []string{"bal tool remove openapi"},
	},
}

func NewToolCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "tool <command> [args]",
		Short:         "Extend the Ballerina CLI with custom commands",
		Long:          toolLongDesc,
		Example:       templates.FormatExamples(toolExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runTool,
	}

	return cmd
}

func runTool(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'tool' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
