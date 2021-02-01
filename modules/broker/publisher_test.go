package broker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wgentry22/agora/modules/broker"
	"github.com/wgentry22/agora/types/config"
)

var _ = Describe("Publisher", func() {

	Context("when vendor is `kafka`", func() {
		var (
			publisher broker.Publisher
			conf      = config.Broker{
				ID:         "testId",
				Role:       config.BrokerRoleProducer,
				Vendor:     config.BrokerVendorKafka,
				Servers:    []string{"localhost:9092"},
				BufferSize: 1234,
				Timeout:    1234,
				Args:       map[string]interface{}{},
			}
		)

		BeforeEach(func() {
			publisher = broker.NewPublisher(conf)
		})

		It("should be generated by broker.NewPublisher", func() {
			Expect(publisher).ToNot(BeNil())
		})
	})
})
