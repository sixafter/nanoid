// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"math/bits"
	"testing"

	"github.com/sixafter/nanoid/x/crypto/prng"
	"github.com/stretchr/testify/assert"
)

// TestGetConfig tests the Config() method of the generator.
func TestGetConfig(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	gen, err := NewGenerator()
	is.NoError(err, "NewGenerator() should not return an error with the default alphabet")

	// Assert that generator implements Configuration interface
	config, ok := gen.(Configuration)
	is.True(ok, "Generator should implement Configuration interface")

	runtimeConfig := config.Config()

	is.Equal(DefaultAlphabet, string(runtimeConfig.RuneAlphabet()), "Config.RuneAlphabet should match the default alphabet")
	is.Equal(uint16(len([]rune(DefaultAlphabet))), runtimeConfig.AlphabetLen(), "Config.AlphabetLen should match the default alphabet length")

	// Update expectedMask calculation based on RuntimeConfig
	expectedMask := uint((1 << bits.Len(uint(runtimeConfig.AlphabetLen()-1))) - 1)
	is.Equal(expectedMask, runtimeConfig.Mask(), "Config.Mask should be correctly calculated")

	is.Equal((runtimeConfig.AlphabetLen()&(runtimeConfig.AlphabetLen()-1)) == 0, runtimeConfig.IsPowerOfTwo(), "Config.IsPowerOfTwo should be correct")
	is.Positive(runtimeConfig.BaseMultiplier(), "Config.BaseMultiplier should be a positive integer")
	is.Positive(runtimeConfig.BitsNeeded(), "Config.BitsNeeded should be a positive integer")
	is.Positive(runtimeConfig.BufferMultiplier(), "Config.BufferMultiplier should be a positive integer")
	is.Positive(runtimeConfig.BufferSize(), "Config.BufferSize should be a positive integer")
	is.NotNil(runtimeConfig.ByteAlphabet(), "Config.ByteAlphabet should not be nil")
	is.Positive(runtimeConfig.BytesNeeded(), "Config.BytesNeeded should be a positive integer")
	is.Equal(true, runtimeConfig.IsASCII(), "Config.IsASCII should be true by default")
	is.Equal(true, runtimeConfig.IsPowerOfTwo(), "Config.IsPowerOfTwo should be true by default")
	is.Positive(runtimeConfig.LengthHint(), "Config.LengthHint should be a positive integer")
	is.Equal(1, runtimeConfig.MaxBytesPerRune(), "Config.MaxBytesPerRune should be 1 by default")
	is.Equal(prng.Reader, runtimeConfig.RandReader(), "Config.RandReader should be rand.Reader by default")
	is.NotNil(runtimeConfig.RuneAlphabet(), "Config.RuneAlphabet should not be nil")
	is.Positive(runtimeConfig.ScalingFactor(), "Config.ScalingFactor should be a positive integer")
}
