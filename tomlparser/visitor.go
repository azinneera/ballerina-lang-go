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

// TomlNodeVisitor is the visitor interface for TOML AST nodes.
// Maps to: io.ballerina.toml.semantic.ast.TomlNodeVisitor
type TomlNodeVisitor interface {
	// VisitTomlTableNode maps to: visit(TomlTableNode)
	VisitTomlTableNode(node *TomlTableNode)
	// VisitTomlTableArrayNode maps to: visit(TomlTableArrayNode)
	VisitTomlTableArrayNode(node *TomlTableArrayNode)
	// VisitTomlKeyValueNode maps to: visit(TomlKeyValueNode)
	VisitTomlKeyValueNode(node *TomlKeyValueNode)
	// VisitTomlKeyNode maps to: visit(TomlKeyNode)
	VisitTomlKeyNode(node *TomlKeyNode)
	// VisitTomlValueNode maps to: visit(TomlValueNode)
	VisitTomlValueNode(node TomlValueNode)
	// VisitTomlArrayValueNode maps to: visit(TomlArrayValueNode)
	VisitTomlArrayValueNode(node *TomlArrayValueNode)
	// VisitTomlKeyEntryNode maps to: visit(TomlKeyEntryNode)
	VisitTomlKeyEntryNode(node *TomlKeyEntryNode)
	// VisitTomlStringValueNode maps to: visit(TomlStringValueNode)
	VisitTomlStringValueNode(node *TomlStringValueNode)
	// VisitTomlDoubleValueNodeNode maps to: visit(TomlDoubleValueNodeNode)
	VisitTomlDoubleValueNodeNode(node *TomlDoubleValueNodeNode)
	// VisitTomlLongValueNode maps to: visit(TomlLongValueNode)
	VisitTomlLongValueNode(node *TomlLongValueNode)
	// VisitTomlBooleanValueNode maps to: visit(TomlBooleanValueNode)
	VisitTomlBooleanValueNode(node *TomlBooleanValueNode)
	// VisitTomlInlineTableValueNode maps to: visit(TomlInlineTableValueNode)
	VisitTomlInlineTableValueNode(node *TomlInlineTableValueNode)
}

// BaseVisitor provides default implementations that traverse children.
// Extend this by embedding and override specific methods as needed.
type BaseVisitor struct{}

// VisitTomlTableNode visits a table node and its children.
func (v *BaseVisitor) VisitTomlTableNode(node *TomlTableNode) {
	for _, child := range node.Entries() {
		child.Accept(v)
	}
}

// VisitTomlTableArrayNode visits a table array node and its children.
func (v *BaseVisitor) VisitTomlTableArrayNode(node *TomlTableArrayNode) {
	for _, child := range node.Children() {
		child.Accept(v)
	}
}

// VisitTomlKeyValueNode visits a key-value node and its value.
func (v *BaseVisitor) VisitTomlKeyValueNode(node *TomlKeyValueNode) {
	if node.Value() != nil {
		node.Value().Accept(v)
	}
}

// VisitTomlKeyNode visits a key node and its entries.
func (v *BaseVisitor) VisitTomlKeyNode(node *TomlKeyNode) {
	for _, entry := range node.Keys() {
		entry.Accept(v)
	}
}

// VisitTomlValueNode visits a generic value node.
func (v *BaseVisitor) VisitTomlValueNode(node TomlValueNode) {}

// VisitTomlArrayValueNode visits an array value node and its elements.
func (v *BaseVisitor) VisitTomlArrayValueNode(node *TomlArrayValueNode) {
	for _, elem := range node.Elements() {
		elem.Accept(v)
	}
}

// VisitTomlKeyEntryNode visits a key entry node (leaf node).
func (v *BaseVisitor) VisitTomlKeyEntryNode(node *TomlKeyEntryNode) {}

// VisitTomlStringValueNode visits a string value node (leaf node).
func (v *BaseVisitor) VisitTomlStringValueNode(node *TomlStringValueNode) {}

// VisitTomlDoubleValueNodeNode visits a double value node (leaf node).
func (v *BaseVisitor) VisitTomlDoubleValueNodeNode(node *TomlDoubleValueNodeNode) {}

// VisitTomlLongValueNode visits a long value node (leaf node).
func (v *BaseVisitor) VisitTomlLongValueNode(node *TomlLongValueNode) {}

// VisitTomlBooleanValueNode visits a boolean value node (leaf node).
func (v *BaseVisitor) VisitTomlBooleanValueNode(node *TomlBooleanValueNode) {}

// VisitTomlInlineTableValueNode visits an inline table node and its elements.
func (v *BaseVisitor) VisitTomlInlineTableValueNode(node *TomlInlineTableValueNode) {
	for _, elem := range node.Elements() {
		elem.Accept(v)
	}
}
