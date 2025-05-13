package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/royroki/services/chatting-service/internal/constants"
	"github.com/royroki/services/chatting-service/internal/domain/entity"
	"github.com/sony/sonyflake"
)

type WCServer struct {
	clients   map[int64]*entity.User
	paired    map[int64]int64
	mu        sync.Mutex
	sonyFlake *sonyflake.Sonyflake
}

// NewWCServer initializes the WebSocket server with Sonyflake for generating unique IDs
func NewWCServer() *WCServer {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	if sf == nil {
		log.Fatal("Failed to initialize Sonyflake")
	}
	return &WCServer{
		clients:   make(map[int64]*entity.User),
		paired:    make(map[int64]int64),
		sonyFlake: sf,
	}
}

// Upgrade upgrades the HTTP connection to WebSocket and generates a unique user ID
func (s *WCServer) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, int64, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins, but we should secure this in production
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, 0, err
	}

	// Generate a unique user ID using Sonyflake
	id, err := s.sonyFlake.NextID()
	if err != nil {
		conn.Close()
		return nil, 0, err
	}

	return conn, int64(id), nil
}

// Register adds a new user to the active clients map
func (s *WCServer) Register(user *entity.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[user.ID] = user
}

// IsActive checks if a user is still connected by looking up in the clients map
func (s *WCServer) IsActive(userID int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.clients[userID]
	return exists
}

// Get User Entity by Id
func (s *WCServer) GetUser(userID int64) *entity.User {
	s.mu.Lock()
	defer s.mu.Unlock()
	user := s.clients[userID]
	return user
}

// Get Partner ID
func (s *WCServer) GetPartnerID(userID int64) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	partnerID, ok := s.paired[userID]
	if !ok {
		return 0
	}
	return partnerID
}

// Send sends a message to a specific user by user ID
func (s *WCServer) Send(userID int64, message entity.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if user, ok := s.clients[userID]; ok {
		err := user.Conn.WriteJSON(message)
		if err != nil {
			log.Printf("Failed to send message to user %d: %v", userID, err)
			user.Conn.Close()
			delete(s.clients, userID)
		}
	}
}

// Pair connects two users by updating the paired map
func (s *WCServer) Pair(userA, userB int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.paired[userA] = userB
	s.paired[userB] = userA
}

// Unregister removes a user from the clients and paired maps
func (s *WCServer) Unregister(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Close the connection if it's still open
	if user, ok := s.clients[userID]; ok {
		_ = user.Conn.Close()
		delete(s.clients, userID)
	}

	// Remove pairing if exists
	if partnerID, ok := s.paired[userID]; ok {
		delete(s.paired, userID)
		delete(s.paired, partnerID)

		// Optionally notify the partner
		if partner, exists := s.clients[partnerID]; exists {
			partner.Conn.WriteJSON(entity.Message{
				From:    0,
				Content: "Your partner has disconnected. Re-matching...",
				Status:  constants.PartnerDisconnected,
			})
		}
	}
}
