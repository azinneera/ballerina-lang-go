package cmd

import (
	"ballerina-lang-go/cli/bal/pkg/utils"

	"github.com/spf13/cobra"
)

func newDummyCmd(use, short string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Run: func(cmd *cobra.Command, args []string) {
			_ = utils.ExecuteBallerinaCommand(javaCmdPass, cmdLineArgsPass)
		},
	}
}

// Core Commands

func buildCmd() *cobra.Command {
	return newDummyCmd("build", "Compile the current package")
}

func runCmd() *cobra.Command {
	return newDummyCmd("run", "Compile and run the current package")
}

func testCmd() *cobra.Command {
	return newDummyCmd("test", "Run package tests")
}

func docCmd() *cobra.Command {
	return newDummyCmd("doc", "Generate current package's documentation")
}

func packCmd() *cobra.Command {
	return newDummyCmd("pack", "Create distribution format of the current package")
}

// Package Commands

func newCmd() *cobra.Command {
	return newDummyCmd("new", "Create a new Ballerina package")
}

func addCmd() *cobra.Command {
	return newDummyCmd("add", "Add a new module to the current package")
}

func pullCmd() *cobra.Command {
	return newDummyCmd("pull", "Pull a package from Ballerina Central")
}

func pushCmd() *cobra.Command {
	return newDummyCmd("push", "Publish a package to Ballerina Central")
}

func searchCmd() *cobra.Command {
	return newDummyCmd("search", "Search Ballerina Central for packages")
}

func semverCmd() *cobra.Command {
	return newDummyCmd("semver", "Show SemVer compatibility and local changes against published packages")
}

func graphCmd() *cobra.Command {
	return newDummyCmd("graph", "Print the dependency graph")
}

func deprecateCmd() *cobra.Command {
	return newDummyCmd("deprecate", "Deprecate a package in Ballerina Central")
}

// Other Commands

func cleanCmd() *cobra.Command {
	return newDummyCmd("clean", "Clean the artifacts generated during the build")
}

func formatCmd() *cobra.Command {
	return newDummyCmd("format", "Format Ballerina source files")
}

func grpcCmd() *cobra.Command {
	return newDummyCmd("grpc", "Generate Ballerina sources for the given protobuf definition")
}

func graphqlCmd() *cobra.Command {
	return newDummyCmd("graphql", "Generate Ballerina client sources for a GraphQL config")
}

func openapiCmd() *cobra.Command {
	return newDummyCmd("openapi", "Generate Ballerina sources for an OpenAPI contract")
}

func asyncapiCmd() *cobra.Command {
	return newDummyCmd("asyncapi", "Generate a Ballerina listener from an AsyncAPI contract")
}

func persistCmd() *cobra.Command {
	return newDummyCmd("persist", "Manage data persistence")
}

func bindgenCmd() *cobra.Command {
	return newDummyCmd("bindgen", "Generate Ballerina bindings for Java APIs")
}

func shellCmd() *cobra.Command {
	return newDummyCmd("shell", "Run Ballerina interactive shell")
}

func toolCmd() *cobra.Command {
	return newDummyCmd("tool", "Manage Ballerina tool commands")
}

func versionCmd() *cobra.Command {
	return newDummyCmd("version", "Print the Ballerina version")
}

func profileCmd() *cobra.Command {
	return newDummyCmd("profile", "Profile Ballerina programs")
}
