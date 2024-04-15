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

	sndr, err := sender.InitSender()
	if err != nil {
		logger.Fatal(err)
	}

	messageSendingTicker := time.NewTicker(time.Second / 30)
	messageToSend := &messageModule.Message{}

	go func() {
		for range messageSendingTicker.C {
			messageToSend = messageQueue.Dequeue()
			if messageToSend != nil {
				err = sndr.SendMessage(sender.SendMessageParams{
					Url:         messageToSend.Url,
					Data:        messageToSend.Data,
					RequestType: messageToSend.RequestType,
				})
				if err != nil {
					logger.Errorf("Ошибка отправки сообщения: %v", err)
					var err *sender.SendError
					if errors.As(err, &err) {
						logger.Warningf("Сервис вернул ошибку, повторяем отправку")
						time.Sleep(5 * time.Second)
						messageQueue.AddToTheBeginningEnqueue(*messageToSend) // Помещаем сообщение в начало очереди
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
