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

var graphLongDesc = `Resolve the dependencies of the current package and print the dependency graph in the console.

This produces the textual representation of the dependency graph using the DOT graph language.`

var graphExamples = []templates.Example{
	{
		Description: "Print the dependency graph for the current package.",
		Commands:    []string{"bal graph"},
	},
}

func NewGraphCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "graph [OPTIONS] [<package>|<source-file>]",
		Short:         "Print the dependency graph",
		Long:          graphLongDesc,
		Example:       templates.FormatExamples(graphExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runGraph,
	}

	return cmd
}

func runGraph(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'graph' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
