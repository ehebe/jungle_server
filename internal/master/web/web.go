package web

import "github.com/gofiber/fiber/v2"

type JungleServer struct {
	HttpServer *fiber.App
}

func NewJungleServer() *JungleServer {

	jserver := &JungleServer{}
	jserver.HttpServer = fiber.New(fiber.Config{
    Prefork: true,
    ServerHeader: "Jungle",
})
	return jserver
}
