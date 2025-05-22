package master

import (
	"log"

	"github.com/ehebe/jungle/internal/master/ws"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func Start() {
	app := fiber.New()
	hub := ws.NewAgentHub()
	hub.OnConnect(func(ctx *ws.WSContext) {
		log.Printf("Agent %s connected", ctx.ID)
	})

	hub.OnMessage(func(ctx *ws.WSContext, msg ws.Message) {
		log.Printf("[%s] type=%s payload=%v", ctx.ID, msg.Type, msg.Payload)
		if msg.Type == "response" {
			ctx.Conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"ack"}`))
		}
	})

	hub.OnDisconnect(func(ctx *ws.WSContext) {
		log.Printf("Agent %s disconnected", ctx.ID)
	})

	app.Get("/ws/:id", hub.Handle())

	log.Fatal(app.Listen(":8080"))
}
