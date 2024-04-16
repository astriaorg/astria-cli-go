package safebuffer

import (
	"bytes"
	"strings"
	"sync"
	"sync/atomic"
)

// SafeBuffer is a thread/goroutine safe buffer that keeps track of the number of lines in the buffer.
type SafeBuffer struct {
	buf bytes.Buffer
	mu  sync.Mutex

	// lineCount is the number of lines in the buffer. This is not the same as buf.Len()
	lineCount int64
}

// Write writes p to the buffer.
func (sb *SafeBuffer) Write(p []byte) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	n, err = sb.buf.Write(p)
	if err == nil {
		atomic.AddInt64(&sb.lineCount, int64(bytes.Count(p, []byte("\n"))))
	}

	return n, err
}

// WriteString writes s to the buffer.
func (sb *SafeBuffer) WriteString(s string) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	n, err = sb.buf.WriteString(s)
	if err == nil {
		atomic.AddInt64(&sb.lineCount, int64(strings.Count(s, "\n")))
	}

	return n, err
}

// String returns the contents of the buffer as a string.
func (sb *SafeBuffer) String() string {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.String()
}

// GetLineCount returns the number of lines in the buffer.
func (sb *SafeBuffer) GetLineCount() int64 {
	return atomic.LoadInt64(&sb.lineCount)
}

// Reset resets the buffer and line count.
func (sb *SafeBuffer) Reset() {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.buf.Reset()
	atomic.StoreInt64(&sb.lineCount, 0)
}
