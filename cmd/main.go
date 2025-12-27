package main

import (
	"flag"
	"time"

	"github.com/ac0d3r/hyuga/internal/app"
	"github.com/ac0d3r/hyuga/internal/config"
	"github.com/ac0d3r/hyuga/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	buildstamp string = time.Now().Format("2006-01-02 15:04:05")
	githash    string = "dev"
)

func main() {
	var (
		configPath string
	)

	flag.StringVar(&configPath, "config", "../configs/config.toml", "hyuga config path")
	flag.Parse()

	cnf, err := config.Load(configPath)
	if err != nil {
		panic(err)
	}

	if err := logger.Init(cnf.Logger); err != nil {
		panic(err)
	}

	logrus.Infof("build info, build-stamp:%s build-hash:%s", buildstamp, githash)

	app, err := app.New(cnf)
	if err != nil {
		panic(err)
	}

	if err = app.Run(); err != nil {
		panic(err)
	}
}
