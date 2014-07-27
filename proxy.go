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
	RequestConverter func(originalRequest, proxyRequest *http.Request)
	Transport http.RoundTripper
}

// Create a host-based reverse-proxy.
func NewProxyWithHostConverter(hostConverter func(string) string) *Proxy {
	return &Proxy{
		RequestConverter: func(originalRequest, proxyRequest *http.Request) {
			proxyRequest.URL.Host = hostConverter(originalRequest.Host)
		},
		Transport: http.DefaultTransport,
	}
}

// Create a request-based reverse-proxy.
func NewProxyWithRequestConverter(requestConverter func(*http.Request, *http.Request)) *Proxy {
	return &Proxy{
		RequestConverter: requestConverter,
		Transport: http.DefaultTransport,
	}
}

func (proxy *Proxy) ServeHTTP(writer http.ResponseWriter, originalRequest *http.Request) {
	// Create a new proxy request object by coping the original request.
	proxyRequest := proxy.copyRequest(originalRequest)

	// Convert an original request into another proxy request.
	proxy.RequestConverter(originalRequest, proxyRequest)

	// Convert a request into a response by using its Transport.
	response, err := proxy.Transport.RoundTrip(proxyRequest)
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

// Create a new proxy request with some modifications from an original request.
func (proxy *Proxy) copyRequest(originalRequest *http.Request) *http.Request {
	proxyRequest := new(http.Request)
	*proxyRequest = *originalRequest
	proxyRequest.Proto = "HTTP/1.1"
	proxyRequest.ProtoMajor = 1
	proxyRequest.ProtoMinor = 1
	proxyRequest.Close = false
	proxyRequest.Header = make(http.Header)
	proxyRequest.URL.Scheme = "http"
	proxyRequest.URL.Path = originalRequest.URL.Path

	// Copy all header fields.
	for key, values := range originalRequest.Header {
		for _, value := range values {
			proxyRequest.Header.Add(key, value)
		}
	}

	// Remove ignored header fields.
	for _, headerName := range ignoredHeaderNames {
		proxyRequest.Header.Del(headerName)
	}

	// Append this machine's host name into X-Forwarded-For.
	if requestHost, _, err := net.SplitHostPort(originalRequest.RemoteAddr); err == nil {
		if originalValues, ok := proxyRequest.Header["X-Forwarded-For"]; ok {
			requestHost = strings.Join(originalValues, ", ") + ", " + requestHost
		}
		proxyRequest.Header.Set("X-Forwarded-For", requestHost)
	}

	return proxyRequest
}
