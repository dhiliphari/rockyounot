package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/dhiliphari/rockyounot/data"
	"github.com/dhiliphari/rockyounot/internal/stores"
)

type config struct {
	redisAddr string
	seedFile  string
	capacity  int
	errorRate float64
	chunkSize int
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := loadCfg()

	store, rdb := createStore(cfg)
	defer rdb.Close()

	logger.Info("seeding passwords", "data", cfg.seedFile)
	count, err := data.SeedPasswords(cfg.seedFile, store, cfg.chunkSize, time.Minute)
	if err != nil {
		logger.Error("seeding failed", "count", count, "error", err.Error())
		os.Exit(1)
	}
	logger.Info("seed completed", "count", count)
}

func loadCfg() *config {
	var cfg config

	flag.IntVar(&cfg.capacity, "capacity", 15_000_000, "capacity to reserve")
	flag.Float64Var(&cfg.errorRate, "errorRate", 0.001, "acceptable error rate")
	flag.IntVar(&cfg.chunkSize, "chunkSize", 1000, "reader chunk size")
	flag.StringVar(&cfg.redisAddr, "redisAddr", "localhost:6379", "redis connection address")
	flag.StringVar(&cfg.seedFile, "seedFile", "./data/rockyou.txt", "file containing passwords to seed into the filter")

	flag.Parse()

	return &cfg
}

func createStore(cfg *config) (stores.PasswordStore, *redis.Client) {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.redisAddr,
	})

	store := &stores.RedisStore{
		RDB: rdb,
	}
	store.Initialize(context.Background(), cfg.errorRate, cfg.capacity)

	return store, rdb
}
