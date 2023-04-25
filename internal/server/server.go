package server

import (
	"context"

	"github.com/ac0d3r/hyuga/internal/config"
	"github.com/ac0d3r/hyuga/internal/db"
	"github.com/ac0d3r/hyuga/internal/handler"
	"github.com/ac0d3r/hyuga/pkg/httpx"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, g *errgroup.Group, cnf *config.Web, db *db.DB) {
	gin.SetMode(gin.ReleaseMode)

	g.Go(func() error {
		logrus.Infof("[server] web listen at '%s'", cnf.Address)
		web := httpx.NewBaseGinServer(
			handler.NewRESTfulHandler(db, cnf).RegisterHandler)

		return web.Start(ctx, false, cnf.Address)
	})
}
