package HttpAbout

import (
	"net/http"
	"time"
)

type SetHttpConfig struct {
	TimeOut time.Duration
	Method  string
	Body    string
	Uri     string
	Client  *http.Client
	Header  map[string]string
}

func NewHttpConfig() SetHttpConfig {
	var config SetHttpConfig
	config.Header = make(map[string]string)
	return config
}
