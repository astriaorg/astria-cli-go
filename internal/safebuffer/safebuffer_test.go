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

func TestSafeBuffer_Write(t *testing.T) {
	sb := &SafeBuffer{}

	// Write a single line
	n, err := sb.Write([]byte("This is a test\n"))
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != 15 {
		t.Errorf("Write() got = %v, want %v", n, 15)
	}
	if count := sb.GetLineCount(); count != 1 {
		t.Errorf("GetLineCount() got = %v, want %v", count, 1)
	}

	// Write multiple lines
	_, _ = sb.Write([]byte("Second line.\nThird line.\n"))
	if count := sb.GetLineCount(); count != 3 {
		t.Errorf("GetLineCount() after writing more lines got = %v, want %v", count, 3)
	}
}

func TestSafeBuffer_WriteString(t *testing.T) {
	sb := &SafeBuffer{}

	// WriteString with a single line
	_, _ = sb.WriteString("First line\n")
	if count := sb.GetLineCount(); count != 1 {
		t.Errorf("GetLineCount() got = %v, want %v", count, 1)
	}

	// WriteString with multiple lines
	_, _ = sb.WriteString("Second line\nThird line\n")
	if count := sb.GetLineCount(); count != 3 {
		t.Errorf("GetLineCount() after writing more lines got = %v, want %v", count, 3)
	}
}

func TestSafeBuffer_LineCountAccuracy(t *testing.T) {
	sb := &SafeBuffer{}

	// Write a string without a newline
	_, _ = sb.WriteString("This string does not end with a newline")
	if count := sb.GetLineCount(); count != 0 {
		t.Errorf("GetLineCount() got = %v, want %v for string without newline", count, 0)
	}

	// Write a string that ends with a newline
	_, _ = sb.WriteString("This string ends with a newline\n")
	if count := sb.GetLineCount(); count != 1 {
		t.Errorf("GetLineCount() got = %v, want %v for string with newline", count, 1)
	}

	// Write multiple lines at once without a trailing newline
	_, _ = sb.Write([]byte("Multi-line\nInput\nWithout trailing newline"))
	if count := sb.GetLineCount(); count != 3 {
		t.Errorf("GetLineCount() got = %v, want %v for multi-line input without trailing newline", count, 3)
	}

	// Ensure line count does not decrease
	_, _ = sb.Write([]byte(""))
	if count := sb.GetLineCount(); count != 3 {
		t.Errorf("GetLineCount() got = %v, want %v after writing empty string", count, 3)
	}
}
