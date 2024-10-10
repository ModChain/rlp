package rlp

import (
	"encoding/binary"
	"io"
)

// DecodeUint64 returns a uint64 value for a 0~8 bytes long value, and will panic if
// the buffer is longer than that
func DecodeUint64(buf []byte) uint64 {
	if len(buf) > 8 {
		panic("DecodeUint64 input buffer length higher than 8")
	}
	// buf is a trimmed bigendian uint64
	var tmp [8]byte
	copy(tmp[8-len(buf):], buf)
	return binary.BigEndian.Uint64(tmp[:])
}

// Decode returns an array of values for the given RLP array.
func Decode(buf []byte) ([]any, error) {
	var res []any

	for len(buf) > 0 {
		c := buf[0]
		buf = buf[1:]

		if c <= 0x7f {
			res = append(res, []byte{c})
			continue
		}

		isArray := c&0x40 == 0x40
		ln := uint64(c & 0x3f) // 0~55 = as is, 56~63 = actual value comes after
		if len(buf) < int(ln) {
			return nil, io.ErrUnexpectedEOF
		}

		if ln > 55 {
			ln -= 55
			ln, buf = DecodeUint64(buf[:ln]), buf[ln:]
			if len(buf) < int(ln) {
				return nil, io.ErrUnexpectedEOF
			}
		}
		v := buf[:ln]
		buf = buf[ln:]

		if isArray {
			tmp, err := Decode(v)
			if err != nil {
				return nil, err
			}
			res = append(res, tmp)
		} else {
			res = append(res, v)
		}
	}
	return res, nil
}
