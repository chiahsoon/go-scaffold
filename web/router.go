package web

import (
	"log"

	"github.com/chiahsoon/go_scaffold/web/handlers"
	"github.com/chiahsoon/go_scaffold/web/helper"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Router struct {
	*gin.Engine
}

// @BasePath /apis/v1
func NewRouter() *Router {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Middlewares
	pprof.Register(router)
	// router.Use(
	//	 middlewares.Logger(),
	//	 middlewares.Recovery(version),
	//	 middlewares.Jsonifier(version),
	// )

	// API Endpoints
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1API := router.Group("/apis/v1")
	v1API.GET("home", helper.IsAuthorized, handlers.Home)
	v1API.POST("login", handlers.Login)
	v1API.POST("signup", handlers.Signup)
	v1API.GET("logout", helper.IsAuthorized, handlers.Logout)

	return &Router{
		router,
	}
}

func (r *Router) Run() {
	port := ":" + viper.GetString("port")

	if err := r.Engine.Run(port); err != nil {
		log.Fatal("failed to start server: \n", err)
	}
}

func Run() {
	NewRouter().Run()
}
