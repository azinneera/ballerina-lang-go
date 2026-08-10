// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
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

package splice

import (
	"bytes"
	"debug/pe"
	"encoding/binary"
	"fmt"
)

// PESectionName is the PE section the payload is stored under.
// Duplicated as windowsSectionName in
// cli/internal/executable/embed_windows.go. 7 chars + NUL padding fits
// the fixed 8-byte Name field, no COFF long-name indirection needed.
const PESectionName = ".balexe"

const (
	peDosELfanewOffset  = 0x3C // offset of the 4-byte LE PE header pointer
	peSignatureSize     = 4    // "PE\x00\x00"
	peCoffHeaderSize    = 20   // pe.FileHeader's on-disk size
	peSectionHeaderSize = 40   // pe.SectionHeader32's on-disk size
)

// splicePESection returns a copy of stub with a new section named
// sectionName appended, holding payload verbatim. stub must be 64-bit PE
// with header slack for one more section entry (Go's linker always
// reserves room for 16, regardless of actual count).
//
// Unlike ELF, the PE loader maps every section unconditionally, so the
// new one is non-executable, non-writable, and unreferenced by any Data
// Directory — mapped but inert. SizeOfImage is grown to cover its VA
// range, which the loader requires.
func splicePESection(stub []byte, sectionName string, payload []byte) ([]byte, error) {
	f, err := pe.NewFile(bytes.NewReader(stub))
	if err != nil {
		return nil, fmt.Errorf("parsing PE stub: %w", err)
	}
	defer func() { _ = f.Close() }()

	opt, ok := f.OptionalHeader.(*pe.OptionalHeader64)
	if !ok {
		return nil, fmt.Errorf("unsupported PE optional header %T: only 64-bit PE stubs are supported", f.OptionalHeader)
	}
	if int(f.SizeOfOptionalHeader) != binary.Size(pe.OptionalHeader64{}) {
		return nil, fmt.Errorf("unsupported PE optional header size %d (expected %d)",
			f.SizeOfOptionalHeader, binary.Size(pe.OptionalHeader64{}))
	}
	for _, s := range f.Sections {
		if s.Name == sectionName {
			return nil, fmt.Errorf("stub already has a %q section", sectionName)
		}
	}
	if opt.FileAlignment == 0 || opt.SectionAlignment == 0 {
		return nil, fmt.Errorf("invalid PE alignment: FileAlignment=%d SectionAlignment=%d", opt.FileAlignment, opt.SectionAlignment)
	}

	if len(stub) < peDosELfanewOffset+4 {
		return nil, fmt.Errorf("stub too small to contain a DOS header")
	}
	elfanew := int64(binary.LittleEndian.Uint32(stub[peDosELfanewOffset : peDosELfanewOffset+4]))
	coffHeaderOff := elfanew + peSignatureSize
	optHeaderOff := coffHeaderOff + peCoffHeaderSize
	sectionTableOff := optHeaderOff + int64(f.SizeOfOptionalHeader)

	oldNumSections := int64(f.NumberOfSections)
	if coffHeaderOff < 0 || sectionTableOff+oldNumSections*peSectionHeaderSize > int64(len(stub)) {
		return nil, fmt.Errorf("PE section header table out of range")
	}

	newTableEnd := sectionTableOff + (oldNumSections+1)*peSectionHeaderSize
	if newTableEnd > int64(opt.SizeOfHeaders) {
		return nil, fmt.Errorf("insufficient PE header slack to add a new section header: need %d bytes, have %d",
			newTableEnd-sectionTableOff, int64(opt.SizeOfHeaders)-sectionTableOff)
	}

	// New-section geometry is the max over every existing section, not
	// just the table's last entry — section order isn't guaranteed to
	// match address/offset order.
	var maxVAEnd, maxFileEnd uint32
	for _, s := range f.Sections {
		if vaEnd := s.VirtualAddress + s.VirtualSize; vaEnd > maxVAEnd {
			maxVAEnd = vaEnd
		}
		if fileEnd := s.Offset + s.Size; fileEnd > maxFileEnd {
			maxFileEnd = fileEnd
		}
	}
	if uint64(maxFileEnd) != uint64(len(stub)) {
		return nil, fmt.Errorf("stub has %d unexpected trailing bytes past its last section (expected exactly %d bytes)",
			int64(len(stub))-int64(maxFileEnd), maxFileEnd)
	}

	rawDataOff := alignUp32(maxFileEnd, opt.FileAlignment)
	rawDataSize := alignUp32(uint32(len(payload)), opt.FileAlignment)
	newVA := alignUp32(maxVAEnd, opt.SectionAlignment)
	newVASize := alignUp32(uint32(len(payload)), opt.SectionAlignment)

	out := append([]byte(nil), stub...)
	if pad := int64(rawDataOff) - int64(len(out)); pad > 0 {
		out = append(out, make([]byte, pad)...)
	}
	out = append(out, payload...)
	if pad := int64(rawDataOff) + int64(rawDataSize) - int64(len(out)); pad > 0 {
		out = append(out, make([]byte, pad)...)
	}

	var nameBytes [8]byte
	copy(nameBytes[:], sectionName)
	newHeader := pe.SectionHeader32{
		Name:             nameBytes,
		VirtualSize:      uint32(len(payload)),
		VirtualAddress:   newVA,
		SizeOfRawData:    rawDataSize,
		PointerToRawData: rawDataOff,
		Characteristics:  uint32(pe.IMAGE_SCN_CNT_INITIALIZED_DATA | pe.IMAGE_SCN_MEM_READ),
	}
	var hdrBuf bytes.Buffer
	if err := binary.Write(&hdrBuf, binary.LittleEndian, newHeader); err != nil {
		return nil, fmt.Errorf("encoding new section header: %w", err)
	}
	newEntryOff := sectionTableOff + oldNumSections*peSectionHeaderSize
	copy(out[newEntryOff:newEntryOff+peSectionHeaderSize], hdrBuf.Bytes())

	// Patch NumberOfSections via a read-modify-write of the whole COFF
	// header struct, rather than hand-computed byte offsets.
	var coff pe.FileHeader
	if err := binary.Read(bytes.NewReader(stub[coffHeaderOff:]), binary.LittleEndian, &coff); err != nil {
		return nil, fmt.Errorf("reading COFF header: %w", err)
	}
	coff.NumberOfSections++
	var coffBuf bytes.Buffer
	if err := binary.Write(&coffBuf, binary.LittleEndian, coff); err != nil {
		return nil, fmt.Errorf("encoding COFF header: %w", err)
	}
	copy(out[coffHeaderOff:coffHeaderOff+peCoffHeaderSize], coffBuf.Bytes())

	// Patch SizeOfImage the same way: read-modify-write the whole
	// optional header struct. Every other field (SizeOfHeaders, CheckSum,
	// all 16 Data Directory entries) is carried through unchanged.
	newOpt := *opt
	newOpt.SizeOfImage = newVA + newVASize
	var optBuf bytes.Buffer
	if err := binary.Write(&optBuf, binary.LittleEndian, newOpt); err != nil {
		return nil, fmt.Errorf("encoding optional header: %w", err)
	}
	copy(out[optHeaderOff:optHeaderOff+int64(f.SizeOfOptionalHeader)], optBuf.Bytes())

	return out, nil
}

func alignUp32(v, align uint32) uint32 {
	if align == 0 {
		return v
	}
	return (v + align - 1) / align * align
}
