package router

import (
	"net/http"

	"hyuga/conf"
	"hyuga/internal/core"
	"hyuga/router/api"
	v1 "hyuga/router/api/v1"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func registerRoutes(e *echo.Echo) {
	e.GET("/", api.Index)
	e.POST("/v1/users", v1.CreateUser)
	e.GET("/v1/records", v1.GetRecords)
	e.GET("/v1/users/:identity/dns-rebinding", v1.GetUserDNSRebinding)
	e.POST("/v1/users/:identity/dns-rebinding", v1.SetUserDNSRebinding)
}

func registerCORS(e *echo.Echo) {
	origins := []string{"http://" + conf.Domain.Main}
	if conf.AppEnv == "development" {
		origins[0] = "*"
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))
}

// New create hyuga api router
func New() (app *echo.Echo) {
	app = echo.New()
	app.Pre(core.HTTPDog())
	registerRoutes(app)
	registerCORS(app)
	return
}
