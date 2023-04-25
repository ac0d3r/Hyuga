package httpx

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HttpServer interface {
	Start(ctx context.Context, debug bool, addr string) error
}

type BaseGin struct {
	RegisterRoute func(engine *gin.Engine)
}

func NewBaseGinServer(r func(engine *gin.Engine)) HttpServer {
	return &BaseGin{RegisterRoute: r}
}

var _ HttpServer = (*BaseGin)(nil)

func (s *BaseGin) route(debug bool) http.Handler {
	var engine *gin.Engine

	if debug {
		engine = gin.Default()
	} else {
		engine = gin.New()
		engine.Use(gin.Recovery())
	}

	if s.RegisterRoute != nil {
		s.RegisterRoute(engine)
	}
	return engine
}

func (s *BaseGin) Start(ctx context.Context, debug bool, addr string) (err error) {
	srv := http.Server{
		Addr:    addr,
		Handler: s.route(debug),
	}

	processed := make(chan struct{})
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		err = srv.Shutdown(ctx)
		if err := srv.Shutdown(ctx); nil != err {
			logrus.Warnf("server shutdown failed, err: %v\n", err)
		}
		logrus.Infof("server gracefully shutdown")

		close(processed)
	}()

	err = srv.ListenAndServe()
	if http.ErrServerClosed != err {
		logrus.Warnf("server not gracefully shutdown, err :%v\n", err)
	}
	<-processed
	return err
}
