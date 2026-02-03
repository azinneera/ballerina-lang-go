# TOML Parser Migration Plan

Migration from Java to Go for the Ballerina TOML parser.

## Migration Status Overview

| Phase | Component | Files | Status |
|-------|-----------|-------|--------|
| 1 | Foundation | 5 | Complete |
| 2 | Lexer | 10 | Not Started |
| 3 | Internal Syntax Tree | 25 | Not Started |
| 4 | Parser | 8 | Not Started |
| 5 | Public Syntax Tree | 30 | Not Started |
| 6 | Semantic AST | 22 | Not Started |
| 7 | Transformer | 3 | Not Started |
| 8 | Public API | 2 | Not Started |
| 9 | Validator | 18 | Not Started |
| 10 | Tests | 21 | Not Started |

**Legend**: Not Started | In Progress | Complete

---

## Phase 1: Foundation (Priority: Critical)

Core types and utilities needed by all other components.

### Files to Create

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `semantic/toml_type.go` | `TomlType.java` | Complete |
| `syntax/tree/syntax_kind.go` | `SyntaxKind.java` | Complete |
| `internal/diagnostics/diagnostic_error_code.go` | `DiagnosticErrorCode.java` | Complete |
| `internal/diagnostics/diagnostic_message_helper.go` | `DiagnosticMessageHelper.java` | Complete |
| `internal/parser/parser_mode.go` | `ParserMode.java` | Complete |

### Key Decisions
- [x] Define Go interfaces for node hierarchies (DiagnosticCode interface defined)
- [x] Choose error handling strategy (error values with DiagnosticErrorCode struct)
- [ ] Define position/location structures (deferred to Phase 3/4)

---

## Phase 2: Lexer (Priority: Critical)

Tokenizes input into a stream of tokens.

### Files to Create

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `internal/parser/lexer_terminals.go` | `LexerTerminals.java` | Not Started |
| `internal/parser/input_reader.go` | `InputReader.java` | Not Started |
| `internal/parser/abstract_lexer.go` | `AbstractLexer.java` | Not Started |
| `internal/parser/toml_lexer.go` | `TomlLexer.java` | Not Started |
| `internal/parser/abstract_token_reader.go` | `AbstractTokenReader.java` | Not Started |
| `internal/parser/token_reader.go` | `TokenReader.java` | Not Started |

### Dependencies
- Phase 1: `SyntaxKind`, `ParserMode`

---

## Phase 3: Internal Syntax Tree (Priority: Critical)

Internal representation of parsed TOML.

### Files to Create

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `internal/parser/tree/st_node.go` | `STNode.java` | Not Started |
| `internal/parser/tree/st_token.go` | `STToken.java` | Not Started |
| `internal/parser/tree/st_minutiae.go` | `STMinutiae.java` | Not Started |
| `internal/parser/tree/st_minutiae_list.go` | `STMinutiaeList.java` | Not Started |
| `internal/parser/tree/st_abstract_node_factory.go` | `STAbstractNodeFactory.java` | Not Started |
| `internal/parser/tree/st_node_factory.go` | `STNodeFactory.java` | Not Started |
| `internal/parser/tree/st_node_list.go` | `STNodeList.java` | Not Started |
| `internal/parser/tree/st_node_visitor.go` | `STNodeVisitor.java` | Not Started |
| `internal/parser/tree/st_node_transformer.go` | `STNodeTransformer.java` | Not Started |
| `internal/parser/tree/st_document_node.go` | `STDocumentNode.java` | Not Started |
| `internal/parser/tree/st_table_node.go` | `STTableNode.java` | Not Started |
| `internal/parser/tree/st_table_array_node.go` | `STTableArrayNode.java` | Not Started |
| `internal/parser/tree/st_key_value_node.go` | `STKeyValueNode.java` | Not Started |
| `internal/parser/tree/st_key_node.go` | `STKeyNode.java` | Not Started |
| `internal/parser/tree/st_array_node.go` | `STArrayNode.java` | Not Started |
| `internal/parser/tree/st_inline_table_node.go` | `STInlineTableNode.java` | Not Started |
| `internal/parser/tree/st_basic_literal_node.go` | `STBasicLiteralNode.java` | Not Started |
| `internal/parser/tree/st_string_literal_node.go` | `STStringLiteralNode.java` | Not Started |
| `internal/parser/tree/st_numeric_literal_node.go` | `STNumericLiteralNode.java` | Not Started |
| `internal/parser/tree/st_boolean_literal_node.go` | `STBooleanLiteralNode.java` | Not Started |
| `internal/parser/tree/st_identifier_token.go` | `STIdentifierToken.java` | Not Started |
| `internal/parser/tree/st_literal_value_token.go` | `STLiteralValueToken.java` | Not Started |
| `internal/parser/tree/st_missing_token.go` | `STMissingToken.java` | Not Started |
| `internal/parser/tree/st_invalid_node_minutiae.go` | `STInvalidNodeMinutiae.java` | Not Started |

