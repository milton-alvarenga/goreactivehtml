package handle

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/milton-alvarenga/goreactivehtml/internal/server/types"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // permitir qualquer origem
}

type PID uint8

type ProcessorQueue map[PID]types.ClientOutput
