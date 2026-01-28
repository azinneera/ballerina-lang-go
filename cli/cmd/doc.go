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

var docLongDesc = `Generates the documentation of the current package and writes it to the 'target/apidocs' directory.`

var docExamples = []templates.Example{
	{
		Description: "Generate API documentation for the current package.",
		Commands:    []string{"bal doc"},
	},
}

func NewDocCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "doc [OPTIONS] [<package>]",
		Short:         "Build the documentation of a Ballerina package",
		Long:          docLongDesc,
		Example:       templates.FormatExamples(docExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runDoc,
	}

	return cmd
}

func runDoc(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'doc' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
