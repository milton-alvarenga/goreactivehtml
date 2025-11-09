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
		case websocket.TextMessage:
			message = string(msg)

			//TODO
			//Method not supported
		case websocket.BinaryMessage:
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

		// Handle the message in a separate goroutine
		go handleMessage(wsc, msg)
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
		//RPC
		case 2:
			fmt.Println("Received message type RPC")
			input := rpc.ClientInputRPC{
				WSConn: wsc,
			}
		//ENDPOINT
		case 3:
			fmt.Println("Received message type ENDPOINT")
			input := rest.ClientInputRest{
				WSConn: wsc,
			}
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

	input.Unmarshal(message)

	err := input.IsValidMessage()

	if err != nil {
		errorMsg := err.Error()
		outputMsg := types.ClientOutput{
			MsgType:     types.WSTypeErrorOutputMessage,
			Destination: input.Origin,
			Data:        errorMsg,
		}
		input.SendToClient(outputMsg)
		log.Println(errorMsg)
	}

	switch input.Endpoint {
	case "/book":
		// Handle book message
		service.BookClientInputProcessor(input)
	case "/trades":
		// Handle trades message
	case "/broker":
		// Handle broker message
	case "/stats/brokers":
	case "/stats/symbol":
	case "/stats/asset":
	case "/stats/exchange":
	case "/pingpong":
		service.PingPongProcessor(input)
	}
	log.Println("End of handleMessage...")
}
