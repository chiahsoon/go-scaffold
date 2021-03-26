package web

import (
	"log"

	"github.com/chiahsoon/go_scaffold/web/handlers"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Router struct {
	*gin.Engine
}

// @title golang-scaffold
// @description This is a API server scaffold written in Go.

// @BasePath /apis/v1
func NewRouter() *Router {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Middlewares
	pprof.Register(router)

	// API Endpoints
	v1APIWithoutAuth := router.Group("/apis/v1")
	v1APIWithoutAuth.POST("signup", handlers.Signup)
	v1APIWithoutAuth.POST("login", handlers.Login)

	v1API := v1APIWithoutAuth.Use(handlers.IsAuthenticated)
	v1API.GET("home", handlers.Home)
	v1API.GET("logout", handlers.Logout)
	v1API.GET("current_user", handlers.CurrentUser)

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
