package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/KekemonBS/p2pSharer-back/infrastructure/env"
	"github.com/KekemonBS/p2pSharer-back/router"
	"github.com/KekemonBS/p2pSharer-back/storage/redis"
	"github.com/KekemonBS/p2pSharer-back/tracker"
	rds "github.com/redis/go-redis/v9"
)

func main() {
	//Init config
	cfg, err := env.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	//Init logger
	logger := log.New(os.Stdout, "log:", log.Lshortfile)
	logger.Println("TEST")
	//Open redis connection
	opt, err := rds.ParseURL(cfg.RedisURI)
	if err != nil {
		panic(err)
	}
	client := rds.NewClient(opt)

	cacheImpl := redis.New(client)

	//Init handlers
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		logger.Printf("got signal: %v", <-ch)
		cancel()
	}()
	//Init tracker handlers
	handlers := tracker.New(ctx, logger, cacheImpl)

	router := router.New(handlers)
	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		err = s.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			s.Close()
			return
		}
	}
}
