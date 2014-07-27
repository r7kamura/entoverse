package main

import(
	"net/http"
	"github.com/r7kamura/entoverse"
)

func main() {
	hostConverter := func(originalHost string) string {
		return ":9292"
	}
	proxy := entoverse.NewProxyWithHostConverter(hostConverter)
	http.ListenAndServe(":3000", proxy)
}
