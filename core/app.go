package core

import (
	"Hyuga/core/api"
	v1 "Hyuga/core/api/v1"
	"Hyuga/core/conf"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func registerRoutes(e *echo.Echo) {
	e.GET("/", api.Index)
	e.POST("/v1/users", v1.CreateUser)
	e.GET("/v1/records", v1.GetRecords)
}

func registerCORS(e *echo.Echo) {
	origins := []string{"http://" + conf.Domain}
	if conf.AppENV == "development" {
		origins[0] = "*"
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))
}

// Create Hyuga echo application
func Create() (app *echo.Echo) {
	app = echo.New()
	app.Pre(httpDog())
	registerRoutes(app)
	registerCORS(app)
	return
}
