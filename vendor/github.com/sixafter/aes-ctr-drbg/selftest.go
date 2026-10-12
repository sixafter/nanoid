// Copyright (c) 2024-2026 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package ctrdrbg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"os"
	"sync"
	"time"
)

var (
	selfTestOnce sync.Once
	selfTestErr  error

	drbgSelfTestOnce sync.Once
	drbgSelfTestErr  error
)

// ErrSelfTestFailed indicates the FIPS 140-2 power-on self-test did not pass.
var ErrSelfTestFailed = errors.New("ctrdrbg: FIPS 140-2 self-test failed")

// ErrDRBGSelfTestFailed indicates that the NIST SP 800-90A §11.3 known-answer self-test of the instantiate, reseed,
// and generate functions did not pass. NewReader returns it when BacktrackingResistance is enabled. The result is
// kept for the life of the process, so every later NewReader call with BacktrackingResistance also returns it.
var ErrDRBGSelfTestFailed = errors.New("ctrdrbg: NIST SP 800-90A known-answer self-test failed")

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

// runDRBGSelfTests runs the NIST SP 800-90A §11.3 known-answer tests once per process and returns the result.
func runDRBGSelfTests() error {
	drbgSelfTestOnce.Do(func() {
		drbgSelfTestErr = runDRBGKAT()
	})
	return drbgSelfTestErr
}

// drbgKATVector is a NIST CAVP CTR_DRBG known-answer vector ("no df", with a personalization string and additional
// input of key size + 16 bytes). Without prediction resistance, the vector instantiates, reseeds with
// entropyInputReseed and additionalInputReseed, then generates twice. With prediction resistance, each of the two
// generate requests first reseeds with the next entropyInputPR value.
type drbgKATVector struct {
	entropyInput          string
	personalization       string
	entropyInputReseed    string
	additionalInputReseed string
	entropyInputPR        [2]string
	additionalInput       [2]string
	returnedBits          string
	keySize               KeySize
	predictionResistance  bool
}

