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

			for len(chunk.Lines) < chunkSize && scanner.Scan() {
				line := getLine()

				line.Data = line.Data[:0]
				line.Data = append(line.Data[:0], scanner.Bytes()...)

				chunk.Lines = append(chunk.Lines, line)
			}

			if len(chunk.Lines) > 0 {
				out <- ChunkResult{Chunk: chunk}
			} else {
				chunk.Release()
				break
			}
		}

		if err := scanner.Err(); err != nil {
			out <- ChunkResult{Err: err}
		}
	}()
	return out
}
