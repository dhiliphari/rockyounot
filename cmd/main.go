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

	lines := reader.Lines(file)

	for line := range lines {
		if line.Err != nil {
			log.Fatal(line.Err.Error())
		}
		err := rdb.Do(context.Background(), "BF.ADD", "passwords", line.Line).Err()
		if err != nil {
			log.Printf("error adding password to filter: %v", err)
		}
	}

	log.Printf("Processing completed in %v", time.Since(startTime))
}
