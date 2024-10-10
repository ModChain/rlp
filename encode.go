package rlp

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

func rlpEncodeSigned[T int | int8 | int16 | int32 | int64](v T) ([]byte, error) {
	if v < 0 {
		return nil, ErrRlpMinusValue
	}
	return rlpEncodeUnsigned(uint64(v))
}

func rlpEncodeUnsigned[T uint | uint8 | uint16 | uint32 | uint64 | uintptr](v T) ([]byte, error) {
	return EncodeValue(rlpTrim(binary.BigEndian.AppendUint64(nil, uint64(v))))
}

// rlpEncodeLen returns a buffer that encodes the length of the incoming buffer or array
func rlpEncodeLen(ln int, isArray bool) []byte {
	if ln <= 55 {
		if isArray {
			return []byte{byte(0xc0 + ln)}
		} else {
			return []byte{byte(0x80 + ln)}
		}
	}
	buf := binary.BigEndian.AppendUint64(nil, uint64(ln))
	for len(buf) > 0 && buf[0] == 0 {
		buf = buf[1:]
	}
	if isArray {
		return append([]byte{0xf7 + byte(len(buf))}, buf...)
	} else {
		return append([]byte{0xb7 + byte(len(buf))}, buf...)
	}
}

func rlpTrim(in []byte) []byte {
	for len(in) > 0 && in[0] == 0 {
		in = in[1:]
	}
	return in
}

// EncodeValue encodes a single value into rlp format
func EncodeValue(v any) ([]byte, error) {
	switch in := v.(type) {
	case int8:
		return rlpEncodeSigned(in)
	case int16:
		return rlpEncodeSigned(in)
	case int32:
		return rlpEncodeSigned(in)
	case int64:
		return rlpEncodeSigned(in)
	case int:
		return rlpEncodeSigned(in)
	case uint:
		return rlpEncodeUnsigned(in)
	case uint8:
		return rlpEncodeUnsigned(in)
	case uint16:
		return rlpEncodeUnsigned(in)
	case uint32:
		return rlpEncodeUnsigned(in)
	case uint64:
		return rlpEncodeUnsigned(in)
	case uintptr:
		return rlpEncodeUnsigned(in)
	case *big.Int:
		return EncodeValue(rlpTrim(in.Bytes()))
	case []byte:
		if len(in) == 1 && in[0] <= 0x7f {
			return []byte{byte(in[0])}, nil
		}
		return append(rlpEncodeLen(len(in), false), in...), nil
	case string:
		// 0x...
		var ok bool
		if in, ok = strings.CutPrefix(in, "0x"); !ok {
			return nil, ErrStringPrefix
		}
		if len(in) == 0 {
			return EncodeValue([]byte{})
		}
		if len(in)&1 == 1 {
			in = "0" + in
		}
		buf, err := hex.DecodeString(in)
		if err != nil {
			return nil, err
		}
		return EncodeValue(buf)
	case []any:
		return Encode(in...)
	default:
		return nil, fmt.Errorf("unsupported type %T", in)
	}
}

// Encode encodes a number of arguments into a RLP array
func Encode(in ...any) ([]byte, error) {
	var buf []byte
	for _, v := range in {
		t, err := EncodeValue(v)
		if err != nil {
			return nil, err
		}
		buf = append(buf, t...)
	}
	return append(rlpEncodeLen(len(buf), true), buf...), nil
}
