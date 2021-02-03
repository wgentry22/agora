package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wgentry22/agora/modules/auth"
	"github.com/wgentry22/agora/types/config"
)

func testHandler(c *gin.Context) {
	user, ok := c.Get("subject")
	if ok {
		c.JSON(http.StatusOK, map[string]string{"hello": user.(string)})
	} else {
		c.AbortWithStatus(http.StatusForbidden)
	}
}

var _ = Describe("AuthMiddleware", func() {

	var (
		router *gin.Engine
	)

	BeforeEach(func() {
		router = gin.New()
	})

	Context("when no config.Auth is provided to auth.Use", func() {
		It("should panic", func() {
			defer func() {
				if r := recover(); r != nil {
					err, isErr := r.(error)
					if !isErr {
						Fail("expected to panic with error")
					}

					Expect(err).To(Equal(auth.ErrAuthConfigurationRequired))
				}
			}()

			router.GET("/test", auth.RequiresTokenMiddleware, testHandler)
		})
	})

	Context("when a config.Auth is provided to auth.Use", func() {
		var (
			conf = config.Auth{
				Vendor: config.AuthVendorMock,
			}
			request *http.Request
		)

		JustBeforeEach(func() {
			auth.Use(conf)
			router.GET("/test", auth.RequiresTokenMiddleware, testHandler)

			req, err := http.NewRequest(http.MethodGet, "/test", nil)
			if err != nil {
				panic(err)
			}

			request = req
		})

		It("should return 401 UNAUTHORIZED when Bearer Token is missing", func() {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should reach the handler when Bearer Token is present", func() {
			w := httptest.NewRecorder()

			request.Header.Set("Authorization", "Bearer test")

			router.ServeHTTP(w, request)

			Expect(w.Code).To(Equal(http.StatusOK))

			var body map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &body)
			Expect(err).To(BeNil())

			Expect(body["hello"]).To(Equal("mock"))
		})
	})
})
