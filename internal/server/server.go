package server

import (
	"context"

	"github.com/ac0d3r/hyuga/internal/config"
	"github.com/ac0d3r/hyuga/internal/db"
	"github.com/ac0d3r/hyuga/internal/handler"
	"github.com/ac0d3r/hyuga/internal/oob"
	"github.com/ac0d3r/hyuga/internal/record"
	"github.com/ac0d3r/hyuga/pkg/event"
	"github.com/ac0d3r/hyuga/pkg/httpx"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context,
	g *errgroup.Group,
	db *db.DB,
	cnf *config.Web,
	dns *config.DNS,
	eventbus *event.EventBus,
	recorder *record.Recorder,
) {

	gin.SetMode(gin.ReleaseMode)

	g.Go(func() error {
		logrus.Infof("[server] web listen at '%s'", cnf.Address)
		web := httpx.NewBaseGinServer(
			handler.NewRESTfulHandler(db, cnf, dns, eventbus, recorder).RegisterHandler)

		return web.Start(ctx, false, cnf.Address)
	})

	// notify
	g.Go(func() error {
		s := eventbus.Subscribe("*")
		defer eventbus.Unsubscribe(s)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case msg := <-s.Out():
				r, ok := msg.(oob.OOBRecord)
				if !ok {
					continue
				}
				logrus.Infof("[server][notify] eventbus '*' receive msg: %v", r)

				u, err := db.GetUserBySid(r.Sid)
				if err == nil && u != nil && u.Notify.Enable {
					// TODO: send notify
					logrus.Infof("[server][notify] TODO msg: %v", r)
				}
			}
		}
	})
}
