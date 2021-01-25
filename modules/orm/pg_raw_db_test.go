package orm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wgentry22/agora/modules/orm"
)

var _ = Describe("PostgreSQL RawDB", func() {

	Context("when config has not been used", func() {
		It("should panic", func() {
			defer func() {
				if r := recover(); r != nil {
					err, isErr := r.(error)
					if !isErr {
						Fail("expected to panic with error since DB has not used a config")
					}

					Expect(err).To(Equal(orm.ErrDBConnectionNotInitialized))
				}
			}()

			_ = orm.Get()
		})
	})

	Context("when config has been used", func() {
		JustBeforeEach(func() {
			config := mustParseConfig()
			orm.UseConfig(config)
		})

		It("should ping successfully", func() {
			orm := orm.Get()

			raw, err := orm.DB()
			Expect(err).To(BeNil())

			err = raw.Ping()
			Expect(err).To(BeNil())
		})
	})
})
