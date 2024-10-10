package rlp_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
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
		&encodeTestVector{
			In:  []any{[]byte("cat"), []any{[]byte("puppy"), []byte("cow")}, []byte("horse"), []any{[]any{}}, []byte("pig"), []any{"0x"}, []byte("sheep")},
			Out: "e383636174ca85707570707983636f7785686f727365c1c083706967c180857368656570",
			Dec: "[[636174 [7075707079 636f77] 686f727365 [[]] 706967 [] 7368656570]]",
		},
		&encodeTestVector{
			In:  []any{},
			Out: "c0", // the empty list = [ 0xc0 ]
			Dec: "[[]]",
		},
		&encodeTestVector{
			In:  []any{[]byte{}},
			Out: "c180", // the empty string ('null') = [ 0x80 ]
			Dec: "[[]]",
		},
		&encodeTestVector{
			In:  []any{[]byte{0x42}},
			Out: "c142",
			Dec: "[[42]]",
		},
		&encodeTestVector{
			In:  []any{[]byte("dog")},
			Out: "c483646f67", // the string "dog" = [ 0x83, 'd', 'o', 'g' ]
			Dec: "[[646f67]]",
		},
		&encodeTestVector{
			In:  []any{[]byte("cat"), []byte("dog")},
			Out: "c88363617483646f67", // the list [ "cat", "dog" ] = [ 0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g' ]
			Dec: "[[636174 646f67]]",
		},
		&encodeTestVector{
			In:  []any{new(big.Int)},
			Out: "c180", // the integer 0 = [ 0x80 ]
			Dec: "[[]]",
		},
		&encodeTestVector{
			In:  []any{[]byte{0}},
			Out: "c100", // the byte '\x00' = [ 0x00 ]
			Dec: "[[00]]",
		},
		&encodeTestVector{
			In:  []any{[]byte{0x0f}},
			Out: "c10f", // the byte '\x0f' = [ 0x0f ]
			Dec: "[[0f]]",
		},
		&encodeTestVector{
			In:  []any{[]byte{0x04, 0x00}},
			Out: "c3820400", // the bytes '\x04\x00' = [ 0x82, 0x04, 0x00 ]
			Dec: "[[0400]]",
		},
		&encodeTestVector{
			In:  []any{[]any{}, []any{[]any{}}, []any{[]any{}, []any{[]any{}}}},
			Out: "c7c0c1c0c3c0c1c0", // the set theoretical representation(opens in a new tab) of three, [ [], [[]], [ [], [[]] ] ] = [ 0xc7, 0xc0, 0xc1, 0xc0, 0xc3, 0xc0, 0xc1, 0xc0 ]
			Dec: "[[[] [[]] [[] [[]]]]]",
		},
		&encodeTestVector{
			In:  []any{[]byte("Lorem ipsum dolor sit amet, consectetur adipisicing elit")}, //
			Out: "f83ab8384c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e7365637465747572206164697069736963696e6720656c6974",
			Dec: "[[4c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e7365637465747572206164697069736963696e6720656c6974]]",
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