// drbgKATVectors are COUNT = 0 of the "no df" sections with personalization and additional input in the NIST CAVP
// files drbgvectors_pr_false/CTR_DRBG.rsp and drbgvectors_pr_true/CTR_DRBG.rsp, for each key size.
var drbgKATVectors = []drbgKATVector{
	{
		keySize:               KeySize128,
		entropyInput:          "289e5c8283cbd7dbe707255cb3cf2907d8a5ce5b347314966f9b2bebb1a1e200",
		personalization:       "7f7b59f23510b976fe155d047525c94e2dacb30d77ac8b09281544dd815d5293",
		entropyInputReseed:    "98c522028f36fc6b85a8f3c003efd4b130dd90180ec81cf7c67d4c53d10f0022",
		additionalInputReseed: "f7a0378328d939f0f8521e39409d7175d87319c7597a9050414f7adc392a328d",
		additionalInput: [2]string{
			"19c286f5b36194d1cc62c0188140bc9d61d2a9c5d88bb5aebc224bfb04dfca83",
			"820650c3201d347f5b20d3d25d1c8c7bef4d9f66a5a04c7dd9d669e95182a0c4",
		},
		returnedBits: "79a79d44edada58e3fc12a4e36ae900eeace290265f01262f40f2958a70dcbd4" +
			"d4185f708c088ede7ff8c8375f44f4012f2512d38328a5df171a17029d90f185",
	},
	{
		keySize:               KeySize192,
		entropyInput:          "4b58271b116237eedd4e9ff9360382a59f3e2a173d860f2bbd8b2bace142b2395c67cf5a513f06f3",
		personalization:       "cf76c16cd5d270707ea9acc39744db69bfac63e566256fd6917bf9819679840f3fea2aa535d8df01",
		entropyInputReseed:    "1867f371a345eef98b2d70fc1960397892645b7b29a4ead252e8835e0b600618a9bd6ff99785d890",
		additionalInputReseed: "6d44839aff8b7165deebd489ad088ecb7dcec11c32b1e747dba8f0e8a0b89f74a84ea8a05586fe9e",
		additionalInput: [2]string{
			"42248fce0994e0e63504209d629a6943eb3e2ad512f03f79cbd5102928392bce1cacbba056ac6ca9",
			"bd529b600273329423a58d6f8a12be0f17989a02e73e347bc7d49d9169337a6cff7c07e8a807a80a",
		},
		returnedBits: "02486d32cd55954f406ba55705f1460d384439592dede81a84fda221fd45c0d6" +
			"51d67ec4a81a8b404151a643f331ad051cb004352289de37bca71e8cc0a6aeab",
	},
	{
		keySize:               KeySize256,
		entropyInput:          "ae7ebe062971f5eb32e5b21444750785de816595ad2cbe80a209c8f8ab04b5468166de8c6ae522d8f10b56386a3b424f",
		personalization:       "55860dae57fcac297087c137efb796878a75868f6e7681114e9b73ed0c67e3c62bfc9f5d77e8caa59bcdb223f4ffd247",
		entropyInputReseed:    "a42407931bfeca70e6ee5dd197021a129525051c07468e8b25587c5ad50abe9204e882fe847b8fd47cf7b4360e5aa034",
		additionalInputReseed: "ee4c88d1eb05f4853663eada501d2fc4b4984b283a88db579af2113031e03d9bc570de943dd168918f3ba8065581fea7",
		additionalInput: [2]string{
			"4b4b03ef19b0f259dca2b3ee3ae4cd86c3895a784b3d8eee043a2003c08289f8fffdad141e6b1ab2174d8d5d79c1e581",
			"3062b33f116b46e20fe3c354726ae9b2a3a4c51922c8107863cb86f1f0bdad7554075659d91c371e2b11b1e8106a1ed5",
		},
		returnedBits: "0d270518baeafac160ff1cb28c11ef68712c764c0c01674e6c9ca2cc9c7e0e8a" +
			"ccfd3c753635ee070081eee7628af6187fbc2854b3c204461a796cf3f3fcb092",
	},
	{
		keySize:              KeySize128,
		predictionResistance: true,
		entropyInput:         "f623911bb7efe5f6c4e969799724b317d6aa2107a356ab6a05143eee99e675b3",
		personalization:      "19733edbc20bedc0bc269c4837a1a471377cfeb6cd60039e4944b31af73bc4c0",
		entropyInputPR: [2]string{
			"60a809d63f6d5efeb4335b819e1d0a06b7be668765ca72f5404992a0e90b7a70",
			"530e7a6dfa6339aaa6661bc4fb13ff1f351c79e9d3415e2d0d570b8f9907c8b4",
		},
		additionalInput: [2]string{
			"0591feffa77de223d6fb1a91579f766c4c09f3e9540d14fd2ee98894f7a80053",
			"c5f7b282a516a2f831d668ca4448662c92907cbbfa9b90c3745ea2e79b28ec6d",
		},
		returnedBits: "e108b1bd65928d8097190c106321e9fcb3b930ef3b74421ed8e70802d3f06ebe" +
			"7338948cbf08ab91cdf71dd398ddc4752b6d3dbc265fe836588783fbddd509fb",
	},
	{
		keySize:              KeySize192,
		predictionResistance: true,
		entropyInput:         "feea3b6964664c58bc58f5ac8f0e9d35d3ccf87220bb5e7a84c13c68cd9a28f1ff576675224c84c2",
		personalization:      "eebac711edbb4034a7d7fd17c301467bb1369b85aeff93dbe82bce803ca059492caf3e96b5f01d83",
		entropyInputPR: [2]string{
			"ab87b81b15e4203ced5dcfa7a2e92f71784832a5c47f0054850ff32a06c6874747446c69cb481ef4",
			"29de4f15b489c4e8cdc73244570012b443a1242f28f9dd7fa43f98293a0e490aef281eea13d1e136",
		},
		additionalInput: [2]string{
			"a5500c3f37b20cec3391effcac45676674ab8aba5f9690517b80eb6d2610f832564128f87c4f88b2",
			"61384443e07a5659c7bdb9bf9f3fde677ad95737e9260b53635fd7b9d4b63a7995861285ee613265",
		},
		returnedBits: "2f19570a4eebe46b24104c44290a1c206daaaaf99d43f0530df4183322bd798d" +
			"08d795c28e9f67ae24507be481deabb4f760d2f693cd0af236b0ba5996806e34",
	},
	{
		keySize:              KeySize256,
		predictionResistance: true,
		entropyInput:         "c54805274bde00aa5289e0513579019707666d2fa7a1c8908865891c87c0c652335a4d3cc415bc30742b164647f8820f",
		personalization:      "d63fb5afa2101fa4b8a6c3b89d9c250ac728fc1ddad0e7585b5d54728ed20c2f940e89155596e3b963635b6d6088164b",
		entropyInputPR: [2]string{
			"b2ad31d1f20dcf30dd526ec9156c07f270216bdb59197325bab180675929888ab699c54fb21819b7d921d6346bff2f7f",
			"eca449048d26fd38f8ca435237dce66eadec7069ee5dd0b70084b819a711c0820a7556bbd0ae20f06e5169278b593b71",
		},
		additionalInput: [2]string{
			"744bfae3c23a5cc9a3b373b6c50795068d35eb8a339746ac810d16f864e880061082edf9d2687c211960aa83400f85f9",
			"ad55c682962aa4fe9ebc227c9402e79b0aa7874844d33eaee7e2d15baf81d9d33936e4d93f28ad109657b512aee115a5",
		},
		returnedBits: "f08fdfc1775b6feb5a4177110bf29d7c3ab715dfdc4b27200359288c0624bd5c" +
			"1028acc9914d88a82b09f5eaafdc3bca8547b98481df39b86504314221cbdc3c",
	},
}

