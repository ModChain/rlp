package rlp

import "errors"

var (
	ErrRlpMinusValue = errors.New("rlp: negative values are not encodable")
	ErrStringPrefix  = errors.New("rlp: strings must start with 0x")
)
