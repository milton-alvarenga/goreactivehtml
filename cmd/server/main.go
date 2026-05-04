package main

import (
    "flag"
	"log"
	"net/http"

	"github.com/milton-alvarenga/goreactivehtml/internal/server/handle"
)

/*
Msg pattern

	INPUT
		connId
		namespace
		function name
		data
	OUTPUT
		connID
		status (E/S) Error or Success
		data
*/
func main() {
    var version = "dev"
    var buildDate = "unknown"

	showVersion := flag.Bool("v", false, "Show current version (git hash)")
	flag.Parse()

	if *showVersion {
		log.Printf("Ganha1000 Version: %s\n", version)
		log.Printf("Build Date: %s\n", buildDate)
		return
	}


	http.HandleFunc("/ws", handle.WS)

	// Serve static files from the "static" directory
	/*
	project/
	├── main.go
	├── handle/
	│   └── ws.go
	└── static/
		├── index.html
		├── style.css
		└── app.js
	

	http://localhost:8080/
	Go will serve static/index.html automatically.
	*/
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
