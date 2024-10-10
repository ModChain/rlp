[![GoDoc](https://godoc.org/github.com/ModChain/rlp?status.svg)](https://godoc.org/github.com/ModChain/rlp)

# Recursive Length Prefix encoding (RLP)

See: https://ethereum.org/en/developers/docs/data-structures-and-encoding/rlp

This is a very lightweight implementation of rlp decoding & encoding, depending on no external library.

a typical transaction: RLP([nonce, gasPrice, gasLimit, to, value, data, v, r, s])

## Encoding an array of values

This encoder supports a number of value types, including integer values, `big.Int`, byte arrays and strings.

```go
import "github.com/ModChain/rlp"

buf, err := rlp.Encode(nonce, gasPrice, 21000, "0x123456...", valueBig, []byte{})
```

## Decoding

When decoding, all values will be decoded as []byte. Decode() will return a `[]any` that can have values that
are either []byte or another []any (recursive).
