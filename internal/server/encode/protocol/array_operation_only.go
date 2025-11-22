package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//
// ┌────────┬───────────┬──────────┬──────────┬────────────┐
// │ bit 7  │ bit 6     │ bits 4-5 │ bits 2-3 │ bits 0-1   │
// │ bulk   │ partial   │ dataSize │ posSize  │ operation  │
// └────────┴───────────┴──────────┴──────────┴────────────┘
//
// Operation:
//   00 = DELETE
//   01 = UPDATE
//   11 = INSERT
//   10 = reserved
//

type OperationType byte

const (
	OpDelete OperationType = 0b00
	OpUpdate OperationType = 0b01
	OpInsert OperationType = 0b11
)

type Encoder struct{}

//
// ─────────────────────────────────────────────────────────────
//  INT ENCODING HELPERS
// ─────────────────────────────────────────────────────────────
//

// v encoded using sizeIndicator bytes: 0–3 bytes
func encodeIntWithSize(v uint32, sizeIndicator uint8) ([]byte, error) {
	switch sizeIndicator {
	case 0:
		return []byte{}, nil

	case 1:
		if v > 0xFF {
			return nil, fmt.Errorf("value %d too large for 1 byte", v)
		}
		return []byte{byte(v)}, nil

	case 2:
		if v > 0xFFFF {
			return nil, fmt.Errorf("value %d too large for 2 bytes", v)
		}
		out := make([]byte, 2)
		binary.BigEndian.PutUint16(out, uint16(v))
		return out, nil

	case 3:
		if v > 0xFFFFFF {
			return nil, fmt.Errorf("value %d too large for 3 bytes", v)
		}
		return []byte{
			byte((v >> 16) & 0xFF),
			byte((v >> 8) & 0xFF),
			byte(v & 0xFF),
		}, nil

	default:
		return nil, fmt.Errorf("invalid size indicator %d", sizeIndicator)
	}
}

// Minimal byte size for an integer
func autoSizeIndicator(v uint32) uint8 {
	switch {
	case v <= 0xFF:
		return 1
	case v <= 0xFFFF:
		return 2
	default:
		return 3
	}
}

// Same logic for data length
func autoLengthIndicator(v uint32) uint8 {
	return autoSizeIndicator(v)
}

//
// ─────────────────────────────────────────────────────────────
//  HEADER BUILDER — fixed bit layout
// ─────────────────────────────────────────────────────────────
//

func buildHeader(op OperationType, bulk bool, partial bool, posSize, dataSize uint8) byte {
	var header byte

	// bits 0-1: operation
	header |= byte(op & 0b11)

	// bits 2-3: posSize
	header |= (byte(posSize) & 0b11) << 2

	// bits 4-5: dataSize
	header |= (byte(dataSize) & 0b11) << 4

	// bit 6: partial
	if partial {
		header |= 1 << 6
	}

	// bit 7: bulk
	if bulk {
		header |= 1 << 7
	}

	return header
}

//
// ─────────────────────────────────────────────────────────────
//  PUBLIC API
// ─────────────────────────────────────────────────────────────
//

// ─── Single DELETE ───────────────────────────────────────────
func (e *Encoder) EncodeDelete(pos uint32) ([]byte, error) {

	posSize := autoSizeIndicator(pos)
	dataSize := uint8(0)

	header := buildHeader(OpDelete, false, false, posSize, dataSize)

	buf := bytes.Buffer{}
	buf.WriteByte(header)

	encodedPos, err := encodeIntWithSize(pos, posSize)
	if err != nil {
		return nil, err
	}
	buf.Write(encodedPos)

	return buf.Bytes(), nil
}

// ─── Bulk DELETE (dense) ─────────────────────────────────────
func (e *Encoder) EncodeDeleteRange(start, end uint32) ([]byte, error) {

	posSize := autoSizeIndicator(max(start, end))
	dataSize := uint8(0)

	header := buildHeader(OpDelete, true, false, posSize, dataSize)

	buf := bytes.Buffer{}
	buf.WriteByte(header)

	encA, _ := encodeIntWithSize(start, posSize)
	encB, _ := encodeIntWithSize(end, posSize)
	buf.Write(encA)
	buf.Write(encB)

	return buf.Bytes(), nil
}

// ─── Single INSERT (full) ────────────────────────────────────
func (e *Encoder) EncodeInsert(pos uint32, data []byte) ([]byte, error) {

	posSize := autoSizeIndicator(pos)
	dataSize := autoLengthIndicator(uint32(len(data)))

	header := buildHeader(OpInsert, false, false, posSize, dataSize)

	buf := bytes.Buffer{}
	buf.WriteByte(header)

	encPos, err := encodeIntWithSize(pos, posSize)
	if err != nil {
		return nil, err
	}
	buf.Write(encPos)

	encLen, err := encodeIntWithSize(uint32(len(data)), dataSize)
	if err != nil {
		return nil, err
	}
	buf.Write(encLen)

	buf.Write(data)

	return buf.Bytes(), nil
}

