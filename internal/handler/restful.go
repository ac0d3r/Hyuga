package handler

import (
	"embed"
	"io/fs"
	"net"
	"net/http"

	"github.com/ac0d3r/hyuga/internal/config"
	"github.com/ac0d3r/hyuga/internal/db"
	"github.com/ac0d3r/hyuga/internal/oob"
	"github.com/ac0d3r/hyuga/internal/record"
	"github.com/ac0d3r/hyuga/pkg/event"
	"github.com/gin-contrib/static"
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

//go:embed dist
var dist embed.FS

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	return err == nil
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}

func (r *restfulHandler) RegisterHandler(g *gin.Engine) {
	g.Use(r.oobHttp())                                  // oob http log
	g.Use(static.Serve("/", EmbedFolder(dist, "dist"))) // static file
	g.NoRoute(func(c *gin.Context) { c.Status(http.StatusNotFound) })

	v2 := g.Group("/api/v2")
	v2.GET("/login", r.login)

	user := v2.Group("user")
	user.Use(r.userToken())
	{
		user.GET("/info", r.info)      // get user info
		user.Any("/record", r.record)  // stream record
		user.POST("/notify", r.notify) // set user notify
		user.POST("/reset", r.reset)   // reset user api token
		user.POST("/logout", r.logout) // logout
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
