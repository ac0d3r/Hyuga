package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"hyuga/config"
	"hyuga/database"
	"hyuga/handler/router"
	"hyuga/oob"
)

func main() {
	if err := config.SetFromYaml("config.yaml"); err != nil {
		log.Fatalln(err)
	}
	if err := database.Init(config.RedisDsn); err != nil {
		log.Fatal(err)
	}

	http_ := &http.Server{
		Addr:    ":8000",
		Handler: router.Router(),
	}

	dns := oob.NewDnsServer(":53")
	jndi := oob.NewJndiServer(":8881")
	go dns.ListenAndServe()
	go http_.ListenAndServe()
	go jndi.ListenAndServe()

	defer func() {
		log.Println("[dns] shutdown", dns.Shutdown())
		log.Println("[http] shutdown", http_.Shutdown(context.Background()))
		log.Println("[jndi] shutdown", jndi.Shutdown())
		jndi.Wait()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
