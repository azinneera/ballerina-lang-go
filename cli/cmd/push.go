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

var pushLongDesc = `Push the Ballerina archive (.bala) of the current package to Ballerina Central, local or a
custom remote repository. Once the package is pushed to Ballerina Central, it becomes public and sharable
and will be permanent.

To be able to publish a package to Ballerina Central, you should sign in to Ballerina Central and obtain
an access token.`

var pushExamples = []templates.Example{
	{
		Description: "Push the current package to Ballerina Central.",
		Commands:    []string{"bal push"},
	},
	{
		Description: "Push a .bala file to a local repository.",
		Commands:    []string{"bal push --repository=local"},
	},
}

func NewPushCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "push [OPTIONS] [<bala-path>]",
		Short:         "Push the Ballerina Archive (BALA) of the current package to a package repository",
		Long:          pushLongDesc,
		Example:       templates.FormatExamples(pushExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runPush,
	}

	return cmd
}

func runPush(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'push' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
