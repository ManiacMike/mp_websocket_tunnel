package main

import (
	// "encoding/json"
	// "fmt"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.

type TunnelId struct {

	//used by any clients
	active bool

	//created time
	createTime int64

	//updated when disconnected or created
	lastActiveTime int64
}

type Hub struct {
	tunnelIdPool map[string]*TunnelId

	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
		tunnelIdPool: make(map[string]*TunnelId),
	}
}

func (h *Hub) addTunnelId(tunnelIdStr string){
	tunnelId := &TunnelId{
		active:  false,
		createTime: time.Now().Unix(),
		lastActiveTime: time.Now().Unix(),
	}
	h.tunnelIdPool[tunnelIdStr] = tunnelId
}

func (h *Hub) checkTunnelId(tunnelId string) int {
	if h.tunnelIdPool[tunnelId] == nil{
		return 0
	}else if h.tunnelIdPool[tunnelId].active == false{
		return 1
	}else{
		return 2
	}

}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// welcomeMessage := []byte("hello")
			// client.send <- welcomeMessage
			h.clients[client.tunnelId] = client
			h.tunnelIdPool[client.tunnelId].active = true
			// userCountMessage := map[string]interface{}{
			// 	"type":       "user_count",
			// 	"user_count": len(h.clients),
			// }
			// userCountMessagebody, _ := json.Marshal(userCountMessage)
			// fmt.Println(userCountMessage)
			// for client := range h.clients {
			// 	select {
			// 	case client.send <- userCountMessagebody:
			// 	default:
			// 		close(client.send)
			// 		delete(h.clients, client)
			// 	}
			// }
		case client := <-h.unregister:
			tunnelId := client.tunnelId
			if _, ok := h.clients[tunnelId]; ok {
				h.tunnelIdPool[tunnelId].active = false
				delete(h.clients, tunnelId)
				close(client.send)
			}
			// userCountMessage := map[string]interface{}{
			// 	"type":       "user_count",
			// 	"user_count": len(h.clients),
			// }
			// userCountMessagebody, _ := json.Marshal(userCountMessage)
			// fmt.Println(userCountMessage)
			// for client := range h.clients {
			// 	select {
			// 	case client.send <- userCountMessagebody:
			// 	default:
			// 		close(client.send)
			// 		delete(h.clients, client)
			// 	}
			// }
		case message := <-h.broadcast:
			for _,client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.tunnelId)
				}
			}
		}
	}
}
