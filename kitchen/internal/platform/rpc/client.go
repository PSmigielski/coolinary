package rpc

import (
	"net/http"
	"time"

	"connectrpc.com/connect"
)

// NewHTTPClient returns an HTTP client that speaks HTTP/2 without TLS (h2c),
// with keepalive pings so idle connections are detected and replaced.
func NewHTTPClient() *http.Client {
	var protocols http.Protocols
	protocols.SetUnencryptedHTTP2(true)
	return &http.Client{
		Transport: &http.Transport{
			Protocols: &protocols,
			HTTP2: &http.HTTP2Config{
				SendPingTimeout: 30 * time.Second,
				PingTimeout:     10 * time.Second,
			},
		},
	}
}

// ClientOptions makes generated clients use the gRPC protocol plus our interceptor.
func ClientOptions() []connect.ClientOption {
	return []connect.ClientOption{
		connect.WithGRPC(),
		connect.WithInterceptors(NewClientInterceptor()),
	}
}

// BaseURL turns "notifier:50051" into "http://notifier:50051".
func BaseURL(addr string) string { return "http://" + addr }
