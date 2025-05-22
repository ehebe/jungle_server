package network

import (
	"encoding/json"
	"fmt"
	"net"
)

type Context struct {
	Data  []byte
	Reply func([]byte) error
}

type HandlerFunc func(ctx *Context)

type ProtocolHandler struct {
	handlers map[string]HandlerFunc
}

func NewProtocolHandler() *ProtocolHandler {
	return &ProtocolHandler{handlers: make(map[string]HandlerFunc)}
}

func (ph *ProtocolHandler) Handle(event string, fn HandlerFunc) {
	ph.handlers[event] = fn
}

func (ph *ProtocolHandler) Dispatch(event string, ctx *Context) {
	if h, ok := ph.handlers[event]; ok {
		h(ctx)
	} else {
		fmt.Println("[WARN] Unknown event:", event)
	}
}

type TransportLayer struct {
	role string
	conn net.Conn
}

func NewTransportLayer(role string) *TransportLayer {
	return &TransportLayer{role: role}
}

func (t *TransportLayer) Listen(dispatch func(event string, ctx *Context)) {
	if t.role == "master" {
		ln, _ := net.Listen("tcp", ":9000")
		fmt.Println("Master listening on :9000")
		conn, _ := ln.Accept()
		t.conn = conn
	} else {
		conn, _ := net.Dial("tcp", "127.0.0.1:9000")
		t.conn = conn
		fmt.Println("Agent connected to master")
	}

	for {
		buf := make([]byte, 1024)
		n, err := t.conn.Read(buf)
		if err != nil {
			fmt.Println("[ERR] read:", err)
			return
		}
		var msg map[string]json.RawMessage
		json.Unmarshal(buf[:n], &msg)
		event := string(msg["event"])
		ctx := &Context{
			Data: msg["payload"],
			Reply: func(resp []byte) error {
				respMsg := map[string]any{
					"event":   "reply",
					"payload": resp,
				}
				out, _ := json.Marshal(respMsg)
				_, err := t.conn.Write(out)
				return err
			},
		}
		dispatch(event, ctx)
	}
}

func (t *TransportLayer) Send(data []byte) {
	t.conn.Write(data)
}
