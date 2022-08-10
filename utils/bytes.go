package utils

import (
	"crypto/rand"
	"math/big"
)

// random slice of specified size
func GenRandom(size int) []byte {
	randomBytes := make([]byte, size)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return randomBytes
}

// big.Int into a slice of 16 bytes
func To16Bytes(x *big.Int) []byte {
	buf := make([]byte, 16)
	x.FillBytes(buf)
	return buf
}

// big.Int into a slice of 32 bytes
func To32Bytes(x *big.Int) []byte {
	buf := make([]byte, 32)
	x.FillBytes(buf)
	return buf
}
