# [红蓝工具]GoPoc

**基于 Json 、自定义Go脚本、fofa的快速验证扫描引擎，可用于快速验证目标是否存在该漏洞。**

你甚至可以拿来写自动化脚本来加快你的工作效率,比如从海量url中获取到符合条件的url,从而避免重复的点鼠标环节

使用场景：

- 你是否发现一个新漏洞但苦于无法快速编写成脚本？这款工具可能适合你，只要你能看http请求包/响应包就能上手写poc。
- 想更高深的利用？比如文件上传写冰蝎、哥斯拉。那么代码模式适合你。
- 和 nuclei 的区别: nuclei有非常庞大的用户社区,可以比较快速的提供非常新的poc,而GoPoc侧重在poc/exp编写,
  通过断点调试的方法可以帮助我们快速编写poc/exp , 最后再结合fofa来快速打点

# 注意事项:
- User文件夹下只能放一个自己写的go poc文件,这是由于go语言的特性,个人暂时没有解决办法
- 想使用 Poc文件夹下的 go文件时,请"复制"过去,而不是直接拖过去,因为拖过去会弹个窗让你重构,你点击重构代码就乱套了,就无法使用 (你要是会开发可以直接拖过去)

# 免责声明

```md
本工具仅面向合法授权的企业安全建设行为，如您需要测试本工具的可用性，请自行搭建靶机环境。

为避免被恶意使用，本项目不提供任何poc，不存在漏洞利用过程，不会对目标发起真实攻击和漏洞利用。

在使用本工具进行检测时，您应确保该行为符合当地的法律法规，并且已经取得了足够的授权。请勿对非授权目标进行扫描。

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任。

在安装并使用本工具前，请您务必审慎阅读、充分理解各条款内容，限制、免责条款或者其他涉及您重大权益的条款可能会以加粗、加下划线等形式提示您重点注意。 除非您已充分阅读、完全理解并接受本协议所有条款，否则，请您不要安装并使用本工具。您的使用行为或者您以其他任何明示或者默示方式表示接受本协议的，即视为您已阅读并同意本协议的约束。
```



# 使用教程

- 1. 请 **下载源码** 使用 
  2. 使用ide打开项目，比如 goland
  3. 在 User 包下新建一个 "文件名.go"
  4. 创建一个配置文件,比如 config.ini
  5. 使用 -ini 参数加载配置文件.
- config.ini 内容如下：

```md
-email // fofa的email (必须)
-key // fofa的key (必须)
-url // 扫单个url (非必须)
-file // 扫url文件中的每一个url (非必须)
-vul // poc/exp文件,文件后缀为.go (必须)
-mod // 指定poc/exp这两种模式 (必须)
-proxy // burpsuite 代理,用于方便写poc/exp (必须)
-maxConcurrentLevel // 最大并发量,越大扫描速度越快 (必须)
-maxFofaSize	   // fofa最大检索数 (必须)
------------------------------------
例如
-email
xxxxxxxxxxx@gamil.com
-key
fdgdfhfgdhdfgdhfghfdg
-vul
D:\Coding\Github\GoPoc\main\User\test.go
-mod
poc
-url
http://127.0.0.1
-proxy
http://127.0.0.1:8082
-maxConcurrentLevel
3
-maxFofaSize
300
```

- 两种利用模式：

  - json模式 [非常不推荐,因为不够灵活,且后续版本不会对其进行优化,所以会非常有可能会遇到各种奇奇怪怪的问题]

    新建一个 **文件名.go** 文件，输入以下内容：

  ```go
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
  	// 有代码使用代码,无代码使用json
  	// 如果存在代码,可以不写Json格式(即Json格式有架构,但内容为空).但必须存在 fofa语句
  	// 此处的json只是说明json的使用方式,与代码模式并无关联
  	Json = `{
      // 必须,表明想要查找的fofa语句.
      "fofa":"body=\"hello world\"",
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
  }
  
  ```

  

  - go 模式

    在 main\User 文件夹下新建一个 **文件名.go** 文件，然后 config 文件指定该 go 文件即可,go 文件的内容如下[只做样例说明,真正的文件请看 dvwaSqlScan.go 文件]：

  ```go
  package User

  import (
  	"GoPoc/main/Developer/AllFormat"
  	"GoPoc/main/Developer/Http"
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
  		config := Http.NewHttpConfig()
  		config.Uri = "/dvwa/"
  		config.Client = client
  		resp, err := Http.SendHttpRequest(hostInfo, config)
  		if err != nil {
  			return "", err
  		}
  		config.Header["Cookie"] = resp.Header["Set-Cookie"][0]
  		config.Body = "username=admin&password=password&Login=Login&user_token=" + Utils.RandomStringByModule(24, 1)
  		config.Uri = "/dvwa/login.php"
  		resp, err = Http.SendHttpRequest(hostInfo, config)
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
  		config := Http.NewHttpConfig()
  		config.Header["Cookie"], err = getSessionByLogin(hostInfo, client) // (非强制) 自定义Header头
  		if err != nil {
  			return customResponse, nil
  		}
  		config.Header["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"
  		config.Header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:129.0) Gecko/20100101 Firefox/129.0" // (非强制) 如果不写的话随机从默认URL列表上选取一个
  		config.TimeOut = 5                                                                                               // (非强制) 如果不写的话默认值为 5秒
  		config.Method = "GET"                                                                                            // (非强制) 如果不写的话默认值为 GET方式
  		//config.Body = `123` 			// (非强制) 如果不写的话默认值为 ""
  		config.Uri = "/dvwa/index.php" // (非强制) 如果不写的话默认值为 ""
  		config.Client = client         // (强制) 因为这个 client 挂上了burpSuite代理,如果你使用自己的client可能会因为没有挂代理而无法得知利用过程,会不好写poc/exp
  		return Http.SendHttpRequest(hostInfo, config)
  	}

  	// 建议: 函数名+随机命名
  	sendSqlPayload5251552 := func(hostInfo string, client *http.Client) (Format.CustomResponseFormat, error) {
  		config := Http.NewHttpConfig()
  		config.Header["Cookie"] = "security=low; PHPSESSID=1abcbe73869e90560e9061ca636c813e"
  		config.Uri = "/dvwa/vulnerabilities/sqli/?id=%27&Submit=Submit#"
  		config.Client = client
  		return Http.SendHttpRequest(hostInfo, config)
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

  
  ```

  

- 书写规范:	

  ```md
  1. 使用go模式时必须存在 Poc/Exp 函数,如果是使用json模式不写 Poc/Exp 函数
  2. 如果存在代码,可以不写Json格式(即Json格式有架构,但内容为空).但必须存在 fofa语句
  3. json 模式只支持 poc,不支持 exp
  4. json 模式能书写的复杂度不如 go 模式,如果想更深层次利用请使用 go 模式
  ```

# 效果展示：

![演示](pic\演示.gif)

# 特性

```md
1. 基于 fofa 规则匹配对应产品,匹配成功后才开始使用POC,避免发送无用包
2. POC 易编写,只需要会看http响应包和http回显包即可
```
