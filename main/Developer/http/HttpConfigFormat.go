package Http

import (
	"io"
	"net/http"
	"time"
)

type SetHttpConfig struct {
	TimeOut time.Duration
	Method  string
	Body    io.Reader
	Uri     string
	Client  *http.Client
}
