# GoPoc

基于 Json文件 、自定义脚本的快速验证扫描器，用于快速验证目标是否存在该漏洞。

使用场景：

你是否发现一个新漏洞但苦于无法快速编写成脚本？这款工具可能适合你，只要你能看http请求包/响应包就能上手写poc。

# 免责声明

```md
本工具仅面向合法授权的企业安全建设行为，如您需要测试本工具的可用性，请自行搭建靶机环境。

为避免被恶意使用，本项目不提供任何poc，不存在漏洞利用过程，不会对目标发起真实攻击和漏洞利用。

在使用本工具进行检测时，您应确保该行为符合当地的法律法规，并且已经取得了足够的授权。请勿对非授权目标进行扫描。

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任。

在安装并使用本工具前，请您务必审慎阅读、充分理解各条款内容，限制、免责条款或者其他涉及您重大权益的条款可能会以加粗、加下划线等形式提示您重点注意。 除非您已充分阅读、完全理解并接受本协议所有条款，否则，请您不要安装并使用本工具。您的使用行为或者您以其他任何明示或者默示方式表示接受本协议的，即视为您已阅读并同意本协议的约束。
```



# 使用教程

- 请到此处下载发行版：[releases](https://github.com/TonyD0g/GoPoc/releases)

- 程序使用

```md
第1步: 创建一个配置文件,比如 config.txt
第2步:使用 -ini 参数加载配置文件.
------------------------------------
config.txt内容如下:

-email // fofa的email (必须)
-key // fofa的key (必须)
-url // 扫单个url (非必须)
-file // 扫url文件中的每一个url (非必须)
-pocJson // poc的json文件 (必须)
-proxy // burpsuite 代理 (必须)
-maxConcurrentLevel // 最大并发量,越大速度越快 (必须)
-maxFofaSize	   // 最大检索数 (必须)
------------------------------------
例如
-email
21212121
-key
212121212
-pocJson
C:\Users\xxx\Desktop\1.json
-proxy
http://127.0.0.1:8082
-maxConcurrentLevel
3
```

- 书写规范:	

  ```md
  1. 使用代码模式时必须存在 Poc/Exp 函数,如果是使用json模式不写 Poc/Exp 函数
  2. json 中必须存在 fofa语句
  ```

- Poc 编写

```json
{
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
}
```



# 特性

```md
1. 基于 fofa 规则匹配对应产品,匹配成功后才开始使用POC,避免发送无用包
2. POC 易编写,只需要会看http响应包和http回显包即可
```

# TODO

```md
进一步封装代码模式
```



【建议先别二开，代码写的很狗屎，先让我简单优化下。】