### Dependencies
- Phase 1: `SyntaxKind`

---

## Phase 4: Parser (Priority: Critical)

Recursive descent parser that builds syntax tree.

### Files to Create

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `internal/parser/parser_rule_context.go` | `ParserRuleContext.java` | Not Started |
| `internal/parser/abstract_parser.go` | `AbstractParser.java` | Not Started |
| `internal/parser/toml_parser.go` | `TomlParser.java` | Not Started |
| `internal/parser/abstract_parser_error_handler.go` | `AbstractParserErrorHandler.java` | Not Started |
| `internal/parser/toml_parser_error_handler.go` | `TomlParserErrorHandler.java` | Not Started |
| `internal/parser/syntax_errors.go` | `SyntaxErrors.java` | Not Started |
| `internal/parser/parser_factory.go` | `ParserFactory.java` | Not Started |
| `internal/diagnostics/syntax_diagnostic.go` | `SyntaxDiagnostic.java` | Not Started |

### Dependencies
- Phase 2: Lexer
- Phase 3: Internal Syntax Tree

---

## Phase 5: Public Syntax Tree (Priority: High)

User-facing syntax tree API.

### Files to Create

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `syntax/tree/node.go` | `Node.java` | Not Started |
| `syntax/tree/non_terminal_node.go` | `NonTerminalNode.java` | Not Started |
| `syntax/tree/token.go` | `Token.java` | Not Started |
| `syntax/tree/syntax_tree.go` | `SyntaxTree.java` | Not Started |
| `syntax/tree/node_factory.go` | `NodeFactory.java` | Not Started |
| `syntax/tree/node_visitor.go` | `NodeVisitor.java` | Not Started |
| `syntax/tree/node_transformer.go` | `NodeTransformer.java` | Not Started |
| `syntax/tree/node_list.go` | `NodeList.java` | Not Started |
| `syntax/tree/separated_node_list.go` | `SeparatedNodeList.java` | Not Started |
| `syntax/tree/minutiae.go` | `Minutiae.java` | Not Started |
| `syntax/tree/minutiae_list.go` | `MinutiaeList.java` | Not Started |
| `syntax/tree/document_node.go` | `DocumentNode.java` | Not Started |
| `syntax/tree/table_node.go` | `TableNode.java` | Not Started |
| `syntax/tree/table_array_node.go` | `TableArrayNode.java` | Not Started |
| `syntax/tree/key_value_node.go` | `KeyValueNode.java` | Not Started |
| `syntax/tree/key_node.go` | `KeyNode.java` | Not Started |
| `syntax/tree/array_node.go` | `ArrayNode.java` | Not Started |
| `syntax/tree/inline_table_node.go` | `InlineTableNode.java` | Not Started |
| `syntax/tree/value_node.go` | `ValueNode.java` | Not Started |
| `syntax/tree/string_literal_node.go` | `StringLiteralNode.java` | Not Started |
| `syntax/tree/numeric_literal_node.go` | `NumericLiteralNode.java` | Not Started |
| `syntax/tree/boolean_literal_node.go` | `BooleanLiteralNode.java` | Not Started |
| `syntax/tree/identifier_literal_node.go` | `IdentifierLiteralNode.java` | Not Started |
| `internal/syntax/external_tree_node_list.go` | `ExternalTreeNodeList.java` | Not Started |
| `internal/syntax/node_list_utils.go` | `NodeListUtils.java` | Not Started |
| `internal/syntax/syntax_utils.go` | `SyntaxUtils.java` | Not Started |

