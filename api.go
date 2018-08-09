package main

import (
	"fmt"
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
	// "crypto/sha1"
	// "encoding/hex"
	// "strings"
)


type ApiServer struct {
	hub *Hub
	apiName string
	tcId    string
	tcKey   string
}

type ApiParams struct{
	data string
	tcId string
	tcKey string
	signature string
}

func (this *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err, apiParams := this.CheckParams(r);

	if  err != nil {
		returnMsg := fmt.Sprintf("{\"code\":400,\"msg\":\"%v\",\"time\":%v}", err.Error(), time.Now().Unix())
		fmt.Fprint(w, returnMsg)
		return
	}

	switch this.apiName {
	case "get-wsurl":
		err = this.GetWsurl(w, apiParams)
	case "ws-push":
		err = this.WsPush(w, apiParams)
	default:
		fmt.Fprint(w, "Invalid api")
	}

	if  err != nil {
		returnMsg := fmt.Sprintf("{\"code\":400,\"msg\":\"%v\",\"time\":%v}", err.Error(), time.Now().Unix())
		fmt.Fprint(w, returnMsg)
		return
	}
}

func (this *ApiServer) CheckParams(r *http.Request) (error, *ApiParams) {
	result, _:= ioutil.ReadAll(r.Body)
	r.Body.Close()
	fmt.Println(result)

	var f interface{}
	json.Unmarshal(result, &f) 
	m := f.(map[string]interface{})

	if m["tcId"] == nil || m["tcId"].(string) == ""{
		return Error("tcId missing"), nil
	}
	tcId := m["tcId"].(string)
	if m["signature"] == nil || m["signature"].(string) == ""{
		return Error("signature missing"), nil
	}
	signature := m["signature"].(string)

	tcKey := ""
	if m["tcKey"] != nil{
		tcKey = m["tcKey"].(string) 
	}

	data := ""
	if m["data"] != nil{
		data = m["data"].(string)
	}

	signatureCompute := sha1Encode(data + this.tcKey)

	if signatureCompute != signature{
		return Error("signature error"), nil
	}

	apiParams := &ApiParams{tcId: tcId, signature: signature, tcKey: tcKey, data: data}
	return nil, apiParams
}


func (this *ApiServer) GetWsurl(w http.ResponseWriter, r *ApiParams) error {

	dataNode := JsonDecode(r.data)
	data := dataNode.(map[string]interface{})
	protocol := data["protocolType"].(string)
	// receiveUrl := data["receiveUrl"].(string)
	tunnelId := GenerateUnixNanoId()
	fmt.Println("tunnelId: ", tunnelId)
	url := fmt.Sprintf("%v://"+ *wsdomain +"/?tunnelId=%v", protocol, tunnelId)
	this.hub.addTunnelId(tunnelId)

	returnDataMap := map[string]string{"tunnelId": tunnelId, "connectUrl": url}
	returnData := JsonEncode(returnDataMap)
	result := map[string]interface{}{"code": 0, "data": returnData, "signature": sha1Encode(returnData + this.tcKey)}
	this.Success(result, w)
	return nil
}

func (this *ApiServer) WsPush(w http.ResponseWriter, r *ApiParams) error {
	dataNode := JsonDecode(r.data)
	data := dataNode.([]interface{})
	fmt.Println("debug WsPush \n")
	invalidTunnelIds := []string{}
	for _, v := range data {
		vv := v.(map[string]interface{})
		if vv["type"].(string) == "message"{
			tunnelIds := vv["tunnelIds"].([]string)
			for _, tunnelId := range tunnelIds {
				err := this.hub.sendByTunnelId(tunnelId, "message:" + vv["content"].(string))
				if err != nil{
					invalidTunnelIds = append(invalidTunnelIds, tunnelId)
				}
			}
		}
	}
	
	result := map[string]interface{}{"code": 0, "data": map[string]interface{}{"invalidTunnelIds": invalidTunnelIds}}
	this.Success(result, w)
	return nil
}



func (this *ApiServer) Success(result interface{}, w http.ResponseWriter) {
	returnMsg := JsonEncode(result)
	fmt.Fprint(w, returnMsg)
}
