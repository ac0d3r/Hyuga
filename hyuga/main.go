package main

import (
	"log"
	"net/http"

	"hyuga/config"
	"hyuga/database"
	"hyuga/handler/frontend"
	"hyuga/oob"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.SetFromYaml("config.yaml"); err != nil {
		log.Fatalln(err)
	}
	if err := database.Init(config.RedisDsn); err != nil {
		log.Fatal(err)
	}

	var engine *gin.Engine
	if config.DebugMode {
		engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
	}
	engine.SetTrustedProxies(nil)
	addRoute(engine)

	dns := oob.NewDnsServer("")
	go func() {
		log.Println("Listening and serving Dns on :53")
		dns.ListenAndServe()
	}()
	defer func() {
		if err := dns.Shutdown(); err != nil {
			log.Println(err)
		}
	}()

	if err := engine.Run(":8000"); err != nil {
		log.Printf("Could not serve http on port 8000: %s\n", err)
	}
}

func addRoute(engine *gin.Engine) {
	engine.Use(frontend.MiddlewareForwardLog(),
		static.Serve("/", static.LocalFile("dist", false)))
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
