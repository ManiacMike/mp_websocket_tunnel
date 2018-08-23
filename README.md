## mp_websocket_tunnel 腾讯云小程序web socket信道服务golang版，与后端http服务器配合使用


## 主要实现以下业务

* http接口 /get/wsurl 生成一个tunnel id，与ws connection绑定，返回ws服务器的url
~~~ 
    参数：
        data:
            protocolType: wss
            receiveUrl : https://www.24dota.com/tunnel/
        tcId
        tcKey
        signature

    返回：
        data : {connectUrl : wss://ws.24dota.com, tunnelId: xxx}
        signature : xxx
~~~

* 接受客户端的包文转发到后端服务器
~~~ 
    connect -> 转发
    ping -> 不转发，自己处理
    close -> 转发
    message:{type: "speak", content: {word: "I say something at Thu Aug 09 2018 11:17:00 GMT+0800 (CST)"}} -> 转发

    post请求后端服务器的payload
    {
        data{
            tunnelId
            type : connect|message|close
            content: {"type":"speak","content":{"word":"I say something at Thu Aug 09 2018 11:17:00 GMT+0800 (CST)"}}
        }
        signature
    }
~~~ 


* http接口 /ws/push 根据tunnel id，推送内容给ws connection
~~~ 
    参数：
        data:
            [tunnelIds: xxx
            type : message
            content : hello world]
        tcId
        signature

    返回：
        code : 0
        data:
            invalidTunnelIds : []

~~~ 


## 命令
~~~ 
    -d websocket gateway域名，默认127.0.0.1，与小程序配置中"socket合法域名"一致
    -h 后端服务器域名 ，默认127.0.0.1
    -p websocket gateway端口，默认8002
    -k 腾讯云api的加密秘钥，配置文件中的tcKey字段，用于生成和后端服务器交互的请求签名
    -r 后端服务器接收gateway服务的url
    -e 无用的 tunnel id 过期时间（秒），默认3600秒

    使用示例：mp_websocket_tunnel -r https://www.mydomain.com/tunnel -k xxxxxxxxxx -d ws.mydomain.com
~~~ 
