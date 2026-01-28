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

var pullLongDesc = `Download the specified package from Ballerina Central along with its dependencies and cache
it in the '.ballerina' directory in the user home.

Ballerina Central is a package repository hosted at https://central.ballerina.io/. A package repository
organizes packages into a three-level hierarchy: organization, package name, and version.`

var pullExamples = []templates.Example{
	{
		Description: "Pull a package from Ballerina Central.",
		Commands:    []string{"bal pull ballerina/io"},
	},
	{
		Description: "Pull a specific version of a package.",
		Commands:    []string{"bal pull ballerina/io:1.0.0"},
	},
}

func NewPullCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "pull <org-name>/<package-name>[:<version>]",
		Short:         "Fetch packages from Ballerina Central or a custom package repository",
		Long:          pullLongDesc,
		Example:       templates.FormatExamples(pullExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runPull,
	}

	return cmd
}

func runPull(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'pull' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
