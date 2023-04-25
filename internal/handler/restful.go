package handler

import (
	"github.com/ac0d3r/hyuga/internal/config"
	"github.com/ac0d3r/hyuga/internal/db"
	"github.com/gin-gonic/gin"
)

type restfulHandler struct {
	db  *db.DB
	cnf *config.Web
}

var _ Register = (*restfulHandler)(nil)

func NewRESTfulHandler(db *db.DB, cnf *config.Web) Register {
	return &restfulHandler{db: db, cnf: cnf}
}

func (r *restfulHandler) RegisterHandler(g *gin.Engine) {
	api := g.Group("/api/v2")
	api.GET("/login", r.login)

	user := api.Group("user")
	user.Use(r.userToken())
	{
		// TODO
		user.GET("/info", r.info)
		user.POST("/reset_token", r.reset_token)
		user.POST("/logout", r.logout)
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

		c.Set("sid", user.ID)
		c.Set("token", user.Token)
		c.Next()
	}
}
