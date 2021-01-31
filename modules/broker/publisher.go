package broker

import (
  "errors"
  "time"

  "github.com/hashicorp/errwrap"
  "github.com/wgentry22/agora/modules/logg"
  "github.com/wgentry22/agora/types/config"
  "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
  ErrFailedToDeliverMessage = errors.New("delivery of message failed")
  logger = logg.Root()
)

type Publisher interface {
  Publish(event Event)
  Errors() <-chan error
}

func NewPublisher(conf config.Broker) Publisher {
  return &kafkaPublisher{
    publisher: conf.NewPublisher(),
    events:    make(chan kafka.Event, conf.BufferSize),
    errors:    make(chan error),
  }
}

type kafkaPublisher struct {
  publisher *kafka.Producer
  events    chan kafka.Event
  errors    chan error
}

func (k *kafkaPublisher) Publish(event Event) {
  if err := k.publisher.Produce(&kafka.Message{
    TopicPartition: kafka.TopicPartition{ //nolint:exhaustivestruct
      Topic:     event.Topic(),
      Partition: kafka.PartitionAny,
    },
    Value:         event.Payload(),
    Key:           event.Kind(),
    Timestamp:     time.Now(),
    TimestampType: 0,
    Opaque:        nil,
    Headers:       nil,
  }, k.events); err != nil {
    logger.
      WithError(err).
      Warning("Failed to produce message")

    k.errors <- err
  }

  go func(ec chan kafka.Event) {
    e := <-ec
    if msg, ok := e.(*kafka.Message); ok {
      if msg.TopicPartition.Error != nil {
        k.errors <- errwrap.Wrap(ErrFailedToDeliverMessage, msg.TopicPartition.Error)
      } else {
        logger.
          WithField("topic", *msg.TopicPartition.Topic).
          WithField("delivered", msg.Timestamp).
          Debug("Message sent successfully")
      }
    }
  }(k.events)
}

func (k *kafkaPublisher) Errors() <-chan error {
  return k.errors
}
