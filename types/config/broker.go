package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type BrokerRole int8

const (
	BrokerRoleUnknown BrokerRole = iota
	BrokerRoleProducer
	BrokerRoleConsumer
)

type BrokerVendor int8

const (
	BrokerVendorUnknown BrokerVendor = iota
	BrokerVendorKafka
)

var (
	bvDisplay = []string{"unknown", "kafka"}
	brDisplay = []string{"unknown", "producer", "consumer"}
	bvLookup  = map[string]BrokerVendor{
		"unknown": BrokerVendorUnknown,
		"kafka":   BrokerVendorKafka,
	}
	brLookup = map[string]BrokerRole{
		"unknown":  BrokerRoleUnknown,
		"producer": BrokerRoleProducer,
		"consumer": BrokerRoleConsumer,
	}
	ErrBrokerRoleRequired      = errors.New("value for `broker.role` is expected")
	ErrBrokerVendorRequired    = errors.New("value for `broker.vendor` is expected")
	ErrBrokerIDRequired        = errors.New("value for `broker.id` is expected")
	ErrBrokerServersRequired   = errors.New("values for `broker.servers` must be present")
	ErrFailedToSetBrokerConfig = func(key string, value interface{}) error {
		return fmt.Errorf("failed to set `%s:%s`", key, value)
	}
)

func (br BrokerRole) String() string {
	return brDisplay[br]
}

func (bv BrokerVendor) String() string {
	return bvDisplay[bv]
}

type Broker struct {
	ID         string                 `toml:"id"`
	Vendor     BrokerVendor           `toml:"vendor"`
	Role       BrokerRole             `toml:"role"`
	Servers    []string               `toml:"servers"`
	BufferSize int                    `toml:"buffer_size"`
	Timeout    int                    `toml:"timeout"`
	Args       map[string]interface{} `toml:"args"`
}

func (b *Broker) UnmarshalTOML(data interface{}) (err error) {
	dataMap := data.(map[string]interface{})

	if id, ok := dataMap["id"].(string); ok {
		b.ID = id
	} else {
		err = errwrap.Wrap(ErrBrokerIDRequired, err)
	}

	if role, ok := dataMap["role"].(string); ok {
		found, roleFound := brLookup[role]
		if roleFound {
			b.Role = found
		}
	}

	if vals, ok := dataMap["servers"].([]interface{}); ok {
		servers := make([]string, len(vals))
		for i, val := range vals {
			servers[i] = val.(string)
		}

		b.Servers = servers
	} else {
		err = errwrap.Wrap(ErrBrokerServersRequired, err)
	}

	if args, ok := dataMap["args"].(map[string]interface{}); ok {
		b.Args = args
	}

	if buffer, ok := dataMap["buffer_size"].(int64); ok {
		b.BufferSize = int(buffer)
	} else {
		b.BufferSize = 100
	}

	if timeout, ok := dataMap["timeout"].(int64); ok {
		b.Timeout = int(timeout)
	} else {
		b.Timeout = 150
	}

	if role, ok := dataMap["role"].(string); ok {
		r, ok := brLookup[role]
		if ok {
			b.Role = r
		} else {
			b.Role = BrokerRoleUnknown
		}
	} else {
		err = errwrap.Wrap(ErrBrokerRoleRequired, err)
	}

	if vendor, ok := dataMap["vendor"].(string); ok {
		if val, valOk := bvLookup[vendor]; valOk {
			b.Vendor = val
		} else {
			b.Vendor = BrokerVendorUnknown
		}
	} else {
		err = errwrap.Wrap(ErrBrokerVendorRequired, err)
	}

	return err
}

func (b Broker) ForSubscriber() *kafka.ConfigMap {
	if len(b.Servers) == 0 {
		panic(ErrBrokerServersRequired)
	}

	if b.ID == "" {
		panic(ErrBrokerIDRequired)
	}

	servers := strings.Join(b.Servers, ",")

	confMap := &kafka.ConfigMap{
		"bootstrap.servers": servers,
		"group.id":          b.ID,
	}

	for k, v := range b.Args {
		if err := confMap.SetKey(k, v); err != nil {
			panic(errwrap.Wrap(ErrFailedToSetBrokerConfig(k, v), err))
		}
	}

	return confMap
}

func (b Broker) ForPublisher() *kafka.ConfigMap {
	if len(b.Servers) == 0 {
		panic(ErrBrokerServersRequired)
	}

	if b.ID == "" {
		panic(ErrBrokerIDRequired)
	}

	servers := strings.Join(b.Servers, ",")

	confMap := &kafka.ConfigMap{
		"bootstrap.servers": servers,
		"client.id":         b.ID,
	}

	for k, v := range b.Args {
		if err := confMap.SetKey(k, v); err != nil {
			panic(errwrap.Wrap(ErrFailedToSetBrokerConfig(k, v), err))
		}
	}

	return confMap
}
