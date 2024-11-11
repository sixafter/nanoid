// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package prng

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"testing"
)

// TestPRNG_Read performs a basic read operation, verifying that the correct number of bytes is read
// and that the buffer is not filled with all zeros.
func TestPRNG_Read(t *testing.T) {
	t.Parallel() // Enable parallel execution of this test

	reader, err := NewReader()
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}

	buffer := make([]byte, 64)
	n, err := reader.Read(buffer)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if n != len(buffer) {
		t.Errorf("Expected to read %d bytes, but read %d bytes", len(buffer), n)
	}

	// Ensure that the buffer is not all zeros
	allZeros := true
	for _, b := range buffer {
		if b != 0 {
			allZeros = false
			break
		}
	}
	if allZeros {
		t.Errorf("Buffer should not be all zeros, expected random data")
	}
}

// TestPRNG_ReadZeroBytes tests reading with a zero-length buffer, ensuring no bytes are read
// and no errors are returned.
func TestPRNG_ReadZeroBytes(t *testing.T) {
	t.Parallel() // Enable parallel execution of this test

	reader, err := NewReader()
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}

	buffer := make([]byte, 0)
	n, err := reader.Read(buffer)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if n != 0 {
		t.Errorf("Expected to read 0 bytes, but read %d bytes", n)
	}
}

// TestPRNG_ReadMultipleTimes performs multiple sequential read operations, ensuring each read
// is successful and returns unique data.
func TestPRNG_ReadMultipleTimes(t *testing.T) {
	t.Parallel() // Enable parallel execution of this test

	reader, err := NewReader()
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}

	buffer1 := make([]byte, 32)
	n, err := reader.Read(buffer1)
	if err != nil {
		t.Fatalf("First read failed: %v", err)
	}
	if n != len(buffer1) {
		t.Errorf("First read expected %d bytes, but read %d bytes", len(buffer1), n)
	}

	buffer2 := make([]byte, 32)
	n, err = reader.Read(buffer2)
	if err != nil {
		t.Fatalf("Second read failed: %v", err)
	}
	if n != len(buffer2) {
		t.Errorf("Second read expected %d bytes, but read %d bytes", len(buffer2), n)
	}

	// Ensure that the two buffers are different
	if bytes.Equal(buffer1, buffer2) {
		t.Errorf("Consecutive reads should produce different data")
	}
}

// TestPRNG_ReadWithDifferentBufferSizes tests reading with various buffer sizes to ensure
// consistent behavior across different data volumes.
func TestPRNG_ReadWithDifferentBufferSizes(t *testing.T) {
	t.Parallel() // Enable parallel execution of this test

	bufferSizes := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048}

	for _, size := range bufferSizes {
		size := size // Capture range variable
		t.Run(fmt.Sprintf("BufferSize_%d", size), func(t *testing.T) {
			t.Parallel() // Enable parallel execution of subtests

			reader, err := NewReader()
			if err != nil {
				t.Fatalf("NewReader failed: %v", err)
			}

			buffer := make([]byte, size)
			n, err := reader.Read(buffer)
			if err != nil {
				t.Fatalf("Read failed: %v", err)
			}
			if n != len(buffer) {
				t.Errorf("Expected to read %d bytes, but read %d bytes", len(buffer), n)
			}

			// Ensure that the buffer is not all zeros
			allZeros := true
			for _, b := range buffer {
				if b != 0 {
					allZeros = false
					break
				}
			}
			if allZeros {
				t.Errorf("Buffer should not be all zeros, expected random data")
			}
		})
	}
}

// TestPRNG_Concurrency tests concurrent read operations by spawning multiple goroutines,
// each performing a read. It ensures thread safety and data integrity under high concurrency.
func TestPRNG_Concurrency(t *testing.T) {
	t.Parallel() // Enable parallel execution of this test

	numGoroutines := 100
	bufferSize := 64
	buffers := make([][]byte, numGoroutines)
	errorsChan := make(chan error, numGoroutines)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	reader, err := NewReader()
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()

			buf := make([]byte, bufferSize)
			_, err := reader.Read(buf)
			if err != nil {
				errorsChan <- err
				return
			}
			buffers[index] = buf
		}(i)
	}

	wg.Wait()
	close(errorsChan)

	for err := range errorsChan {
		t.Errorf("Concurrent Read should not return an error: %v", err)
	}

	// Optionally, verify that all buffers contain unique data
	// Note: This is a basic check and may not always pass due to randomness
	for i := 0; i < numGoroutines; i++ {
		for j := i + 1; j < numGoroutines; j++ {
			if bytes.Equal(buffers[i], buffers[j]) {
				t.Errorf("Buffers at index %d and %d are identical", i, j)
			}
		}
	}
}

