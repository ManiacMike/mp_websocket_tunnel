package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"strings"
)

//handle api request from api
type ApiServer struct {
	ApiName string
	AppId   string
}

//TODO create a error standard

func (this *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := this.CheckParams(r); err != nil {
		returnMsg := fmt.Sprintf("{\"code\":400,\"msg\":\"%v\",\"time\":%v}", err.Error(), time.Now().Unix())
		fmt.Fprint(w, returnMsg)
		return
	}

	switch this.ApiName {
	case "create-channel":
		this.CreateChannel(w, r)
	case "push":
		this.Push(w, r)
	case "broadcast":
		this.Broadcast(w, r)
	case "get-channel":
		this.GetChannel(w, r)
	case "close-channel":
		this.CloseChannel(w, r)
	case "app-status":
		this.AppStatus(w, r)
	default:
		fmt.Fprint(w, "Invalid api")
	}
}

func (this *ApiServer) CheckParams(r *http.Request) error {
	appId := r.PostFormValue("app_id")
	if appId == "" {
		return Error("app_id missing")
	}
	appSecret := r.PostFormValue("app_secret")
	if appSecret == "" {
		return Error("app_secret missing")
	}
	config, ok := applications_config[appId]
	if ok == false {
		return Error("app_id invalid")
	}
	if config.AppSecret != appSecret {
		return Error("app_secret invalid")
	}
	this.AppId = appId
	return nil
}

func (this *ApiServer) CreateChannel(w http.ResponseWriter, r *http.Request) error {
	uid := r.PostFormValue("uid")
	if uid == "" {
		return Error("uid missing")
	}
	appId := this.AppId
	c, ok := applications[appId].Services[uid]
	if ok == true {
		msg := fmt.Sprintf("{\"uid\":\"%v\",\"token\":\"%v\"}", c.Uid, c.Token)
		this.Success(msg, w)
		return nil
	}

	token := GenerateId()
	fmt.Println("token: ", token)
	channelService := ChannelService{Uid: uid, Token: token}
	applications[appId].Services[uid] = channelService
	msg := fmt.Sprintf("{\"uid\":\"%v\",\"token\":\"%v\"}", channelService.Uid, channelService.Token)
	this.Success(msg, w)
	return nil
}

func (this *ApiServer) Push(w http.ResponseWriter, r *http.Request) error {
	uidStr := r.PostFormValue("uid")
	if uidStr == "" {
		return Error("uid missing")
	}
	appId := this.AppId
	uidSlice := strings.Split(uidStr, ",")
	for _,uid := range uidSlice{
		channelService, ok := applications[appId].Services[uid]
		if ok == false {
			continue
		}
		msg := r.PostFormValue("message")
		for _, conn := range channelService.Conns {
			if err := websocket.Message.Send(conn, msg); err != nil {
				applications.removeConn(appId, uid, conn)
			}
		}
	}
	this.Success("", w)
	return nil
}

func (this *ApiServer) Broadcast(w http.ResponseWriter, r *http.Request) error {
	appId := this.AppId
	channelServices := applications[appId]
	msg := r.PostFormValue("message")
	for uid, cs := range channelServices.Services {
		for _, conn := range cs.Conns {
			if err := websocket.Message.Send(conn, msg); err != nil {
				applications.removeConn(appId, uid, conn)
			}
		}
	}
	this.Success("", w)
	return nil
}

func (this *ApiServer) GetChannel(w http.ResponseWriter, r *http.Request) error {
	uid := r.PostFormValue("uid")
	if uid == "" {
		return Error("uid missing")
	}
	appId := this.AppId
	channelService, ok := applications[appId].Services[uid]
	if ok == false {
		return Error("channel not created")
	}
	msg := fmt.Sprintf("{\"uid\":\"%v\",\"token\":\"%v\",\"connections\":%d}", channelService.Uid, channelService.Token, len(channelService.Conns))
	this.Success(msg, w)
	return nil
}

func (this *ApiServer) CloseChannel(w http.ResponseWriter, r *http.Request) error {
	uid := r.PostFormValue("uid")
	if uid == "" {
		return Error("uid missing")
	}
	appId := this.AppId
	_, ok := applications[appId].Services[uid]
	if ok == false {
		return Error("uid invalid")
	}
	applications.removeChannel(appId, uid)
	this.Success("", w)
	return nil
}

func (this *ApiServer) AppStatus(w http.ResponseWriter, r *http.Request) error {
	appId := this.AppId
	channelNum := len(applications[appId].Services)
	msg := fmt.Sprintf("{\"channelNum\":%d}", channelNum)
	this.Success(msg, w)
	return nil
}

func (this *ApiServer) Success(msg string, w http.ResponseWriter) {
	if msg == "" {
		msg = "\"success\""
	}
	returnMsg := fmt.Sprintf("{\"code\":0,\"result\":%v}", msg)
	fmt.Fprint(w, returnMsg)
}
