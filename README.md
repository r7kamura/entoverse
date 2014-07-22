# Entoverse
A library to implement host-based HTTP reverse-proxy in Golang.

## Usage
Here are some example uses of entoverse.

### Example1 - single host
`:3000` -> `:4000`

```go
package main

import(
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
	http.ListenAndServe("localhost:3000", proxy)
}
```

### Example2 - multi hosts
`:3000` -> `:4000, :5000, :6000`

```go
package main

import(
	"container/ring"
	"net/http"
	"sync"
	"github.com/r7kamura/entoverse"
)

func main() {
	// Creates a circular list for round-robin HTTP roundtrip.
	hosts := []string{
		":4000",
		":5000",
		":6000",
	}
	hostRing := ring.New(len(hosts))
	for _, host := range hosts {
		hostRing.Value = host
		hostRing = hostRing.Next()
	}

	// Locks by mutex because hostConverter will be executed in parallel.
	mutex := sync.Mutex{}
	hostConverter := func(originalHost string) string {
		mutex.Lock()
		defer mutex.Unlock()
		host := hostRing.Value.(string)
		hostRing = hostRing.Next()
		return host
	}

	// Runs a reverse-proxy server on http://localhost:3000/
	proxy := entoverse.NewProxy(hostConverter)
	http.ListenAndServe(":3000", proxy)
}
```
