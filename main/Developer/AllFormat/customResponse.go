package Format

import (
	"crypto/tls"
	"io"
	"net/http"
)

type CustomResponseFormat struct {
	Body             string
	Status           string // e.g. "200 OK"
	StatusCode       int    // e.g. 200
	Proto            string // e.g. "HTTP/1.0"
	ProtoMajor       int    // e.g. 1
	ProtoMinor       int    // e.g. 0
	Header           map[string][]string
	RawBody          io.ReadCloser
	ContentLength    int64
	TransferEncoding []string
	Close            bool
	Uncompressed     bool
	Trailer          map[string][]string
	Request          *http.Request
	TLS              *tls.ConnectionState
}
