package rpc

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/milton-alvarenga/goreactivehtml/internal/server/types"
)

type ClientInputRPC struct {
	ReqId  *uint8                     // Unique request ID
	Class  string                     // What bff to process and receive the requests. 2 bytes for lenght, followed by the string
	Method string                     // RPC method name (e.g., "getUser", "updateData"). 1 byte for length, followed by the ascii
	Params map[string]interface{}     // Parameters for the RPC call. 4 bytes for lenght, followed by the JSON string
	Header map[string]string          // Optional headers (for metadata or authentication). 2 bytes for lenght, followed by the JSON string
	WSConn *types.WebSocketConnection // WebSocket connection for communication
}

func (c ClientInputRPC) SendToClient(ClientOutput types.ClientOutput) bool {
	message, err := ClientOutput.Marshal()
	if err != nil {
		log.Println("Could not marshal Client Output to send to client")
		return false
	}
	err = c.WSConn.Write(websocket.TextMessage, message)
	if err != nil {
		if websocket.IsCloseError(err) {
			log.Println("Client closed the connection")
		} else {
			log.Println("Error reading message:", err)
		}
		return false
	}
	return true
}

func (c ClientInputRPC) IsValidMessage() error {

	if !c.IsValidExecutor() {
		return errors.New("invalid endpoint")
	}

	if !c.IsValidOperation(c.Method) {
		return errors.New("invalid operation")
	}

	return nil
}

func (c *ClientInputRPC) Close() error {
	return c.WSConn.Conn.Close()
}

func (c ClientInputRPC) IsValidExecutor() bool {
}

func (c ClientInputRPC) IsValidOperation(operation string) bool {

}

func (c *ClientInputRPC) Unmarshal(message []byte) error {
	offset := 0
	if len(message) < 1+1+2+1+4+2 {
		return errors.New("message too short")
	}

	// --- 1. Skip conn_type (1 byte)
	offset++

	// --- 2. reqId (1 byte)
	reqId := message[offset]
	c.ReqId = &reqId
	offset++

	// --- 3. bff length (2 bytes, big endian)
	if offset+2 > len(message) {
		return errors.New("missing bff length")
	}
	bffLen := int(binary.BigEndian.Uint16(message[offset : offset+2]))
	offset += 2

	if offset+bffLen > len(message) {
		return errors.New("invalid bff length")
	}
	c.Class = string(message[offset : offset+bffLen])
	offset += bffLen

	// --- 4. method length (1 byte)
	if offset+1 > len(message) {
		return errors.New("missing method length")
	}
	methodLen := int(message[offset])
	offset++

	if offset+methodLen > len(message) {
		return errors.New("invalid method length")
	}
	c.Method = string(message[offset : offset+methodLen])
	offset += methodLen

	// --- 5. params length (4 bytes, big endian)
	if offset+4 > len(message) {
		return errors.New("missing params length")
	}
	paramsLen := int(binary.BigEndian.Uint32(message[offset : offset+4]))
	offset += 4

	if offset+paramsLen > len(message) {
		return errors.New("invalid params length")
	}
	paramsBytes := message[offset : offset+paramsLen]
	offset += paramsLen

	// --- 6. Unmarshal params JSON into map
	if len(paramsBytes) > 0 {
		var params map[string]interface{}
		if err := json.Unmarshal(paramsBytes, &params); err != nil {
			return fmt.Errorf("invalid params JSON: %w", err)
		}
		c.Params = params
	} else {
		c.Params = make(map[string]interface{})
	}

	// --- 7. header length (2 bytes, big endian)
	if offset+2 > len(message) {
		return errors.New("missing header length")
	}
	headerLen := int(binary.BigEndian.Uint16(message[offset : offset+2]))
	offset += 2

	if offset+headerLen > len(message) {
		return errors.New("invalid header length")
	}
	headerBytes := message[offset : offset+headerLen]

	// --- 8. Unmarshal header JSON into map
	if len(headerBytes) > 0 {
		var header map[string]string
		if err := json.Unmarshal(headerBytes, &header); err != nil {
			return fmt.Errorf("invalid header JSON: %w", err)
		}
		c.Header = header
	} else {
		c.Header = make(map[string]string)
	}

	return nil
}
