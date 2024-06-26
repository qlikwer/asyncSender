package server

import (
	"asyncSender/cmd/app/handlers"
	"asyncSender/pkg/logger"
	"asyncSender/pkg/message"
	"flag"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func StartServer(messageQueue *message.Queue) {
	var port string
	flag.StringVar(&port, "port", "8065", "Port on which the server will be launched")
	flag.Parse()

	app := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
		AppName:      "Async Sender App v0.0.1",
	})

	apiRouter := app.Group("/api")
	apiRouter.Get("/health", func(c *fiber.Ctx) error {
		return handlers.HealthCheckerHandler(c, messageQueue)
	})
	apiRouter.Post("/send", func(c *fiber.Ctx) error {
		return handlers.SendMessageHandler(c, messageQueue)
	})

	logger.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
