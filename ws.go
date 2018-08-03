package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	// "log"
	"net/http"
	"net/url"
	// "time"
)

func WsServer(ws *websocket.Conn) {
	var err error
	var appId, uid string
	if appId, uid, err = acceptClientToken(ws); err != nil {
		errMsg := err.Error()
		websocket.Message.Send(ws, errMsg)
		ws.Close()
		return
	}
	config := applications_config[appId]
	for {
		var receiveMsg string

		if err = websocket.Message.Receive(ws, &receiveMsg); err != nil {
			applications.removeConn(appId, uid, ws)
			break
		}
		if config.MessageTransferApi != "" {
			fmt.Println(config.MessageTransferApi,receiveMsg)
			go http.PostForm(config.MessageTransferApi, url.Values{"uid": {uid}, "data": {receiveMsg}, "event": {"message"}})
		}
	}
}


func acceptClientToken(ws *websocket.Conn) (string, string, error) {
	appId := ws.Request().FormValue("app_id")
	if appId == "" {
		return "", "", Error("app_id missing")
	}
	config, ok := applications_config[appId]
	if ok == false {
		return "", "", Error("app_id invalid")
	}
	uid := ws.Request().FormValue("uid")
	if uid == "" {
		return "", "", Error("uid missing")
	}
	channelService, ok := applications[appId].Services[uid]
	if ok == false {
		return "", "", Error("uid invalid")
	}
	token := ws.Request().FormValue("token")
	if token == "" {
		return "", "", Error("token missing")
	}
	if token != channelService.Token {
		return "", "", Error("invalid token")
	}
	conns := channelService.Conns
	if len(conns) > (config.MaxClientConn - 1) {
		//close the first conn
		conns[0].Close()
		conns = conns[1:]
	}
	conns = append(conns, ws)
	if len(conns) == 1 && config.GetConnectApi != "" {
		go http.PostForm(config.GetConnectApi, url.Values{"uid": {uid}, "event": {"connect"}})
	}
	channelService.Conns = conns
	applications[appId].Services[uid] = channelService
	return appId, uid, nil
}
