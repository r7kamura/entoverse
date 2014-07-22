package main

import(
	"log"
	"net/http"
	"github.com/r7kamura/entoverse"
)

func main() {
	hostConverter := func(originalHost string) string {
		return ":9292"
	}
	proxy := entoverse.NewProxy(hostConverter)
	serverError := http.ListenAndServe(":3000", proxy)
	log.Fatal(serverError)
}
