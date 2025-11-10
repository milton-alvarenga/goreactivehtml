package procedures

import (
	"github.com/milton-alvarenga/goreactivehtml/internal/server/types"
)

// HandleFunc is a function signature that all handler functions should match
// It takes a request (represented as an interface{} for flexibility) and returns a result (interface{}) and an error
// type HandleFunc func(request interface{}) (interface{}, error)

// ClientInput and ClientOutput
type HandleFunc func(*types.ClientInputInterface) *types.ClientOutput

type Class string
type Method string

type Procedure map[Class]map[Method]HandleFunc

var routes = make(Procedure)

func Register(class Class, method Method, handler HandleFunc) {
	routes[class] = make(map[Method]HandleFunc)
	routes[class][method] = handler
}

func IsValidClass(class Class) bool {
	_, ok := routes[class]
	return ok
}

func IsValid(class Class, method Method) bool {
	_, ok := routes[class][method]
	return ok
}

func Exec(ClientInput types.ClientInputInterface) *types.ClientOutput {
	// Type assertion to access fields on the concrete type
	clientInputRPC, ok := ClientInput.(*types.ClientInputRPC)
	if !ok {
		// If ClientInput is not of type *types.ClientInputRPC, return nil or handle the error
		return nil // or handle the error appropriately
	}

	return routes[Class(clientInputRPC.Class)][Method(clientInputRPC.Method)](clientInputRPC)
}
