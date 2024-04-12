package handlers

import (
	"asyncSender/pkg/logger"
	messageModule "asyncSender/pkg/message"
	"github.com/gofiber/fiber/v2"
	"time"
)

func HealthCheckerHandler(c *fiber.Ctx) error {
	requestDateTime := time.Now().Format(logger.DateTimeFormat)
	logger.Info("Проверка работоспособности")

	return c.JSON(fiber.Map{
		"requestDateTine": requestDateTime,
	})
}

func SendMessageHandler(c *fiber.Ctx, messageQueue *messageModule.Queue) error {
	requestDateTime := time.Now().Format(logger.DateTimeFormat)

	Url := c.Get("Url")
	Data := c.Body()
	RequestType := c.Get("RequestType")
	logger.Errorf(Url)
	logger.Errorf(string(Data))
	logger.Errorf(RequestType)

	if len(Data) == 0 {
		logger.Errorf("Ошибка чтения тела запроса")
		return c.Status(fiber.StatusBadRequest).SendString("Ошибка чтения тела запроса")
	}

	newMessage := messageModule.Message{
		Url:         Url,
		Data:        Data,
		RequestType: RequestType,
	}

	messageQueue.Enqueue(newMessage) // Добавляем сообщение в очередь

	return c.JSON(fiber.Map{
		"requestDateTime": requestDateTime,
		"queueSize":       messageQueue.Size(),
	})
}
