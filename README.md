# Entoverse
A library to implement host-based HTTP reverse-proxy in Golang.

## Usage
e.g. `localhost:3000` -> `localhost:4000`

```go
package main

import(
	"log"
	"net/http"
	"github.com/r7kamura/entoverse"
)

func main() {
	// Please implement your host converter function.
	// This example always delegates HTTP requests to localhost:4000.
	hostConverter := func(originalHost string) string {
		return "localhost:4000"
	}

	// Creates an entoverse.Proxy object as an HTTP handler.
	proxy := entoverse.NewProxy(hostConverter)

	// Runs a reverse-proxy server on http://localhost:3000/
	serverError := http.ListenAndServe("localhost:3000", proxy)
	log.Fatal(serverError)
}
```
