package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/royroki/LetsGo/internal/modules/chat/infrastructure/waitingqueue"
)

// WebSocketHub holds all active WebSocket connections for paired users.
type WebSocketHub struct {
	clients map[string]*websocket.Conn // Stores user connections by userID
	queue   *waitingqueue.WaitingQueue // Waiting queue for unpaired users
}

// NewWebSocketHub initializes a new WebSocketHub.
func NewWebSocketHub(wq *waitingqueue.WaitingQueue) *WebSocketHub {
	return &WebSocketHub{
		clients: make(map[string]*websocket.Conn),
		queue:   wq,
	}
}

// HandleWebSocket manages WebSocket connections, pairing users, and messaging.
func (hub *WebSocketHub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins
	}

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading WebSocket:", err)
		return
	}
	defer conn.Close()

	// Generate or get user ID (You can pass userID as a query param or generate it)
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		userID = generateUniqueUserID() // Implement generateUniqueUserID()
	}

	// Add the user to the queue
	log.Printf("User %s added to queue", userID)
	if err := hub.queue.AddUserToQueue(userID); err != nil {
		log.Println("Error adding user to queue:", err)
		return
	}

	// Try to pair this user with another one from the queue
	pairedUserID, err := hub.queue.GetNextUser(userID)
	if err != nil {
		log.Println("Error getting next user from queue:", err)
		return
	}

	// If the current user is paired with someone, start the chat
	if pairedUserID != userID {
		// Pair them up
		log.Printf("Pairing users: %s and %s", userID, pairedUserID)
		hub.clients[userID] = conn // Store the connection
		hub.clients[pairedUserID].WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("You are paired with %s", userID)))
		hub.clients[userID].WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("You are paired with %s", pairedUserID)))

		// Now listen for messages from both clients
		go hub.listenForMessages(userID, conn)
		go hub.listenForMessages(pairedUserID, hub.clients[pairedUserID])
	} else {
		// If no pair is found, keep the user in the queue for now
		log.Printf("Only one user found, waiting for another user to connect: %s", userID)
	}
}

// listenForMessages listens for incoming messages from the WebSocket connection.
func (hub *WebSocketHub) listenForMessages(userID string, conn *websocket.Conn) {
	for {
		// Read incoming message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %s: %v", userID, err)
			delete(hub.clients, userID) // Remove from active clients if disconnected
			break
		}

		// Broadcast the message to the paired user
		pairedUserID, err := hub.queue.GetNextUser(userID) // Get paired user from the queue
		if err != nil {
			log.Printf("Error retrieving paired user for %s: %v", userID, err)
			continue
		}

		// Send message to paired user
		if err := hub.clients[pairedUserID].WriteMessage(messageType, message); err != nil {
			log.Printf("Error sending message to paired user %s: %v", pairedUserID, err)
		}
	}
}

func generateUniqueUserID() string {
	return uuid.New().String()
}
