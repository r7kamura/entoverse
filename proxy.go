package entoverse

import(
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

// These headers won't be copied from original request to proxy request.
var ignoredHeaderNames = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

// Provide host-based proxy server.
type Proxy struct {
	HostConverter func(string) string
}

// Create a entoverse.Proxy object with a default round tripper.
func NewProxy(hostConverter func(string) string) *Proxy {
	return &Proxy{
		HostConverter: hostConverter,
	}
}

func (proxy *Proxy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// Create a http.Request to be sent to an upstream server by shallow-coping.
	proxyRequest := new(http.Request)
	*proxyRequest = *request
	proxyRequest.Proto = "HTTP/1.1"
	proxyRequest.ProtoMajor = 1
	proxyRequest.ProtoMinor = 1
	proxyRequest.Close = false
	proxyRequest.Header = make(http.Header)
	proxyRequest.URL.Scheme = "http"
	proxyRequest.URL.Path = request.URL.Path
	proxyRequest.URL.Host = proxy.HostConverter(request.Host)
	if proxyRequest.URL.Host == "" {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// Copy all header fields.
	for key, values := range request.Header {
		for _, value := range values {
			proxyRequest.Header.Add(key, value)
		}
	}

	// Remove ignored header fields.
	for _, headerName := range ignoredHeaderNames {
		proxyRequest.Header.Del(headerName)
	}

	// Append this machine's host name into X-Forwarded-For.
	if requestHost, _, err := net.SplitHostPort(request.RemoteAddr); err == nil {
		if originalValues, ok := proxyRequest.Header["X-Forwarded-For"]; ok {
			requestHost = strings.Join(originalValues, ", ") + ", " + requestHost
		}
		proxyRequest.Header.Set("X-Forwarded-For", requestHost)
	}

	// Convert a request into a response by using its Transport.
	response, err := http.DefaultTransport.RoundTrip(proxyRequest)
	if err != nil {
		log.Printf("ErrorFromProxy: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Ensure a response body from upstream will be always closed.
	defer response.Body.Close()

	// Copy all header fields.
	for key, values := range response.Header {
		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}

	// Copy a status code.
	writer.WriteHeader(response.StatusCode)

	// Copy a response body.
	io.Copy(writer, response.Body)
}
