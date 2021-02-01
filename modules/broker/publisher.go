package broker

import (
	"errors"

	"github.com/wgentry22/agora/types/config"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
	ErrProducerConfigurationExpected = errors.New("expected configuration with broker role `producer`")
	ErrFailedToDeliverMessage        = errors.New("delivery of message failed")
)

type Publisher interface {
	Publish(event Event)
	Errors() <-chan error
}

func NewPublisher(conf config.Broker) Publisher {
	if conf.Role.String() == "producer" {
		if conf.Vendor.String() == "kafka" {
			return newKafkaPublisher(conf)
		}

		return nil
	}

	panic(ErrProducerConfigurationExpected)
}

func newKafkaPublisher(conf config.Broker) Publisher {
	pub, err := kafka.NewProducer(conf.ForPublisher())
	if err != nil {
		panic(err)
	}

	return &kafkaPublisher{
		publisher: pub,
		events:    make(chan kafka.Event),
		errors:    make(chan error, 1),
	}
}

type kafkaPublisher struct {
	publisher *kafka.Producer
	events    chan kafka.Event
	errors    chan error
}

func (k *kafkaPublisher) Publish(event Event) {
	if event == nil {
		panic(errors.New("cannot publish nil event"))
	}

	panic(event)
}

func (k *kafkaPublisher) Errors() <-chan error {
	return k.errors
}
