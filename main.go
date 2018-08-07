package main

import (
	"fmt"
	// "github.com/gorilla/websocket"
	"log"
	"net/http"
	// "net/url"
	// "github.com/larspensjo/config"
	// "strconv"
)

// var applications ApplicationGroup
// var applications_config ApplicationGroupConfig

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
	http.ServeFile(w, req, "home.html")
	return
}

func main() {

	var err error

	fmt.Println("浏览器访问 http://yourhost:port/chat")
	http.HandleFunc("/chat", StaticServer)

	hub := newHub()
	go hub.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	
	tcId := "id_test"
	tcKey := "key_test"

	http.Handle("/get/wsurl", &ApiServer{apiName: "get-wsurl", tcId: tcId, tcKey: tcKey})
	http.Handle("/ws/push", &ApiServer{apiName: "ws-push", tcId: tcId, tcKey: tcKey})

	fmt.Println("listen on port 8002")
	//TODO offer a init commad to reload application info file
	// applications = make(ApplicationGroup)
	// applications_config = make(ApplicationGroupConfig)

	// if err = initServer(); err != nil {
	// 	panic(err.Error())
	// }

	if err = http.ListenAndServe(":8002", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}