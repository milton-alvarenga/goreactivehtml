package types

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	Conn *websocket.Conn
	mu   sync.Mutex
}

func (wsc *WebSocketConnection) Write(messageType int, data []byte) error {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()
	return wsc.Conn.WriteMessage(messageType, data)
}

type PID uint8

type ProcessorQueue map[PID]ClientOutput

type WSTypeOutputMessage byte

type WSOperation string

type WSDestination string

const (
	WSOperationSelect WSOperation = "S"
	WSOperationInsert WSOperation = "I"
	WSOperationUpdate WSOperation = "U"
	WSOperationDelete WSOperation = "D"

	WSTypeErrorOutputMessage   WSTypeOutputMessage = 'E'
	WSTypeSuccessOutputMessage WSTypeOutputMessage = 'S'

	WSDestinationUnknown WSDestination = "*unknown"
)

func GetValidOperations() map[WSOperation]bool {
	return map[WSOperation]bool{
		WSOperationSelect: true,
		WSOperationInsert: true,
		WSOperationUpdate: true,
		WSOperationDelete: true,
	}
}

type Subscription struct {
	MsgType     WSTypeOutputMessage
	Destination string
	Callback    func(ClientOutput) // The function to execute
}

type WebSocketClient struct {
	Subscriptions []Subscription
}

func (client *WebSocketClient) Subscribe(subscription Subscription) {
	client.Subscriptions = append(client.Subscriptions, subscription)
}

func (client *WebSocketClient) ProcessMessage(msg ClientOutput) {
	for _, sub := range client.Subscriptions {
		if sub.MsgType == msg.MsgType && sub.Destination == msg.Destination {
			sub.Callback(msg) // Execute the callback
		}
	}
}
