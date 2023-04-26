package handler

import (
	"net"
	"net/http"

	"github.com/ac0d3r/hyuga/internal/config"
	"github.com/ac0d3r/hyuga/internal/db"
	"github.com/ac0d3r/hyuga/internal/event"
	"github.com/ac0d3r/hyuga/internal/oob"
	"github.com/ac0d3r/hyuga/internal/record"
	"github.com/gin-gonic/gin"
)

type restfulHandler struct {
	db       *db.DB
	cnf      *config.Web
	dns      *config.DNS
	eventbus *event.EventBus
	recorder *record.Recorder
}

var _ Register = (*restfulHandler)(nil)

func NewRESTfulHandler(db *db.DB,
	cnf *config.Web,
	dns *config.DNS,
	eventbus *event.EventBus,
	recorder *record.Recorder) Register {

	return &restfulHandler{
		db:       db,
		cnf:      cnf,
		dns:      dns,
		eventbus: eventbus,
		recorder: recorder,
	}
}

func (r *restfulHandler) RegisterHandler(g *gin.Engine) {
	g.Use(r.oobHttp())
	g.NoRoute(func(c *gin.Context) { c.Status(http.StatusNotFound) })

	api := g.Group("/api/v2")
	api.GET("/login", r.login)

	user := api.Group("user")
	user.Use(r.userToken())
	{
		// TODO
		user.GET("/info", r.info)
		user.Any("/record", r.record)
		user.POST("/info", r.setInfo)
		user.POST("/reset_token", r.resetToken)
		user.POST("/logout", r.logout)
	}
}

func (r *restfulHandler) oobHttp() gin.HandlerFunc {
	httplog := oob.NewHTTP(r.dns, r.recorder)

	return func(c *gin.Context) {
		host, _, _ := net.SplitHostPort(c.Request.Host)
		if host != r.dns.Main {
			httplog.Record(c)
			c.Abort()
		}
	}
}

func (r *restfulHandler) userToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			ReturnUnauthorized(c, errUnauthorizedAccess)
			c.Abort()
			return
		}

		user, err := r.db.GetUserByToken(token)
		if err != nil || user == nil {
			ReturnUnauthorized(c, errUnauthorizedAccess)
			c.Abort()
			return
		}

		c.Set("sid", user.Sid)
		c.Set("token", user.Token)
		c.Next()
	}
}
