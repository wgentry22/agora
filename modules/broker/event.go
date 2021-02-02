package broker

import (
	"fmt"
)

type Event interface {
	Topic() *string
	Key() []byte
	Payload() []byte
}

func EventFactory(topic string) func(key, payload []byte) Event {
	return func(key, payload []byte) Event {
		return NewEventBuilder(topic).WithKey(key).WithPayload(payload).Build()
	}
}

type EventHandler func([]byte) error

type EventBuilder interface {
	WithKey([]byte) EventBuilder
	WithPayload([]byte) EventBuilder
	Build() Event
}

type simpleEvent struct {
	topic   *string
	kind    []byte
	payload []byte
}

func (s *simpleEvent) Topic() *string {
	return s.topic
}

func (s *simpleEvent) Key() []byte {
	return s.kind
}

func (s *simpleEvent) Payload() []byte {
	return s.payload
}

func (s *simpleEvent) WithKey(data []byte) EventBuilder {
	s.kind = data

	return s
}

func (s *simpleEvent) WithPayload(data []byte) EventBuilder {
	s.payload = data

	return s
}

func (s *simpleEvent) Build() Event {
	return s
}

func NewEventBuilder(topic string) EventBuilder {
	var t = &topic

	return &simpleEvent{
		topic: t,
	}
}

func (s *simpleEvent) String() string {
	return fmt.Sprintf("Event[Topic: %s, Payload: %s, Key: %s]", *s.topic, string(s.payload), string(s.kind))
}
