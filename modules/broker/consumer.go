package broker

import (
  "sync"

  "github.com/wgentry22/agora/types/config"
  "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type Consumer interface {
  Start()
  RegisterHandler(topic string, handler EventHandler)
  Errors() <-chan error
}

func NewConsumer(conf config.Broker) Consumer {
  return &kafkaConsumer{
    timeout:  conf.Timeout,
    consumer: conf.NewSubscriber(),
    handlers: sync.Map{},
    errc:     make(chan error),
  }
}

type kafkaConsumer struct {
  timeout  int
  consumer *kafka.Consumer
  handlers sync.Map
  errc     chan error
}

func (k *kafkaConsumer) Start() {
  go func(ec chan error) {
    run := true

    logger.Info("Successfully started broker.Consumer")

    for run {
      event := k.consumer.Poll(k.timeout)
      switch e := event.(type) {
      case *kafka.Message:
        if handler, ok := k.handlers.Load(*e.TopicPartition.Topic); ok {
          if eventHandler, ok := handler.(EventHandler); ok {
            if err := eventHandler(e.Value); err != nil {
              ec <- err
            }
          }
        }
      case kafka.PartitionEOF:
        logger.Warning("reached end of partition")
      case kafka.Error:
        logger.WithError(e).Warning("Stopping consumer")

        run = false
      }
    }

    if err := k.consumer.Close(); err != nil {
      ec <- err
    }
  }(k.errc)
}

func (k *kafkaConsumer) RegisterHandler(topic string, handler EventHandler) {
  k.handlers.Store(topic, handler)
}

func (k *kafkaConsumer) Errors() <-chan error {
  return k.errc
}
