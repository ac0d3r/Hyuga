package main

import (
	"flag"
	"fmt"
	"time"

	"hyuga/internal/app"
	"hyuga/internal/config"
	"hyuga/pkg/logger"
)

var (
	buildstamp string = time.Now().Format("2006-01-02_15:04:05")
	githash    string = "dev"
	lasttag    string = "dev"
)

func main() {
	fmt.Printf(`anylog build info:
	buildstamp: %s 
	githash: %s
	lasttag: %s 
`, buildstamp, githash, lasttag)

	var (
		configPath string
	)

	flag.StringVar(&configPath, "config", "../configs/config.yaml", "hyuga config path")
	flag.Parse()

	cnf, err := config.Load(configPath)
	if err != nil {
		panic(err)
	}

	if err := logger.Init(cnf.Logger); err != nil {
		panic(err)
	}

	a, err := app.New(cnf)
	if err != nil {
		panic(err)
	}

	if err = a.Run(); err != nil {
		panic(err)
	}
}
