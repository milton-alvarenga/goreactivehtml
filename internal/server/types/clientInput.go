package types

type ClientInputInterface interface {
	SendToClient(ClientOutput ClientOutput) bool
	IsValidMessage() error
	Close() error
	IsValidExecutor(string, string) bool
	IsValidOperation(operation string) bool
	Unmarshal(message byte)
}
