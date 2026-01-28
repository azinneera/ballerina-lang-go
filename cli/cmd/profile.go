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

var profileLongDesc = `Compile the current package and run the program with Ballerina Profiler. This command
generates an html file with name 'ProfilerReport.html' in the target directory.

The generated 'ProfilerReport.html' file contains the flame graph that visualizes the distributed
Ballerina functions with execution details.

Note: This is an experimental feature, which supports only a limited set of functionality.`

var profileExamples = []templates.Example{
	{
		Description: "Profile the current package.",
		Commands:    []string{"bal profile"},
	},
	{
		Description: "Profile a specific source file.",
		Commands:    []string{"bal profile main.bal"},
	},
}

func NewProfileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "profile [OPTIONS] [<package>|<source-file>]",
		Short:         "Run Ballerina Profiler on the source and generate flame graph",
		Long:          profileLongDesc,
		Example:       templates.FormatExamples(profileExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runProfile,
	}

	return cmd
}

func runProfile(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'profile' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
