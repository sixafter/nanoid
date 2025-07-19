// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"math/bits"
	"testing"

	"github.com/sixafter/prng-chacha"
	"github.com/stretchr/testify/assert"
)

// Test_Config tests the Config() method of the generator.
func Test_Config(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	gen, err := NewGenerator()
	is.NoError(err, "NewGenerator() should not return an error with the default alphabet")

	config := gen.Config()

	is.Equal(DefaultAlphabet, string(config.RuneAlphabet()), "Config.RuneAlphabet should match the default alphabet")
	is.Equal(uint16(len([]rune(DefaultAlphabet))), config.AlphabetLen(), "Config.AlphabetLen should match the default alphabet length")

	// Update expectedMask calculation based on RuntimeConfig
	expectedMask := uint((1 << bits.Len(uint(config.AlphabetLen()-1))) - 1)
	is.Equal(expectedMask, config.Mask(), "Config.Mask should be correctly calculated")

	is.Equal((config.AlphabetLen()&(config.AlphabetLen()-1)) == 0, config.IsPowerOfTwo(), "Config.IsPowerOfTwo should be correct")
	is.Positive(config.BaseMultiplier(), "Config.BaseMultiplier should be a positive integer")
	is.Positive(config.BitsNeeded(), "Config.BitsNeeded should be a positive integer")
	is.Positive(config.BufferMultiplier(), "Config.BufferMultiplier should be a positive integer")
	is.Positive(config.BufferSize(), "Config.BufferSize should be a positive integer")
	is.NotNil(config.ByteAlphabet(), "Config.ByteAlphabet should not be nil")
	is.Positive(config.BytesNeeded(), "Config.BytesNeeded should be a positive integer")
	is.Equal(true, config.IsASCII(), "Config.IsASCII should be true by default")
	is.Equal(true, config.IsPowerOfTwo(), "Config.IsPowerOfTwo should be true by default")
	is.Positive(config.LengthHint(), "Config.LengthHint should be a positive integer")
	is.Equal(1, config.MaxBytesPerRune(), "Config.MaxBytesPerRune should be 1 by default")
	is.Equal(prng.Reader, config.RandReader(), "Config.RandReader should be rand.Reader by default")
	is.NotNil(config.RuneAlphabet(), "Config.RuneAlphabet should not be nil")
	is.Positive(config.ScalingFactor(), "Config.ScalingFactor should be a positive integer")
}
