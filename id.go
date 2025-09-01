// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"strings"
)

// ID represents a Nano ID as a string.
type ID string

// EmptyID represents an empty Nano ID.
var EmptyID = ID("")

// IsEmpty returns true if the ID is an empty ID (EmptyID)
func (id *ID) IsEmpty() bool {
	if id == nil {
		return true
	}

	return id.Compare(EmptyID) == 0
}

// Compare compares two IDs lexicographically and returns an integer.
// The result will be 0 if id==other, -1 if id < other, and +1 if id > other.
//
// Parameters:
//   - other ID: The ID to compare against.
//
// Returns:
//   - int: An integer indicating the comparison result.
//
// Usage:
//
//	id1 := ID("V1StGXR8_Z5jdHi6B-myT")
//	id2 := ID("V1StGXR8_Z5jdHi6B-myT")
//	result := id1.Compare(id2)
//	fmt.Println(result) // Output: 0
func (id ID) Compare(other ID) int {
	return strings.Compare(string(id), string(other))
}

// String returns the string representation of the ID.
// It implements the fmt.Stringer interface, allowing the ID to be
// used seamlessly with fmt package functions like fmt.Println and fmt.Printf.
//
// Example:
//
//	id := Must()
//	fmt.Println(id) // Output: V1StGXR8_Z5jdHi6B-myT
func (id ID) String() string {
	return string(id)
}

// MarshalText converts the ID to a byte slice.
// It implements the encoding.TextMarshaler interface, enabling the ID
// to be marshaled into text-based formats such as XML and YAML.
//
// Returns:
//   - A byte slice containing the ID.
//   - An error if the marshaling fails.
//
// Example:
//
//	id := Must()
//	text, err := id.MarshalText()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(text)) // Output: V1StGXR8_Z5jdHi6B-myT
func (id *ID) MarshalText() ([]byte, error) {
	if id == nil {
		return nil, ErrNilPointer
	}

	return []byte(*id), nil
}

// UnmarshalText parses a byte slice and assigns the result to the ID.
// It implements the encoding.TextUnmarshaler interface, allowing the ID
// to be unmarshaled from text-based formats.
//
// Parameters:
//   - text: A byte slice containing the ID data.
//
// Returns:
//   - An error if the unmarshaling fails.
//
// Example:
//
//	var id ID
//	err := id.UnmarshalText([]byte("new-id"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(id) // Output: new-id
func (id *ID) UnmarshalText(text []byte) error {
	if id == nil {
		return ErrNilPointer
	}

	*id = ID(text)
	return nil
}

// MarshalBinary converts the ID to a byte slice.
// It implements the encoding.BinaryMarshaler interface, enabling the ID
// to be marshaled into binary formats for efficient storage or transmission.
//
// Returns:
//   - A byte slice containing the ID.
//   - An error if the marshaling fails.
//
// Example:
//
//	id := Must()
//	binaryData, err := id.MarshalBinary()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(binaryData) // Output: [86 49 83 116 71 88 82 56 95 90 ...]
func (id *ID) MarshalBinary() ([]byte, error) {
	if id == nil {
		return nil, ErrNilPointer
	}

	return []byte(*id), nil
}

// UnmarshalBinary parses a byte slice and assigns the result to the ID.
// It implements the encoding.BinaryUnmarshaler interface, allowing the ID
// to be unmarshaled from binary formats.
//
// Parameters:
//   - data: A byte slice containing the binary ID data.
//
// Returns:
//   - An error if the unmarshaling fails.
//
// Example:
//
//	var id ID
//	err := id.UnmarshalBinary([]byte{86, 49, 83, 116, 71, 88, 82, 56, 95, 90}) // "V1StGXR8_Z5jdHi6B-myT"
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(id) // Output: V1StGXR8_Z5jdHi6B-myT
func (id *ID) UnmarshalBinary(data []byte) error {
	if id == nil {
		return ErrNilPointer
	}

	*id = ID(data)
	return nil
}
