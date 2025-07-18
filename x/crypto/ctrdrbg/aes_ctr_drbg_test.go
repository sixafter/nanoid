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
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test_CTRDRBG_Read verifies that a single Read operation from a new DRBG instance
// produces a buffer filled with nonzero, apparently random data. The test ensures
// the DRBG is correctly seeded and generating cryptographically strong output on first use.
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

// Test_CTRDRBG_ReadZeroBytes checks that reading into a zero-length buffer
// is a no-op and returns immediately, as required by the io.Reader contract.
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

// Test_CTRDRBG_ReadMultipleTimes validates that consecutive Read calls from a DRBG
// instance yield different outputs, ensuring the internal counter advances and no state is reused.
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

// Test_CTRDRBG_ReadWithDifferentBufferSizes runs Read on a variety of buffer sizes (1–2KiB).
// It ensures the returned buffer is always filled, and that the DRBG supports
// all size requests without error or truncation.
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

// Test_CTRDRBG_Concurrency verifies that the DRBG is safe under heavy concurrency
// by launching 100 goroutines, each reading a buffer in parallel. The test asserts
// all reads succeed and at least two buffers differ, confirming thread safety and uniqueness.
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

// Test_CTRDRBG_Stream validates that reading a large (1 MiB) buffer using io.ReadFull
// from the DRBG fills the entire buffer with nonzero, random data, ensuring correct
// handling of large sequential requests.
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

// Test_CTRDRBG_ReadAll checks that very large reads (10 KiB) succeed and the buffer
// is filled with unique, nonzero data. This protects against length or edge-case errors.
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

// Test_CTRDRBG_ReadConsistency performs 50 sequential reads from the same DRBG instance,
// storing the output from each. It verifies every buffer is nonzero and ensures that
// at least two reads differ, confirming uniqueness and liveness across multiple calls.
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

// Test_CTRDRBG_AsyncRekey validates the asynchronous key rotation mechanism of the DRBG.
// It configures a small MaxBytesPerKey to trigger a rekey, reads enough data to exceed the
// threshold, then waits and verifies that the internal AES block is replaced and usage is reset.
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

	// Get a snapshot of the initial state pointer (or block pointer)
	initialState := d.state.Load()
	initialBlock := initialState.block

	buf := make([]byte, 128) // Exceeds MaxBytesPerKey, triggers rekey
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
			currState := d.state.Load()
			// Compare block pointers (address equality) to detect a swap
			if currState.block != initialBlock && atomic.LoadUint64(&d.usage) == 0 {
				return // success
			}
		case <-wait.C:
			t.Fatal("Timed out waiting for asyncRekey to complete")
		}
	}
}

// Test_CTRDRBG_Personalization_Changes_Stream ensures that two DRBG instances constructed
// with different personalization parameters yield distinct output streams. The test asserts
// that the personalization string directly impacts the stream as required by NIST SP 800-90A.
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

