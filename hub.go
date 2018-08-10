package main

import (
	// "encoding/json"
	"fmt"
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

func (h *Hub) clearExpiredTunnelId(){
	for tunnelIdStr, tunnelId := range h.tunnelIdPool{
		if tunnelId.active == false && (time.Now().Unix() - tunnelId.lastActiveTime > int64(*tunnelIdExpire)){
			delete(h.tunnelIdPool, tunnelIdStr)
		}
	}
}

func (h *Hub) sendByTunnelId(tunnelId string, message string) error{
	if h.tunnelIdPool[tunnelId] == nil || h.tunnelIdPool[tunnelId].active == false{
		return Error("invalid tunnelId")
	}
	h.clients[tunnelId].messageSendChan <- []byte(message)
	return nil
}

func (h *Hub) run() {
	clearExpiredTunnelIdInterval := time.Second * 60
	ticker := time.NewTimer(clearExpiredTunnelIdInterval)

	for {
		select {
		case client := <-h.register:
			h.clients[client.tunnelId] = client
			h.tunnelIdPool[client.tunnelId].active = true

		case client := <-h.unregister:
			tunnelId := client.tunnelId
			if _, ok := h.clients[tunnelId]; ok {
				fmt.Println("unregister")
				h.tunnelIdPool[tunnelId].active = false
				h.tunnelIdPool[tunnelId].lastActiveTime = time.Now().Unix()
				client.postToServerChan <- map[string]string{"packetType": "close"}
				delete(h.clients, tunnelId)
				close(client.messageSendChan)
				close(client.postToServerChan)
			}
		case message := <-h.broadcast:
			for _,client := range h.clients {
				select {
				case client.messageSendChan <- message:
				default:
					tunnelId := client.tunnelId
					h.tunnelIdPool[tunnelId].active = false
					h.tunnelIdPool[tunnelId].lastActiveTime = time.Now().Unix()
					delete(h.clients, tunnelId)
					close(client.messageSendChan)
					close(client.postToServerChan)
				}
			}
		case <-ticker.C:
			h.clearExpiredTunnelId()
			ticker.Reset(clearExpiredTunnelIdInterval)
		}
	}
}
