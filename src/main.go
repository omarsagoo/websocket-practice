package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	allUsers  AllUsers
	clients   = make(map[*websocket.Conn]bool) // client connections map
	broadcast = make(chan Message)             // Broadcast Channel
	// configure the upgrader
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type User struct {
	Username string `json:"user"`
	Avatar   string `json:"avatar"`
}

// AllUsers holds all the users in a string array
type AllUsers struct {
	Users []User `json:"users"`
}

// Message holds the values for each message
type Message struct {
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// upgrade Get request to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("ConnectionHandler:", err)
	}

	// close when this function ends
	defer ws.Close()

	// add client to global clients map
	clients[ws] = true

	for {
		var msg Message

		// read new msg in as JSON, and initialize it as a Message object
		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("ConnectionHandlerMsg: %v", err)

			if websocket.IsCloseError(err, 1001) {
				broadcast <- Message{Username: "ATTENTION!", Avatar: "red-alert.png", Message: "Someone left the chat."}
			}

			delete(clients, ws)
			break
		}
		// send the msg to the broadcast channel, to be handled by handleMessage goroutine
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		if msg.Message == "" {
			usr := User{Username: msg.Username, Avatar: msg.Avatar}
			// users := append(allUsers.Users, usr)
			msg.Message = fmt.Sprintf("NEW USER JOINED!!!!! WELCOME %s", usr.Username)
			msg.Username = "ATTENTION!"
			msg.Avatar = "red-alert.png"
			// fmt.Println(users)
		}
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// create file server
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// start listening for messages
	go handleMessages()

	log.Println("http server started on port :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
