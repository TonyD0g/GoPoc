package User

import (
	"GoPoc/main/Developer/AllFormat"
	"context"
	"io"
	"net/http"
	"strings"
	"time"
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
		// Create an HTTP.Request object
		procedureRequest, err := http.NewRequest("GET", hostInfo+"/robots.txt", nil)
		if err != nil {
			return false
		}
		// Set a timeout for the request
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		procedureRequest = procedureRequest.WithContext(ctx)
		// Send request and obtain response results
		procedureResponse, err := client.Do(procedureRequest)
		if err != nil {
			cancel()
			return false
		}
		cancel()
		bodyOfExecutionResults, err := io.ReadAll(procedureResponse.Body)
		return strings.Contains(string(bodyOfExecutionResults), "Hello World")
	}

	Exp = func(expResult Format.ExpResult, client *http.Client) Format.ExpResult {
		return expResult
	}
}
