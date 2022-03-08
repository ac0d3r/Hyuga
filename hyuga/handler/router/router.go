package router

import (
	"hyuga/config"
	"hyuga/handler/frontend"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	var engine *gin.Engine
	if config.DebugMode {
		engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
	}
	engine.SetTrustedProxies(nil)
	addRoute(engine)
	return engine
}

func addRoute(engine *gin.Engine) {
	engine.Use(frontend.MiddlewareForwardLog())
	engine.Use(static.Serve("/", static.LocalFile("dist", false)))
	engine.NoRoute(func(c *gin.Context) { c.Status(http.StatusNotFound) })

	addAPIRoute(engine.Group("/api"))
}

func addAPIRoute(group *gin.RouterGroup) {
	userGroup := group.Group("/user")
	userGroup.POST("create", frontend.CreateUser)

	userGroup.Use(frontend.MiddlewareUserToken())
	{
		userGroup.GET("/dns-rebinding", frontend.GetUserDnsRebinding)
		userGroup.POST("/dns-rebinding", frontend.UpdateUserDnsRebinding)
		userGroup.POST("/delete", frontend.DeleteUser)
	}

	recordsGroup := group.Group("/record")
	recordsGroup.Use(frontend.MiddlewareUserToken())
	{
		recordsGroup.GET("/list", frontend.GetRecords)
		recordsGroup.POST("/clean", frontend.CleanRecords)
	}
}
