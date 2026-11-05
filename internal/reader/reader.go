package reader

import (
	"bufio"
	"io"
)

type LineResult struct {
	Line []byte
	Err  error
}

func Lines(file io.Reader) <-chan LineResult {
	out := make(chan LineResult, 100)
	go func() {
		defer close(out)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// create a copy to avoid data corruption
			b := scanner.Bytes()
			line := make([]byte, len(b))
			copy(line, b)

			out <- LineResult{Line: line}
		}

		if err := scanner.Err(); err != nil {
			out <- LineResult{Err: err}
		}
	}()
	return out
}
