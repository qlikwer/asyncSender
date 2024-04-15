package server

import (
	"asyncSender/cmd/app/handlers"
	"asyncSender/pkg/logger"
	messageModule "asyncSender/pkg/message"
	"github.com/gofiber/fiber/v2"
)

func StartServer(messageQueue *messageModule.Queue) {
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

	logger.Fatal(app.Listen(":8080"))
}
