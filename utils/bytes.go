package utils

import (
	"crypto/rand"
	"math"
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

func Concat(slices ...[]byte) []byte {
	totalSize := 0
	for _, v := range slices {
		totalSize += len(v)
	}
	buf := make([]byte, totalSize)
	offset := 0
	for _, v := range slices {
		copy(buf[offset:offset+len(v)], v)
		offset += len(v)
	}
	return buf
}

func BitsToBytes(b []int) []byte {
	bigint := new(big.Int)
	for i := 0; i < len(b); i++ {
		bigint.SetBit(bigint, i, uint(b[i]))
	}
	len := int(math.Ceil(float64(len(b)) / 8))
	buf := make([]byte, len)
	bigint.FillBytes(buf)
	return buf
}

func XorBytes(a, b []byte) []byte {
	if len(a) != len(b) {
		panic("len(a) != len(b)")
	}
	buf := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		buf[i] = a[i] ^ b[i]
	}
	return buf
}

func Flatten(matrix [][]byte) []byte {
	var buf []byte
	for i := 0; i < len(matrix); i++ {
		buf = append(buf, matrix[i]...)
	}
	return buf
}

func Slice(data []byte, chunkSize int) [][]byte {
	if len(data)%chunkSize != 0 {
		panic("len(data) % chunkSize != 0")
	}
	size := len(data) / chunkSize
	chunks := make([][]byte, size)
	for i := 0; i < size; i++ {
		chunks[i] = data[i*chunkSize : (i+1)*chunkSize]
	}
	return chunks
}
