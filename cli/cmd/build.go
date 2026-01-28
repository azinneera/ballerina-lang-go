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

var buildLongDesc = `Compile a package and its dependencies to produce a standalone executable or compile a
workspace and generate executables.

Build writes the resulting executable to the 'target/bin' directory when compiling the current package.

Build writes the resulting executable to the current directory when compiling a single '.bal' file. The
generated executable is named after the source file.

Note: Building individual '.bal' files of a package is not allowed. For workspaces, the executable is
generated only for packages that do not have dependents except when the package is built explicitly.`

var buildExamples = []templates.Example{
	{
		Description: "Build the current package or the workspace. This will generate an executable in " +
			"the 'target/bin' directory of the package. For workspaces, this will generate executables " +
			"for all packages that does not have dependants in the workspace.",
		Commands: []string{"bal build"},
	},
	{
		Description: "Build a single '.bal' file. This will generate an executable in the current directory.",
		Commands:    []string{"bal build app.bal"},
	},
	{
		Description: "Build a specific package in the workspace. This will generate an executable for " +
			"the specified package.",
		Commands: []string{"bal build <package-path-in-workspace>"},
	},
	{
		Description: "Build the 'app' package from a different directory.",
		Commands:    []string{"bal build <app-package-path>"},
	},
	{
		Description: "Build the package with additional GraalVM native image options.",
		Commands:    []string{`bal build --graalvm --graalvm-build-options="--static --enable-monitoring"`},
	},
}

func NewBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "build [OPTIONS] [<package>|<source-file>]",
		Short:         "Compiles the current package",
		Long:          buildLongDesc,
		Example:       templates.FormatExamples(buildExamples),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runBuild,
	}

	addBuildFlags(cmd)

	return cmd
}

func addBuildFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.Bool("export-openapi", false,
		"Generate OpenAPI definitions for service declarations in the current package.")

	flags.String("cloud", "",
		"Generate cloud artifacts: '--cloud=k8s' for Kubernetes and '--cloud=docker' for Docker.")
	templates.SetPlaceholder(cmd, "cloud", "provider")

	flags.Bool("list-conflicted-classes", false,
		"List the conflicting classes of conflicting JARs in the package.")

	flags.StringP("output", "o", "",
		"Write the output to the given file. The option only works for the single '.bal' file scenario.")

	flags.Bool("observability-included", false,
		"Include the dependencies that are required to enable observability.")

	flags.Bool("offline", false,
		"Proceed without accessing the network. Attempt to proceed with the previously downloaded "+
			"dependencies in local caches, will fail otherwise.")

	flags.Bool("sticky", false,
		"Attempt to stick to the dependency versions available in the 'Dependencies.toml' file. "+
			"If the file doesn't exist, this option is ignored.")

	flags.String("target-dir", "", "Target directory path.")
	templates.SetPlaceholder(cmd, "target-dir", "path")

	flags.Bool("graalvm", false,
		"Generate a GraalVM native image. Native image generation is an experimental feature which "+
			"supports only a limited set of functionality.")

	flags.String("graalvm-build-options", "",
		"Additional build options to be passed to the GraalVM native image.")

	flags.Bool("remote-management", false,
		"Include the dependencies that are required to enable remote package management service.")

	flags.Bool("show-dependency-diagnostics", false,
		"Print the diagnostics that are related to the dependencies. By default, these diagnostics "+
			"are not printed to the console.")

	flags.Bool("optimize-dependency-compilation", false,
		"[EXPERIMENTAL] Enables memory-efficient compilation of package dependencies using separate "+
			"processes. This can help prevent out-of-memory issues during the initial compilation "+
			"with a clean central cache.")

	flags.Bool("experimental", false, "Enable experimental language features.")
}

func runBuild(cmd *cobra.Command, args []string) error {
	err := fmt.Errorf("command 'build' is not yet implemented")
	printError(err, "", false, cmd.Name())
	return err
}
