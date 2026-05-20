package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	ChatID int64
}

type Hub struct {
	mu      sync.RWMutex
	clients map[int64]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int64]map[*Client]struct{}),
	}
}

func (h *Hub) AddClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client.ChatID]; !ok {
		h.clients[client.ChatID] = make(map[*Client]struct{})
	}

	h.clients[client.ChatID][client] = struct{}{}
}

func (h *Hub) RemoveClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	chatClients, ok := h.clients[client.ChatID]
	if !ok {
		return
	}

	delete(chatClients, client)

	if len(chatClients) == 0 {
		delete(h.clients, client.ChatID)
	}
}

func (h *Hub) Broadcast(chatID int64, payload []byte) {
	h.mu.RLock()
	chatClients := h.clients[chatID]
	h.mu.RUnlock()

	for client := range chatClients {
		if err := client.Conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			_ = client.Conn.Close()
			h.RemoveClient(client)
		}
	}
}
