package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"hyuga/internal/config"
	"hyuga/internal/db"
	"hyuga/internal/handler"
	"hyuga/internal/oob"

	"github.com/sirupsen/logrus"
)

type App struct {
	ctx    context.Context
	cancel func()

	cnf  *config.Config
	sigs []os.Signal
}

func New(c *config.Config) (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		ctx:    ctx,
		cancel: cancel,
		cnf:    c,
		sigs:   []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}, nil
}

func (a *App) Run() error {
	var (
		err error
		wg  sync.WaitGroup
	)

	// db client
	dbc, err := db.New(a.cnf)
	if err != nil {
		return err
	}

	// http server
	h := handler.New(a.cnf, dbc)
	engine := h.Route()
	api := &http.Server{
		Addr:    a.cnf.Api.Address,
		Handler: engine,
	}
	dns := oob.NewDns(&a.cnf.OOB, dbc)
	jndi := oob.NewJndi(&a.cnf.OOB, dbc)

	wg.Add(4)
	go func() {
		defer wg.Done()
		if err = api.ListenAndServe(); err != nil {
			logrus.Warnf("api server listen error", err)
		}
	}()
	go func() {
		defer wg.Done()
		if err = dns.ListenAndServe(); err != nil {
			logrus.Warnf("dns server listen error", err)
		}
	}()
	go func() {
		defer wg.Done()
		err = jndi.ListenAndServe()
	}()

	go func() {
		defer wg.Done()
		<-a.ctx.Done()
		if err = api.Shutdown(a.ctx); err != nil {
			logrus.Warnf("api server shurdown error", err)
		}
		if err = dns.Shutdown(); err != nil {
			logrus.Warnf("dns server shurdown error", err)
		}
		if err = jndi.Shutdown(); err != nil {
			logrus.Warnf("jndi server shurdown error", err)
		}
	}()

	// signel
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.sigs...)
	go func() {
		<-c
		if a.cancel != nil {
			a.cancel()
		}
	}()

	wg.Wait()
	return err
}
