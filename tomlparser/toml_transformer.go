// Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package tomlparser

import (
	"fmt"
	"time"

	"ballerina-lang-go/tools/diagnostics"
)

// TomlTransformer transforms parsed TOML document into a semantic AST.
// Maps to: io.ballerina.toml.semantic.ast.TomlTransformer
type TomlTransformer struct {
	// diagnostics collects errors during transformation
	// Maps to: DiagnosticLog dlog in Java
	diagnostics []Diagnostic
}

// NewTomlTransformer creates a new transformer instance.
// Maps to: TomlTransformer() constructor [line 66-68]
func NewTomlTransformer() *TomlTransformer {
	return &TomlTransformer{
		diagnostics: make([]Diagnostic, 0),
	}
}

// Transform converts a Toml document into a semantic AST.
// Maps to: transform(DocumentNode) [line 77-114]
// Returns TomlNode interface (concrete type is *TomlTableNode for root)
func (t *TomlTransformer) Transform(toml *Toml) TomlNode {
	rootTable := t.createRootTable()
	data := toml.ToMap()
	t.transformMembers(rootTable, data)

	for _, diag := range toml.Diagnostics() {
		rootTable.AddDiagnostic(diag)
	}

	return rootTable
}

// createRootTable creates the root table node for the document.
// Maps to: createRootTable(DocumentNode) [line 116-125]
func (t *TomlTransformer) createRootTable() *TomlTableNode {
	return NewTomlTableNode("__root", nil, false)
}

// transformMembers processes all top-level members of the document.
// Maps to: the loop in transform(DocumentNode) [line 81-112]
func (t *TomlTransformer) transformMembers(parent *TomlTableNode, data map[string]any) {
	for key, value := range data {
		switch v := value.(type) {
		case map[string]any:
			// Maps to: case TABLE [line 91-94]
			tableChild := t.transformTable(key, v)
			t.addChildTableToParent(parent, tableChild)

		case []any:
			if t.isTableArray(v) {
				// Maps to: case TABLE_ARRAY [line 98-101]
				tableArrayChild := t.transformTableArray(key, v)
				t.addChildParentArrayToParent(parent, tableArrayChild)
			} else {
				// Maps to: case KEY_VALUE [line 103-106]
				kvChild := t.transformKeyValue(key, value)
				t.addChildKeyValueToParent(parent, kvChild)
			}

		default:
			// Maps to: case KEY_VALUE [line 103-106]
			kvChild := t.transformKeyValue(key, value)
			t.addChildKeyValueToParent(parent, kvChild)
		}
	}
}

// isTableArray checks if an array is an array of tables.
func (t *TomlTransformer) isTableArray(arr []any) bool {
	if len(arr) == 0 {
		return false
	}
	for _, elem := range arr {
		if _, ok := elem.(map[string]any); !ok {
			return false
		}
	}
	return true
}

// transformTable transforms a map into a TomlTableNode.
// Maps to: transform(TableNode) [line 312-317]
func (t *TomlTransformer) transformTable(key string, data map[string]any) *TomlTableNode {
	tableNode := NewTomlTableNode(key, nil, false)
	t.addChildToTable(tableNode, data)
	return tableNode
}

// addChildToTable adds all children to a table node.
// Maps to: addChildToTable(TableNode, TomlTableNode) [line 319-336]
func (t *TomlTransformer) addChildToTable(tableNode *TomlTableNode, data map[string]any) {
	for key, value := range data {
		switch v := value.(type) {
		case map[string]any:
			childTable := t.transformTable(key, v)
			t.addChildTableToParent(tableNode, childTable)

		case []any:
			if t.isTableArray(v) {
				tableArrayChild := t.transformTableArray(key, v)
				t.addChildParentArrayToParent(tableNode, tableArrayChild)
			} else {
				kvChild := t.transformKeyValue(key, value)
				t.addChildKeyValueToParent(tableNode, kvChild)
			}

		default:
			kvChild := t.transformKeyValue(key, value)
			t.addChildKeyValueToParent(tableNode, kvChild)
		}
	}
}

