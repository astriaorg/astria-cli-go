package safebuffer

import (
	"strconv"
	"sync"
	"testing"
)

func TestSafeBuffer_ConcurrentWritesAndReads(t *testing.T) {
	var wg sync.WaitGroup
	sb := &SafeBuffer{}

	numWriters := 10
	numReaders := 5
	writeCount := 100

	// start writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(writerID int) {
			defer wg.Done()
			for j := 0; j < writeCount; j++ {
				_, err := sb.Write([]byte(strconv.Itoa(writerID) + " "))
				if err != nil {
					t.Errorf("Failed to write to safeBuffer: %v", err)
				}
			}
		}(i)
	}

	// start readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = sb.String()
		}()
	}

	wg.Wait()

	// verify buffer content
	finalContent := sb.String()
	expectedLength := numWriters * writeCount * 2 // each write includes an ID and a space
	if len(finalContent) != expectedLength {
		t.Errorf("Expected buffer length %d, got %d", expectedLength, len(finalContent))
	}
}
