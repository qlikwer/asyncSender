package handlers

import (
	"asyncSender/pkg/logger"
	messageModule "asyncSender/pkg/message"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func HealthCheckerHandler(c *fiber.Ctx, messageQueue *messageModule.Queue) error {
	requestDateTime := time.Now().Format(logger.DateTimeFormat)
	logger.Info("Health check. Queue size - " + fmt.Sprintf("%d", messageQueue.Size()))

	return c.JSON(fiber.Map{
		"requestDateTine": requestDateTime,
		"queueSize":       messageQueue.Size(),
	})
}

func SendMessageHandler(c *fiber.Ctx, messageQueue *messageModule.Queue) error {
	requestDateTime := time.Now().Format(logger.DateTimeFormat)

	allHeaders := c.GetReqHeaders()
	RequestType := c.Method()

	var Url string

	for key, values := range allHeaders {
		for _, value := range values {
			switch key {
			case "Url":
				Url = value
				delete(allHeaders, key)
			}
		}
	}

	Data := c.Body()

	if len(Data) == 0 {
		logger.Errorf("Error reading request body")
		return c.Status(fiber.StatusBadRequest).SendString("Error reading request body")
	}

	newMessage := messageModule.Message{
		Url:         Url,
		Data:        string(Data),
		RequestType: RequestType,
		Iteration:   0,
		Headers:     allHeaders,
	}

	messageQueue.Enqueue(newMessage) // Добавляем сообщение в очередь

	return c.JSON(fiber.Map{
		"requestDateTime": requestDateTime,
		"queueSize":       messageQueue.Size(),
	})
}
