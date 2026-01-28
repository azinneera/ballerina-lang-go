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
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const maxLineWidth = 80
const indentWidth = 7   // 7 spaces for main content
const optionIndent = 11 // 11 spaces for option descriptions

type SubCommandInfo struct {
	Name  string
	Short string
}

type CommandData struct {
	Name        string
	Short       string
	Synopsis    string
	Description string
	Options     string
	Arguments   string
	Commands    []SubCommandInfo
	Examples    string
}

const subcommandTemplate = `NAME
       ballerina-{{.Name}} - {{.Short}}

SYNOPSIS
       bal {{.Synopsis}}

{{if .Description}}
DESCRIPTION
{{.Description}}

{{end}}{{if .Options}}
OPTIONS
{{.Options}}
{{end}}{{if .Arguments}}
ARGUMENTS
{{.Arguments}}
{{end}}{{if .Commands}}
SUBCOMMANDS
{{range .Commands}}       {{.Name | printf "%-15s"}} {{.Short}}
{{end}}
{{end}}{{if .Examples}}
EXAMPLES
{{.Examples}}
{{end}}{{if .Commands}}
Use 'bal {{.Name}} <subcommand> --help' for more information on a specific subcommand.
{{end}}`

// FillTemplate renders the help template to stdout
func FillTemplate(commandData CommandData) {
	FillTemplateToWriter(os.Stdout, commandData)
}

// FillTemplateToWriter renders the help template to the given writer
func FillTemplateToWriter(w io.Writer, commandData CommandData) {
	tmpl, err := template.New("subHelpTemplate").Parse(subcommandTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template: %v\n", err)
		return
	}

	err = tmpl.Execute(w, commandData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template: %v\n", err)
	}
}

// GetSubCommands extracts subcommand info from a cobra command
func GetSubCommands(cmd *cobra.Command) []SubCommandInfo {
	if len(cmd.Commands()) == 0 {
		return nil
	}

	var commands []SubCommandInfo
	for _, subCmd := range cmd.Commands() {
		if subCmd.Hidden {
			continue
		}
		commands = append(commands, SubCommandInfo{
			Name:  subCmd.Name(),
			Short: sanitizeText(subCmd.Short),
		})
	}

	return commands
}

// GetCommandDataFromCobra builds CommandData from cobra.Command properties
func GetCommandDataFromCobra(cmd *cobra.Command) CommandData {
	return CommandData{
		Name:        getFullCommandName(cmd),
		Short:       sanitizeText(cmd.Short),
		Synopsis:    getSynopsis(cmd),
		Description: formatDescription(cmd.Long),
		Options:     formatOptions(cmd),
		Arguments:   formatArguments(cmd),
		Commands:    GetSubCommands(cmd),
		Examples:    formatExamples(cmd.Example),
	}
}

// getFullCommandName returns the full command name (e.g., "tool-list" for "bal tool list")
func getFullCommandName(cmd *cobra.Command) string {
	names := []string{}
	for c := cmd; c != nil && c.Name() != "bal"; c = c.Parent() {
		names = append([]string{c.Name()}, names...)
	}
	return strings.Join(names, "-")
}

// getSynopsis extracts the synopsis from the Use field
func getSynopsis(cmd *cobra.Command) string {
	// Build the full command path
	var cmdPath []string
	for c := cmd; c != nil && c.Name() != "bal"; c = c.Parent() {
		cmdPath = append([]string{c.Name()}, cmdPath...)
	}

	// The Use field typically contains "command [flags] [args]"
	// We want to extract just the arguments/flags part
	use := sanitizeText(cmd.Use)
	parts := strings.SplitN(use, " ", 2)
	if len(parts) > 1 {
		// Combine command path with the rest of the Use string
		return strings.Join(cmdPath, " ") + " " + parts[1]
	}
	return strings.Join(cmdPath, " ")
}

// sanitizeText normalizes whitespace and line endings in text.
// This handles cases where third-party developers may have inconsistent formatting.
func sanitizeText(text string) string {
	if text == "" {
		return ""
	}
	// Replace Windows line endings
	text = strings.ReplaceAll(text, "\r\n", "\n")
	// Replace carriage returns
	text = strings.ReplaceAll(text, "\r", "\n")
	// Trim leading/trailing whitespace
	text = strings.TrimSpace(text)
	return text
}

