package logg_test

import (
	"bytes"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wgentry22/agora/modules/logg"
	"github.com/wgentry22/agora/types/config"
)

type LoggingTestStruct struct {
	From    string `json:"from"`
	Message string `json:"msg"`
}

var _ = Describe("Logger", func() {

	Describe("configured using config.Logging", func() {
		Context("when the configuration data is valid", func() {
			var (
				validConfig = config.Logging{
					Level:       "trace",
					OutputPaths: []string{"stdout"},
					Fields: map[string]interface{}{
						"from": "test",
					},
				}
				formattedConfig = config.Logging{
					Level:       "trace",
					OutputPaths: []string{"stdout"},
					Fields: map[string]interface{}{
						"from": "formatted_config",
					},
				}
				panicConfig = config.Logging{
					Level:       "trace",
					OutputPaths: []string{"stdout"},
					Fields: map[string]interface{}{
						"from": "panic",
					},
				}
				panicFormattedConfig = config.Logging{
					Level:       "trace",
					OutputPaths: []string{"stdout"},
					Fields: map[string]interface{}{
						"from": "formatted_panic",
					},
				}
			)

			It("should instantiate properly", func() {
				logger := logg.NewLogrusLogger(validConfig)
				Expect(logger).ToNot(BeNil())
			})

			It("should log", func() {
				var buf bytes.Buffer

				logger := logg.NewLogrusLogger(validConfig).WithWriter(&buf)
				Expect(logger).ToNot(BeNil())

				logger.Trace("Booyah!")
				buf.Write([]byte{','})
				logger.Debug("Booyah!")
				buf.Write([]byte{','})
				logger.Info("Booyah!")
				buf.Write([]byte{','})
				logger.Warning("Booyah!")
				buf.Write([]byte{','})
				logger.Warn("Booyah!")
				buf.Write([]byte{','})
				logger.Error("Booyah!")

				jsonArrayPre := append([]byte{'['}, buf.Bytes()...)
				jsonArray := append(jsonArrayPre, ']')

				var entries []LoggingTestStruct
				err := json.Unmarshal(jsonArray, &entries)
				Expect(err).To(BeNil())

				Expect(entries).To(HaveLen(6))
				for _, entry := range entries {
					Expect(entry.From).To(Equal(validConfig.Fields["from"]))
					Expect(entry.Message).To(Equal("Booyah!"))
				}
			})

			It("should log formatted", func() {
				var buf bytes.Buffer

				logger := logg.NewLogrusLogger(formattedConfig).WithWriter(&buf)
				Expect(logger).ToNot(BeNil())

				logger.Tracef("%s - formatted", "Booyah!")
				buf.Write([]byte{','})
				logger.Debugf("%s - formatted", "Booyah!")
				buf.Write([]byte{','})
				logger.Infof("%s - formatted", "Booyah!")
				buf.Write([]byte{','})
				logger.Warningf("%s - formatted", "Booyah!")
				buf.Write([]byte{','})
				logger.Warnf("%s - formatted", "Booyah!")
				buf.Write([]byte{','})
				logger.Errorf("%s - formatted", "Booyah!")

				jsonArrayPre := append([]byte{'['}, buf.Bytes()...)
				jsonArray := append(jsonArrayPre, ']')

				var entries []LoggingTestStruct
				err := json.Unmarshal(jsonArray, &entries)
				Expect(err).To(BeNil())

				Expect(entries).To(HaveLen(6))
				for _, entry := range entries {
					Expect(entry.From).To(Equal(formattedConfig.Fields["from"]))
					Expect(entry.Message).To(Equal("Booyah! - formatted"))
				}
			})

			It("should panic", func() {
				var buf bytes.Buffer

				logger := logg.NewLogrusLogger(panicConfig).WithWriter(&buf)
				Expect(logger).ToNot(BeNil())

				defer func() {
					if r := recover(); r != nil {
						var entry LoggingTestStruct

						err := json.Unmarshal(buf.Bytes(), &entry)
						Expect(err).To(BeNil())

						Expect(entry.Message).To(Equal("Booyah!"))
						Expect(entry.From).To(Equal(panicConfig.Fields["from"]))
					} else {
						Fail("expected to panic when logg.Logger.Panic is invoked")
					}
				}()

				logger.Panic("Booyah!")
			})

			It("should panic formatted", func() {
				var buf bytes.Buffer

				logger := logg.NewLogrusLogger(panicFormattedConfig).WithWriter(&buf)
				Expect(logger).ToNot(BeNil())

				defer func() {
					if r := recover(); r != nil {
						var entry LoggingTestStruct

						err := json.Unmarshal(buf.Bytes(), &entry)
						Expect(err).To(BeNil())

						Expect(entry.Message).To(Equal("Booyah! - formatted"))
						Expect(entry.From).To(Equal(panicFormattedConfig.Fields["from"]))
					} else {
						Fail("expected to panic when logg.Logger.Panicf is invoked")
					}
				}()

				logger.Panicf("%s - formatted", "Booyah!")
			})
		})
	})
})
