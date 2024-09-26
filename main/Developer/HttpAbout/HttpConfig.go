package HttpAbout

import (
	"net/http"
	"time"
)

type SetHttpConfig struct {
	TimeOut    time.Duration
	Method     string
	Body       string
	Uri        string
	Client     *http.Client
	Header     map[string]string
	IsRedirect bool // 是否重定向
}

func NewHttpConfig() SetHttpConfig {
	var config SetHttpConfig
	config.Header = make(map[string]string)
	return config
}
