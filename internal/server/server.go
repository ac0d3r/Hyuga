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
	notifier "github.com/ac0d3r/hyuga/thirdparty/notify"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context,
	g *errgroup.Group,
	db *db.DB,
	cnf *config.Web,
	oob_ *config.OOB,
	eventbus *event.EventBus,
	recorder *record.Recorder,
) {

	gin.SetMode(gin.ReleaseMode)

	g.Go(func() error {
		logrus.Infof("[server] web listen at '%s'", cnf.Address)
		web := httpx.NewBaseGinServer(
			handler.NewRESTfulHandler(db, cnf, oob_, eventbus, recorder).RegisterHandler)

		return web.Start(ctx, false, cnf.Address)
	})

	// notify
	g.Go(func() error {
		// subscribe all event
		s := eventbus.Subscribe("*")
		defer eventbus.Unsubscribe(s)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case msg := <-s.Out():
				r, ok := msg.(oob.Record)
				if !ok {
					continue
				}
				logrus.Infof("[server][notify] eventbus '*' receive msg: %v", r)

				u, err := db.GetUserBySid(r.Sid)
				if err == nil && u != nil && u.Notify.Enable {
					logrus.Infof("[server][notify] msg: '%s', '%s'", r.Type.String(), r.Name)

					if u.Notify.Bark.Key != "" {
						notifier.WithBark(u.Notify.Bark.Key, u.Notify.Bark.Server, r.Type.String(), r.Name)
					}
					if u.Notify.Dingtalk.Token != "" && u.Notify.Dingtalk.Secret != "" {
						notifier.WithDingTalk(u.Notify.Dingtalk.Token, u.Notify.Dingtalk.Secret, r.Type.String(), r.Name)
					}
					if u.Notify.Lark.Token != "" && u.Notify.Lark.Secret != "" {
						notifier.WithLark(u.Notify.Lark.Token, u.Notify.Lark.Secret, r.Type.String(), r.Name)
					}
					if u.Notify.Feishu.Token != "" && u.Notify.Feishu.Secret != "" {
						notifier.WithFeishu(u.Notify.Feishu.Token, u.Notify.Feishu.Secret, r.Type.String(), r.Name)
					}
					if u.Notify.ServerChan.UserID != "" && u.Notify.ServerChan.SendKey != "" {
						notifier.WithServerChan(u.Notify.ServerChan.UserID, u.Notify.ServerChan.SendKey, r.Type.String(), r.Name)
					}
				}
			}
		}
	})
}
