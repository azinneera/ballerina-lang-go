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

import "strings"

// Example represents a single command example with description
type Example struct {
	Description string   // What this example does
	Commands    []string // One or more command lines (without "$ " prefix)
}

// FormatExamples converts structured examples to the expected string format
func FormatExamples(examples []Example) string {
	var parts []string
	for _, ex := range examples {
		parts = append(parts, ex.Description)
		for _, cmd := range ex.Commands {
			parts = append(parts, "$ "+cmd)
		}
		parts = append(parts, "") // blank line between examples
	}
	// Trim trailing blank line
	result := strings.Join(parts, "\n")
	return strings.TrimSuffix(result, "\n")
}
