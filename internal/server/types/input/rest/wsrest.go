package rest

// Define a custom type for REST methods
type RESTMethod string

// Define constants for each HTTP method
const (
	GET     RESTMethod = "GET"
	POST    RESTMethod = "POST"
	PUT     RESTMethod = "PUT"
	PATCH   RESTMethod = "PATCH"
	DELETE  RESTMethod = "DELETE"
	HEAD    RESTMethod = "HEAD"
	OPTIONS RESTMethod = "OPTIONS"
	TRACE   RESTMethod = "TRACE"
	CONNECT RESTMethod = "CONNECT"
)