### Dependencies
- Phase 3: Internal Syntax Tree
- Phase 4: Parser

---

## Phase 6: Semantic AST (Priority: High)

High-level typed representation for querying.

### Files to Create

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `semantic/ast/toml_node.go` | `TomlNode.java` | Not Started |
| `semantic/ast/top_level_node.go` | `TopLevelNode.java` | Not Started |
| `semantic/ast/toml_value_node.go` | `TomlValueNode.java` | Not Started |
| `semantic/ast/toml_basic_value_node.go` | `TomlBasicValueNode.java` | Not Started |
| `semantic/ast/toml_string_value_node.go` | `TomlStringValueNode.java` | Not Started |
| `semantic/ast/toml_long_value_node.go` | `TomlLongValueNode.java` | Not Started |
| `semantic/ast/toml_double_value_node.go` | `TomlDoubleValueNodeNode.java` | Not Started |
| `semantic/ast/toml_boolean_value_node.go` | `TomlBooleanValueNode.java` | Not Started |
| `semantic/ast/toml_array_value_node.go` | `TomlArrayValueNode.java` | Not Started |
| `semantic/ast/toml_inline_table_value_node.go` | `TomlInlineTableValueNode.java` | Not Started |
| `semantic/ast/toml_key_node.go` | `TomlKeyNode.java` | Not Started |
| `semantic/ast/toml_unquoted_key_node.go` | `TomlUnquotedKeyNode.java` | Not Started |
| `semantic/ast/toml_key_entry_node.go` | `TomlKeyEntryNode.java` | Not Started |
| `semantic/ast/toml_key_value_node.go` | `TomlKeyValueNode.java` | Not Started |
| `semantic/ast/toml_table_node.go` | `TomlTableNode.java` | Not Started |
| `semantic/ast/toml_table_array_node.go` | `TomlTableArrayNode.java` | Not Started |
| `semantic/ast/toml_node_visitor.go` | `TomlNodeVisitor.java` | Not Started |
| `semantic/diagnostics/toml_diagnostic.go` | `TomlDiagnostic.java` | Not Started |
| `semantic/diagnostics/toml_node_location.go` | `TomlNodeLocation.java` | Not Started |
| `semantic/diagnostics/diagnostic_log.go` | `DiagnosticLog.java` | Not Started |
| `semantic/diagnostics/diagnostic_comparator.go` | `DiagnosticComparator.java` | Not Started |

### Dependencies
- Phase 1: `TomlType`
- Phase 5: Public Syntax Tree

---

## Phase 7: Transformer (Priority: High)

Converts syntax tree to semantic AST.

### Files to Create

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `semantic/ast/toml_transformer.go` | `TomlTransformer.java` | Not Started |
| `semantic/ast/toml_table_value_factory.go` | Extracted from TomlTransformer | Not Started |

### Dependencies
- Phase 5: Public Syntax Tree
- Phase 6: Semantic AST

---

## Phase 8: Public API (Priority: High)

Main entry point for users.

### Files to Create

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `api/toml.go` | `Toml.java` | Not Started |

