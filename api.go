package main

import (
	"fmt"
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
	// "strings"
)


type ApiServer struct {
	apiName string
	tcId    string
	tcKey   string
}

// type ApiRequest struct{
// 	data string
// 	tcId string
// 	tcKey string
// 	signature string
// }

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
	result, _:= ioutil.ReadAll(r.Body)
	r.Body.Close()
	fmt.Println(result)


	var f interface{}
	json.Unmarshal(result, &f) 
	m := f.(map[string]interface{})

	fmt.Println(m)

	if m["tcId"] == nil || m["tcId"].(string) == ""{
		return Error("tcId missing")
	}
	tcId := m["tcId"].(string)
	if m["signature"] == nil || m["signature"].(string) == ""{
		return Error("signature missing")
	}
	fmt.Println(tcId)
	signature := m["signature"].(string)
	fmt.Println(signature)
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
