package app

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/ac0d3r/hyuga/internal/config"
	"github.com/ac0d3r/hyuga/internal/db"
	"github.com/ac0d3r/hyuga/internal/oob"
	"github.com/ac0d3r/hyuga/internal/server"
	"golang.org/x/sync/errgroup"
)

type App struct {
	db       *db.DB
	recorder *db.Recorder
	cnf      *config.Config
}

var errOSSignal = errors.New("os signal")

func New(cnf *config.Config) (*App, error) {
	db_, err := db.NewDB(cnf.DB)
	if err != nil {
		return nil, err
	}

	return &App{
		db:       db_,
		recorder: db.NewRecorder(),
		cnf:      cnf,
	}, nil
}

func (a *App) Run() (err error) {
	defer func() {
		a.db.Close()
	}()

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c:
			return errOSSignal
		}
	})

	server.Run(ctx, g, a.db, a.cnf.Web, &a.cnf.OOB.DNS, a.recorder)
	oob.Run(ctx, g, a.db, a.cnf.OOB, a.recorder)

	err = g.Wait()
	if errors.Is(err, errOSSignal) {
		return nil
	}
	return err
}
