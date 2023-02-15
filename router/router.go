package router

import (
	"github.com/romberli/go-template-web/api/v1/health"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/pingcap/errors"
	"github.com/romberli/go-template-web/config"
	"github.com/romberli/log"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/romberli/go-template-web/docs"
)

const (
	defaultStatusURL     = "/status"
	defaultBaseURL       = "/api/v1"
	defaultSwaggerURL    = "/swagger"
	defaultSwaggerDocURL = "/swagger/doc.json"
	defaultSwaggerAnyURL = "/*any"
)

type Router interface {
	http.Handler
	Use(middleware ...gin.HandlerFunc)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	Register()
	Run(addr ...string) error
}

type GinRouter struct {
	Engine *gin.Engine
}

func NewGinRouter() Router {
	if log.GetLevel() != zapcore.DebugLevel {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(Logger(), Recovery())

	return &GinRouter{engine}
}

func (gr *GinRouter) Use(middleware ...gin.HandlerFunc) {
	gr.Engine.Use(middleware...)
}

func (gr *GinRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	gr.Engine.ServeHTTP(w, req)
}

func (gr *GinRouter) Register() {
	// pprof
	if viper.GetBool(config.ServerPProfEnabledKey) {
		pprof.Register(gr.Engine)
	}

	// swagger
	gr.Swagger()
	gr.Engine.GET(defaultStatusURL, health.Status)

	apiV1 := gr.Engine.Group(defaultBaseURL)
	RegisterAll(apiV1)
}

func (gr *GinRouter) Run(addr ...string) error {
	return errors.Trace(gr.Engine.Run(addr...))
}

func (gr *GinRouter) Swagger() {
	swaggerGroup := gr.Engine.Group(defaultSwaggerURL)
	{
		url := ginSwagger.URL(defaultSwaggerDocURL)
		swaggerGroup.GET(defaultSwaggerAnyURL, ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}
}

func RegisterAll(group *gin.RouterGroup) {
	// health
	RegisterHealth(group)
}
