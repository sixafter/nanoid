package nanoid

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	id, err := New()
	is.NoError(err)
	is.Equal(DefaultSize, utf8.RuneCountInString(id))
}

func TestNewSize(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	size := 30
	id, err := NewSize(size)
	is.NoError(err)
	is.Equal(size, utf8.RuneCountInString(id))
}

func TestNewCustom(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	size := 15
	customAlphabet := "abcdef123456"
	id, err := NewCustom(size, customAlphabet)
	is.NoError(err)
	is.Equal(size, utf8.RuneCountInString(id))

	// Ensure that the ID contains only characters from the custom alphabet
	for _, r := range id {
		is.Contains(customAlphabet, string(r))
	}
}

func TestNewCustomWithUnicodeAlphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	size := 10
	unicodeAlphabet := "あいうえお漢字🙂🚀"
	id, err := NewCustom(size, unicodeAlphabet)
	is.NoError(err)
	is.Equal(size, utf8.RuneCountInString(id))

	// Ensure that the ID contains only characters from the Unicode alphabet
	for _, r := range id {
		is.Contains(unicodeAlphabet, string(r))
	}
}

func TestNewCustomReader(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	size := 10
	customAlphabet := "abcd"
	// Use a deterministic random source for testing
	randomData := bytes.Repeat([]byte{0xFF}, 10) // Simulate random bytes
	rnd := bytes.NewReader(randomData)

	id, err := NewCustomReader(size, customAlphabet, rnd)
	is.NoError(err)
	is.Equal(size, utf8.RuneCountInString(id))

	// Since we used 0xFF, the index will be masked accordingly
	for _, r := range id {
		is.Contains(customAlphabet, string(r))
	}
}

func TestNewCustomReaderNil(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	size := 10
	customAlphabet := DefaultAlphabet

	_, err := NewCustomReader(size, customAlphabet, nil)
	is.Error(err)
	is.EqualError(err, "random source cannot be nil")
}

func TestNewCustomEmptyAlphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	size := 10
	customAlphabet := ""

	_, err := NewCustom(size, customAlphabet)
	is.Error(err)
	is.EqualError(err, "alphabet must not be empty")
}

func TestNewCustomNegativeSize(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	size := -5
	customAlphabet := DefaultAlphabet

	_, err := NewCustom(size, customAlphabet)
	is.Error(err)
	is.EqualError(err, "size must be greater than zero")
}

func TestNewCustomSingleCharacterAlphabet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	size := 10
	customAlphabet := "X"

	id, err := NewCustom(size, customAlphabet)
	is.NoError(err)
	is.Equal(strings.Repeat(customAlphabet, size), id)
}

func TestThreadSafety(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	const numGoroutines = 100
	const idSize = 21

	ids := make(chan string, numGoroutines)
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := New()
			is.NoError(err)
			is.Equal(idSize, utf8.RuneCountInString(id))
			ids <- id
		}()
	}

	wg.Wait()
	close(ids)

	idSet := make(map[string]struct{})
	for id := range ids {
		// Ensure uniqueness
		is.NotContains(idSet, id)
		idSet[id] = struct{}{}
	}
}
