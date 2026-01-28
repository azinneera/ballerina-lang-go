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

var formatLongDesc = `Formats the Ballerina source files. Formatting can be performed on a Ballerina package,
module, or source file.

The formatted content will be written to the original files. By using the 'dry run' option, you will
be able to check which files will be formatted after the execution.

If the Ballerina sources contain syntax errors, they will be notified and formatting will not be
proceeded until they are fixed.`

var formatExamples = []templates.Example{
	{
		Description: "Format the current package.",
		Commands:    []string{"bal format"},
	},
	{
		Description: "Format a specific source file.",
		Commands:    []string{"bal format main.bal"},
	},
	{
		Description: "Dry run to check which files would be formatted.",
		Commands:    []string{"bal format --dry-run"},
	},
}

func NewFormatCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "format [OPTIONS] [<package>|<module>|<source-file>]",
		Short:         "Format the Ballerina source files",
		Long:          formatLongDesc,
		Example:       templates.FormatExamples(formatExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runFormat,
	}

	return cmd
}

func runFormat(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'format' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
