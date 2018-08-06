#mp_websocket_channel


1./get/wsurl 生成一个tunnel id，与ws connection绑定，返回ws服务器的url
    参数：
        data:
            protocolType: wss
            receiveUrl : https://www.24dota.com/tunnel/
        tcId
        tcKey
        signature

    返回：
        wss://ws.24dota.com


2. /ws/push 根据tunnel id，推送内容给ws connection

    参数：
        data:
            tunnelIds: wss
            receiveUrl : https://www.24dota.com/tunnel/
        tcId
        signature

    返回：
        code : 0
        data:
            invalidTunnelIds : []

