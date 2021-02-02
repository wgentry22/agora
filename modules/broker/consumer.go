package broker

import (
	"errors"
	"sync"

	"github.com/hashicorp/errwrap"
	"github.com/wgentry22/agora/types/config"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
	ErrConsumerConfigurationExpected = errors.New("expected configuration with broker role `consumer`")
)

type Consumer interface {
	Start()
	RegisterHandler(topic string, handler EventHandler)
	Errors() <-chan error
}

func NewConsumer(conf config.Broker) Consumer {
	if conf.Role.String() == "consumer" {
		if conf.Vendor.String() == "kafka" {
			return newKafkaConsumer(conf)
		}

		return nil
	}

	panic(ErrConsumerConfigurationExpected)
}

func newKafkaConsumer(conf config.Broker) Consumer {
	consumer, err := kafka.NewConsumer(conf.ForSubscriber())
	if err != nil {
		panic(err)
	}

	return &kafkaConsumer{
		timeout:  conf.Timeout,
		consumer: consumer,
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
	run := true

	topics := make([]string, 0)
	
	k.handlers.Range(func(k, v interface{}) bool {
		topics = append(topics, k.(string))
		return true
	})

	if err := k.consumer.SubscribeTopics(topics, nil); err != nil {
		panic(errwrap.Wrap(errors.New("failed to subscribe to topics"), err))
	}

	for run {
		event := k.consumer.Poll(k.timeout)
		switch e := event.(type) {
		case *kafka.Message:
			if handler, ok := k.handlers.Load(*e.TopicPartition.Topic); ok {
				if eventHandler, ok := handler.(EventHandler); ok {
					if err := eventHandler(e.Value); err != nil {
						k.errc <- err
					}
				}
			}
		case kafka.Error:
			run = false
			k.errc <- e
		}
	}

	if err := k.consumer.Close(); err != nil {
		k.errc <- err
	}
}

func (k *kafkaConsumer) RegisterHandler(topic string, handler EventHandler) {
	k.handlers.Store(topic, handler)
}

func (k *kafkaConsumer) Errors() <-chan error {
	return k.errc
}
