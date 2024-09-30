package HttpAbout

import (
	"net/http"
	"time"
)

type SetHttpConfig struct {
	TimeOut    time.Duration     // 请求等待时间
	Method     string            // 请求方法
	Body       string            // 请求体
	Uri        string            // 请求的路径
	Client     *http.Client      // 代理选项
	Header     map[string]string //请求头
	IsRedirect bool              // 是否重定向
}

func NewHttpConfig() SetHttpConfig {
	var config SetHttpConfig
	config.Header = make(map[string]string)
	return config
}
