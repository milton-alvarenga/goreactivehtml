package rest

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/milton-alvarenga/goreactivehtml/internal/server/types"
	"github.com/milton-alvarenga/goreactivehtml/internal/server/types/input/rest/endpoints"
)

type ClientInputRest struct {
	ReqId    *uint8
	Method   RESTMethod
	Endpoint string
	Data     string
	Header   map[string]string
	WSConn   *types.WebSocketConnection
}

func (c ClientInputRest) SendToClient(ClientOutput types.ClientOutput) bool {
	message, err := ClientOutput.Marshal()
	if err != nil {
		log.Println("Could not marshal client output")
		return false
	}
	err = c.WSConn.Write(websocket.BinaryMessage, message)
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

func (c ClientInputRest) IsValidMessage() error {

	if !c.IsValidEndpoint(c.Endpoint) {
		return errors.New("invalid endpoint")
	}

	if !c.IsValidOperation(string(c.Method)) {
		return errors.New("invalid method operation")
	}

	return nil
}

func (c *ClientInputRest) Close() error {
	return c.WSConn.Conn.Close()
}

func (c ClientInputRest) IsValidExecutor(endpoint_identification string, method string) bool {
	return endpoints.IsValid(
		endpoints.Endpoint(endpoint_identification),
		endpoints.Method(method),
	)
}

func (c ClientInputRest) IsValidEndpoint(endpoint_identification string) bool {
	return endpoints.IsValidEndpoint(endpoints.Endpoint(endpoint_identification))
}

func (c ClientInputRest) IsValidOperation(operation string) bool {
	switch RESTMethod(operation) {
	case GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS, TRACE, CONNECT:
		return true
	default:
		return false
	}
}

func (c *ClientInputRest) Unmarshal(message []byte) error {
	offset := 0
	if len(message) < 1+1+1 {
		return errors.New("message too short")
	}

	// --- 1. Skip conn_type (1 byte)
	offset++

	// --- 2. reqId (1 byte)
	reqId := message[offset]
	c.ReqId = &reqId
	offset++

	// --- 3. method length (1 byte)
	methodLen := int(message[offset])
	offset++

	if offset+methodLen > len(message) {
		return errors.New("invalid method length")
	}
	c.Method = RESTMethod(string(message[offset : offset+methodLen]))
	offset += methodLen

	// --- 4. endpoint length (2 bytes, big endian)
	if offset+2 > len(message) {
		return errors.New("missing endpoint length")
	}
	endpointLen := int(binary.BigEndian.Uint16(message[offset : offset+2]))
	offset += 2

	if offset+endpointLen > len(message) {
		return errors.New("invalid endpoint length")
	}
	c.Endpoint = string(message[offset : offset+endpointLen])
	offset += endpointLen

	// --- 5. payload length (4 bytes, big endian)
	if offset+4 > len(message) {
		return errors.New("missing payload length")
	}
	payloadLen := int(binary.BigEndian.Uint32(message[offset : offset+4]))
	offset += 4

	if offset+payloadLen > len(message) {
		return errors.New("invalid payload length")
	}
	c.Data = string(message[offset : offset+payloadLen])
	offset += payloadLen

	// --- 6. header length (2 bytes, big endian)
	if offset+2 > len(message) {
		return errors.New("missing header length")
	}
	headerLen := int(binary.BigEndian.Uint16(message[offset : offset+2]))
	offset += 2

	if offset+headerLen > len(message) {
		return errors.New("invalid header length")
	}
	headerBytes := message[offset : offset+headerLen]

	// --- 7. parse header JSON
	var header map[string]string
	if len(headerBytes) > 0 {
		if err := json.Unmarshal(headerBytes, &header); err != nil {
			return fmt.Errorf("invalid header JSON: %w", err)
		}
	}
	c.Header = header

	return nil
}
