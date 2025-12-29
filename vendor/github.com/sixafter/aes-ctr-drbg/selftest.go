// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package ctrdrbg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"sync"
)

var (
	selfTestOnce sync.Once
	selfTestErr  error
)

// ErrSelfTestFailed indicates the FIPS 140-2 power-on self-test did not pass.
var ErrSelfTestFailed = errors.New("ctrdrbg: FIPS 140-2 self-test failed")

// RunSelfTests executes FIPS 140-2 Known Answer Tests (KAT) for AES-CTR.
//
// This function verifies that the AES-CTR cryptographic primitive is functioning
// correctly by comparing output against known NIST CAVP test vectors.
//
// RunSelfTests is safe for concurrent use and executes only once per process
// via sync.Once. Subsequent calls return the cached result.
//
// Returns nil on success, or ErrSelfTestFailed if the test fails.
func RunSelfTests() error {
	selfTestOnce.Do(func() {
		selfTestErr = runKAT()
	})
	return selfTestErr
}

// runKAT performs the Known Answer Test using NIST CAVP test vectors.
//
// Test vector source: NIST SP 800-38A, §F.5.5
// AES-256-CTR test vector
func runKAT() error {
	// NIST SP 800-38A F.5.5 AES-256-CTR test vector
	key := [32]byte{
		0x60, 0x3d, 0xeb, 0x10, 0x15, 0xca, 0x71, 0xbe,
		0x2b, 0x73, 0xae, 0xf0, 0x85, 0x7d, 0x77, 0x81,
		0x1f, 0x35, 0x2c, 0x07, 0x3b, 0x61, 0x08, 0xd7,
		0x2d, 0x98, 0x10, 0xa3, 0x09, 0x14, 0xdf, 0xf4,
	}
	iv := [16]byte{
		0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7,
		0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff,
	}
	plaintext := [16]byte{
		0x6b, 0xc1, 0xbe, 0xe2, 0x2e, 0x40, 0x9f, 0x96,
		0xe9, 0x3d, 0x7e, 0x11, 0x73, 0x93, 0x17, 0x2a,
	}
	expected := [16]byte{
		0x60, 0x1e, 0xc3, 0x13, 0x77, 0x57, 0x89, 0xa5,
		0xb7, 0xa7, 0xf5, 0x04, 0xbb, 0xf3, 0xd2, 0x28,
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return ErrSelfTestFailed
	}

	stream := cipher.NewCTR(block, iv[:])
	var out [16]byte
	stream.XORKeyStream(out[:], plaintext[:])

	if !bytes.Equal(out[:], expected[:]) {
		return ErrSelfTestFailed
	}
	return nil
}
