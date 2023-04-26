package event

import (
	"strings"
	"sync"
)

type EventBus struct {
	mux         sync.RWMutex
	subscribers []*Subscriber
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make([]*Subscriber, 0),
	}
}

func (e *EventBus) Publish(topic string, msg any) int {
	e.mux.RLock()
	defer e.mux.RUnlock()

	if topic == "" {
		return 0
	}

	count := 0
	for _, s := range e.subscribers {
		if s.Has(topic) {
			s.out <- msg
			count++
		}
	}
	return count
}

func (e *EventBus) Subscribe(topic string) *Subscriber {
	e.mux.Lock()
	defer e.mux.Unlock()

	s := NewSubscriber(topic)
	e.subscribers = append(e.subscribers, s)
	return s
}

func (e *EventBus) Unsubscribe(s *Subscriber) {
	e.mux.Lock()
	defer e.mux.Unlock()

	close(s.out)
	for i := range e.subscribers {
		if e.subscribers[i] == s {
			e.subscribers = append(e.subscribers[:i], e.subscribers[i+1:]...)
			break
		}
	}
}

type Subscriber struct {
	topic string
	out   chan interface{}
}

func NewSubscriber(topic string) *Subscriber {
	return &Subscriber{topic: topic,
		out: make(chan any, 1e2)}
}

func (s *Subscriber) Out() <-chan interface{} {
	return s.out
}

func (s *Subscriber) Has(topic string) bool {
	return strings.EqualFold(s.topic, topic)
}
