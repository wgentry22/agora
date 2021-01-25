package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wgentry22/agora/types/config"
	"net/http"
	"sync"
)

var (
	m sync.Mutex
)

type Router struct {
	api    config.API
	router *gin.Engine
}

func (r *Router) Server() *http.Server {
	return &http.Server{
		Addr:         r.api.ListenAddr(),
		Handler:      r.router,
		ReadTimeout:  r.api.Timeout.Read,
		WriteTimeout: r.api.Timeout.Write,
	}
}

func (r *Router) Handler() http.Handler {
	return r.router
}

type Controller struct {
	uri    string
	routes []Route
}

type Route struct {
	handler func(ctx *gin.Context)
	subPath string
	method  string
}

func NewRouter(config config.API) Router {
	router := &Router{
		api:    config,
		router: gin.Default(),
	}

	infoController := NewInfoController(config.Info())

	router.Register(infoController)

	return *router
}

func NewController(uri string) Controller {
	return Controller{uri, make([]Route, 0)}
}

func (r *Router) Register(controller Controller) {
	m.Lock()
	defer m.Unlock()

	rg := r.routerGroup().Group(controller.uri)

	for _, route := range controller.routes {
		rg.Handle(route.method, route.subPath, route.handler)
	}
}

func (r *Router) routerGroup() *gin.RouterGroup {
	return r.router.Group(r.api.PathPrefix)
}

func (c *Controller) Register(route Route) {
	c.routes = append(c.routes, route)
}

func NewGETRoute(uri string, handler func(c *gin.Context)) Route {
	return newRoute(http.MethodGet, uri, handler)
}

func NewPOSTRoute(uri string, handler func(c *gin.Context)) Route {
	return newRoute(http.MethodPost, uri, handler)
}

func NewPUTRoute(uri string, handler func(c *gin.Context)) Route {
	return newRoute(http.MethodPut, uri, handler)
}

func NewPATCHRoute(uri string, handler func(c *gin.Context)) Route {
	return newRoute(http.MethodPatch, uri, handler)
}

func NewDELETERoute(uri string, handler func(c *gin.Context)) Route {
	return newRoute(http.MethodDelete, uri, handler)
}

func newRoute(method, uri string, handler func(c *gin.Context)) Route {
	return Route{
		handler: handler,
		subPath: uri,
		method:  method,
	}
}

func NewInfoController(info config.Info) Controller {
	versionHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"version": info.Version.String()})
	}

	infoController := NewController("/info")

	infoController.Register(NewGETRoute("/version", versionHandler))

	return infoController
}