// runDRBGKAT runs every DRBG known-answer vector and returns ErrDRBGSelfTestFailed if any result differs from the
// known answer. Each test instance is private to the test, so no test output is returned to callers.
func runDRBGKAT() error {
	for _, v := range drbgKATVectors {
		if !v.passes() {
			return ErrDRBGSelfTestFailed
		}
	}
	return nil
}

// passes runs one vector through the instantiate, reseed, and generate functions used with
// BacktrackingResistance and reports whether the second generated output equals the known answer.
func (v drbgKATVector) passes() bool {
	valid := true
	decode := func(s string) []byte {
		b, err := hex.DecodeString(s)
		if err != nil {
			valid = false
		}
		return b
	}

	cfg := DefaultConfig()
	cfg.KeySize = v.keySize
	cfg.Personalization = decode(v.personalization)
	cfg.BacktrackingResistance = true
	cfg.PredictionResistance = v.predictionResistance

	// Entropy for reseeding is read in the order the vector consumes it.
	var reseedEntropy []byte
	if v.predictionResistance {
		reseedEntropy = append(decode(v.entropyInputPR[0]), decode(v.entropyInputPR[1])...)
	} else {
		reseedEntropy = decode(v.entropyInputReseed)
	}

	d := &drbg{
		config:         &cfg,
		pid:            os.Getpid(),
		lastReseedTime: time.Now(),
		entropySource:  bytes.NewReader(reseedEntropy),
	}

	seedMaterial := decode(v.entropyInput)
	mixSeed(seedMaterial, cfg.Personalization, nil)
	if err := d.instantiate(seedMaterial); err != nil {
		return false
	}
	if !v.predictionResistance {
		if err := d.Reseed(decode(v.additionalInputReseed)); err != nil {
			return false
		}
	}

	want := decode(v.returnedBits)
	if !valid {
		return false
	}
	got := make([]byte, len(want))
	for _, ai := range v.additionalInput {
		if _, err := d.ReadWithAdditionalInput(got, decode(ai)); err != nil {
			return false
		}
	}
	defer clear(got)
	return bytes.Equal(got, want)
}
