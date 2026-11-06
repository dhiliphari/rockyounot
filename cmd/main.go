package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/dhiliphari/rockyounot/data"
	"github.com/dhiliphari/rockyounot/internal/stores"
)

type config struct {
	port      int
	redisAddr string
	seedFile  string
	capacity  int
	errorRate float64
	chunkSize int
}

type application struct {
	config *config
	store  stores.PasswordStore
	logger *slog.Logger
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := loadCfg()

	rdb, err := openRDB(ctx, cfg.redisAddr)
	if err != nil {
		logger.Error("error connecting to redis", "error", err.Error())
		os.Exit(1)
	}
	defer rdb.Close()

	store := &stores.RedisStore{
		RDB: rdb,
	}
	err = store.Initialize(ctx, cfg.errorRate, cfg.capacity)
	if err != nil {
		logger.Error("error initializing filter", "errors", err.Error())
		os.Exit(1)
	}

	logger.Info("seeding passwords", "data", cfg.seedFile)

	count, err := data.SeedPasswords(ctx, cfg.seedFile, store, cfg.chunkSize, time.Minute)
	if err != nil {
		logger.Error("seeding failed", "count", count, "error", err.Error())
		os.Exit(1)
	}

	logger.Info("seed completed", "count", count)

	app := &application{
		logger: logger,
		config: cfg,
		store:  store,
	}

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr)

	ch := make(chan error, 1)

	go func() {
		err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
		if err != nil {
			ch <- err
		}
		close(ch)
	}()

	select {
	case err = <-ch:
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		err = srv.Shutdown(timeout)
	}

	if err != nil {
		logger.Error(err.Error())
	}

	os.Exit(1)
}

func loadCfg() *config {
	var cfg config

	flag.IntVar(&cfg.port, "addr", 5001, "server port")

	flag.IntVar(&cfg.chunkSize, "chunkSize", 1000, "reader chunk size")
	flag.IntVar(&cfg.capacity, "capacity", 15_000_000, "capacity to reserve")
	flag.Float64Var(&cfg.errorRate, "errorRate", 0.001, "acceptable error rate")
	flag.StringVar(&cfg.redisAddr, "redisAddr", "localhost:6379", "redis connection address")
	flag.StringVar(&cfg.seedFile, "seedFile", "./data/rockyou.txt", "file containing passwords to seed into the filter")

	flag.Parse()

	return &cfg
}

func openRDB(ctx context.Context, addr string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
