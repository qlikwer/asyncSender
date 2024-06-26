package message

import (
	"sync"
)

type Message struct {
	Url         string              `json:"Url"`
	Data        string              `json:"json data"`
	RequestType string              `json:"Request type"`
	Iteration   int                 `json:"Current iteration"`
	Headers     map[string][]string `json:"Array of headers"`
}

type Queue struct {
	sync.Mutex
	messages []Message
}

// AddToTheBeginningEnqueue Помещает сообщение в начало очереди
func (q *Queue) AddToTheBeginningEnqueue(message Message) {
	q.Lock()
	defer q.Unlock()
	q.messages = append([]Message{message}, q.messages...)
}

// Enqueue Добавляет сообщение в очередь
func (q *Queue) Enqueue(message Message) {
	q.Lock()
	defer q.Unlock()
	q.messages = append(q.messages, message)
}

// Dequeue извлекает первое сообщение из очереди
func (q *Queue) Dequeue() *Message {
	q.Lock()
	defer q.Unlock()
	if len(q.messages) == 0 {
		return nil
	}
	message := q.messages[0]
	q.messages = q.messages[1:]

	return &message
}

// Size Возвращает размер очереди
func (q *Queue) Size() int {
	q.Lock()
	defer q.Unlock()

	return len(q.messages)
}
