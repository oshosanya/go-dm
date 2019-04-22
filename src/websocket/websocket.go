package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/oshosanya/go-dm/src/data"
)

// Define our message object
type Message struct {
	Action string `json:"action"`
	Args   string `json:"args"`
}

type ResponseMessage struct {
	Action  string `json:"action"`
	Payload string `json:"payload"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[*websocket.Conn]bool) // connected clients
var inbound = make(chan Message)
var outbound = make(chan ResponseMessage)
var responseData string = ""
var downloads []data.Download

func Start() {
	// http.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	conn, _ := upgrader.Upgrade(w, r, nil)
	// 	go func(conn *websocket.Conn) {
	// 		for {
	// 			mType, msg, err := conn.ReadMessage()
	// 			if err != nil {
	// 				conn.Close()
	// 			}
	// 			println(msg)
	// 			conn.WriteMessage(mType, msg)
	// 		}
	// 	}(conn)
	// })
	http.HandleFunc("/v1/ws", handleConnections)
	go handleMessages()
	http.ListenAndServe(":3000", nil)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	allDownloads := data.GetAllDownloads()
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(allDownloads)
	if err != nil {
		log.Printf("Could not serialize message: %v", err)
	}
	response := ResponseMessage{
		Action:  "all",
		Payload: string(b),
	}
	ws.WriteJSON(response)
	// Make sure we close the connection when the function returns
	defer ws.Close()
	clients[ws] = true

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		fmt.Printf("%+v\n", msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		inbound <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-inbound
		// Loop through inbound and handle actions
		for _ = range clients {
			callDynamically(msg)
			// err := client.WriteJSON(responseData)
			// println("response")
			// println(responseData)
			// err := client.WriteMessage(websocket.TextMessage, []byte(responseData))
			// if err != nil {
			// 	log.Printf("error: %v", err)
			// 	client.Close()
			// 	delete(clients, client)
			// }
		}
	}
}

func SendMessageToClient(message interface{}) error {
	// var raw map[string]interface{}

	b, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("Could not serialize message: %v", err)
	}
	response := ResponseMessage{
		Action:  "update",
		Payload: string(b),
	}
	for client := range clients {
		// err := client.WriteMessage(websocket.TextMessage, []byte())
		err := client.WriteJSON(response)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
	return nil
}
