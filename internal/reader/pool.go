package reader

import "sync"

type Chunk struct {
	buf  []byte
	ends []int
}

func (c *Chunk) Len() int {
	return len(c.ends)
}

func (c *Chunk) Line(i int) []byte {
	start := 0
	if i > 0 {
		start = c.ends[i-1]
	}
	end := c.ends[i]
	return c.buf[start:end:end] // cap the slice so appends by the caller can't override the next line
}

func (c *Chunk) add(b []byte) {
	c.buf = append(c.buf, b...)
	c.ends = append(c.ends, len(c.buf))
}

var chunkPool = sync.Pool{
	New: func() any {
		return &Chunk{}
	},
}

func getChunk() *Chunk {
	return chunkPool.Get().(*Chunk)
}

func (c *Chunk) Release() {
	c.buf = c.buf[:0]
	c.ends = c.ends[:0]
	chunkPool.Put(c)
}
