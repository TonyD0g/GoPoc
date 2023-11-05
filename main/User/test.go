package User

import (
	"GoPoc/main/Developer/AllFormat"
	"GoPoc/main/Developer/Http"
	"net/http"
	"strings"
)

var Json string
var Poc func(hostInfo string, client *http.Client) bool
var Exp func(expResult Format.ExpResult, client *http.Client) Format.ExpResult

func init() {
	Json = `{
    "fofa": "body=\"Hello World\"",
    "Request": {
        "Method": "GET",
        "Uri": [
            "/robots.txt"
        ],
        "Header": {
            "Accept-Encoding": "gzip"
        }
    },
    "Response": {
        "Operation": "OR",
        "Group": [
            {
                "Body": [
                    "Hello World"
                ]
            }
        ]
    }
}`

	Poc = func(hostInfo string, client *http.Client) bool {
		var config Http.SetHttpConfig
		config.TimeOut = 5
		config.Method = "GET"
		config.Body = nil
		config.Uri = "/robots.txt"
		config.Client = client
		resp, err := Http.SendHttpRequest(hostInfo, config)
		return err == nil && strings.Contains(resp.Body, "Hello World21211221")
	}

	Exp = func(expResult Format.ExpResult, client *http.Client) Format.ExpResult {
		return expResult
	}
}
