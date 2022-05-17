package broker

import (
	"effie/logger"
	"sync"
)

var log = logger.Get("broker").Sugar()

const (
	chBufferSize = 128
)

type Handler func(topic string, message interface{})

type Message struct {
	topic string
	msg   interface{}
}

var ch chan *Message

var handlers map[string][]Handler
var handlersMutex sync.Mutex

func init() {
	handlers = make(map[string][]Handler)
	ch = make(chan *Message, chBufferSize)

	go func() {
		for message := range ch {
			log.Debugw("handing message to subscribers", "topic", message.topic)

			handlersMutex.Lock()
			chachedHandlers, exists := handlers[message.topic]
			handlersMutex.Unlock()
			if !exists {
				continue
			}

			for _, handler := range chachedHandlers {
				handler(message.topic, message.msg)
			}
		}
	}()
}

func Publish(topic string, msg interface{}) error {
	ch <- &Message{
		topic: topic,
		msg:   msg,
	}
	log.Debugw("published message to topic", "topic", topic)

	return nil
}

func Subscribe(topic string, handler Handler) error {
	handlersMutex.Lock()
	defer handlersMutex.Unlock()

	handlerList, exists := handlers[topic]
	if !exists {
		handlers[topic] = make([]Handler, 1)
		handlers[topic][0] = handler
		return nil
	}

	handlers[topic] = append(handlerList, handler)
	log.Debug("registered new subscribe handler", "topic", topic)
	return nil
}
