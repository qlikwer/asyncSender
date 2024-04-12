package main

import (
	"asyncSender/cmd/app/server"
	"asyncSender/pkg/logger"
	messageModule "asyncSender/pkg/message"
	"asyncSender/pkg/sender"
	"errors"
	"sync"
	"time"
)

var (
	messageQueue = &messageModule.Queue{}
	waitGroup    sync.WaitGroup
)

func main() {

	bot, err := sender.InitSender()
	if err != nil {
		logger.Fatal(err)
	}

	messageSendingTicker := time.NewTicker(time.Second / 30)
	messageToSend := &messageModule.Message{}

	go func() {
		for range messageSendingTicker.C {
			messageToSend = messageQueue.Dequeue()
			if messageToSend != nil {
				err = bot.SendMessage(sender.SendMessageParams{
					Url:         messageToSend.Url,
					Data:        messageToSend.Data,
					RequestType: messageToSend.RequestType,
				})
				if err != nil {
					logger.Errorf("Ошибка отправки сообщения: %v", err)
					var telegramErr *sender.SendError
					if errors.As(err, &telegramErr) {
						retryAfter, _ := sender.ParseRetryAfter(telegramErr)
						if retryAfter != 0 {
							logger.Warningf("Сервис прилёг отдохнуть на %d %s",
								retryAfter, sender.Pluralize(retryAfter, "секунда", "секунды", "секунд"))
							time.Sleep(time.Duration(retryAfter) * time.Second)
							messageQueue.AddToTheBeginningEnqueue(*messageToSend) // Помещаем сообщение в начало очереди
						}
					}
				} else {
					logger.Info("Запрос успешно отправлен")
				}
			}
		}
	}()

	server.StartServer(messageQueue)

	waitGroup.Wait()
}
