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
	"os"
	"path/filepath"
)

// ScanSourceFiles scans .bal files in the given directory and returns document configurations.
func ScanSourceFiles(dir string, pkgId PackageId) ([]DocumentConfig, error) {
	modId := NewModuleId(pkgId.name, pkgId)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var docs []DocumentConfig
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".bal" {
			continue
		}

		docId := NewDocumentId(entry.Name(), modId)
		filePath := filepath.Join(dir, entry.Name())

		// Create closure that captures filePath
		capturedPath := filePath
		contentFn := func() string {
			data, _ := os.ReadFile(capturedPath)
			return string(data)
		}

		docs = append(docs, NewDocumentConfig(docId, entry.Name(), contentFn))
	}

	return docs, nil
}
