package endpoints

import (
	"github.com/milton-alvarenga/goreactivehtml/internal/server/types"
)

// ClientInput and ClientOutput
type HandleFunc func(*types.ClientInputInterface) *types.ClientOutput

type Endpoint string
type Method string

type Router map[Endpoint]map[Method]HandleFunc

var routes = make(Router)

func Register(endpoint Endpoint, method Method, handler HandleFunc) {
	routes[endpoint] = make(map[Method]HandleFunc)
	routes[endpoint][method] = handler
}

func IsValidEndpoint(endpoint Endpoint) bool {
	_, ok := routes[endpoint]
	return ok
}

func IsValid(endpoint Endpoint, method Method) bool {
	_, ok := routes[endpoint][method]
	return ok
}

func Exec(endpoint Endpoint, method Method, ClientInput *types.ClientInputInterface) *types.ClientOutput {
	return routes[endpoint][method](ClientInput)
}
