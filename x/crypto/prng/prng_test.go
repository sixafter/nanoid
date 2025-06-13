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

	"github.com/stretchr/testify/assert"
)

// TestPRNG_Read verifies that a single Read call fills the buffer completely
// and produces non-zero (random) data.
func TestPRNG_Read(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err, "NewReader should not error")

	buffer := make([]byte, 64)
	n, err := rdr.Read(buffer)
	is.NoError(err, "Read should not error")
	is.Equal(len(buffer), n, "Read should return full buffer length")

	allZeros := true
	for _, b := range buffer {
		if b != 0 {
			allZeros = false
			break
		}
	}
	is.False(allZeros, "Buffer should not be all zeros")
}

// TestPRNG_ReadZeroBytes ensures that reading into a zero-length slice
// returns immediately with a count of zero and no error.
func TestPRNG_ReadZeroBytes(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err)

	buffer := make([]byte, 0)
	n, err := rdr.Read(buffer)
	is.NoError(err, "Reading zero-length buffer should not error")
	is.Equal(0, n, "Reading zero-length buffer should return 0")
}

// TestPRNG_ReadMultipleTimes checks that consecutive Read calls produce
// different outputs, confirming the generator does not repeat data immediately.
func TestPRNG_ReadMultipleTimes(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err)

	buf1 := make([]byte, 32)
	n, err := rdr.Read(buf1)
	is.NoError(err)
	is.Equal(len(buf1), n)

	buf2 := make([]byte, 32)
	n, err = rdr.Read(buf2)
	is.NoError(err)
	is.Equal(len(buf2), n)

	is.False(bytes.Equal(buf1, buf2), "Consecutive reads should differ")
}

// TestPRNG_ReadWithDifferentBufferSizes runs Read with various sizes to
// ensure correctness across a range of buffer lengths.
func TestPRNG_ReadWithDifferentBufferSizes(t *testing.T) {
	t.Parallel()

	sizes := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048}
	for _, size := range sizes {
		size := size
		t.Run(fmt.Sprintf("Size_%d", size), func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			rdr, err := NewReader()
			is.NoError(err)

			buf := make([]byte, size)
			n, err := rdr.Read(buf)
			is.NoError(err)
			is.Equal(size, n)

			allZeros := true
			for _, b := range buf {
				if b != 0 {
					allZeros = false
					break
				}
			}
			is.False(allZeros, "Buffer of size %d should not be all zeros", size)
		})
	}
}

// TestPRNG_Concurrency spawns many goroutines performing Read concurrently
// to verify thread safety and data integrity under parallel usage.
func TestPRNG_Concurrency(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	const (
		numGoroutines = 100
		bufferSize    = 64
	)
	rdr, err := NewReader()
	is.NoError(err)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	errCh := make(chan error, numGoroutines)
	buffers := make([][]byte, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			buf := make([]byte, bufferSize)
			if _, err := rdr.Read(buf); err != nil {
				errCh <- err
				return
			}
			buffers[i] = buf
		}(i)
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		is.NoError(err, "Concurrent Read should not error")
	}

	// Optional uniqueness check (best-effort for randomness)
	for i := 0; i < numGoroutines; i++ {
		for j := i + 1; j < numGoroutines; j++ {
			is.False(bytes.Equal(buffers[i], buffers[j]), "Buffers %d and %d should differ", i, j)
		}
	}
}

// TestPRNG_Stream uses io.ReadFull to read a large stream of data,
// ensuring the Reader can handle substantial volumes without error.
func TestPRNG_Stream(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err)

	const total = 1 << 20 // 1 MiB
	buf := make([]byte, total)
	n, err := io.ReadFull(rdr, buf)
	is.NoError(err)
	is.Equal(total, n)

	allZeros := true
	for _, b := range buf {
		if b != 0 {
			allZeros = false
			break
		}
	}
	is.False(allZeros, "Stream buffer should not be all zeros")
}

// TestPRNG_ReadUnique reads twice into buffers and verifies the outputs differ,
// providing additional confidence in randomness between successive calls.
func TestPRNG_ReadUnique(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err)

	b1 := make([]byte, 64)
	_, err = rdr.Read(b1)
	is.NoError(err)

	b2 := make([]byte, 64)
	_, err = rdr.Read(b2)
	is.NoError(err)

	is.False(bytes.Equal(b1, b2), "Consecutive reads should produce unique data")
}

// TestPRNG_NewReader ensures NewReader returns a non-nil Reader that
// can successfully produce random bytes on Read.
func TestPRNG_NewReader(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err)
	is.NotNil(rdr, "NewReader should return non-nil Reader")

	buf := make([]byte, 32)
	n, err := rdr.Read(buf)
	is.NoError(err)
	is.Equal(len(buf), n)

	allZeros := true
	for _, b := range buf {
		if b != 0 {
			allZeros = false
			break
		}
	}
	is.False(allZeros, "NewReader buffer should not be all zeros")
}

// TestPRNG_ReadAll reads a larger buffer in one call to exercise
// the Reader's ability to fill arbitrary-length slices.
func TestPRNG_ReadAll(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err)

	buf := make([]byte, 10*1024) // 10 KiB
	n, err := rdr.Read(buf)
	is.NoError(err)
	is.Equal(len(buf), n)

	allZeros := true
	for _, b := range buf {
		if b != 0 {
			allZeros = false
			break
		}
	}
	is.False(allZeros, "ReadAll buffer should not be all zeros")
}

// TestPRNG_ReadConsistency performs multiple reads of the same size
// and validates each buffer is filled and buffers differ from one another.
func TestPRNG_ReadConsistency(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	const (
		numReads   = 50
		bufferSize = 128
	)
	rdr, err := NewReader()
	is.NoError(err)

	buffers := make([][]byte, numReads)
	for i := 0; i < numReads; i++ {
		buf := make([]byte, bufferSize)
		n, err := rdr.Read(buf)
		is.NoError(err, "Read %d should not error", i)
		is.Equal(bufferSize, n, "Read %d should fill the buffer", i)

		allZeros := true
		for _, b := range buf {
			if b != 0 {
				allZeros = false
				break
			}
		}
		is.False(allZeros, "Buffer %d should not be all zeros", i)
		buffers[i] = buf
	}

	for i := 0; i < numReads; i++ {
		for j := i + 1; j < numReads; j++ {
			is.False(bytes.Equal(buffers[i], buffers[j]), "Buffers %d and %d should differ", i, j)
		}
	}
}
