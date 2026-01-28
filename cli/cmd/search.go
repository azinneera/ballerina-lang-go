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

var searchLongDesc = `Search Ballerina Central for the given keyword appearing in the organization name or package name.

The keyword is treated as case-insensitive when the search is executed.`

var searchExamples = []templates.Example{
	{
		Description: "Search for packages containing 'http'.",
		Commands:    []string{"bal search http"},
	},
}

func NewSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "search <keyword>",
		Short:         "Search Ballerina Central for packages",
		Long:          searchLongDesc,
		Example:       templates.FormatExamples(searchExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runSearch,
	}

	return cmd
}

func runSearch(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'search' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
