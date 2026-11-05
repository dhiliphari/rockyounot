package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/dhiliphari/rockyounot/internal/reader"
)

const CAPACITY = 14_344_392 // equivalent to total passwords exposed by rockyou.txt

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// Setting a 0.1% error rate for 14 million items
	err := rdb.Do(context.Background(), "BF.RESERVE", "passwords", 0.001, CAPACITY).Err()
	if err != nil && err.Error() != "ERR item exists" { // Ignore error if it already exists
		log.Fatalf("error setting up filter: %v", err)
	}

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

		args := make([]any, 0, len(r.Chunk.Lines)+2)
		args = append(args, "BF.MADD", "passwords")
		for _, line := range r.Chunk.Lines {
			args = append(args, line.Data)
		}

		err := rdb.Do(context.Background(), args...).Err()
		if err != nil {
			r.Chunk.Release()
			log.Fatalf("error adding %d passwords to filter: %v", len(r.Chunk.Lines), err)
		}

		r.Chunk.Release()
	}

	log.Printf("Processing completed in %v", time.Since(startTime))
}
