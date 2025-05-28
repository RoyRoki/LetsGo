package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/royroki/services/chatting-service/internal/constants"
	"github.com/royroki/services/chatting-service/internal/domain/entity"
	"github.com/sony/sonyflake"
)

type WCServer struct {
	clients    map[int32]*entity.User
	paired     map[int32]int32
	mu         sync.Mutex
	sonyFlake  *sonyflake.Sonyflake
	OnMessage  func(user *entity.User, tags []string, module string) // callback
}


func NewWCServer() *WCServer {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	if sf == nil {
		log.Fatal("Failed to initialize Sonyflake")
	}
	return &WCServer{
		clients:   make(map[int32]*entity.User),
		paired:    make(map[int32]int32),
		sonyFlake: sf,
	}
}


func (s *WCServer) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, int32, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, 0, err
	}

	id, err := s.sonyFlake.NextID()
	if err != nil {
		conn.Close()
		return nil, 0, err
	}
	safeID := int32(id & 0x7FFFFFFF)

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}()

	return conn, safeID, nil
}

func (s *WCServer) Register(user *entity.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[user.ID] = user
}

func (s *WCServer) IsActive(userID int32) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.clients[userID]
	return exists
}

func (s *WCServer) GetUser(userID int32) *entity.User {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clients[userID]
}

func (s *WCServer) GetPartnerID(userID int32) int32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.paired[userID]
}

func (s *WCServer) Send(userID int32, msg entity.Message) {
	s.mu.Lock()
	user, ok := s.clients[userID]
	s.mu.Unlock()
	if !ok || user.Conn == nil {
		log.Printf("Send failed: user %d not connected", userID)
		return
	}

	if err := user.Conn.WriteJSON(msg); err != nil {
		log.Printf("Send failed to user %d: %v. Retrying...", userID, err)
		if err := user.Conn.WriteJSON(msg); err != nil {
			log.Printf("Retry failed to user %d: %v", userID, err)
		}
	}
}

func (s *WCServer) Pair(userA, userB int32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.paired[userA] = userB
	s.paired[userB] = userA
}

func (s *WCServer) Unregister(userID int32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if partnerID, ok := s.paired[userID]; ok {
		delete(s.paired, userID)
		delete(s.paired, partnerID)
		if partner, exists := s.clients[partnerID]; exists && partner.Conn != nil {
			_ = partner.Conn.WriteJSON(entity.Message{
				From:    0,
				Content: "Your partner has disconnected. Re-matching...",
				Status:  constants.PartnerDisconnected,
			})
		}
	}

	if user, ok := s.clients[userID]; ok {
		if user.Conn != nil {
			_ = user.Conn.Close()
		}
		delete(s.clients, userID)
	}
}

// HandleMessage reads, unmarshals, and processes a single WebSocket message
func (s *WCServer) HandleMessage(user *entity.User, msgBytes []byte) {
	var chatReq struct {
		Tags   []string `json:"tags"`
		Module string   `json:"module"`
	}

	if err := json.Unmarshal(msgBytes, &chatReq); err != nil {
		log.Printf("Failed to unmarshal chat message from user %d: %v", user.ID, err)
		return
	}

	if chatReq.Tags != nil {
		s.mu.Lock()
		tagSet := make(map[string]struct{})
		for _, tag := range user.Tags {
			tagSet[tag] = struct{}{}
		}
		for _, tag := range chatReq.Tags {
			tagSet[tag] = struct{}{}
		}

		merged := make([]string, 0, len(tagSet))
		for tag := range tagSet {
			merged = append(merged, tag)
		}
		user.Tags = merged
		user.Module = chatReq.Module
		s.mu.Unlock()

		log.Printf("User %d tags: %v, module: %s", user.ID, user.Tags, user.Module)
		go s.OnMessage(user, user.Tags, user.Module)
	}
}
	