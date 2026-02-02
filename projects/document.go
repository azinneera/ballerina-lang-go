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

package projects

import (
	"ballerina-lang-go/parser"
	"ballerina-lang-go/parser/tree"
	"ballerina-lang-go/tools/text"
)

// Document represents a single .bal source file.
type Document struct {
	module     *Module
	id         DocumentId
	name       string
	filePath   string // full path to the file (for parsing)
	contentFn  func() string
	syntaxTree *tree.SyntaxTree
	content    string
}

// NewDocumentFromConfig creates a Document from a DocumentConfig.
func NewDocumentFromConfig(mod *Module, cfg DocumentConfig) *Document {
	return &Document{
		module:    mod,
		id:        cfg.DocumentId(),
		name:      cfg.Name(),
		contentFn: cfg.contentFn,
	}
}

// DocumentId returns the document's unique identifier.
func (d *Document) DocumentId() DocumentId {
	return d.id
}

// Name returns the document file name.
func (d *Document) Name() string {
	return d.name
}

// Module returns the parent module.
func (d *Document) Module() *Module {
	return d.module
}

// Content returns the document content (lazy loaded).
func (d *Document) Content() string {
	if d.content == "" && d.contentFn != nil {
		d.content = d.contentFn()
	}
	return d.content
}

// SyntaxTree returns the parsed syntax tree (lazy loaded).
func (d *Document) SyntaxTree() (*tree.SyntaxTree, error) {
	if d.syntaxTree == nil {
		content := d.Content()
		syntaxTree, err := parseContent(content, d.name)
		if err != nil {
			return nil, err
		}
		d.syntaxTree = syntaxTree
	}
	return d.syntaxTree, nil
}

// parseContent parses Ballerina source content and returns a syntax tree.
func parseContent(content, fileName string) (*tree.SyntaxTree, error) {
	// Create CharReader from content
	reader := text.CharReaderFromText(content)

	// Create Lexer
	lexer := parser.NewLexer(reader, nil)

	// Create TokenReader from Lexer
	tokenReader := parser.CreateTokenReader(*lexer, nil)

	// Create Parser from TokenReader
	ballerinaParser := parser.NewBallerinaParserFromTokenReader(tokenReader, nil)

	// Parse the content
	rootNode := ballerinaParser.Parse().(*tree.STModulePart)

	// Create syntax tree
	moduleNode := tree.CreateUnlinkedFacade[*tree.STModulePart, *tree.ModulePart](rootNode)
	syntaxTree := tree.NewSyntaxTreeFromNodeTextDocumentStringBool(moduleNode, nil, fileName, false)
	return &syntaxTree, nil
}
