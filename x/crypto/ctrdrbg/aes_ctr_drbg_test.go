// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.
//
// Tests for ctrdrbg: validates AES-CTR-DRBG output, uniqueness, concurrency, async rekey, personalization.

package ctrdrbg

import (
	"bytes"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_CTRDRBG_Read validates that a single call to Read fills the buffer with nonzero, random data.
func Test_CTRDRBG_Read(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err)

	buf := make([]byte, 64)
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
	is.False(allZeros, "Buffer should not be all zeros")
}

// Test_CTRDRBG_ReadZeroBytes ensures that reading into a zero-length slice is a no-op.
func Test_CTRDRBG_ReadZeroBytes(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err)

	buf := make([]byte, 0)
	n, err := rdr.Read(buf)
	is.NoError(err)
	is.Equal(0, n)
}

// Test_CTRDRBG_ReadMultipleTimes checks that consecutive Read calls yield different results.
func Test_CTRDRBG_ReadMultipleTimes(t *testing.T) {
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

// Test_CTRDRBG_ReadWithDifferentBufferSizes ensures correct output across a range of buffer sizes.
func Test_CTRDRBG_ReadWithDifferentBufferSizes(t *testing.T) {
	t.Parallel()

	sizes := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048}
	for _, size := range sizes {
		size := size
		t.Run("Size_"+string(rune(size)), func(t *testing.T) {
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

// Test_CTRDRBG_Concurrency verifies thread safety under heavy concurrency.
func Test_CTRDRBG_Concurrency(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	const numGoroutines = 100
	const bufferSize = 64

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

	// Optional uniqueness check: at least two buffers should differ
	unique := false
outer:
	for i := 0; i < numGoroutines; i++ {
		for j := i + 1; j < numGoroutines; j++ {
			if !bytes.Equal(buffers[i], buffers[j]) {
				unique = true
				break outer
			}
		}
	}
	is.True(unique, "At least two buffers should differ")
}

// Test_CTRDRBG_Stream tests reading a large (1 MiB) buffer with io.ReadFull.
func Test_CTRDRBG_Stream(t *testing.T) {
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

// Test_CTRDRBG_ReadAll validates that very large requests succeed and return random data.
func Test_CTRDRBG_ReadAll(t *testing.T) {
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

// Test_CTRDRBG_ReadConsistency checks that multiple reads are unique and filled.
func Test_CTRDRBG_ReadConsistency(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	const numReads = 50
	const bufferSize = 128

	rdr, err := NewReader()
	is.NoError(err)

	buffers := make([][]byte, numReads)
	for i := 0; i < numReads; i++ {
		buf := make([]byte, bufferSize)
		n, err := rdr.Read(buf)
		is.NoError(err)
		is.Equal(bufferSize, n)

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
	// Ensure at least two reads differ
	unique := false
outer:
	for i := 0; i < numReads; i++ {
		for j := i + 1; j < numReads; j++ {
			if !bytes.Equal(buffers[i], buffers[j]) {
				unique = true
				break outer
			}
		}
	}
	is.True(unique, "At least two buffers should differ")
}

// Test_CTRDRBG_AsyncRekey validates that async rekeying occurs when MaxBytesPerKey is exceeded.
func Test_CTRDRBG_AsyncRekey(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	cfg.MaxBytesPerKey = 64                  // Small threshold to trigger rekey
	cfg.RekeyBackoff = 10 * time.Millisecond // Speed up test
	cfg.MaxRekeyAttempts = 3
	cfg.MaxInitRetries = 3
	cfg.EnableKeyRotation = true

	// DRBG instance
	d, err := newDRBG(&cfg)
	is.NoError(err)

	initialBlock := d.block

	buf := make([]byte, 128) // Exceeds MaxBytesPerKey
	_, err = d.Read(buf)
	is.NoError(err)

	// Wait for async rekey to finish
	wait := time.NewTimer(500 * time.Millisecond)
	tick := time.NewTicker(10 * time.Millisecond)
	defer wait.Stop()
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if d.block != initialBlock && d.usage == 0 {
				return // success
			}
		case <-wait.C:
			t.Fatal("Timed out waiting for asyncRekey to complete")
		}
	}
}

// Test_CTRDRBG_Personalization_Changes_Stream verifies that two DRBGs constructed with
// different personalization parameters yield distinct output streams.
func Test_CTRDRBG_Personalization_Changes_Stream(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1, err := NewReader(WithPersonalization([]byte("foo")))
	is.NoError(err)
	r2, err := NewReader(WithPersonalization([]byte("bar")))
	is.NoError(err)

	buf1 := make([]byte, 64)
	buf2 := make([]byte, 64)

	_, err = r1.Read(buf1)
	is.NoError(err)
	_, err = r2.Read(buf2)
	is.NoError(err)

	is.False(bytes.Equal(buf1, buf2), "Personalization should affect output")
}