// normalizeWhitespace collapses multiple spaces into single space
func normalizeWhitespace(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

// wrapText wraps text at the specified width, preserving the indent for continuation lines
func wrapText(text string, firstLineIndent, continuationIndent int, maxWidth int) string {
	if text == "" {
		return ""
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var lines []string
	currentLine := strings.Repeat(" ", firstLineIndent)
	currentWidth := firstLineIndent

	for i, word := range words {
		wordLen := len(word)

		// Check if adding this word would exceed max width
		if currentWidth+wordLen+1 > maxWidth && currentWidth > firstLineIndent {
			// Start a new line
			lines = append(lines, currentLine)
			currentLine = strings.Repeat(" ", continuationIndent) + word
			currentWidth = continuationIndent + wordLen
		} else {
			// Add word to current line
			if i == 0 || currentWidth == firstLineIndent || currentWidth == continuationIndent {
				currentLine += word
				currentWidth += wordLen
			} else {
				currentLine += " " + word
				currentWidth += wordLen + 1
			}
		}
	}

	// Don't forget the last line
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

// formatDescription formats the long description with proper indentation and line wrapping
func formatDescription(desc string) string {
	desc = sanitizeText(desc)
	if desc == "" {
		return ""
	}

	// Split into paragraphs (separated by blank lines)
	paragraphs := strings.Split(desc, "\n\n")
	var formatted []string

	for _, para := range paragraphs {
		// Normalize whitespace within paragraph
		para = strings.TrimSpace(para)
		para = normalizeWhitespace(para)

		if para == "" {
			formatted = append(formatted, "")
			continue
		}

		// Wrap the paragraph
		wrapped := wrapText(para, indentWidth, indentWidth, maxLineWidth)
		formatted = append(formatted, wrapped)
	}

	return strings.Join(formatted, "\n\n")
}

// formatOptions formats the command flags in the Java help style
func formatOptions(cmd *cobra.Command) string {
	var result strings.Builder

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}

		// Build the flag line: "       -s <value>, --long <value>" or "       --long"
		var flagLine strings.Builder
		flagLine.WriteString("       ")

		// Get the value placeholder from annotations or use flag name
		valuePlaceholder := f.Name
		if f.Annotations != nil {
			if placeholder, ok := f.Annotations["placeholder"]; ok && len(placeholder) > 0 {
				valuePlaceholder = placeholder[0]
			}
		}

		if f.Shorthand != "" {
			if f.Value.Type() != "bool" {
				flagLine.WriteString(fmt.Sprintf("-%s <%s>, ", f.Shorthand, valuePlaceholder))
			} else {
				flagLine.WriteString(fmt.Sprintf("-%s, ", f.Shorthand))
			}
		}

		flagLine.WriteString("--")
		flagLine.WriteString(f.Name)

		if f.Value.Type() != "bool" {
			flagLine.WriteString(fmt.Sprintf(" <%s>", valuePlaceholder))
		}

		result.WriteString(flagLine.String())
		result.WriteString("\n")

		// Format usage with proper indentation and line wrapping
		usage := sanitizeText(f.Usage)
		if usage != "" {
			usage = normalizeWhitespace(usage)
			wrapped := wrapText(usage, optionIndent, optionIndent, maxLineWidth)
			result.WriteString(wrapped)
			result.WriteString("\n")
		}

		result.WriteString("\n")
	})

	return result.String()
}

// formatArguments formats the arguments section
func formatArguments(cmd *cobra.Command) string {
	if cmd.Annotations != nil {
		if args, ok := cmd.Annotations["arguments"]; ok {
			return formatDescription(args)
		}
	}
	return ""
}

// formatExamples formats the examples section with proper indentation
func formatExamples(examples string) string {
	examples = sanitizeText(examples)
	if examples == "" {
		return ""
	}

	lines := strings.Split(examples, "\n")
	var formatted []string
	var currentDesc []string

	flushDesc := func() {
		if len(currentDesc) > 0 {
			// Join description lines and wrap them
			desc := strings.Join(currentDesc, " ")
			desc = normalizeWhitespace(desc)
			wrapped := wrapText(desc, indentWidth, indentWidth, maxLineWidth)
			formatted = append(formatted, wrapped)
			currentDesc = nil
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			flushDesc()
			formatted = append(formatted, "")
		} else if strings.HasPrefix(trimmed, "$") {
			flushDesc()
			// Command lines get extra indent
			formatted = append(formatted, "           "+trimmed)
		} else {
			// Accumulate description lines
			currentDesc = append(currentDesc, trimmed)
		}
	}
	flushDesc()

	return strings.Join(formatted, "\n")
}

// PrintCommandHelp prints the help for a command using the custom template
func PrintCommandHelp(cmd *cobra.Command) {
	data := GetCommandDataFromCobra(cmd)
	FillTemplate(data)
}

// PrintCommandHelpToWriter prints the help for a command to the given writer
func PrintCommandHelpToWriter(w io.Writer, cmd *cobra.Command) {
	data := GetCommandDataFromCobra(cmd)
	FillTemplateToWriter(w, data)
}
