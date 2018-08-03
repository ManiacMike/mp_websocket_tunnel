package main

import (
	// "fmt"
	//"errors"
	// "github.com/larspensjo/config"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	// "strconv"
)

const MAX_CLIENT_CONN, DEFAULT_CLIENT_CONN = 5, 1

//const TOKEN_METHOD_COOKIE,TOKEN_METHOD_GET = 1,2

//refer to one user with multiple client connections
type ChannelService struct {
	Uid   string
	Token string
	Conns []*websocket.Conn
}

type ApplicationConfig struct {
	AppId     string
	AppSecret string
	//TokenMethod int
	MaxClientConn      int
	GetConnectApi      string
	LoseConnectApi     string
	MessageTransferApi string
}

//refer to one application with multiple channel services
type Application struct {
	Services map[string]ChannelService
	Config   ApplicationConfig
}

type ApplicationGroup map[string]Application
type ApplicationGroupConfig map[string]ApplicationConfig


//remove lost conns
func (this *ApplicationGroup) removeConn(appId, uid string, ws *websocket.Conn) error {
	cs, ok := (*this)[appId].Services[uid]
	if ok == false {
		return nil
	}
	config := applications_config[appId]
	for i, conn := range cs.Conns {
		if ws == conn {
			cs.Conns = append((cs.Conns)[:i], (cs.Conns)[i+1:]...)
			break
		}
	}
	(*this)[appId].Services[uid] = cs
	if len(cs.Conns) == 0 && config.LoseConnectApi != "" {
		go http.PostForm(config.LoseConnectApi, url.Values{"uid": {uid}, "event": {"disconnect"}})
	}
	return nil
}

func (this *ApplicationGroup) removeChannel(appId, uid string) error {
	cs := (*this)[appId].Services[uid]
	for _, conn := range cs.Conns {
		conn.Close()
	}
	delete((*this)[appId].Services, uid)
	return nil
}
