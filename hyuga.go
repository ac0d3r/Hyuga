package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"hyuga/conf"
	"hyuga/db"
	"hyuga/internal/core"
	"hyuga/router"
)

func main() {
	var (
		err error
		ctx context.Context = context.Background()
	)
	// conf
	if err = conf.SetFromFile("config.yml"); err != nil {
		log.Fatalln(err)
	}
	db.New()
	// redis
	times := 5
	for {
		if _, err := db.Ping(ctx); err != nil {
			times--
			log.Println(err)
			if times == 0 {
				log.Fatalln("redis connect failed")
			}
			time.Sleep(time.Second)
			continue
		}
		break
	}

	r := router.New()
	dns, _ := core.NewDNSDog(":53")
	go r.Start(":5000")
	go dns.ListenAndServe()

	defer func() {
		db.Close()
		dns.Shutdown()
		r.Shutdown(ctx)
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
