package safebuffer

import (
	"bytes"
	"sync"
)

// SafeBuffer is a thread/goroutine safe buffer that keeps track of the number of lines in the buffer.
type SafeBuffer struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

// Write writes p to the buffer.
func (sb *SafeBuffer) Write(p []byte) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	n, err = sb.buf.Write(p)
	return n, err
}

// WriteString writes s to the buffer.
func (sb *SafeBuffer) WriteString(s string) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	n, err = sb.buf.WriteString(s)
	return n, err
}

// String returns the contents of the buffer as a string.
func (sb *SafeBuffer) String() string {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.String()
}

// Reset resets the buffer and line count.
func (sb *SafeBuffer) Reset() {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	sb.buf.Reset()
}
