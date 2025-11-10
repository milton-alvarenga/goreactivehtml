package handle

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/milton-alvarenga/goreactivehtml/internal/server/handle/auth"
	"github.com/milton-alvarenga/goreactivehtml/internal/server/types"
	"github.com/milton-alvarenga/goreactivehtml/internal/server/types/input/rest"
)

var connections = make(map[*types.WebSocketConnection]bool)

func WS(w http.ResponseWriter, r *http.Request) {
	if !auth.Check(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	wsc := &types.WebSocketConnection{
		Conn: c,
	}

	connections[wsc] = true

	var message string
loop:
	for {
		msgType, msg, err := wsc.Conn.ReadMessage()
		log.Println("Msg received...")
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				delete(connections, wsc)

				log.Println("Client closed the connection")
				break
			} else {
				log.Println("Error reading message:", err)
			}
			return
		}

		switch msgType {
		//Method not supported
		case websocket.TextMessage:
			message = string(msg)

			output := shared.ClientOutput{
				MsgType:     shared.WSTypeErrorOutputMessage,
				Destination: string(shared.WSDestinationUnknown),
				Data:        "Websocket TextMessage not yeat supported. Use websocket BinaryMessage mode",
			}
			//we need a specific binary unmarshal code
			wsc.Write(websocket.TextMessage, []byte(output.Marshal()))
			continue
		case websocket.BinaryMessage:
			// Handle the message in a separate goroutine
			go handleMessage(wsc, msg)
		case websocket.CloseMessage:
			//TODO
			//ADD CONTEXT AND SEND KILL SIGNAL TO GO ROUTINE
			log.Println("Client send close message for this connection")
			break loop
		//Automaticly treat by gorilla websocket library
		//case websocket.PingMessage:
		case websocket.PongMessage:
			continue
		}
	}
	log.Println("End of EntryConnections...")
}

func handleMessage(wsc *types.WebSocketConnection, message []byte) {

	log.Println("Start of handleMessage...")

	switch message[0] {
		//SUBSCRIBE
		case 1:
			fmt.Println("Received message type SUBSCRIBE")
			input := subscribe.ClientInputSubscription{
				WSConn: wsc,
			}
			input.Unmarshal(message)
			//
		//RPC
		case 2:
			fmt.Println("Received message type RPC")
			input := rpc.ClientInputRPC{
				WSConn: wsc,
			}
			input.Unmarshal(message)
			//Check if class exists
			//Check if method exists
			//Execute
			//Get the result
			//Return the response
		//ENDPOINT
		case 3:
			fmt.Println("Received message type ENDPOINT")
			input := rest.ClientInputRest{
				WSConn: wsc,
			}
			input.Unmarshal(message)
			//Check if endpoint exists
			//Check if the methd exists
			//Execute
			//Return the response
		default:
			msg := "Unknown message type"
			fmt.Println(msg)
			output := types.ClientOutput{
				ReqId: message[0],
				MsgType: types.WSTypeErrorOutputMessage,
				Data: msg,
			}
			wsc.Write(websocket.BinaryMessage,byte(output.Marshal()))
			return
		}
	}
	log.Println("End of handleMessage...")
}
