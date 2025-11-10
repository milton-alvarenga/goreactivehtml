package topics

import (
	"github.com/milton-alvarenga/goreactivehtml/internal/server/types"
)

// ClientInput and ClientOutput
type HandleFunc func(*types.ClientInputInterface) *types.ClientOutput

type Topic string

type Topics map[Topic]HandleFunc

var topics = make(Topics)

func Register(topic Topic, handler HandleFunc) {
	topics[topic] = handler
}

func IsValidEndpoint(topic Topic) bool {
	_, ok := topics[topic]
	return ok
}

func IsValid(topic Topic) bool {
	return IsValidEndpoint(topic)
}

func Exec(topic Topic, ClientInput *types.ClientInputInterface) *types.ClientOutput {
	return topics[topic](ClientInput)
}
