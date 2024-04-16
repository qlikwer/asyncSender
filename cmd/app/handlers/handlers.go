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
	logger.Info("Проверка работоспособности. Размер очереди - " + fmt.Sprintf("%d", messageQueue.Size()))

	return c.JSON(fiber.Map{
		"requestDateTine": requestDateTime,
		"queueSize":       messageQueue.Size(),
	})
}

func SendMessageHandler(c *fiber.Ctx, messageQueue *messageModule.Queue) error {
	requestDateTime := time.Now().Format(logger.DateTimeFormat)

	Url := c.Get("Url")
	Data := c.Body()
	RequestType := c.Get("RequestType")

	if len(Data) == 0 {
		logger.Errorf("Ошибка чтения тела запроса")
		return c.Status(fiber.StatusBadRequest).SendString("Ошибка чтения тела запроса")
	}

	newMessage := messageModule.Message{
		Url:         Url,
		Data:        string(Data),
		RequestType: RequestType,
		Iteration:   0,
	}

	messageQueue.Enqueue(newMessage) // Добавляем сообщение в очередь

	return c.JSON(fiber.Map{
		"requestDateTime": requestDateTime,
		"queueSize":       messageQueue.Size(),
	})
}
