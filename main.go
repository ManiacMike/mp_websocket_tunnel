package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	// "net/url"
	"github.com/larspensjo/config"
	"strconv"
)

var applications ApplicationGroup
var applications_config ApplicationGroupConfig

type ServiceError struct {
	Msg string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("%s",e.Msg)
}

func Error(msg string) error {
	return &ServiceError{msg}
}


func StaticServer(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "demo/demo.html")
	// staticHandler := http.FileServer(http.Dir("./"))
	// staticHandler.ServeHTTP(w, req)
	return
}

func initServer() error {
	configMap := make(map[string]ApplicationConfig)

	cfg, err := config.ReadDefault("config.ini")
	if err != nil {
		return Error("unable to open config file or wrong fomart")
	}
	sections := cfg.Sections()
	if len(sections) == 0 {
		return Error("no app config")
	}

	for _, section := range sections {
		if section != "DEFAULT" {
			sectionData, _ := cfg.SectionOptions(section)
			tmp := make(map[string]string)
			for _, key := range sectionData {
				value, err := cfg.String(section, key)
				if err == nil {
					tmp[key] = value
				}
			}
			maxClientConn, _ := strconv.Atoi(tmp["MaxClientConn"])
			configMap[section] = ApplicationConfig{tmp["AppId"], tmp["AppSecret"], maxClientConn, tmp["GetConnectApi"], tmp["LoseConnectApi"], tmp["MessageTransferApi"]}
		}
	}
	fmt.Println(configMap)

	valid_config := make(map[string]ApplicationConfig)

	for appid, appconfig := range configMap {
		// if appconfig.TokenMethod != TOKEN_METHOD_GET && appconfig.TokenMethod != TOKEN_METHOD_COOKIE{
		//   return Error("invalid TokenMethod appid: " + appid )
		// }
		if appconfig.MaxClientConn < 1 || appconfig.MaxClientConn > MAX_CLIENT_CONN {
			return Error("invalid MaxClientConn appid: " + appid)
		}
		channelGroup := make(map[string]ChannelService)
		app := Application{channelGroup, appconfig}

		applications[appid] = app
		valid_config[appid] = appconfig
	}
	applications_config = valid_config
	return nil
}

func main() {

	var err error

	http.Handle("/", websocket.Handler(WsServer))
	http.Handle("/api/create-channel", &ApiServer{ApiName: "create-channel"}) //create a ChannelService
	http.Handle("/api/push", &ApiServer{ApiName: "push"})
	http.Handle("/api/broadcast", &ApiServer{ApiName: "broadcast"})
	http.Handle("/api/get-channel", &ApiServer{ApiName: "get-channel"})
	http.Handle("/api/close-channel", &ApiServer{ApiName: "close-channel"}) //close a specific ChannelService
	http.Handle("/api/app-status", &ApiServer{ApiName: "app-status"})       //online num and live connection num

	http.HandleFunc("/demo", StaticServer)

	fmt.Println("listen on port 8002")
	//TODO offer a init commad to reload application info file
	applications = make(ApplicationGroup)
	applications_config = make(ApplicationGroupConfig)

	if err = initServer(); err != nil {
		panic(err.Error())
	}

	if err = http.ListenAndServe(":8002", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}