### API to Implement
```go
// Parsing
func Read(path string) (*Toml, error)
func ReadString(content, fileName string) (*Toml, error)
func ReadWithSchema(path string, schema *Schema) (*Toml, error)

// Querying
func (t *Toml) Get(dottedKey string) (TomlValueNode, bool)
func (t *Toml) GetTable(dottedKey string) (*Toml, bool)
func (t *Toml) GetTables(dottedKey string) []*Toml

// Conversion
func (t *Toml) ToMap() map[string]interface{}
func (t *Toml) Unmarshal(v interface{}) error

// Diagnostics
func (t *Toml) Diagnostics() []Diagnostic
func (t *Toml) RootNode() *TomlTableNode
```

### Dependencies
- All previous phases

---

## Phase 9: Validator (Priority: Medium)

JSON Schema validation for TOML files.

### Files to Create

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `validator/toml_validator.go` | `TomlValidator.java` | Not Started |
| `validator/schema_validator.go` | `SchemaValidator.java` | Not Started |
| `validator/boilerplate_generator.go` | `BoilerplateGenerator.java` | Not Started |
| `validator/sample_node_generator.go` | `SampleNodeGenerator.java` | Not Started |
| `validator/validation_util.go` | `ValidationUtil.java` | Not Started |
| `validator/schema/schema.go` | `Schema.java` | Not Started |
| `validator/schema/abstract_schema.go` | `AbstractSchema.java` | Not Started |
| `validator/schema/string_schema.go` | `StringSchema.java` | Not Started |
| `validator/schema/numeric_schema.go` | `NumericSchema.java` | Not Started |
| `validator/schema/boolean_schema.go` | `BooleanSchema.java` | Not Started |
| `validator/schema/array_schema.go` | `ArraySchema.java` | Not Started |
| `validator/schema/object_schema.go` | `ObjectSchema.java` | Not Started |
| `validator/schema/composition_schema.go` | `CompositionSchema.java` | Not Started |
| `validator/schema/primitive_value_schema.go` | `PrimitiveValueSchema.java` | Not Started |
| `validator/schema/schema_deserializer.go` | `SchemaDeserializer.java` | Not Started |
| `validator/schema/schema_visitor.go` | `SchemaVisitor.java` | Not Started |
| `validator/schema/type.go` | `Type.java` | Not Started |

### Dependencies
- Phase 6: Semantic AST
- Phase 8: Public API

---

## Phase 10: Tests (Priority: High)

Migrate all Java test files to Go.

**See [TEST_MIGRATION.md](TEST_MIGRATION.md) for detailed test file tracking.**

### Summary
- **21 Java test files** to migrate
- **87 test resource files** (47 TOML + 40 JSON) to copy
- Test categories: API, Syntax, Diagnostics, Modifier, Validator

### Dependencies
- All previous phases (1-9) must be complete before tests can run

---

## Milestone Summary

### Milestone 1: Parsing (Phases 1-4)
- Foundation types
- Lexer
- Internal syntax tree
- Parser
- **Deliverable**: Can parse TOML files to syntax tree

### Milestone 2: Public API (Phases 5-8)
- Public syntax tree
- Semantic AST
- Transformer
- Public API
- **Deliverable**: Full TOML parsing and querying API

### Milestone 3: Validation (Phase 9)
- Schema validation
- **Deliverable**: JSON Schema validation support

### Milestone 4: Tests (Phase 10)
- Migrate all 21 Java test files
- Copy 87 test resource files (47 TOML + 40 JSON)
- **Deliverable**: Production-ready library with comprehensive tests

---

## Notes

### Complexity Estimates
- **High**: Lexer, Parser, Transformer (careful translation required)
- **Medium**: Syntax trees, Semantic AST (structural translation)
- **Lower**: Types, Enums, API facade

### Key Differences: Java vs Go
1. No inheritance → Use interfaces and embedding
2. No generics until Go 1.18 → Use interface{} or type assertions
3. No exceptions → Return errors
4. No method overloading → Use different function names or variadic args
5. Package-level visibility → Capitalize for public

### Reference Files
- Java source: `/Users/asmaj/ballerina/ballerina-lang/misc/toml-parser/`
- Go target: `/Users/asmaj/ballerina/ballerina-lang-go/toml-parser/`
