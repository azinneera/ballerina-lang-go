# TOML Parser Architecture

This document describes the architecture of the Java TOML parser being migrated to Go.

## Overview

The Ballerina TOML parser is a production-grade, multi-layer implementation with:
- LL(k) lexer with mode-based state machine
- Recursive-descent parser with error recovery
- Two-tree architecture (Syntax Tree + Semantic AST)
- Schema validation via JSON Schema
- Rich diagnostics and error reporting

## Source Location

- **Java Source**: `/Users/asmaj/ballerina/ballerina-lang/misc/toml-parser/src/main/java/io/ballerina/toml/`
- **Go Target**: `/Users/asmaj/ballerina/ballerina-lang-go/toml-parser/`

## Package Structure

### Java Packages (11 core packages, ~155 files)

```
io.ballerina.toml/
├── api/                          # Public API (1 file)
│   └── Toml.java                 # Main entry point
├── internal/
│   ├── diagnostics/              # Internal diagnostics (3 files)
│   │   ├── DiagnosticErrorCode.java
│   │   ├── DiagnosticMessageHelper.java
│   │   └── SyntaxDiagnostic.java
│   ├── parser/                   # Lexer & Parser (14 files)
│   │   ├── AbstractLexer.java
│   │   ├── AbstractParser.java
│   │   ├── TomlLexer.java
│   │   ├── TomlParser.java
│   │   ├── TomlParserErrorHandler.java
│   │   ├── ParserMode.java
│   │   ├── LexerTerminals.java
│   │   └── ...
│   ├── parser/tree/              # Internal Syntax Tree (34 files)
│   │   ├── STNode.java
│   │   ├── STToken.java
│   │   ├── STNodeFactory.java
│   │   └── [ST* node implementations]
│   └── syntax/                   # Syntax utilities (5 files)
├── semantic/
│   ├── TomlType.java             # Type enum
│   ├── ast/                      # Semantic AST (22 files)
│   │   ├── TomlNode.java
│   │   ├── TomlValueNode.java
│   │   ├── TomlTableNode.java
│   │   ├── TomlTransformer.java
│   │   └── [value/structure nodes]
│   └── diagnostics/              # Semantic diagnostics (4 files)
├── syntax/tree/                  # Public Syntax Tree API (35 files)
│   ├── Node.java
│   ├── SyntaxKind.java
│   ├── SyntaxTree.java
│   └── [public node implementations]
└── validator/                    # Schema validation (18 files)
    ├── TomlValidator.java
    ├── SchemaValidator.java
    └── schema/
        ├── Schema.java
        └── [schema type implementations]
```

## Parsing Pipeline

```
Input String
    │
    ▼
┌─────────────────────────────────────────────────────────┐
│ [1] LEXER (TomlLexer)                                   │
│     - Character-by-character scanning                    │
│     - Mode-based state machine                           │
│     - Produces tokens with leading/trailing minutiae     │
└─────────────────────────────────────────────────────────┘
    │
    ▼
Token Stream (with trivia/whitespace)
    │
    ▼
┌─────────────────────────────────────────────────────────┐
│ [2] PARSER (TomlParser)                                 │
│     - Recursive descent LL(k) parser                     │
│     - Error recovery via TomlParserErrorHandler          │
│     - Produces internal syntax tree (STNode)             │
└─────────────────────────────────────────────────────────┘
    │
    ▼
Internal Syntax Tree (STNode hierarchy)
    │
    ▼
┌─────────────────────────────────────────────────────────┐
│ [3] PUBLIC FACADE (Node classes)                        │
│     - Wraps STNode with user-friendly API               │
│     - Lazy evaluation of line/column ranges             │
│     - Parent/child navigation                            │
└─────────────────────────────────────────────────────────┘
    │
    ▼
Public Syntax Tree (Node hierarchy)
    │
    ▼
┌─────────────────────────────────────────────────────────┐
│ [4] SEMANTIC TRANSFORMER (TomlTransformer)              │
│     - Converts syntax tree to semantic AST              │
│     - Flattens dotted keys to nested tables             │
│     - Validates uniqueness and structure                │
│     - Produces typed value nodes                        │
└─────────────────────────────────────────────────────────┘
    │
    ▼
Semantic AST (TomlNode hierarchy)
    │
    ▼
┌─────────────────────────────────────────────────────────┐
│ [5] PUBLIC API (Toml class)                             │
│     - Query: get(), getTable(), getTables()             │
│     - Conversion: toMap(), to(Class<T>)                 │
│     - Diagnostics: diagnostics()                        │
└─────────────────────────────────────────────────────────┘
```

