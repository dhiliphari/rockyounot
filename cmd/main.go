package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/dhiliphari/rockyounot/internal/reader"
	"github.com/dhiliphari/rockyounot/internal/stores"
)

const CAPACITY = 14_344_392 // equivalent to total passwords exposed by rockyou.txt

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	store := &stores.RedisStore{
		RDB: rdb,
	}
	// Setting a 0.1% error rate for 14 million items
	store.Initialize(context.Background(), 0.001, CAPACITY)

	startTime := time.Now()

	file, err := os.Open("./data/rockyou.txt")
	if err != nil {
		log.Fatalf("error reading passwords file: %v", err)
	}
	defer file.Close()

	results := reader.Chunks(file, 1000)

	for r := range results {
		if r.Err != nil {
			log.Fatal(r.Err.Error())
		}

		passwords := make([]any, 0, r.Chunk.Len())
		for i := range r.Chunk.Len() {
			passwords = append(passwords, r.Chunk.Line(i))
		}

		err := store.Add(context.Background(), passwords)
		if err != nil {
			r.Chunk.Release()
			log.Fatalf("error adding %d passwords to filter: %v", r.Chunk.Len(), err)
		}

		r.Chunk.Release()
	}

	log.Printf("Processing completed in %v", time.Since(startTime))
}
