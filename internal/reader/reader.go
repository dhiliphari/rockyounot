package reader

import (
	"bufio"
	"io"
)

type ChunkResult struct {
	Chunk *Chunk
	Err   error
}

func Chunks(file io.Reader, chunkSize int) <-chan ChunkResult {
	if chunkSize <= 0 {
		panic("chunkSize cannot be zero")
	}
	out := make(chan ChunkResult, 100)
	go func() {
		defer close(out)

		scanner := bufio.NewScanner(file)

		for {
			chunk := getChunk()

			for chunk.Len() < chunkSize && scanner.Scan() {
				chunk.add(scanner.Bytes())
			}

			if chunk.Len() == 0 {
				chunk.Release()
				break
			}
			out <- ChunkResult{Chunk: chunk}
		}

		if err := scanner.Err(); err != nil {
			out <- ChunkResult{Err: err}
		}
	}()
	return out
}
