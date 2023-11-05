package Http

import (
	"GoPoc/main/Developer/AllFormat"
	"context"
	"io"
	"net/http"
	"time"
)

func SendHttpRequest(hostInfo string, config SetHttpConfig) (Format.CustomResponseFormat, error) {
	var customResponse Format.CustomResponseFormat
	// Create an HTTP.Request object
	procedureRequest, err := http.NewRequest(config.Method, hostInfo+config.Uri, config.Body)
	if err != nil {
		return customResponse, err
	}

	// Set a timeout for the request
	if config.TimeOut == 0 {
		config.TimeOut = 5
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*config.TimeOut)
	// todo 添加各种header头
	procedureRequest = procedureRequest.WithContext(ctx)
	// Send request and obtain response results
	procedureResponse, err := config.Client.Do(procedureRequest)
	if err != nil {
		cancel()
		return customResponse, err
	}

	// 对自定义Response进行赋值
	bodyOfExecutionResults, err := io.ReadAll(procedureResponse.Body)
	if err != nil {
		cancel()
		return customResponse, err
	}
	customResponse.Body = string(bodyOfExecutionResults)
	customResponse.Request = procedureResponse.Request
	customResponse.Header = procedureResponse.Header
	customResponse.Status = procedureResponse.Status
	customResponse.RawBody = procedureResponse.Body
	customResponse.Close = procedureResponse.Close
	customResponse.ContentLength = procedureResponse.ContentLength
	customResponse.Proto = procedureResponse.Proto
	customResponse.ProtoMajor = procedureResponse.ProtoMajor
	customResponse.ProtoMinor = procedureResponse.ProtoMinor
	customResponse.TLS = procedureResponse.TLS
	customResponse.StatusCode = procedureResponse.StatusCode
	customResponse.Trailer = procedureResponse.Trailer
	customResponse.TransferEncoding = procedureResponse.TransferEncoding
	customResponse.Uncompressed = procedureResponse.Uncompressed

	cancel()
	return customResponse, nil
}
