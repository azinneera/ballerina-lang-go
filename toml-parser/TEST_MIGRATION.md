# Test Migration Plan

Detailed tracking of Java test file migration to Go.

**Java Source**: `/Users/asmaj/ballerina/ballerina-lang/misc/toml-parser/src/test/java/toml/parser/test/`
**Go Target**: `/Users/asmaj/ballerina/ballerina-lang-go/toml-parser/test/`

**Legend**: Not Started | In Progress | Complete

---

## Status Overview

| Category | Files | Status |
|----------|-------|--------|
| Test Utilities | 3 | Not Started |
| Core API Tests | 2 | Not Started |
| API Core Tests | 2 | Not Started |
| API Error Tests | 2 | Not Started |
| API Object Tests | 1 | Not Started |
| Syntax Tests | 5 | Not Started |
| Diagnostics Tests | 1 | Not Started |
| Modifier Tests | 1 | Not Started |
| Validator Tests | 4 | Not Started |
| **Total Test Files** | **21** | Not Started |
| TOML Test Resources | 47 | Not Started |
| JSON Test Resources | 40 | Not Started |
| **Total Resources** | **87** | Not Started |

---

## Test Utility Files

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `test/parser_test_constants.go` | `ParserTestConstants.java` | Not Started |
| `test/parser_test_utils.go` | `ParserTestUtils.java` | Not Started |
| `test/api/errors/error_test_utils.go` | `api/errors/ErrorTestUtils.java` | Not Started |

---

## Core API Tests

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `test/test_toml.go` | `TestToml.java` | Not Started |
| `test/test_toml_validator.go` | `TestTomlValidator.java` | Not Started |

---

## API Core Tests

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `test/api/core/key_value_test.go` | `api/core/KeyValueTest.java` | Not Started |
| `test/api/core/table_test.go` | `api/core/TableTest.java` | Not Started |

---

## API Error Tests

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `test/api/errors/key_value_pair_test.go` | `api/errors/KeyValuePairTest.java` | Not Started |
| `test/api/errors/table_test.go` | `api/errors/TableTest.java` | Not Started |

---

## API Object Tests

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `test/api/object/to_object_test.go` | `api/object/ToObjectTest.java` | Not Started |

---

## Syntax Tests

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `test/syntax/abstract_toml_parser_test.go` | `syntax/AbstractTomlParserTest.java` | Not Started |
| `test/syntax/key_value_test.go` | `syntax/KeyValueTest.java` | Not Started |
| `test/syntax/key_value_negative_test.go` | `syntax/KeyValueNegetiveTest.java` | Not Started |
| `test/syntax/table_test.go` | `syntax/TableTest.java` | Not Started |
| `test/syntax/table_negative_test.go` | `syntax/TableNegetiveTest.java` | Not Started |

---

## Diagnostics Tests

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `test/diagnostics/diagnostic_code_test.go` | `diagnostics/DiagnosticCodeTest.java` | Not Started |

---

## Modifier Tests

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `test/modifier/modifier_test.go` | `modifier/ModifierTest.java` | Not Started |

---

## Validator Tests

| Go File | Java Equivalent | Status |
|---------|-----------------|--------|
| `test/validator/boilerplate_generator_test.go` | `validator/BoilerplateGeneratorTest.java` | Not Started |
| `test/validator/custom_error_test.go` | `validator/CustomErrorTest.java` | Not Started |
| `test/validator/schema_test.go` | `validator/SchemaTest.java` | Not Started |
| `test/validator/toml_validate_test.go` | `validator/TomlValidateTest.java` | Not Started |

---

## Test Resources - TOML Files (47 files)

**Java Source**: `/Users/asmaj/ballerina/ballerina-lang/misc/toml-parser/src/test/resources/`
**Go Target**: `/Users/asmaj/ballerina/ballerina-lang-go/toml-parser/testdata/`

### Basic TOML

| Resource | Status |
|----------|--------|
| `basic-toml.toml` | Not Started |

### Modifier Tests

| Resource | Status |
|----------|--------|
| `modifier/Dependencies.toml` | Not Started |

### Object Tests

| Resource | Status |
|----------|--------|
| `object/complex.toml` | Not Started |

### Semantic Tests

| Resource | Status |
|----------|--------|
| `semantic/existing-node.toml` | Not Started |

### Syntax Key-Value Tests (16 files)

| Resource | Status |
|----------|--------|
| `syntax/key-value/array.toml` | Not Started |
| `syntax/key-value/array-missing-comma-negative.toml` | Not Started |
| `syntax/key-value/array-missing-value-negative.toml` | Not Started |
| `syntax/key-value/dotted.toml` | Not Started |
| `syntax/key-value/empty-string-key-sem-negative.toml` | Not Started |
| `syntax/key-value/inline-negative.toml` | Not Started |
| `syntax/key-value/inline-tables.toml` | Not Started |
| `syntax/key-value/key-conflict-sem-negative.toml` | Not Started |
| `syntax/key-value/key-value-multi-negative.toml` | Not Started |
| `syntax/key-value/keys.toml` | Not Started |
| `syntax/key-value/missing-equal-negative.toml` | Not Started |
| `syntax/key-value/missing-key-negative.toml` | Not Started |
| `syntax/key-value/missing-new-line-negative.toml` | Not Started |
| `syntax/key-value/missing-value-negative.toml` | Not Started |
| `syntax/key-value/no-newline-end.toml` | Not Started |
| `syntax/key-value/values.toml` | Not Started |

### Syntax Table Tests (9 files)