// transformTableArray transforms an array of tables into TomlTableArrayNode.
// Maps to: transform(TableArrayNode) [line 360-367]
func (t *TomlTransformer) transformTableArray(key string, arr []any) *TomlTableArrayNode {
	tableArrayNode := NewTomlTableArrayNode(key, nil)

	// Maps to: addChildsToTableArray() [line 369-388]
	for _, elem := range arr {
		if tableData, ok := elem.(map[string]any); ok {
			childTable := NewTomlTableNode(key, nil, false)
			t.addChildToTable(childTable, tableData)
			tableArrayNode.AddChild(childTable)
		}
	}

	return tableArrayNode
}

// transformKeyValue transforms a key-value pair into TomlKeyValueNode.
// Maps to: transform(KeyValueNode) [line 398-423]
func (t *TomlTransformer) transformKeyValue(key string, value any) *TomlKeyValueNode {
	valueNode := t.transformValue(value)
	return NewTomlKeyValueNode(key, valueNode, nil)
}

// addChildKeyValueToParent adds a key-value node to a parent table.
// Maps to: addChildKeyValueToParent(TomlTableNode, TomlKeyValueNode) [line 127-163]
func (t *TomlTransformer) addChildKeyValueToParent(parent *TomlTableNode, kv *TomlKeyValueNode) {
	t.addChildToTableAST(parent, kv)
}

// addChildToTableAST adds a child node to a table, checking for duplicates.
// Maps to: addChildToTableAST(TomlTableNode, TopLevelNode) [line 165-174]
func (t *TomlTransformer) addChildToTableAST(parent *TomlTableNode, child TopLevelNode) {
	key := child.Key()
	entries := parent.Entries()

	if _, exists := entries[key]; exists {
		// Maps to: dlog.error(value.location(), DiagnosticErrorCode.ERROR_EXISTING_NODE, key)
		diag := t.createDiagnostic(ErrorExistingNode, child.Location(), key)
		parent.AddDiagnostic(diag)
		return
	}

	entries[key] = child
}

// addChildParentArrayToParent adds a table array to the parent.
// Maps to: addChildParentArrayToParent(TomlTableNode, TomlTableArrayNode) [line 188-211]
func (t *TomlTransformer) addChildParentArrayToParent(parent *TomlTableNode, tableArray *TomlTableArrayNode) {
	key := tableArray.Key()
	entries := parent.Entries()

	existing, exists := entries[key]
	if !exists {
		t.addChildToTableAST(parent, tableArray)
		return
	}

	if existingArray, ok := existing.(*TomlTableArrayNode); ok {
		// Maps to: ((TomlTableArrayNode) topLevelNode).addChild(...)
		for _, child := range tableArray.Children() {
			existingArray.AddChild(child)
		}
	} else {
		diag := t.createDiagnostic(ErrorExistingNode, tableArray.Location(), key)
		parent.AddDiagnostic(diag)
	}
}

// addChildTableToParent adds a table node to the parent.
// Maps to: addChildTableToParent(TomlTableNode, TomlTableNode) [line 243-274]
func (t *TomlTransformer) addChildTableToParent(parent *TomlTableNode, tableChild *TomlTableNode) {
	key := tableChild.Key()
	entries := parent.Entries()

	existing, exists := entries[key]
	if !exists {
		t.addChildToTableAST(parent, tableChild)
		return
	}

	if existingTable, ok := existing.(*TomlTableNode); ok {
		if existingTable.Generated() {
			// Maps to: parentTable.replaceGeneratedTable(newTableNode)
			parent.ReplaceGeneratedTable(tableChild)
		} else {
			diag := t.createDiagnostic(ErrorExistingNode, tableChild.Location(), key)
			parent.AddDiagnostic(diag)
		}
	} else {
		diag := t.createDiagnostic(ErrorExistingNode, tableChild.Location(), key)
		parent.AddDiagnostic(diag)
	}
}

