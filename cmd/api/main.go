package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Define the upgrader for WebSocket connection (HTTP to WebSocket)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (you may want to restrict this in production)
	},
}

// Store all active connections (you can implement a better mechanism later)
var clients = make(map[*websocket.Conn]bool)

// Function to handle incoming WebSocket connections
func handleConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade the connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Add the new connection to the clients map
	clients[conn] = true

	for {
		// Read incoming message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		// Broadcast the message to all connected clients
		for c := range clients {
			if err := c.WriteMessage(messageType, p); err != nil {
				fmt.Println("Error writing message:", err)
				c.Close()
				delete(clients, c)
			}
		}
	}
}

func main() {
	http.HandleFunc("/chat", handleConnection) // Handle connections to the /chat endpoint

	// Start the HTTP server
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
