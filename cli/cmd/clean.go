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

var cleanLongDesc = `Remove the 'target' directory created during the build. This will delete all files and
directories in the target directory.

The default target directory that will be cleaned is the '<project-root>/target'. You may also
provide a custom target to be cleaned.

Additionally, the caches of the dependencies of the current package or workspace can also be
cleaned.`

var cleanExamples = []templates.Example{
	{
		Description: "Clean the target directory.",
		Commands:    []string{"bal clean"},
	},
	{
		Description: "Clean a custom target directory.",
		Commands:    []string{"bal clean --target-dir=build"},
	},
}

func NewCleanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "clean [OPTIONS]",
		Short:         "Remove the artifacts created during the build",
		Long:          cleanLongDesc,
		Example:       templates.FormatExamples(cleanExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runClean,
	}

	return cmd
}

func runClean(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'clean' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
