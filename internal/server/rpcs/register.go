package rpcs

import (
	"github.com/milton-alvarenga/goreactivehtml/internal/server/types"
)

// ClientInput and ClientOutput
type HandleFunc func(*types.WebSocketConnection, []byte)

type Router map[string]HandleFunc

var routes = Router{}

func Register(endpoint string, handler HandleFunc) {
	routes[endpoint] = handler
}

func Exec(endpoint string) {

}
