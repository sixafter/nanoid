// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestID_String tests the String() method of the ID type.
// It verifies that the String() method returns the underlying string value.
func TestID_String(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Initialize expected using Must()
	expectedID := Must()
	expected := expectedID.String()

	// Actual is obtained by calling String() on the ID
	actual := expectedID.String()

	// Assert that actual equals expected
	is.Equal(expected, actual, "ID.String() should return the underlying string")
}

// TestID_MarshalText tests the MarshalText() method of the ID type.
// It verifies that MarshalText() returns the correct byte slice representation of the ID.
func TestID_MarshalText(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Initialize expected using Must()
	expectedID := Must()
	expectedBytes := []byte(expectedID.String())

	// Actual is obtained by calling MarshalText()
	actualBytes, err := expectedID.MarshalText()

	// Assert no error occurred
	is.NoError(err, "MarshalText() should not return an error")

	// Assert that actual bytes match expected bytes
	is.Equal(expectedBytes, actualBytes, "MarshalText() should return the correct byte slice")
}

// TestID_UnmarshalText tests the UnmarshalText() method of the ID type.
// It verifies that UnmarshalText() correctly parses the byte slice and assigns the value to the ID.
func TestID_UnmarshalText(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Initialize expected using Must()
	expectedID := Must()
	inputBytes := []byte(expectedID.String())

	// Initialize a zero-valued ID
	var actualID ID

	// Call UnmarshalText with inputBytes
	err := actualID.UnmarshalText(inputBytes)

	// Assert no error occurred
	is.NoError(err, "UnmarshalText() should not return an error")

	// Assert that actualID matches expectedID
	is.Equal(expectedID, actualID, "UnmarshalText() should correctly assign the input value to ID")
}

// TestID_MarshalBinary tests the MarshalBinary() method of the ID type.
// It verifies that MarshalBinary() returns the correct byte slice representation of the ID.
func TestID_MarshalBinary(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Initialize expected using Must()
	expectedID := Must()
	expectedBytes := []byte(expectedID.String())

	// Actual is obtained by calling MarshalBinary()
	actualBytes, err := expectedID.MarshalBinary()

	// Assert no error occurred
	is.NoError(err, "MarshalBinary() should not return an error")

	// Assert that actual bytes match expected bytes
	is.Equal(expectedBytes, actualBytes, "MarshalBinary() should return the correct byte slice")
}

// TestID_UnmarshalBinary tests the UnmarshalBinary() method of the ID type.
// It verifies that UnmarshalBinary() correctly parses the byte slice and assigns the value to the ID.
func TestID_UnmarshalBinary(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Initialize expected using Must()
	expectedID := Must()
	inputBytes := []byte(expectedID.String())

	// Initialize a zero-valued ID
	var actualID ID

	// Call UnmarshalBinary with inputBytes
	err := actualID.UnmarshalBinary(inputBytes)

	// Assert no error occurred
	is.NoError(err, "UnmarshalBinary() should not return an error")

	// Assert that actualID matches expectedID
	is.Equal(expectedID, actualID, "UnmarshalBinary() should correctly assign the input value to ID")
}

// TestID_Compare tests the Compare() method of the ID type.
// It verifies that Compare() correctly compares two IDs and returns the expected result.
func TestID_Compare(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	id1 := ID("FgEVN8QMTrnKGvBxFjtjw")
	id2 := ID("zTxG5Nl21ZAoM8Fabqk3H")

	// Case 1: id1 < id2
	is.Equal(-1, id1.Compare(id2), "id1 should be less than id2")

	// Case 2: id1 = id2
	is.Equal(0, id1.Compare(id1), "id1 should be equal to id1")

	// Case 3: id1 > id2
	is.Equal(1, id2.Compare(id1), "id2 should be greater than id1")
}

// TestID_IsEmpty tests the IsEmpty() method of the ID type.
// It verifies that IsEmpty() correctly returns true for an empty ID and false for a non-empty ID.
func TestID_IsEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Initialize two IDs using Must()
	id1 := Must()
	id2 := EmptyID

	// Case 1: id1 is not empty
	is.False(id1.IsEmpty(), "id1 should not be empty")

	// Case 2: id2 is empty
	is.True(id2.IsEmpty(), "id2 should be empty")
}

func TestID_IsEmpty_NilReceiver(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	var id *ID // nil pointer
	is.True(id.IsEmpty(), "expected IsEmpty to return true for nil receiver")
}
