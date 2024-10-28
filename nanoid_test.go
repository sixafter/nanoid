// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"bytes"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewGenerator(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	alphabet := "abc123"
	gen, err := New(alphabet, nil)
	is.NoError(err, "Expected no error when creating a new generator")
	is.NotNil(gen, "Generator should not be nil")

	conf, ok := gen.(Configuration)
	is.True(ok, "Expected a Configuration")

	is.Equal(6, conf.GetConfig().AlphabetLen, "AlphabetLen should be 6")
	is.Equal([]byte(alphabet), conf.GetConfig().Alphabet, "Alphabet should match the input")
	is.Equal(byte(7), conf.GetConfig().Mask, "Mask should be 7 for alphabetLen=6")
	is.Equal(341, conf.GetConfig().Step, "Step should be 341 for mask=7 and alphabetLen=6")
}

func TestGenerateDefault(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	id, err := GenerateSize(10)
	is.NoError(err, "GenerateSize should not return an error")
	is.Len(id, 10, "Generated ID should have length 10")
}

func TestGenerateCustomAlphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	alphabet := "ABCDEF"
	gen, err := New(alphabet, nil)
	is.NoError(err, "Creating generator with custom alphabet should not error")
	is.NotNil(gen, "Generator should not be nil")

	id, err := gen.Generate(5)
	is.NoError(err, "GenerateSize should not return an error")
	is.Len(id, 5, "Generated ID should have length 5")

	for _, c := range id {
		is.True(bytes.Contains([]byte(alphabet), []byte{byte(c)}), "Character %c should be in the custom alphabet", c)
	}
}

func TestGenerateZeroLength(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	id, err := GenerateSize(0)
	is.ErrorIs(err, ErrInvalidLength, "Generating ID with zero length should return ErrInvalidLength")
	is.Empty(id, "Generated ID should be empty when length is zero")
}

func TestGenerateWithCustomReader(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Initialize mockReader with sufficient data
	mockReader := &mockReader{
		data:  []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
		reads: 0,
	}
	alphabet := "ABCDEF"
	gen, err := New(alphabet, mockReader)
	is.NoError(err, "Creating generator with custom reader should not error")
	is.NotNil(gen, "Generator should not be nil")

	id, err := gen.Generate(6)
	is.NoError(err, "GenerateSize should not return an error")
	is.Equal("ABCDEF", id, "Generated ID should match expected value 'ABCDEF'")
}

func TestGenerateInsufficientRandomBytes(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Reader that provides only a few bytes
	mockReader := &mockReader{
		data:  []byte{0x00, 0x01}, // Only 2 bytes
		reads: 0,
	}
	alphabet := "ABCDEF"
	gen, err := New(alphabet, mockReader)
	is.NoError(err, "Creating generator with custom reader should not error")
	is.NotNil(gen, "Generator should not be nil")

	id, err := gen.Generate(3) // Requires 3 valid bytes
	is.ErrorIs(err, ErrNoMoreData, "Generating ID with insufficient random bytes should return ErrNoMoreData")
	is.Empty(id, "Generated ID should be empty on error")
}

func TestGenerateConcurrency(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	gen, err := New(DefaultAlphabet, nil)
	is.NoError(err, "Creating generator should not error")
	is.NotNil(gen, "Generator should not be nil")

	if gen == nil {
		return
	}

	const goroutines = 100
	const idsPerGoroutine = 1000
	ids := make(chan string, goroutines*idsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id, err := gen.Generate(10)
				if err != nil {
					t.Errorf("GenerateSize failed: %v", err)
					return
				}
				ids <- id
			}
		}()
	}
	wg.Wait()
	close(ids)

	// Verify uniqueness
	idMap := make(map[string]struct{})
	for id := range ids {
		idMap[id] = struct{}{}
	}
	is.Equal(goroutines*idsPerGoroutine, len(idMap), "All generated IDs should be unique")
}

var ErrNoMoreData = errors.New("no more data")

// mockReader is a mock implementation of io.Reader for testing purposes.
type mockReader struct {
	data  []byte
	reads int
}

func (m *mockReader) Read(p []byte) (int, error) {
	if m.reads >= len(m.data) {
		return 0, ErrNoMoreData
	}
	p[0] = m.data[m.reads]
	m.reads++
	return 1, nil
}

// TestGenerateDoesNotHang tests that GenerateSize does not hang indefinitely.
func TestGenerateDoesNotHang(t *testing.T) {
	gen, err := New("abcdef", nil)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	done := make(chan struct{})
	go func() {
		_, err := gen.Generate(50)
		if err != nil {
			t.Errorf("GenerateSize failed: %v", err)
		}
		close(done)
	}()

	select {
	case <-done:
		// Test passed
	case <-time.After(5 * time.Second):
		t.Error("GenerateSize method is hanging")
	}
}
