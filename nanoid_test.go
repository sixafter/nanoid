// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the MIT License found in the
// LICENSE file in the root directory of this source tree.

// nanoid_test.go
package nanoid_test

import (
	"fmt"
	"github.com/sixafter/nanoid"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

// TestNew verifies that the New function generates an ID of the default size and alphabet.
func TestNew(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	id, err := nanoid.New()
	is.NoError(err, "New() should not return an error")
	is.Equal(nanoid.DefaultSize, len(id), "ID length should match the default size")

	// Check that all characters are within the default alphabet
	for _, char := range id {
		is.Contains(nanoid.DefaultAlphabet, string(char), "Character '%c' should be in the default alphabet", char)
	}
}

// TestNewSize verifies that the NewSize function generates IDs of specified sizes.
func TestNewSize(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	testCases := []struct {
		name string
		size int
	}{
		{"Size1", 1},
		{"Size10", 10},
		{"Size21", 21},
		{"Size50", 50},
		{"Size100", 100},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			id, err := nanoid.NewSize(tc.size)
			is.NoError(err, "NewSize(%d) should not return an error", tc.size)
			is.Equal(tc.size, len(id), "ID length should match the specified size")

			// Check that all characters are within the default alphabet
			for _, char := range id {
				is.Contains(nanoid.DefaultAlphabet, string(char), "Character '%c' should be in the default alphabet", char)
			}
		})
	}
}

// TestNewCustom verifies that the NewCustom function generates IDs using a custom ASCII alphabet.
func TestNewCustom(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	customASCIIAlphabet := "abcdef123456"
	testCases := []struct {
		name     string
		size     int
		alphabet string
	}{
		{"Size10_CustomASCIIAlphabet", 10, customASCIIAlphabet},
		{"Size21_CustomASCIIAlphabet", 21, customASCIIAlphabet},
		{"Size50_CustomASCIIAlphabet", 50, customASCIIAlphabet},
		{"Size100_CustomASCIIAlphabet", 100, customASCIIAlphabet},
		{"Size10_SingleCharacter", 10, "x"}, // Single-character alphabet
		{"Size5_SingleCharacter", 5, "A"},   // Single-character alphabet
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			id, err := nanoid.NewCustom(tc.size, tc.alphabet)
			is.NoError(err, "NewCustom(%d, %s) should not return an error", tc.size, tc.alphabet)
			is.Equal(tc.size, len(id), "ID length should match the specified size")

			// Check that all characters are within the custom alphabet
			for _, char := range id {
				is.Contains(tc.alphabet, string(char), "Character '%c' should be in the custom alphabet", char)
			}
		})
	}
}

// TestErrorHandling verifies that functions return errors for invalid inputs.
func TestErrorHandling(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	testCases := []struct {
		name     string
		function func() (string, error)
	}{
		{
			name: "NewSize with zero size",
			function: func() (string, error) {
				return nanoid.NewSize(0)
			},
		},
		{
			name: "NewSize with negative size",
			function: func() (string, error) {
				return nanoid.NewSize(-10)
			},
		},
		{
			name: "NewCustom with empty alphabet",
			function: func() (string, error) {
				return nanoid.NewCustom(10, "")
			},
		},
		{
			name: "NewCustom with size exceeding MaxUintSize",
			function: func() (string, error) {
				return nanoid.NewCustom(nanoid.MaxUintSize+1, "abcdef")
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			id, err := tc.function()
			is.Error(err, "Expected an error for test case '%s'", tc.name)
			is.Empty(id, "Expected empty ID for test case '%s'", tc.name)
		})
	}
}

// TestUniqueness verifies that multiple generated IDs are unique.
func TestUniqueness(t *testing.T) {
	// Note: Due to the high memory consumption and execution time,
	// this test can be marked as skipped unless specifically needed.
	t.Skip("Skipping TestUniqueness to save resources during regular test runs")

	t.Parallel()

	is := assert.New(t)

	const sampleSize = 100000
	ids := make(map[string]struct{}, sampleSize)

	for i := 0; i < sampleSize; i++ {
		id, err := nanoid.New()
		is.NoError(err, "New() should not return an error")

		if _, exists := ids[id]; exists {
			is.FailNow(fmt.Sprintf("Duplicate ID found: %s", id))
		}
		ids[id] = struct{}{}
	}
}

// TestConcurrencySafety verifies that concurrent ID generation does not produce errors or duplicates.
func TestConcurrencySafety(t *testing.T) {
	t.Parallel()

	is := assert.New(t)

	const (
		concurrency  = 100
		perGoroutine = 1000
		totalSample  = concurrency * perGoroutine
	)
	ids := make(chan string, totalSample)
	errs := make(chan error, totalSample)

	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				id, err := nanoid.New()
				if err != nil {
					errs <- err
					continue
				}
				ids <- id
			}
		}()
	}

	wg.Wait()
	close(ids)
	close(errs)

	// Check for errors
	for err := range errs {
		is.NoError(err, "New() should not return an error in concurrent execution")
	}

	// Check for duplicates
	uniqueIDs := make(map[string]struct{}, totalSample)
	for id := range ids {
		if _, exists := uniqueIDs[id]; exists {
			is.FailNow(fmt.Sprintf("Duplicate ID found: %s", id))
		}
		uniqueIDs[id] = struct{}{}
	}
}
