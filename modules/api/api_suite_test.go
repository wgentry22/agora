package api_test

import (
	"github.com/gin-gonic/gin"
	"github.com/wgentry22/agora/modules/api"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testPathPrefix = "/test"
	getRoute       = api.NewGETRoute("/", func(c *gin.Context) {
		c.String(http.StatusOK, "get all")
	})
	getIdRoute = api.NewGETRoute("/:id", func(c *gin.Context) {
		c.String(http.StatusOK, "get %s", c.Param("id"))
	})
	postRoute = api.NewPOSTRoute("/", func(c *gin.Context) {
		c.Status(http.StatusCreated)
	})
	putRoute = api.NewPUTRoute("/:id", func(c *gin.Context) {
		c.String(http.StatusOK, "put %s", c.Param("id"))
	})
	patchRoute = api.NewPATCHRoute("/:id", func(c *gin.Context) {
		c.String(http.StatusOK, "patch %s", c.Param("id"))
	})
	deleteRoute = api.NewDELETERoute("/:id", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	controller api.Controller
)

var _ = BeforeSuite(func() {
	controller = api.NewController(testPathPrefix)

	controller.Register(getRoute)
	controller.Register(getIdRoute)
	controller.Register(postRoute)
	controller.Register(putRoute)
	controller.Register(patchRoute)
	controller.Register(deleteRoute)
})

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}