| Resource | Status |
|----------|--------|
| `syntax/tables/array-of-tables.toml` | Not Started |
| `syntax/tables/empty-table-close-negative.toml` | Not Started |
| `syntax/tables/empty-table-key-negative.toml` | Not Started |
| `syntax/tables/empty-table-open-negative.toml` | Not Started |
| `syntax/tables/string-missing-close-quotes.toml` | Not Started |
| `syntax/tables/table-key-conflict-sem-negative.toml` | Not Started |
| `syntax/tables/table-key-unordered-conflict-sem-negative.toml` | Not Started |
| `syntax/tables/table.toml` | Not Started |
| `syntax/tables/wrong-closing-brace-negative.toml` | Not Started |

### Validator Basic Tests (12 files)

| Resource | Status |
|----------|--------|
| `validator/sample.toml` | Not Started |
| `validator/basic/Dependencies.toml` | Not Started |
| `validator/basic/additional-field.toml` | Not Started |
| `validator/basic/bal-clean.toml` | Not Started |
| `validator/basic/c2c-clean.toml` | Not Started |
| `validator/basic/c2c-invalid-additional-properties.toml` | Not Started |
| `validator/basic/c2c-invalid-min-max.toml` | Not Started |
| `validator/basic/c2c-invalid-regex.toml` | Not Started |
| `validator/basic/c2c-invalid-type.toml` | Not Started |
| `validator/basic/composition.toml` | Not Started |
| `validator/basic/inline-value.toml` | Not Started |
| `validator/basic/string-length.toml` | Not Started |

### Validator Boilerplate Tests (2 files)

| Resource | Status |
|----------|--------|
| `validator/boilerplate/basic-schema.toml` | Not Started |
| `validator/boilerplate/c2c-schema.toml` | Not Started |

### Validator Custom Error Tests (5 files)

| Resource | Status |
|----------|--------|
| `validator/custom-error/clean.toml` | Not Started |
| `validator/custom-error/regex.toml` | Not Started |
| `validator/custom-error/string-length.toml` | Not Started |
| `validator/custom-error/table.toml` | Not Started |
| `validator/custom-error/type.toml` | Not Started |

---

## Test Resources - JSON Files (40 files)

### Syntax Key-Value Expected (16 files)

| Resource | Status |
|----------|--------|
| `syntax/key-value/array.json` | Not Started |
| `syntax/key-value/array-missing-comma-negative.json` | Not Started |
| `syntax/key-value/array-missing-value-negative.json` | Not Started |
| `syntax/key-value/dotted.json` | Not Started |
| `syntax/key-value/empty-string-key-sem-negative.json` | Not Started |
| `syntax/key-value/inline-negative.json` | Not Started |
| `syntax/key-value/inline-tables.json` | Not Started |
| `syntax/key-value/key-conflict-sem-negative.json` | Not Started |
| `syntax/key-value/key-value-multi-negative.json` | Not Started |
| `syntax/key-value/keys.json` | Not Started |
| `syntax/key-value/missing-equal-negative.json` | Not Started |
| `syntax/key-value/missing-key-negative.json` | Not Started |
| `syntax/key-value/missing-new-line-negative.json` | Not Started |
| `syntax/key-value/missing-value-negative.json` | Not Started |
| `syntax/key-value/no-newline-end.json` | Not Started |
| `syntax/key-value/values.json` | Not Started |

### Syntax Table Expected (9 files)

| Resource | Status |
|----------|--------|
| `syntax/tables/array-of-tables.json` | Not Started |
| `syntax/tables/empty-table-close-negative.json` | Not Started |
| `syntax/tables/empty-table-key-negative.json` | Not Started |
| `syntax/tables/empty-table-open-negative.json` | Not Started |
| `syntax/tables/string-missing-close-quotes.json` | Not Started |
| `syntax/tables/table-key-conflict-sem-negative.json` | Not Started |
| `syntax/tables/table-key-unordered-conflict-sem-negative.json` | Not Started |
| `syntax/tables/table.json` | Not Started |
| `syntax/tables/wrong-closing-brace-negative.json` | Not Started |

### Validator Schemas (15 files)

| Resource | Status |
|----------|--------|
| `validator/sample-schema.json` | Not Started |
| `validator/basic/additional-field.json` | Not Started |
| `validator/basic/c2c-schema.json` | Not Started |
| `validator/basic/composition.json` | Not Started |
| `validator/basic/dep-new.json` | Not Started |
| `validator/basic/dependency-schema.json` | Not Started |
| `validator/basic/inline-value.json` | Not Started |
| `validator/basic/schema.json` | Not Started |
| `validator/boilerplate/basic-schema.json` | Not Started |
| `validator/boilerplate/c2c-schema.json` | Not Started |
| `validator/custom-error/schema.json` | Not Started |
| `validator/schema/additional-properties-boolean.json` | Not Started |
| `validator/schema/min-number.json` | Not Started |
| `validator/schema/pattern-string.json` | Not Started |

---

## Migration Notes

### Go Test Conventions
- Test files must end with `_test.go`
- Test functions must start with `Test`
- Use `testing.T` for test assertions
- Consider using `testify` package for assertions

### Resource Loading
```go
// Go pattern for loading test resources
//go:embed testdata/*
var testdata embed.FS

func loadTestResource(path string) ([]byte, error) {
    return testdata.ReadFile(path)
}
```

### Test Organization
```
test/
├── testdata/           # Test resources (TOML, JSON files)
│   ├── syntax/
│   ├── validator/
│   └── ...
├── api/
│   ├── core/
│   ├── errors/
│   └── object/
├── syntax/
├── diagnostics/
├── modifier/
└── validator/
```
