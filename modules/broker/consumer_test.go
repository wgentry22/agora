package broker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wgentry22/agora/modules/broker"
	"github.com/wgentry22/agora/types/config"
)

var _ = Describe("Consumer", func() {

	Context("when vendor is `kafka`", func() {
		var (
			consumer broker.Consumer
			conf     = config.Broker{
				ID:         "testId",
				Role:       config.BrokerRoleConsumer,
				Vendor:     config.BrokerVendorKafka,
				Servers:    []string{"localhost:9092"},
				BufferSize: 1234,
				Timeout:    1234,
				Args:       map[string]interface{}{},
			}
		)

		BeforeEach(func() {
			consumer = broker.NewConsumer(conf)
		})

		It("should be generated by broker.NewConsumer", func() {
			Expect(consumer).ToNot(BeNil())
		})
	})
})
