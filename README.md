# Entoverse
A library to implement host-based HTTP reverse-proxy in Golang.

## Usage
See [examples](examples) for some example entoverse uses.

* [Simple host-based reverse-proxy (:3000 -> :4000)](examples/single.go)
* [Roundrobin host-based reverse-proxy (:3000 -> :4000, :5000, :6000)](examples/roundrobin.go)
