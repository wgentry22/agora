package broker_test

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/pelletier/go-toml"
  "github.com/wgentry22/agora/types/config"
)

var _ = Describe("Parse", func() {

  Describe("should parse as TOML", func() {

  })

  Context("When data is in correct format", func() {
    var (
      data = []byte(`
id = "testId"
vendor = "kafka"
role = "consumer"
servers = ["localhost:9092"]
buffer_size = 1234
timeout = 1234
[args]
key = "value"
`)
    )

    It("should do so correctly", func() {
      var conf config.Broker

      By("not returning error")
      err := toml.Unmarshal(data, &conf)
      Expect(err).To(BeNil())

      By("correctly parsing")
      Expect(conf).To(Equal(config.Broker{
        ID:         "testId",
        Role:       config.BrokerRoleConsumer,
        Vendor:     config.BrokerVendorKafka,
        Servers:    []string{"localhost:9092"},
        BufferSize: 1234,
        Timeout:    1234,
        Args: map[string]interface{}{
          "key": "value",
        },
      }))
    })
  })
})
