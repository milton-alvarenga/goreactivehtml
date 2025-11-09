package types

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type ClientOutput struct {
	ReqId       uint8
	MsgType     WSTypeOutputMessage
	Destination string
	Data        string
	Header      map[string]string
}

// Marshal serializes the ClientOutput struct into a custom binary format.
/*
Steps for Marshalling to byte:
        Write ReqId (1 byte).
        Write MsgType (1 byte for length, followed by the ascii).
        Write Destination (2 bytes for length, followed by the string).
        Write Data (4 bytes for length, followed by the JSON string).
        Write Header (2 bytes for length, followed by the JSON string).
*/
func (c ClientOutput) Marshal() ([]byte, error) {
	var buf bytes.Buffer

	// ReqId (1 byte) - zero is a subscribed events
	buf.WriteByte(c.ReqId)

	// MsgType (1 byte) - write the length of MsgType followed by the MsgType itself
	buf.WriteByte(byte(c.MsgType)) // Length of MsgType string (1 byte)

	// Destination (2 bytes for length + N bytes for content)
	destLen := uint16(len(c.Destination))
	buf.Write([]byte{byte(destLen >> 8), byte(destLen & 0xFF)}) // 2 bytes length
	buf.WriteString(c.Destination)

	// Data (4 bytes for length + N bytes for JSON serialized content)
	dataJSON, err := json.Marshal(c.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %v", err)
	}
	dataLen := uint32(len(dataJSON))
	buf.Write([]byte{byte(dataLen >> 24), byte(dataLen >> 16), byte(dataLen >> 8), byte(dataLen & 0xFF)}) // 4 bytes length
	buf.Write(dataJSON)

	// Header (3 bytes for length + N bytes for JSON serialized content)
	headerJSON, err := json.Marshal(c.Header)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal header: %v", err)
	}
	headerLen := uint16(len(headerJSON))
	buf.Write([]byte{byte(headerLen >> 8), byte(headerLen & 0xFF)}) // 2 bytes length
	buf.Write(headerJSON)

	return buf.Bytes(), nil
}
