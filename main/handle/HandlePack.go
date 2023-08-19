package handle

import (
	"Scanner/main/format"
	"net/http"
)

// HandlePack TODO 处理http请求包
func HandlePackFunc(resp *http.Request, requestPackage format.RequestPackage) {
	resp.Header.Add("User-Agent", "My Custom User Agent")

}
