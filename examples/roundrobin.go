package main

import(
	"container/ring"
	"log"
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
	serverError := http.ListenAndServe(":3000", proxy)
	log.Fatal(serverError)
}
