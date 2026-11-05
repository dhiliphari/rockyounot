package reader

import "sync"

type Line struct {
	Data []byte
}

type Chunk struct {
	Lines []*Line
}

var linePool = sync.Pool{
	New: func() any {
		return &Line{}
	},
}

var chunkPool = sync.Pool{
	New: func() any {
		return &Chunk{
			Lines: make([]*Line, 0),
		}
	},
}

func getLine() *Line {
	return linePool.Get().(*Line)
}

func getChunk() *Chunk {
	return chunkPool.Get().(*Chunk)
}

func (line *Line) Release() {
	line.Data = line.Data[:0]
	linePool.Put(line)
}

func (chunk *Chunk) Release() {
	for _, line := range chunk.Lines {
		line.Release()
	}
	chunk.Lines = chunk.Lines[:0]
	chunkPool.Put(chunk)
}
