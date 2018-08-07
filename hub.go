// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.

type TunnelId struct {

	//used by any clients
	active bool

	//create time
	createTime uint

	//updated when disconnected or create
	lastActiveTime uint
}

type Hub struct {
	tunnelIdPool map[string]*TunnelId

	// Registered clients.
	clients map[*Client]bool

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
		clients:    make(map[*Client]bool),
		tunnelIdPool: make(map[string]*TunnelId),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			welcomeMessage := []byte("hello")
			client.send <- welcomeMessage
			h.clients[client] = true
			userCountMessage := map[string]interface{}{
				"type":       "user_count",
				"user_count": len(h.clients),
			}
			userCountMessagebody, _ := json.Marshal(userCountMessage)
			fmt.Println(userCountMessage)
			for client := range h.clients {
				select {
				case client.send <- userCountMessagebody:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			userCountMessage := map[string]interface{}{
				"type":       "user_count",
				"user_count": len(h.clients),
			}
			userCountMessagebody, _ := json.Marshal(userCountMessage)
			fmt.Println(userCountMessage)
			for client := range h.clients {
				select {
				case client.send <- userCountMessagebody:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
