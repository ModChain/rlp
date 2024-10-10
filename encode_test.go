package rlp_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ModChain/rlp"
)

type encodeTestVector struct {
	In  []any
	Out string // as hex
	Dec string
}

func TestEncode(t *testing.T) {
	tests := []*encodeTestVector{
		&encodeTestVector{
			In:  []any{42, int32(123456789), 21000, "0xabdef0123456789abcdef0123456789012345789", []byte{1, 2, 3, 4, 5, 6}},
			Out: "e52a84075bcd1582520894abdef0123456789abcdef012345678901234578986010203040506",
			Dec: "[[2a 075bcd15 5208 abdef0123456789abcdef0123456789012345789 010203040506]]",
		},
	}

	for _, v := range tests {
		res, err := rlp.Encode(v.In...)
		if err != nil {
			t.Errorf("encoding error: %s", err)
		}
		resH := hex.EncodeToString(res)
		if resH != v.Out {
			t.Errorf("error in test: expected %s but got %s", v.Out, resH)
		}

		// attempt to decode this
		dec, err := rlp.Decode(res)
		if err != nil {
			t.Errorf("decoding error: %s", err)
		}
		decS := fmt.Sprintf("%x", dec)
		if decS != v.Dec {
			t.Errorf("error in test: expected dec %s but got %s", v.Dec, decS)
		}
	}
}
