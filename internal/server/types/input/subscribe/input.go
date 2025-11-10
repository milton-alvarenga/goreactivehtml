package subscribe

import (
	"errors"
	"log"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/milton-alvarenga/goreactivehtml/internal/server/types"
)

type ClientInputSubscription struct {
	ReqId  string
	Topic  string
	Data   string
	Header map[string]string
	WSConn *types.WebSocketConnection
}

func (c ClientInputSubscription) SendToClient(ClientOutput types.ClientOutput) bool {
	message, err := ClientOutput.Marshal()
	if err != nil {
		log.Println("Could not marshal client output")
		return false
	}
	err = c.WSConn.Write(websocket.TextMessage, []byte(message))
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

func (c ClientInputSubscription) IsValidMessage() error {

	if !c.IsValidTopic() {
		return errors.New("invalid endpoint")
	}

	return nil
}

func (c *ClientInputSubscription) Close() error {
	return c.WSConn.Conn.Close()
}

func (c ClientInputSubscription) IsValidEndpoint() bool {
	validsEndpoints := c.GetValidEndpoints()

	_, found := validsEndpoints[c.Endpoint]
	return found
}

func (c ClientInputSubscription) IsValidOperation(operation WSOperation) bool {
	return GetValidOperations()[operation]
}

func (c *ClientInputSubscription) Unmarshal(message string) {
	validsEndpoints := c.GetValidEndpoints()

	if message[0] != '/' {
		c.ReqId = new(uint8)
		*c.ReqId = uint8(message[0])
		message = message[1:]
	}

	endpoint, found := validsEndpoints[c.Endpoint]
	var operation string
	var origin string
	var data string

	for _, endpoint := range validsEndpoints {
		if strings.HasPrefix(message, endpoint) {
			body := strings.TrimPrefix(message, endpoint)

			if len(body) > 0 {
				if c.IsValidOperation(WSOperation(body[0])) {
					operation = string(body[0])
				}
				parts := strings.SplitN(body[1:], ";", 2)

				if len(parts) == 2 {
					origin = parts[0]
					body = parts[1]
				} else {
					return
				}

				if len(body) > 1 {
					data = body
				}
			}
			c.Endpoint = endpoint
			c.Operation = WSOperation(operation)
			c.Origin = origin
			c.Data = data
		}
	}
}
