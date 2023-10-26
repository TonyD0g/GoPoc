# GoScanner 

【并未发布正式版，暂不可用】

基于 Json文件 、自定义脚本的快速扫描器，用于快速验证目标是否存在该漏洞。

使用场景：

你是否发现一个新漏洞但苦于无法快速编写成脚本？这款工具可能适合你，只要你能看http请求包/响应包就能上手写poc。



# 使用教程

- 程序使用

```md
-email // fofa的email (必须)
-key // fofa的key (必须)
-url // 扫单个url (非必须)
-file // 扫url文件中的每一个url (非必须)
-pocJson // poc的json文件 (必须)
-proxy // burpsuite 代理 (必须)
-maxConcurrentLevel // 最大并发量,越大速度越快 (非必须)
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

- Poc 编写

```json
{
    // 必须,表明想要查找的fofa语句.
    "fofa":"body=\"hello world\"", 
   	// 请求包
    "RequestPackage":{
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
    "ResponsePackage":{
        // 定义多个Group之间的关系,有AND和OR这两种,其中AND是都满足漏洞才存在,OR是其中一个条件满足即可.
		"Operation":"OR",
        // 判断条件
		"Group":[
            	 // 条件1
				{
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
			            },
			        "Body":[
                        				// 支持正则表达式,使用格式为 [regexp]`你所想要的表达式内容`
			        				    "[regexp]`.*?`" 
			            ]  
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
1. 支持 json/go 这两种POC方法
2. 增强fofa
3. 解决【验证是否利用成功】的速度过慢的问题
4. 让header和body一样支持正则表达式
```



【建议先别二开，代码写的很狗屎，先让我简单优化下。】
