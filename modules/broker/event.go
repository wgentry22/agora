package broker

type Event interface {
  Topic() *string
  Kind() []byte
  Payload() []byte
}

type EventHandler func([]byte) error

type EventBuilder interface {
  WithKind([]byte) EventBuilder
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

func (s *simpleEvent) Kind() []byte {
  return s.kind
}

func (s *simpleEvent) Payload() []byte {
  return s.payload
}

func (s *simpleEvent) WithKind(data []byte) EventBuilder {
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
  return &simpleEvent{
    topic: &topic,
  }
}
