package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//////////////////////////////////////////////////
// 🔐 JWT CONFIG
//////////////////////////////////////////////////

var jwtKey = []byte("super-secret-key")

//////////////////////////////////////////////////
// 💬 MESSAGE STRUCT
//////////////////////////////////////////////////

type Message struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

//////////////////////////////////////////////////
// 🌍 GLOBAL VARIABLES
//////////////////////////////////////////////////

var clients = make(map[string]*websocket.Conn)

var mongoClient *mongo.Client
var messageCollection *mongo.Collection

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

//////////////////////////////////////////////////
// 🗄 SAVE MESSAGE TO MONGO
//////////////////////////////////////////////////

func saveMessage(msg Message) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := messageCollection.InsertOne(ctx, msg)
	if err != nil {
		log.Println("Mongo Insert Error:", err)
	}
}

//////////////////////////////////////////////////
// 📜 GET CHAT HISTORY FROM MONGO
//////////////////////////////////////////////////

func getChatHistory(username string) []Message {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := map[string]interface{}{
		"$or": []map[string]interface{}{
			{"sender": username},
			{"receiver": username},
		},
	}

	cursor, err := messageCollection.Find(ctx, filter)
	if err != nil {
		return nil
	}

	var messages []Message
	cursor.All(ctx, &messages)

	return messages
}

//////////////////////////////////////////////////
// 🔥 WEBSOCKET HANDLER
//////////////////////////////////////////////////

func wsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}

	//////////////////////////////////////////////////
	// ✅ STEP 1 — RECEIVE JWT TOKEN
	//////////////////////////////////////////////////

	_, tokenBytes, err := conn.ReadMessage()
	if err != nil {
		log.Println("Token not received")
		conn.Close()
		return
	}

	tokenString := string(tokenBytes)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		log.Println("Invalid Token")
		conn.Close()
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	username := claims["username"].(string)

	log.Println("User Connected:", username)

	// Store connection
	clients[username] = conn

	//////////////////////////////////////////////////
	// ✅ STEP 2 — SEND CHAT HISTORY
	//////////////////////////////////////////////////

	history := getChatHistory(username)

	for _, msg := range history {
		jsonMsg, _ := json.Marshal(msg)
		conn.WriteMessage(websocket.TextMessage, jsonMsg)
	}

	//////////////////////////////////////////////////
	// ✅ STEP 3 — LISTEN FOR NEW MESSAGES
	//////////////////////////////////////////////////

	for {

		_, msgBytes, err := conn.ReadMessage()
		if err != nil {

			delete(clients, username)
			log.Println("User disconnected:", username)
			break
		}

		var msg Message
		json.Unmarshal(msgBytes, &msg)

		msg.Timestamp = time.Now().Format("15:04:05")

		// Save to DB
		saveMessage(msg)

		// 🔥 Send to Receiver ONLY
		if receiverConn, ok := clients[msg.Receiver]; ok {

			jsonMsg, _ := json.Marshal(msg)

			receiverConn.WriteMessage(websocket.TextMessage, jsonMsg)
		}

	}

}

//////////////////////////////////////////////////
// 🚀 MAIN FUNCTION
//////////////////////////////////////////////////

func main() {

	// 🔥 CONNECT TO MONGO
	ctx := context.Background()

	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI("mongodb://mongo:27017"),
	)

	if err != nil {
		log.Fatal("Mongo Connection Failed:", err)
	}

	mongoClient = client
	messageCollection = client.Database("chatdb").Collection("messages")

	http.HandleFunc("/ws", wsHandler)

	log.Println("✅ Secure Chat Service Running on :8002")

	log.Fatal(http.ListenAndServe(":8002", nil))
}
