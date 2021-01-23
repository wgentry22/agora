package orm_test

import (
	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wgentry22/agora/modules/orm"
)

var _ = Describe("PostgreSQL RawDB", func() {

	var (
		rawDB *sql.DB
	)

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

			_ = orm.GetRaw()
		})
	})

	Context("when config has been used", func() {
		JustBeforeEach(func() {
			config := mustParseConfig()
			orm.UseConfig(config)
		})

		It("should ping successfully", func() {
			rawDB = orm.GetRaw()
			Expect(rawDB).ToNot(BeNil())

			err := rawDB.Ping()
			Expect(err).To(BeNil())
		})
	})
})