// ─── Single UPDATE (full replace) ─────────────────────────────
func (e *Encoder) EncodeUpdate(pos uint32, data []byte) ([]byte, error) {

	posSize := autoSizeIndicator(pos)
	dataSize := autoLengthIndicator(uint32(len(data)))

	header := buildHeader(OpUpdate, false, false, posSize, dataSize)

	buf := bytes.Buffer{}
	buf.WriteByte(header)

	encPos, err := encodeIntWithSize(pos, posSize)
	if err != nil {
		return nil, err
	}
	buf.Write(encPos)

	encLen, err := encodeIntWithSize(uint32(len(data)), dataSize)
	if err != nil {
		return nil, err
	}
	buf.Write(encLen)

	buf.Write(data)

	return buf.Bytes(), nil
}

// ─── Bulk INSERT (dense) ─────────────────────────────────────
func (e *Encoder) EncodeInsertRange(start, end uint32, payloads [][]byte) ([]byte, error) {

	if int(end-start)+1 != len(payloads) {
		return nil, fmt.Errorf("payload count must match range size")
	}

	maxPos := max(start, end)
	posSize := autoSizeIndicator(maxPos)

	var maxLen uint32
	for _, p := range payloads {
		if uint32(len(p)) > maxLen {
			maxLen = uint32(len(p))
		}
	}
	dataSize := autoLengthIndicator(maxLen)

	header := buildHeader(OpInsert, true, false, posSize, dataSize)

	buf := bytes.Buffer{}
	buf.WriteByte(header)

	encA, _ := encodeIntWithSize(start, posSize)
	encB, _ := encodeIntWithSize(end, posSize)
	buf.Write(encA)
	buf.Write(encB)

	for _, p := range payloads {
		plen := uint32(len(p))
		encLen, err := encodeIntWithSize(plen, dataSize)
		if err != nil {
			return nil, err
		}
		buf.Write(encLen)
		buf.Write(p)
	}

	return buf.Bytes(), nil
}

// ─── Bulk UPDATE (dense) ─────────────────────────────────────
func (e *Encoder) EncodeUpdateRange(start, end uint32, payloads [][]byte) ([]byte, error) {

	if int(end-start)+1 != len(payloads) {
		return nil, fmt.Errorf("payload count must match range size")
	}

	maxPos := max(start, end)
	posSize := autoSizeIndicator(maxPos)

	var maxLen uint32
	for _, p := range payloads {
		if uint32(len(p)) > maxLen {
			maxLen = uint32(len(p))
		}
	}
	dataSize := autoLengthIndicator(maxLen)

	header := buildHeader(OpUpdate, true, false, posSize, dataSize)

	buf := bytes.Buffer{}
	buf.WriteByte(header)

	encA, _ := encodeIntWithSize(start, posSize)
	encB, _ := encodeIntWithSize(end, posSize)
	buf.Write(encA)
	buf.Write(encB)

	for _, p := range payloads {
		plen := uint32(len(p))
		encLen, err := encodeIntWithSize(plen, dataSize)
		if err != nil {
			return nil, err
		}
		buf.Write(encLen)
		buf.Write(p)
	}

	return buf.Bytes(), nil
}

//
// ─────────────────────────────────────────────────────────────
//   PARTIAL UPDATE SUPPORT (NEW)
// ─────────────────────────────────────────────────────────────
//

// Sparse partial update
func (e *Encoder) EncodePartialUpdate(pos uint32, patch []byte) ([]byte, error) {

	posSize := autoSizeIndicator(pos)
	dataSize := autoLengthIndicator(uint32(len(patch)))

	header := buildHeader(OpUpdate, false, true, posSize, dataSize)

	buf := bytes.Buffer{}
	buf.WriteByte(header)

	encPos, err := encodeIntWithSize(pos, posSize)
	if err != nil {
		return nil, err
	}
	buf.Write(encPos)

	encLen, err := encodeIntWithSize(uint32(len(patch)), dataSize)
	if err != nil {
		return nil, err
	}
	buf.Write(encLen)

	buf.Write(patch)

	return buf.Bytes(), nil
}

// Patch entry for bulk sparse updates
type PartialPatch struct {
	Pos  uint32
	Data []byte
}

// Bulk sparse partial update
func (e *Encoder) EncodePartialUpdateRange(start, end uint32, patches []PartialPatch) ([]byte, error) {

	// posSize must support both range & patch positions
	maxPos := max(start, end)
	for _, p := range patches {
		if p.Pos > maxPos {
			maxPos = p.Pos
		}
	}
	posSize := autoSizeIndicator(maxPos)

	// data size is max of all patches
	var maxLen uint32
	for _, p := range patches {
		if uint32(len(p.Data)) > maxLen {
			maxLen = uint32(len(p.Data))
		}
	}
	dataSize := autoLengthIndicator(maxLen)

	header := buildHeader(OpUpdate, true, true, posSize, dataSize)

	buf := bytes.Buffer{}
	buf.WriteByte(header)

	// bulk range
	encA, _ := encodeIntWithSize(start, posSize)
	encB, _ := encodeIntWithSize(end, posSize)
	buf.Write(encA)
	buf.Write(encB)

	// sparse patches
	for _, p := range patches {
		encPos, err := encodeIntWithSize(p.Pos, posSize)
		if err != nil {
			return nil, err
		}
		buf.Write(encPos)

		plen := uint32(len(p.Data))
		encLen, err := encodeIntWithSize(plen, dataSize)
		if err != nil {
			return nil, err
		}
		buf.Write(encLen)

		buf.Write(p.Data)
	}

	return buf.Bytes(), nil
}

// Utility
func max(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}
