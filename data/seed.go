package data

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dhiliphari/rockyounot/internal/reader"
	"github.com/dhiliphari/rockyounot/internal/stores"
)

func SeedPasswords(seedFile string, store stores.PasswordStore, chunkSize int, timeout time.Duration) (int, error) {
	file, err := os.Open(seedFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var count int
	results := reader.Chunks(file, chunkSize)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case r, ok := <-results:
			if !ok {
				return count, nil
			}

			if r.Err != nil {
				return count, r.Err
			}

			count += r.Chunk.Len()

			passwords := make([]any, 0, r.Chunk.Len())
			for i := range r.Chunk.Len() {
				passwords = append(passwords, r.Chunk.Line(i))
			}

			err := store.Add(ctx, passwords)
			if err != nil {
				r.Chunk.Release()
				return count, err
			}

			r.Chunk.Release()
		case <-ctx.Done():
			return count, fmt.Errorf("seeding timedout: %w", ctx.Err())
		}
	}
}
