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

var testLongDesc = `Compiles and executes test functions and prints a summary of the test results.

Test runs the test functions defined in each module of a package when building the current package.
It runs the test functions defined in the given source file when building a single '.bal' file.

Note: Testing individual '.bal' files of a package is not allowed.`

var testExamples = []templates.Example{
	{
		Description: "Run all tests in the current package.",
		Commands:    []string{"bal test"},
	},
	{
		Description: "Run tests in a single Ballerina file.",
		Commands:    []string{"bal test main_test.bal"},
	},
}

func NewTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "test [OPTIONS] [<package>|<source-file>] [(-Ckey=value)...]",
		Short:         "Run package tests",
		Long:          testLongDesc,
		Example:       templates.FormatExamples(testExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runTest,
	}

	return cmd
}

func runTest(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'test' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
