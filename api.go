package main

import (
	"fmt"
	"net/http"
	"time"
	// "strings"
)


type ApiServer struct {
	apiName string
	tcId    string
	tcKey   string
}

func (this *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := this.CheckParams(r); err != nil {
		returnMsg := fmt.Sprintf("{\"code\":400,\"msg\":\"%v\",\"time\":%v}", err.Error(), time.Now().Unix())
		fmt.Fprint(w, returnMsg)
		return
	}

	switch this.apiName {
	case "get-wsurl":
		this.GetWsurl(w, r)
	case "ws-push":
		this.WsPush(w, r)
	default:
		fmt.Fprint(w, "Invalid api")
	}
}

func (this *ApiServer) CheckParams(r *http.Request) error {
	tcId := r.PostFormValue("tcId")
	if tcId == "" {
		return Error("tcId missing")
	}
	signature := r.PostFormValue("signature")
	if signature == "" {
		return Error("signature missing")
	}
	return nil
}


func (this *ApiServer) GetWsurl(w http.ResponseWriter, r *http.Request) error {
	tcKey := r.PostFormValue("tcKey")
	if tcKey == "" {
		return Error("tcKey missing")
	}
	dataNode := JsonDecode(r.PostFormValue("data"))
	data := dataNode.(map[string]interface{})
	protocol := data["protocol"].(string)
	// receiveUrl := data["receiveUrl"].(string)
	token := GenerateUnixNanoId()
	fmt.Println("token: ", token)
	url := fmt.Sprintf("\"%v://ws.24dota.com/?token=%v\"", protocol, token)
	// channelService := ChannelService{Uid: uid, Token: token}
	// applications[appId].Services[uid] = channelService
	// msg := fmt.Sprintf("{\"uid\":\"%v\",\"token\":\"%v\"}", channelService.Uid, channelService.Token)
	this.Success(url, w)
	return nil
}

func (this *ApiServer) WsPush(w http.ResponseWriter, r *http.Request) error {

	return nil
}



func (this *ApiServer) Success(msg string, w http.ResponseWriter) {
	if msg == "" {
		msg = "\"success\""
	}
	returnMsg := fmt.Sprintf("{\"code\":0,\"result\":%v}", msg)
	fmt.Fprint(w, returnMsg)
}
