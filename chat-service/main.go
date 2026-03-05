package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clients = make(map[string]*websocket.Conn)
	mu      sync.Mutex
)

func wsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// 🔥 Step 1 — First message MUST be username
	_, usernameBytes, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		return
	}

	username := string(usernameBytes)

	mu.Lock()
	clients[username] = conn
	mu.Unlock()

	log.Println("User registered:", username)

	// 🔥 Step 2 — Listen for messages
	for {

		_, msgBytes, err := conn.ReadMessage()
		if err != nil {

			mu.Lock()
			delete(clients, username)
			mu.Unlock()

			log.Println("User disconnected:", username)
			break
		}

		var msg Message
		err = json.Unmarshal(msgBytes, &msg)
		if err != nil {
			continue
		}

		mu.Lock()
		receiverConn := clients[msg.Receiver]
		mu.Unlock()

		if receiverConn != nil {

			jsonMsg, _ := json.Marshal(msg)
			receiverConn.WriteMessage(websocket.TextMessage, jsonMsg)

		}
	}

}

func main() {

	http.HandleFunc("/ws", wsHandler)

	log.Println("Private Chat Running on :8002")

	log.Fatal(http.ListenAndServe(":8002", nil))
}
