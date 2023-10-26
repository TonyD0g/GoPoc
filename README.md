# GoScanner 

【并未发布正式版，暂不可用】

基于 Json文件 、自定义脚本的快速扫描器，用于快速验证目标是否存在该漏洞。

使用场景：

你是否发现一个新漏洞但苦于无法快速编写成脚本？这款工具可能适合你，新手上手写poc速度可以非常快。



# 使用教程

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



# 特性

```md
1. 基于 fofa 规则匹配对应产品,匹配成功后才开始使用POC,避免发送无用包
2. POC 易编写,只需要会看http响应包和http回显包即可
```



# TODO

```md
1. 支持 json/go 这两种POC方法
2. 增强fofa
3. 验证是否利用成功的速度过慢

5. 解决扫描太快,但验证Poc太慢导致提前关闭的问题:remote error: tls: user canceled 
json模式添加 body 正则
```

