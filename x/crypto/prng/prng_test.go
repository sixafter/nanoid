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
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/crypto/chacha20"

	"github.com/stretchr/testify/assert"
)

// Test_PRNG_Read validates that a single call to Read fills the buffer with random (non-zero) data.
// It ensures the Reader returns the expected number of bytes and that output is not all zeros.
func Test_PRNG_Read(t *testing.T) {
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

// Test_PRNG_ReadZeroBytes ensures that reading into a zero-length slice is a no-op:
// it should return immediately with no error and report zero bytes read.
func Test_PRNG_ReadZeroBytes(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rdr, err := NewReader()
	is.NoError(err)

	buffer := make([]byte, 0)
	n, err := rdr.Read(buffer)
	is.NoError(err, "Reading zero-length buffer should not error")
	is.Equal(0, n, "Reading zero-length buffer should return 0")
}

// Test_PRNG_ReadMultipleTimes confirms that consecutive Read calls
// produce different outputs, ensuring the PRNG does not repeat data.
func Test_PRNG_ReadMultipleTimes(t *testing.T) {
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

// Test_PRNG_ReadWithDifferentBufferSizes tests Read with a variety of buffer sizes,
// ensuring the Reader works correctly across a range of input slice lengths.
func Test_PRNG_ReadWithDifferentBufferSizes(t *testing.T) {
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

// Test_PRNG_Concurrency verifies the thread safety of the Reader by launching
// many concurrent goroutines that perform Read operations in parallel.
// It asserts that all goroutines succeed and that output buffers are not all identical.
func Test_PRNG_Concurrency(t *testing.T) {
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

// Test_PRNG_Stream checks that the Reader can handle large requests (e.g., 1 MiB)
// via io.ReadFull, and that the returned buffer contains non-zero (random) data.
func Test_PRNG_Stream(t *testing.T) {
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

// Test_PRNG_ReadUnique verifies that two consecutive Read calls to the Reader
// fill buffers with different random data, reinforcing correct PRNG behavior.
func Test_PRNG_ReadUnique(t *testing.T) {
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

// Test_PRNG_NewReader ensures NewReader returns a non-nil Reader instance
// and that the Reader can fill a buffer with random, non-zero bytes.
func Test_PRNG_NewReader(t *testing.T) {
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

// Test_PRNG_ReadAll reads a large buffer in a single call to ensure the Reader
// can fill arbitrary-length slices and returns random, non-zero data.
func Test_PRNG_ReadAll(t *testing.T) {
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

// Test_PRNG_ReadConsistency performs multiple reads of the same size
// and checks that all buffers are filled and differ from each other,
// verifying both completeness and randomness of output.
func Test_PRNG_ReadConsistency(t *testing.T) {
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

// Test_PRNG_AsyncRekey validates the asynchronous rekeying mechanism of the PRNG implementation.
// It configures a low MaxBytesPerKey to force rekeying, writes enough data to trigger it,
// and then checks that a new cipher is used and usage counter is reset, all within a timeout.
func Test_PRNG_AsyncRekey(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	cfg.MaxBytesPerKey = 64                  // small threshold to trigger rekey
	cfg.RekeyBackoff = 10 * time.Millisecond // speed up test execution
	cfg.MaxRekeyAttempts = 3
	cfg.MaxInitRetries = 3
	cfg.EnableKeyRotation = true

	p, err := newPRNG(&cfg)
	is.NoError(err, "newPRNG should not error")

	initialCipher := p.cipher.Load().(*chacha20.Cipher)

	// Write a large enough buffer to exceed MaxBytesPerKey and trigger rekey
	buf := make([]byte, 128)
	_, err = p.Read(buf)
	is.NoError(err)

	// Wait up to 500ms for the async rekey to complete
	wait := time.NewTimer(500 * time.Millisecond)
	tick := time.NewTicker(10 * time.Millisecond)
	defer wait.Stop()
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			// Check if cipher was replaced and usage reset
			currentCipher := p.cipher.Load().(*chacha20.Cipher)
			currentUsage := atomic.LoadUint64(&p.usage)
			if currentCipher != initialCipher && currentUsage == 0 {
				return // success
			}
		case <-wait.C:
			t.Fatal("Timed out waiting for asyncRekey to complete")
		}
	}
}

// Test_PRNG_Read_Shards verifies that a single call to Read only accesses
// one shard pool out of many, regardless of the pool count. It does not
// assert *which* shard is selected, as shardIndex is intentionally random.
//
// This test is table-driven: it runs the check with a variety of pool counts
// to ensure correct behavior at boundaries and typical values.
func Test_PRNG_Read_Shards(t *testing.T) {
	t.Parallel()

	// Define table of test cases with different pool (shard) counts.
	testCases := []struct {
		name       string
		shardCount int
	}{
		{"SinglePool", 1},
		{"TwoPools", 2},
		{"EightPools", 8},
		{"SixteenPools", 16},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)

			// hit[i] will be set true if pool[i] is accessed
			hit := make([]bool, tc.shardCount)

			// Create sync.Pool array, each tracking access via hit[i]
			pools := make([]*sync.Pool, tc.shardCount)
			for i := 0; i < tc.shardCount; i++ {
				id := i
				pools[i] = &sync.Pool{
					New: func() any {
						// Record that this shard was used.
						hit[id] = true
						cfg := DefaultConfig()
						d, _ := newPRNG(&cfg)
						return d
					},
				}
			}

			r := &reader{
				pools: pools,
			}

			buf := make([]byte, 32)
			_, err := r.Read(buf)
			is.NoError(err)

			// Ensure exactly one shard was accessed.
			used := -1
			for i, v := range hit {
				if v {
					if used != -1 {
						t.Fatalf("multiple pools were accessed: %d and %d", used, i)
					}
					used = i
				}
			}
			is.NotEqual(-1, used, "no pool was used")
			t.Logf("Selected shard: %d (shardCount=%d)", used, tc.shardCount)
		})
	}
}

func Test_PRNG_Reader_Config(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Custom, non-default config to verify values round-trip
	want := Config{
		MaxBytesPerKey:    42,
		MaxInitRetries:    7,
		MaxRekeyAttempts:  8,
		MaxRekeyBackoff:   5 * time.Second,
		RekeyBackoff:      1 * time.Second,
		EnableKeyRotation: true,
		UseZeroBuffer:     true,
		DefaultBufferSize: 128,
		Shards:            4,
	}

	// Construct via functional options
	r, err := NewReader(
		WithMaxBytesPerKey(want.MaxBytesPerKey),
		WithMaxInitRetries(want.MaxInitRetries),
		WithMaxRekeyAttempts(want.MaxRekeyAttempts),
		WithMaxRekeyBackoff(want.MaxRekeyBackoff),
		WithRekeyBackoff(want.RekeyBackoff),
		WithEnableKeyRotation(want.EnableKeyRotation),
		WithZeroBuffer(want.UseZeroBuffer),
		WithDefaultBufferSize(want.DefaultBufferSize),
		WithShards(want.Shards),
	)
	is.NoError(err)

	got := r.Config()

	// All fields must match (deep comparison)
	is.Equal(want, got, "Config() should return the config passed to NewReader")
}
