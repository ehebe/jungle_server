package agent

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ehebe/jungle/internal/collector"
	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string                 `json:"type"`    // command / heartbeat / register / response
	Payload map[string]interface{} `json:"payload"` // depends on type
}

const (
	agentID = "agent-01"
	wsURL   = "ws://localhost:8080/ws/" + agentID
)

func Start() {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatalf("Failed to connect to controller: %v", err)
	}
	defer conn.Close()
	log.Printf("[%s] connected to controller", agentID)

	register := Message{
		Type: "register",
		Payload: map[string]interface{}{
			"version": "v1.0.0",
			"arch":    "x86_64",
		},
	}
	send(conn, register)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			send(conn, Message{Type: "heartbeat", Payload: map[string]interface{}{
				"sysinfo": collector.Collect(),
			}})
		}
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Connection lost: %v", err)
			return
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("Invalid message: %v", err)
			continue
		}

		switch msg.Type {
		case "command":
			cmd := msg.Payload["command"].(string)
			log.Printf("Received command: %s", cmd)

			var result string
			switch cmd {
			case "collect_stats":
				result = collectMockStats()
			default:
				result = "unknown command: " + cmd
			}

			resp := Message{
				Type: "response",
				Payload: map[string]interface{}{
					"command": cmd,
					"data":    result,
				},
			}
			send(conn, resp)

		default:
			log.Printf("Unhandled message type: %s", msg.Type)
		}
	}
}

func collectMockStats() string {
	return "cpu=13.2%, mem=52.1%"
}

func send(conn *websocket.Conn, msg Message) {
	data, _ := json.Marshal(msg)
	conn.WriteMessage(websocket.TextMessage, data)
}
