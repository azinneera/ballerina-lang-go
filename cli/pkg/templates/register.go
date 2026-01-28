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
	"github.com/spf13/cobra"
)

// RegisterCommand adds a command to the parent and sets up the custom help function.
// This ensures all commands (including dynamic/third-party ones) use consistent help formatting.
// Returns the command for chaining.
func RegisterCommand(parent *cobra.Command, cmd *cobra.Command) *cobra.Command {
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		PrintCommandHelp(c)
	})
	parent.AddCommand(cmd)
	return cmd
}

// RegisterCommands registers multiple commands at once with the custom help function.
func RegisterCommands(parent *cobra.Command, cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		RegisterCommand(parent, cmd)
	}
}

// SetPlaceholder sets a custom placeholder for a flag's value in help output.
// For example, SetPlaceholder(cmd, "cloud", "provider") will display "--cloud <provider>"
// instead of the default "--cloud <cloud>".
// This is optional - if not set, the flag name is used as the placeholder.
func SetPlaceholder(cmd *cobra.Command, flagName, placeholder string) {
	flag := cmd.Flags().Lookup(flagName)
	if flag == nil {
		return
	}
	if flag.Annotations == nil {
		flag.Annotations = make(map[string][]string)
	}
	flag.Annotations["placeholder"] = []string{placeholder}
}