## Key Components

### 1. Lexer (TomlLexer)

**Responsibility**: Convert input characters to tokens

**Key Features**:
- LL(k) lookahead for token disambiguation
- Mode-based state machine for different contexts:
  - `DEFAULT` - Normal token recognition
  - `STRING` - Inside double-quoted strings
  - `MULTILINE_STRING` - Inside triple-quoted strings
  - `LITERAL_STRING` - Inside single-quoted strings
  - `MULTILINE_LITERAL_STRING` - Inside triple single-quoted strings
  - `NEW_LINE` - Newline handling

**Token Types** (defined in SyntaxKind):
- Literals: `STRING_LITERAL`, `ML_STRING_LITERAL`, `LITERAL_STRING`, `ML_LITERAL_STRING`
- Numbers: `DEC_INT`, `HEX_INT`, `OCT_INT`, `BIN_INT`, `FLOAT`
- Booleans: `TRUE`, `FALSE`
- Structural: `OPEN_BRACKET`, `CLOSE_BRACKET`, `OPEN_BRACE`, `CLOSE_BRACE`
- Operators: `EQUALS`, `DOT`, `COMMA`
- Keys: `UNQUOTED_KEY`
- Whitespace: `WHITESPACE`, `NEWLINE`, `COMMENT`

### 2. Parser (TomlParser)

**Responsibility**: Build syntax tree from token stream

**Parsing Rules**:
```
Document      → TopLevelNode*
TopLevelNode  → Table | TableArray | KeyValue
Table         → '[' Keys ']' KeyValue*
TableArray    → '[[' Keys ']]' KeyValue*
KeyValue      → Keys '=' Value
Keys          → Key ('.' Key)*
Key           → UNQUOTED_KEY | STRING_LITERAL | LITERAL_STRING
Value         → String | Number | Boolean | Array | InlineTable | DateTime
Array         → '[' (Value (',' Value)*)? ']'
InlineTable   → '{' (KeyValue (',' KeyValue)*)? '}'
```

**Error Recovery**:
- `TomlParserErrorHandler` provides graceful recovery
- Inserts missing tokens
- Skips unexpected tokens
- Maintains context for better error messages

### 3. Internal Syntax Tree (STNode)

**Responsibility**: Preserve all source information

**Base Classes**:
- `STNode` - Abstract base for all nodes
- `STToken` - Lexical tokens with position info
- `STMinutiae` - Whitespace and comments (trivia)

**Node Types**:
- `STDocumentNode` - Root document
- `STTableNode` - `[table]` sections
- `STTableArrayNode` - `[[array.of.tables]]`
- `STKeyValueNode` - Key-value pairs
- `STArrayNode` - Array values
- `STInlineTableNode` - Inline table values
- `STBasicLiteralNode` - Primitive values

**Factory**: `STNodeFactory` creates all node instances

### 4. Public Syntax Tree (Node)

**Responsibility**: User-friendly API over internal tree

**Base Classes**:
- `Node` - Abstract base with parent navigation
- `NonTerminalNode` - Nodes with children
- `Token` - Facade over STToken

**Key API**:
- `parent()` - Navigate to parent node
- `ancestors()` - Navigate to all ancestors
- `children()` - Get child nodes
- `lineRange()` - Get source location
- `toSourceCode()` - Reconstruct source text

### 5. Semantic AST (TomlNode)

**Responsibility**: Typed, navigable representation for querying

**Base Classes**:
- `TomlNode` - Abstract base with location info
- `TopLevelNode` - Entries that can appear at root
- `TomlValueNode` - All value types

**Node Hierarchy**:
```
TomlNode
├── TopLevelNode
│   ├── TomlKeyValueNode
│   ├── TomlTableNode
│   └── TomlTableArrayNode
└── TomlValueNode
    ├── TomlBasicValueNode<T>
    │   ├── TomlStringValueNode
    │   ├── TomlLongValueNode
    │   ├── TomlDoubleValueNodeNode
    │   └── TomlBooleanValueNode
    ├── TomlArrayValueNode
    └── TomlInlineTableValueNode
```

