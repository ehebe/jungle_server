package ws

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Message struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

type WSContext struct {
	ID   string
	Conn *websocket.Conn
}

type WSHandler struct {
	OnConnect    func(ctx *WSContext)
	OnMessage    func(ctx *WSContext, msg Message)
	OnDisconnect func(ctx *WSContext)
}

type AgentHub struct {
	handlers WSHandler
	clients  map[string]*WSContext
	mu       sync.RWMutex
}

func NewAgentHub() *AgentHub {
	return &AgentHub{
		clients: make(map[string]*WSContext),
	}
}

func (h *AgentHub) OnConnect(fn func(ctx *WSContext)) {
	h.handlers.OnConnect = fn
}

func (h *AgentHub) OnMessage(fn func(ctx *WSContext, msg Message)) {
	h.handlers.OnMessage = fn
}

func (h *AgentHub) OnDisconnect(fn func(ctx *WSContext)) {
	h.handlers.OnDisconnect = fn
}

func (h *AgentHub) Handle() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		id := c.Params("id")
		ctx := &WSContext{ID: id, Conn: c}

		h.mu.Lock()
		h.clients[id] = ctx
		h.mu.Unlock()

		if h.handlers.OnConnect != nil {
			h.handlers.OnConnect(ctx)
		}

		defer func() {
			h.mu.Lock()
			delete(h.clients, id)
			h.mu.Unlock()
			if h.handlers.OnDisconnect != nil {
				h.handlers.OnDisconnect(ctx)
			}
			c.Close()
		}()

		for {
			_, msgData, err := c.ReadMessage()
			if err != nil {
				break
			}

			var msg Message
			if err := json.Unmarshal(msgData, &msg); err != nil {
				continue
			}

			if h.handlers.OnMessage != nil {
				h.handlers.OnMessage(ctx, msg)
			}
		}
	})
}

