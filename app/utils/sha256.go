package utils

import (
	"crypto/sha256"
	"encoding"
	"encoding/binary"
)

func ToSha256(data []byte) []byte {
	ret := sha256.Sum256(data)
	return ret[:]
}

// circuit output: outerState = sha256(key xor opad)
// hmac = sha256(key xor opad || "inner hash")
// sha256 shoud start from middle state(outerState)
func ToSha256FromMid(outerState []byte, data []byte) []byte {
	h := sha256.New()

	// Hash implementations in the standard library (e.g. hash/crc32 and crypto/sha256)
	// implement the encoding.BinaryMarshaler and encoding.BinaryUnmarshaler interfaces.
	u, ok := h.(encoding.BinaryUnmarshaler)
	if !ok {
		panic("hash did not implement BinaryUnmarshaler")
	}

	var state []byte

	// Base on the source code of sha256.go, the previous state needs a certain format
	magic256 := "sha\x03"
	state = append(state, magic256...)
	state = append(state, outerState...)

	// set previous chunk to be zeroes
	state = append(state, make([]byte, 64)...)

	// set previous 64 bytes has been processed
	var a [8]byte
	binary.BigEndian.PutUint64(a[:], 64)
	state = append(state, a[:]...)

	if err := u.UnmarshalBinary(state); err != nil {
		panic(err)
	}

	h.Write(data)

	return h.Sum(nil)
}