**Transformer**: `TomlTransformer` converts syntax tree to semantic AST:
- Resolves dotted keys to nested tables
- Creates implicit tables
- Validates key uniqueness
- Parses literal values

### 6. Type System (TomlType)

**Enum Values**:
```java
TABLE, TABLE_ARRAY, KEY_VALUE, ARRAY,
STRING, INTEGER, DOUBLE, BOOLEAN,
INLINE_TABLE, KEY, KEY_ENTRY
```

### 7. Diagnostics

**Two Levels**:

1. **Syntax Diagnostics** (SyntaxDiagnostic):
   - Lexer errors (invalid escape, leading zeros)
   - Parser errors (missing tokens, unexpected input)

2. **Semantic Diagnostics** (TomlDiagnostic):
   - Duplicate keys
   - Type conflicts
   - Structure errors

**Error Codes** (DiagnosticErrorCode):
- `ERROR_MISSING_*` - Missing tokens
- `ERROR_INVALID_*` - Invalid syntax
- `ERROR_EXISTING_NODE` - Duplicate definitions

### 8. Validator (TomlValidator)

**Responsibility**: Validate TOML against JSON Schema

**Components**:
- `SchemaValidator` - Main validation logic
- `Schema` hierarchy - Type-specific schemas
- `SchemaDeserializer` - Parse JSON Schema

**Schema Types**:
- `StringSchema`, `NumericSchema`, `BooleanSchema`
- `ArraySchema`, `ObjectSchema`
- `CompositionSchema` (anyOf, oneOf, allOf)

## Design Patterns

| Pattern | Usage |
|---------|-------|
| Factory | `STNodeFactory`, `NodeFactory` for creating nodes |
| Visitor | `NodeVisitor`, `TomlNodeVisitor` for traversal |
| Transformer | `NodeTransformer`, `TomlTransformer` for conversion |
| Facade | Public `Node` wraps internal `STNode` |
| State Machine | Lexer modes for context-sensitive scanning |

## Public API (Toml class)

```java
// Parsing
Toml.read(Path path)
Toml.read(String content, String fileName)
Toml.read(Path path, Schema schema)

// Querying
<T extends TomlValueNode> Optional<T> get(String dottedKey)
Optional<Toml> getTable(String dottedKey)
List<Toml> getTables(String dottedKey)

// Conversion
Map<String, Object> toMap()
<T> T to(Class<T> targetClass)

// Diagnostics
List<Diagnostic> diagnostics()
TomlTableNode rootNode()
```

## Go Migration Guidelines

### File Naming Convention

| Java Class | Go File |
|------------|---------|
| `TomlLexer.java` | `toml_lexer.go` |
| `STNodeFactory.java` | `st_node_factory.go` |
| `TomlTableNode.java` | `toml_table_node.go` |

### Package Structure

```
toml-parser/
├── api/                    → io.ballerina.toml.api
├── internal/
│   ├── diagnostics/        → io.ballerina.toml.internal.diagnostics
│   ├── parser/             → io.ballerina.toml.internal.parser
│   ├── parser/tree/        → io.ballerina.toml.internal.parser.tree
│   └── syntax/             → io.ballerina.toml.internal.syntax
├── semantic/
│   ├── ast/                → io.ballerina.toml.semantic.ast
│   └── diagnostics/        → io.ballerina.toml.semantic.diagnostics
├── syntax/tree/            → io.ballerina.toml.syntax.tree
└── validator/
    └── schema/             → io.ballerina.toml.validator.schema
```

### Go Idioms

1. **Interfaces over abstract classes**: Use interfaces for `Node`, `STNode`
2. **Embedding over inheritance**: Embed base structs in node types
3. **Options pattern**: Use functional options for factory methods
4. **Error returns**: Return `(result, error)` instead of exceptions

### Comments

Each Go file should include a comment referencing the Java equivalent:
```go
// Package parser implements the TOML lexer and parser.
// Java equivalent: io.ballerina.toml.internal.parser
```
