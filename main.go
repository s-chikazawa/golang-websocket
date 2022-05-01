package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	fmt.Println("Start Websocket")
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", websocketHandler)
	go messageHandler()

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

var upgrader = websocket.Upgrader{}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()
	clients[conn] = true

	for {
		var message Message
		err = conn.ReadJSON(&message)
		log.Print(message)
		if err != nil {
			log.Print(err)
			break
		}
		broadcast <- message
	}
}

func messageHandler() {
	for {
		message := <-broadcast
		for client := range clients {
			err := client.WriteJSON(message)
			log.Print(err)
			if err != nil {
				log.Print(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
