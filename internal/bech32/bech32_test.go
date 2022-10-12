// Copyright (c) 2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package bech32_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/gopherlearning/gophkeeper/internal/bech32"
	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	str   string
	valid bool
}{
	{"AGE-SECRET-KEY-1D0542368K3SPZCNEJY9G2SRZXW43PCYK9UFEVQ9F84AM87VLMGNSQPN3RA", true},
	{"age1wd3kx284gsypv7vpx474rprrnr7w3rv46dt4hns9qda5eh965yfqpx7vma", true},
	{"an83characterlonghumanreadablepartthatcontainsthenumber1andtheexcludedcharactersbio1tt5tgs", true},
	{"abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw", true},
	{"11qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqc8247j", true},
	{"split1checkupstagehandshakeupstreamerranterredcaperred2y9e3w", true},
	{"AGE-SECRET-KEY-1d0542368k3spzcnejy9g2srzxw43pcyk9ufevq9f84am87vlmgnsqpn3ra", false},
	{"split1checkupstagehandshakeupstreamerranterredcaperred2y9e2w", false},                              // invalid checksum
	{"s lit1checkupstagehandshakeupstreamerranterredcaperredp8hs2p", false},                              // invalid character (space) in hrp
	{"spl" + strconv.QuoteRune(127) + "t1checkupstagehandshakeupstreamerranterredcaperred2y9e3w", false}, // invalid character (DEL) in hrp

	{"split1cheo2y9e2w", false}, // invalid character (o) in data part
	{"split1a2y9w", false},      // too short data part
	{"1checkupstagehandshakeupstreamerranterredcaperred2y9e3w", false},                                     // empty hrp
	{"11qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqsqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqc8247j", false}, // too long
}

func TestBech32(t *testing.T) {
	t.Run("empty hrp", func(t *testing.T) {
		_, err := bech32.Encode("", nil)
		assert.Error(t, err)
		_, err = bech32.Encode("dwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwfffffffffffffffffffffwwwwwwwwww", nil)
		assert.Error(t, err)
		_, err = bech32.Encode("hYh", nil)
		assert.Error(t, err)
	})
	t.Run("long", func(t *testing.T) {
		_, err := bech32.Encode("dwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwfffffffffffffffffffffwwwwwwwwww", nil)
		assert.Error(t, err)
	})

	for _, test := range tests {
		str := test.str
		hrp, decoded, err := bech32.Decode(str)

		if !test.valid {
			// Invalid string decoding should result in error.
			if err == nil {
				assert.Errorf(t, err, "expected decoding to fail for "+
					"invalid string %v", test.str)
			}

			continue
		}

		// Valid string decoding should result in no error.
		if err != nil {
			t.Errorf("expected string to be valid bech32: %v", err)
		}

		// Check that it encodes to the same string
		encoded, err := bech32.Encode(hrp, decoded)
		if err != nil {
			t.Errorf("encoding failed: %v", err)
		}

		fmt.Println(encoded)

		// if encoded != strings.ToLower(str) {
		// 	t.Errorf("expected data to encode to %v, but got %v",
		// 		str, encoded)
		// }

		// Flip a bit in the string an make sure it is caught.
		pos := strings.LastIndexAny(str, "1")
		flipped := str[:pos+1] + string((str[pos+1] ^ 1)) + str[pos+2:]
		_, _, err = bech32.Decode(flipped)

		if err == nil {
			t.Error("expected decoding to fail")
		}
	}
}
