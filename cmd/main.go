package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/venomuz/url-shortener-go/config"
	_ "github.com/venomuz/url-shortener-go/docs"
	"github.com/venomuz/url-shortener-go/pkg/logger"
	"github.com/venomuz/url-shortener-go/router"
	"github.com/venomuz/url-shortener-go/storage"
)

var ctx = context.Background()

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "url-short-service")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: "",
		DB:       0,
	})
	redisRepo := storage.NewRedisRepo(rdb)

	server := router.New(router.Option{
		Log:  log,
		Rds:  redisRepo,
		Conf: cfg,
	})
	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", logger.Error(err))
		panic(err)
	}
}