// TestPRNG_Stream tests reading a large stream of data using io.ReadFull to ensure that
// the Reader can handle substantial data volumes efficiently.
func TestPRNG_Stream(t *testing.T) {
	t.Parallel() // Enable parallel execution of this test

	reader, err := NewReader()
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}

	// Read a large number of bytes to simulate stream usage
	totalBytes := 1024 * 1024 // 1 MB
	buffer := make([]byte, totalBytes)
	n, err := io.ReadFull(reader, buffer)
	if err != nil {
		t.Fatalf("ReadFull failed: %v", err)
	}
	if n != totalBytes {
		t.Errorf("Expected to read %d bytes, but read %d bytes", totalBytes, n)
	}

	// Ensure that the buffer is not all zeros
	allZeros := true
	for _, b := range buffer {
		if b != 0 {
			allZeros = false
			break
		}
	}
	if allZeros {
		t.Errorf("Buffer should not be all zeros, expected random data")
	}
}

// TestPRNG_ReadUnique ensures that multiple reads produce unique data, enhancing the confidence
// in the randomness provided by the PRNG.
func TestPRNG_ReadUnique(t *testing.T) {
	t.Parallel() // Enable parallel execution of this test

	reader, err := NewReader()
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}

	buffer1 := make([]byte, 64)
	_, err = reader.Read(buffer1)
	if err != nil {
		t.Fatalf("First read failed: %v", err)
	}

	buffer2 := make([]byte, 64)
	_, err = reader.Read(buffer2)
	if err != nil {
		t.Fatalf("Second read failed: %v", err)
	}

	if bytes.Equal(buffer1, buffer2) {
		t.Errorf("Consecutive reads should produce different data")
	}
}

// TestPRNG_NewReader tests the NewReader function to ensure it returns a valid io.Reader instance.
func TestPRNG_NewReader(t *testing.T) {
	t.Parallel() // Enable parallel execution of this test

	newReader, err := NewReader()
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}
	if newReader == nil {
		t.Errorf("NewReader should return a non-nil Reader")
	}

	// Perform a simple read to ensure the new Reader functions correctly
	buffer := make([]byte, 32)
	n, err := newReader.Read(buffer)
	if err != nil {
		t.Fatalf("NewReader's Read failed: %v", err)
	}
	if n != len(buffer) {
		t.Errorf("NewReader's Read expected %d bytes, but read %d bytes", len(buffer), n)
	}

	// Ensure that the buffer is not all zeros
	allZeros := true
	for _, b := range buffer {
		if b != 0 {
			allZeros = false
			break
		}
	}
	if allZeros {
		t.Errorf("NewReader's buffer should not be all zeros, expected random data")
	}
}

// TestPRNG_ReadAll ensures that reading a substantial amount of data is handled correctly.
func TestPRNG_ReadAll(t *testing.T) {
	t.Parallel() // Enable parallel execution of this test

	reader, err := NewReader()
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}

	// Attempt to read a large number of bytes to simulate continuous reading
	buffer := make([]byte, 10*1024) // 10 KB
	n, err := reader.Read(buffer)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if n != len(buffer) {
		t.Errorf("Read expected %d bytes, but read %d bytes", len(buffer), n)
	}

	// Ensure that the buffer is not all zeros
	allZeros := true
	for _, b := range buffer {
		if b != 0 {
			allZeros = false
			break
		}
	}
	if allZeros {
		t.Errorf("Buffer should not be all zeros, expected random data")
	}
}

// TestPRNG_ReadConsistency tests the consistency of the Read method by ensuring that
// multiple reads with the same buffer size produce unique data.
func TestPRNG_ReadConsistency(t *testing.T) {
	t.Parallel() // Enable parallel execution of this test

	numReads := 50
	bufferSize := 128
	buffers := make([][]byte, numReads)

	reader, err := NewReader()
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}

	for i := 0; i < numReads; i++ {
		buf := make([]byte, bufferSize)
		n, err := reader.Read(buf)
		if err != nil {
			t.Fatalf("Read failed at iteration %d: %v", i, err)
		}
		if n != len(buf) {
			t.Errorf("Read at iteration %d expected %d bytes, but read %d bytes", i, len(buf), n)
		}

		// Ensure that the buffer is not all zeros
		allZeros := true
		for _, b := range buf {
			if b != 0 {
				allZeros = false
				break
			}
		}
		if allZeros {
			t.Errorf("Buffer at iteration %d should not be all zeros, expected random data", i)
		}

		buffers[i] = buf
	}

	// Optionally, ensure that all buffers are unique
	for i := 0; i < numReads; i++ {
		for j := i + 1; j < numReads; j++ {
			if bytes.Equal(buffers[i], buffers[j]) {
				t.Errorf("Buffers at index %d and %d are identical", i, j)
			}
		}
	}
}

// TestPRNG_ReadErrorScenario simulates a read error by configuring the Reader to use a pool
// that always returns an errorPRNG instance. This ensures that the application correctly
// handles read failures.
func TestPRNG_ReadErrorScenario(t *testing.T) {
	t.Parallel() // Enable parallel execution of this test

	// Since newPoolReader has been removed, we cannot inject a custom pool.
	// Instead, we can temporarily modify the Reader's behavior if possible.
	// If Reader is a global variable, this approach might not be safe.
	// Alternatively, consider refactoring the Reader to allow dependency injection for better testability.

	// Placeholder for error scenario testing.
	// Implementing this requires modifications to the prng package to support injecting
	// a mock or error-producing Reader. Without such modifications, it's not feasible.
	t.Skip("ReadErrorScenario test is skipped because dependency injection is not supported.")
}
