package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)

func main() {

	router := gin.Default()

	router.GET("/ws", handleConnections)

	router.Run(":8002")
}

func handleConnections(c *gin.Context) {

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Println(err)
		return
	}

	defer ws.Close()

	clients[ws] = true

	for {

		_, msg, err := ws.ReadMessage()

		if err != nil {
			delete(clients, ws)
			break
		}

		broadcast(msg)
	}
}

func broadcast(message []byte) {

	for client := range clients {

		err := client.WriteMessage(websocket.TextMessage, message)

		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}
