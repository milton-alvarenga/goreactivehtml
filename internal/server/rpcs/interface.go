package rpcs

type RPC interface {
	Marshal(data []byte)
	Unmarshal()
	in
	out
	fns map[string]
}
