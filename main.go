package main

import (
	"fmt"
	"encoding/hex"
	// "github.com/gorilla/websocket"
	"log"
	"net/http"
	"flag"
	"crypto/md5"
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

var host = flag.String("h", "127.0.0.1", "http service host")
var port = flag.String("p", "8002", "http service port")
var tcKey = flag.String("k", "", "sign key")

func main() {

	var err error

	// fmt.Println("浏览器访问 http://yourhost:port/chat")
	// http.HandleFunc("/chat", StaticServer)

	flag.Parse()

	hub := newHub()
	go hub.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	
	tcId := "http://" + *host + ":" + *port
	h := md5.New()
	h.Write([]byte(tcId))
	tcId = hex.EncodeToString(h.Sum(nil))

	fmt.Println("tcKey:" + *tcKey)
	http.Handle("/get/wsurl", &ApiServer{hub: hub, apiName: "get-wsurl", tcId: tcId, tcKey: *tcKey})
	http.Handle("/ws/push", &ApiServer{hub: hub, apiName: "ws-push", tcId: tcId, tcKey: *tcKey})

	fmt.Println("listen on port " + *port)

	if err = http.ListenAndServe(":" + *port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}