// getParentTable navigates to or creates the parent table for a given path.
// Maps to: getParentTable(TomlTableNode, TopLevelNode) [line 217-241]
func (t *TomlTransformer) getParentTable(rootTable *TomlTableNode, keyPath []string) *TomlTableNode {
	if len(keyPath) <= 1 {
		return rootTable
	}

	parentTable := rootTable
	for i := 0; i < len(keyPath)-1; i++ {
		parentKey := keyPath[i]
		entries := parentTable.Entries()
		existing, exists := entries[parentKey]

		if !exists {
			// Maps to: generateTable() [line 291-302]
			implicitTable := t.generateTable(parentTable, parentKey)
			parentTable = implicitTable
			continue
		}

		switch node := existing.(type) {
		case *TomlTableNode:
			parentTable = node
		case *TomlTableArrayNode:
			// Maps to: children.get(children.size() - 1) [line 233]
			children := node.Children()
			if len(children) > 0 {
				parentTable = children[len(children)-1]
			}
		default:
			diag := t.createDiagnostic(ErrorExistingNode, nil, parentKey)
			parentTable.AddDiagnostic(diag)
			return parentTable
		}
	}

	return parentTable
}

// generateTable creates an implicit (generated) table.
// Maps to: generateTable(TomlTableNode, TomlKeyEntryNode, TopLevelNode) [line 291-302]
func (t *TomlTransformer) generateTable(parentTable *TomlTableNode, key string) *TomlTableNode {
	table := NewTomlTableNode(key, nil, true)
	t.addChildToTableAST(parentTable, table)
	return table
}

// transformValue converts a value into a TomlValueNode.
// Maps to: transformValue(ValueNode) [line 435-437]
func (t *TomlTransformer) transformValue(value any) TomlValueNode {
	switch v := value.(type) {
	case string:
		// Maps to: transform(StringLiteralNode) [line 472-488]
		return NewTomlStringValueNode(v, nil)

	case int:
		// Maps to: getTomlNode() returning TomlLongValueNode [line 558-559]
		return NewTomlLongValueNode(int64(v), nil)

	case int64:
		return NewTomlLongValueNode(v, nil)

	case float64:
		// Maps to: getTomlNode() returning TomlDoubleValueNodeNodeNode [line 573-574]
		return NewTomlDoubleValueNodeNode(v, nil)

	case bool:
		// Maps to: transform(BoolLiteralNode) [line 585-590]
		return NewTomlBooleanValueNode(v, nil)

	case time.Time:
		return NewTomlStringValueNode(v.Format(time.RFC3339), nil)

	case []any:
		// Maps to: transform(ArrayNode) [line 446-454]
		elements := make([]TomlValueNode, 0, len(v))
		for _, elem := range v {
			elements = append(elements, t.transformValue(elem))
		}
		return NewTomlArrayValueNode(elements, nil)

	case map[string]any:
		// Maps to: transform(InlineTableNode) [line 600-608]
		elements := make([]TopLevelNode, 0, len(v))
		for k, val := range v {
			valueNode := t.transformValue(val)
			kvNode := NewTomlKeyValueNode(k, valueNode, nil)
			elements = append(elements, kvNode)
		}
		return NewTomlInlineTableValueNode(elements, nil)

	default:
		return NewTomlStringValueNode(fmt.Sprintf("%v", v), nil)
	}
}

// createDiagnostic creates a diagnostic with the given code and arguments.
// Maps to: dlog.error() calls in TomlTransformer
func (t *TomlTransformer) createDiagnostic(code TomlDiagnosticCode, location *Location, args ...any) Diagnostic {
	return Diagnostic{
		Code:     code,
		Message:  code.Format(args...),
		Severity: diagnostics.Error,
		Location: location,
	}
}

// Diagnostics returns all collected diagnostics.
func (t *TomlTransformer) Diagnostics() []Diagnostic {
	return t.diagnostics
}
