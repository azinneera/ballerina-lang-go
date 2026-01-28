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

package templates

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// ValidateArgs validates arguments for help and displays the help information.
func ValidateArgs(args []string, rootCmd *cobra.Command, cmdArr []*cobra.Command) {
	var foundCmd *cobra.Command
	var err error

	if len(args) == 0 {
		// No args, show root help
		Executing_Help_Template(*rootCmd, CommandGroups{})
		return
	}

	// Find the command by traversing the args
	foundCmd, _, err = rootCmd.Find(args)
	if err != nil || foundCmd == nil || foundCmd == rootCmd {
		fmt.Fprintf(os.Stderr, "'%s' is not a valid command.\n", args[len(args)-1])
		os.Exit(1)
	}

	// Check if the found command name matches the last arg
	if foundCmd.Name() != args[len(args)-1] {
		fmt.Fprintf(os.Stderr, "'%s' is not a valid subcommand of '%s'.\n",
			args[len(args)-1], args[len(args)-2])
		os.Exit(1)
	}

	// Print help using cobra command properties
	PrintCommandHelp(foundCmd)
}
