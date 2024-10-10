package rlp

import (
	"encoding/binary"
	"io"
)

func rlpDecodeVarLen(buf []byte) uint64 {
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

		switch {
		case c <= 0x7f:
			res = append(res, []byte{c})
		case c <= 0xb7:
			// value, len <= 55
			ln := c - 0x80
			if len(buf) < int(ln) {
				return nil, io.ErrUnexpectedEOF
			}
			sub := buf[:ln]
			buf = buf[ln:]
			res = append(res, sub)
		case c <= 0xbf:
			// value, len > 55
			lnln := c - 0xb7 // 1~8
			if len(buf) < int(lnln) {
				return nil, io.ErrUnexpectedEOF
			}
			ln := rlpDecodeVarLen(buf[:lnln])
			buf = buf[lnln:]
			if len(buf) < int(ln) {
				return nil, io.ErrUnexpectedEOF
			}
			sub := buf[:ln]
			buf = buf[ln:]
			res = append(res, sub)
		case c <= 0xf7:
			// array, len <= 55
			ln := c - 0xc0
			if len(buf) < int(ln) {
				return nil, io.ErrUnexpectedEOF
			}
			sub := buf[:ln]
			buf = buf[ln:]
			tmp, err := Decode(sub)
			if err != nil {
				return nil, err
			}
			res = append(res, tmp)
		default:
			lnln := c - 0xf7 // 1~8
			if len(buf) < int(lnln) {
				return nil, io.ErrUnexpectedEOF
			}
			ln := rlpDecodeVarLen(buf[:lnln])
			if len(buf) < int(ln) {
				return nil, io.ErrUnexpectedEOF
			}
			buf = buf[lnln:]
			sub := buf[:ln]
			buf = buf[ln:]
			tmp, err := Decode(sub)
			if err != nil {
				return nil, err
			}
			res = append(res, tmp)
		}
	}
	return res, nil
}
