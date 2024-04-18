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
					Headers:     messageToSend.Headers,
				})
				if err != nil {
					if messageToSend.Iteration >= MAX_ITERATION_COUNTER {
						logger.Errorf("The message was not sent in %d iteration. Message discarded from queue",
							MAX_ITERATION_COUNTER)
					} else {
						var err *sender.SendError
						if errors.As(err, &err) {
							logger.Warningf(
								"The service returned an error, current iteration: %d, repeat sending. Data: %s",
								messageToSend.Iteration, messageToSend.Data,
							)
							messageQueue.AddToTheBeginningEnqueue(*messageToSend)
						}
					}
				} else {
					logger.Info("Request sent successfully")
				}
			}
		}
	}()

	server.StartServer(messageQueue)

	waitGroup.Wait()
}
