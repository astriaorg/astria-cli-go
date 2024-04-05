package safebuffer

import (
	"bytes"
	"sync"
)

// SafeBuffer is a thread/goroutine safe buffer.
type SafeBuffer struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

// Write writes p to the buffer.
func (sb *SafeBuffer) Write(p []byte) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.Write(p)
}

// WriteString writes s to the buffer.
func (sb *SafeBuffer) WriteString(s string) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.WriteString(s)
}

// String returns the contents of the buffer as a string.
func (sb *SafeBuffer) String() string {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.String()
}