// Test_CTRDRBG_Read_Shards verifies that a single call to Read only accesses
// one shard pool out of many, regardless of the pool count. It does not
// assert *which* shard is selected, as shardIndex is intentionally random.
//
// This test is table-driven: it runs the check with a variety of pool counts
// to ensure correct behavior at boundaries and typical values.
func Test_CTRDRBG_Read_Shards(t *testing.T) {
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
						d, _ := newDRBG(&cfg)
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

// TestDRBG_FillBlocks_ZeroAlloc_Functional verifies that drbg.fillBlocks produces
// unique, non-zero cryptographically random output without any heap allocations.
//
// The test ensures:
//   - The output buffer is filled with non-zero data (not all zeros).
//   - Output changes across invocations (counter increments, no repeats).
//   - No heap allocations occur per call, asserting high performance for this core routine.
//
// This is required for both performance-critical use and compliance with strict allocation budgets.
//
// It is a functional AND allocation test for fillBlocks, and should remain passing as the code evolves.
func Test_DRBG_FillBlocks_ZeroAlloc(t *testing.T) {
	cfg := DefaultConfig()
	d, _ := newDRBG(&cfg)
	var v [16]byte
	buf := make([]byte, KeySize256)

	st := d.state.Load()

	// Warmup and baseline output
	d.fillBlocks(buf, st, &v)
	baseline := make([]byte, len(buf))
	copy(baseline, buf)

	allocs := testing.AllocsPerRun(10000, func() {
		// Fill should mutate buffer and advance counter
		d.fillBlocks(buf, st, &v)
	})
	if allocs != 0 {
		t.Fatalf("unexpected allocations in fillBlocks: %v", allocs)
	}
	// Functional: Check that the buffer is not all zero
	allZero := true
	for _, b := range buf {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		t.Fatal("fillBlocks wrote only zeros")
	}
	// Functional: Output must differ from baseline (counter advanced)
	if string(baseline) == string(buf) {
		t.Fatal("fillBlocks output did not change across calls (counter not advancing?)")
	}
}

// TestDRBG_Read_Functional_Allow1Alloc verifies that drbg.Read produces non-zero,
// unique cryptographic output, and allocates at most once per call.
//
// The test ensures:
//   - The buffer is always filled with non-zero, apparently random data.
//   - Output changes across subsequent reads (counter is advancing).
//   - Heap allocations are ≤ 1 per call (ideally 0, but up to 1 is accepted to allow sync.Pool/runtime bookkeeping).
//
// This protects against accidental regression in allocation patterns or cryptographic soundness.
func Test_DRBG_Read_OneAlloc(t *testing.T) {
	cfg := DefaultConfig()
	d, _ := newDRBG(&cfg)
	buf := make([]byte, 32)

	// Warm up, baseline
	d.Read(buf)
	baseline := make([]byte, len(buf))
	copy(baseline, buf)

	allocs := testing.AllocsPerRun(10000, func() {
		d.Read(buf)
	})
	if allocs > 1 {
		t.Fatalf("unexpected allocations: %v (expected ≤ 1)", allocs)
	}
	// Buffer filled?
	allZero := true
	for _, b := range buf {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		t.Fatal("Read wrote only zeros")
	}
	// Output differs from baseline?
	if string(baseline) == string(buf) {
		t.Fatal("Read output did not change across calls (counter not advancing?)")
	}
}

func Test_DRBG_Reader_Config(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	want := Config{
		KeySize:           KeySize256,
		MaxBytesPerKey:    1024 * 1024,
		MaxInitRetries:    5,
		MaxRekeyAttempts:  7,
		MaxRekeyBackoff:   100 * time.Millisecond,
		RekeyBackoff:      10 * time.Millisecond,
		EnableKeyRotation: true,
		Personalization:   []byte("reader-domain"),
		UseZeroBuffer:     true,
		DefaultBufferSize: 128,
		Shards:            3,
	}

	rdr, err := NewReader(
		WithKeySize(want.KeySize),
		WithMaxBytesPerKey(want.MaxBytesPerKey),
		WithMaxInitRetries(want.MaxInitRetries),
		WithMaxRekeyAttempts(want.MaxRekeyAttempts),
		WithMaxRekeyBackoff(want.MaxRekeyBackoff),
		WithRekeyBackoff(want.RekeyBackoff),
		WithEnableKeyRotation(want.EnableKeyRotation),
		WithPersonalization(want.Personalization),
		WithUseZeroBuffer(want.UseZeroBuffer),
		WithDefaultBufferSize(want.DefaultBufferSize),
		WithShards(want.Shards),
	)
	is.NoError(err)

	got := rdr.Config()
	is.Equal(want.KeySize, got.KeySize)
	is.Equal(want.MaxBytesPerKey, got.MaxBytesPerKey)
	is.Equal(want.MaxInitRetries, got.MaxInitRetries)
	is.Equal(want.MaxRekeyAttempts, got.MaxRekeyAttempts)
	is.Equal(want.MaxRekeyBackoff, got.MaxRekeyBackoff)
	is.Equal(want.RekeyBackoff, got.RekeyBackoff)
	is.Equal(want.EnableKeyRotation, got.EnableKeyRotation)
	is.True(bytes.Equal(got.Personalization, want.Personalization), "Personalization does not match")
	is.Equal(want.UseZeroBuffer, got.UseZeroBuffer)
	is.Equal(want.DefaultBufferSize, got.DefaultBufferSize)
	is.Equal(want.Shards, got.Shards)
}

// TestDRBG_CounterOverflow Simulate the 128-bit counter rolling over (set `d.v` to `[0xff ... 0xff]`, read one block)
// and ensure it wraps correctly per spec. While extremely unlikely in practice, it’s a security-critical edge case.
func Test_DRBG_CounterOverflow(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	cfg := DefaultConfig()
	d, err := newDRBG(&cfg)
	is.NoError(err)

	// Set the DRBG's working counter to all 0xff (max 128-bit value).
	for i := 0; i < len(d.v); i++ {
		d.v[i] = 0xff
	}

	// Prepare output buffer (block size).
	blockSize := 16 // AES block size
	buf := make([]byte, blockSize)

	// Read a block -- should increment counter and wrap it to zero.
	_, err = d.Read(buf)
	is.NoError(err)

	// After increment, counter should be zero
	expected := make([]byte, 16)
	is.Equal(expected, d.v[:], "Counter should wrap to zero after overflow")

	// Optionally, check that output is nonzero
	allZeros := true
	for _, b := range buf {
		if b != 0 {
			allZeros = false
			break
		}
	}
	is.False(allZeros, "Output block should not be all zeros")
}
