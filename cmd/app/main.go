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

const MAX_ITERATION_COUNTER = 5

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

				messageToSend.Iteration = messageToSend.Iteration + 1
				err = sndr.SendMessage(sender.SendMessageParams{
					Url:         messageToSend.Url,
					Data:        messageToSend.Data,
					RequestType: messageToSend.RequestType,
					Iteration:   messageToSend.Iteration,
				})
				if err != nil {
					if messageToSend.Iteration >= MAX_ITERATION_COUNTER {
						logger.Errorf("Сообщение не было отправлено за %d %s. Сообщение выброшено из очереди",
							MAX_ITERATION_COUNTER, sender.Pluralize(
								MAX_ITERATION_COUNTER, "итерация", "итерации", "итераций",
							))
					} else {
						var err *sender.SendError
						if errors.As(err, &err) {
							logger.Warningf(
								"Сервис вернул ошибку, текущая итерация: %d, повторяем отправку. Data: %s", messageToSend.Iteration, messageToSend.Data,
							)
							time.Sleep(5 * time.Second)
							//messageQueue.AddToTheBeginningEnqueue(*messageToSend) // Помещаем сообщение в начало очереди
							messageQueue.Enqueue(*messageToSend) // Помещаем сообщение в конец очереди
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
