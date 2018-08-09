#mp_websocket_channel



1. http接口 /get/wsurl 生成一个tunnel id，与ws connection绑定，返回ws服务器的url
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


2. 接受客户端的包文转发到后端服务器
         -> 转发
    ping -> 不转发
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



3. http接口 /ws/push 根据tunnel id，推送内容给ws connection

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

