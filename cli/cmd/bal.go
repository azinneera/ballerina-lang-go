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
	"os"

	"ballerina-lang-go/cli/pkg/generate"
	"ballerina-lang-go/cli/pkg/templates"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "bal",
	Short:        "The build system and package manager of Ballerina",
	Long:         `The build system and package manager of Ballerina`,
	SilenceUsage: true,
}

func main() {
	// Core commands
	coreCommands := []*cobra.Command{
		NewBuildCommand(),
		NewRunCommand(),
		NewTestCommand(),
		NewDocCommand(),
		NewPackCommand(),
	}

	// Package commands
	packageCommands := []*cobra.Command{
		NewNewCommand(),
		NewAddCommand(),
		NewPullCommand(),
		NewPushCommand(),
		NewSearchCommand(),
		NewSemverCommand(),
		NewGraphCommand(),
		NewDeprecateCommand(),
	}

	// Other commands
	otherCommands := []*cobra.Command{
		NewCleanCommand(),
		NewFormatCommand(),
		NewBindgenCommand(),
		NewShellCommand(),
		NewToolCommand(),
		NewProfileCommand(),
	}

	// Register all commands with automatic help function setup
	templates.RegisterCommands(rootCmd, coreCommands...)
	templates.RegisterCommands(rootCmd, packageCommands...)
	templates.RegisterCommands(rootCmd, otherCommands...)

	// Load dynamic tool commands
	toolList := generate.GetTools(rootCmd)
	toolsCobra := generate.GetCommandsList(toolList, rootCmd)

	// Build command groups for help display
	commandGroups := templates.CommandGroups{
		{Message: "Core Commands", Commands: coreCommands},
		{Message: "Package Commands", Commands: packageCommands},
		{Message: "Other Commands", Commands: otherCommands},
		{Message: "Tool Commands", Commands: toolsCobra},
	}

	// Create and register help command
	helpCmd := NewHelpCommand(rootCmd, commandGroups, toolsCobra)
	templates.RegisterCommand(rootCmd, helpCmd)

	// Register version command
	templates.RegisterCommand(rootCmd, versionCmd)

	// Set custom help function on root command only.
	// Subcommands use the help function set by RegisterCommand.
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd != rootCmd {
			// Subcommands get their help from RegisterCommand
			templates.PrintCommandHelp(cmd)
			return
		}
		if len(args) <= 1 {
			templates.Executing_Help_Template(*cmd, commandGroups)
		} else {
			templates.ValidateArgs(args[:len(args)-1], rootCmd, toolsCobra)
		}
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
