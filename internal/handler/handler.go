package handler

import (
	"strings"

	"hyuga/internal/config"
	"hyuga/internal/db"
	"hyuga/internal/handler/base"
	"hyuga/internal/oob"

	"github.com/gin-gonic/gin"
)

type handler struct {
	cnf *config.Config
	db  *db.Client
}

func New(c *config.Config, db *db.Client) *handler {
	base.Debug = c.DebugMode
	return &handler{
		cnf: c,
		db:  db,
	}
}

func (h *handler) Route() *gin.Engine {
	var engine *gin.Engine
	if h.cnf.DebugMode {
		engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
	}
	engine.SetTrustedProxies(nil)
	engine.Use(h.middlewareForwardLog())

	user := NewUser(h.db)
	user.Route(engine, h.middlewareUserToken())
	record := NewRecord(h.db)
	record.Route(engine, h.middlewareUserToken())
	return engine
}

func (h *handler) middlewareUserToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			token = c.Query("token")
		} else {
			token = strings.TrimPrefix(authorization, "Bearer ")
		}

		if token == "" {
			base.ReturnUnauthorized(c, 2000)
			c.Abort()
			return
		}

		user, err := h.db.GetUserByToken(c.Request.Context(), token)
		if err != nil {
			base.ReturnUnauthorized(c, 2000)
			c.Abort()
			return
		}

		c.Set("uid", user.ID)
		c.Set("token", user.Token)
		c.Next()
	}
}

func (h *handler) middlewareForwardLog() gin.HandlerFunc {
	httplog := oob.NewHTTPLog(&h.cnf.OOB, h.db)

	return func(c *gin.Context) {
		host := strings.Split(c.Request.Host, ":")[0]
		if host != h.cnf.OOB.Dns.Domain {
			httplog.Record(c)
			c.Abort()
		}
	}
}
