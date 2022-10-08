package utils

import (
	"encoding/binary"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/salsa20/salsa"
)

// replace aes by salsa20, 3x times fast
func randomOracle(msg []byte, nonce uint32) []byte {
	if len(msg) != 16 {
		panic(len(msg) != 16)
	}
	fixedKey := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		20, 21, 22, 23, 24, 25, 26, 27, 28, 0, 0, 0, 0}
	tBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tBytes, nonce)
	copy(fixedKey[28:32], tBytes)
	out := make([]byte, 16)
	var msgArray [16]byte
	copy(msgArray[:], msg)
	salsa.XORKeyStream(out, out, &msgArray, &fixedKey)
	return out
}

// H(A,B) + C
func Encrypt(a, b []byte, nonce uint32, c []byte) []byte {

	a2 := make([]byte, 16)
	copy(a2[:], a[:])
	tailA := make([]byte, 1)
	copy(tailA, a2[0:1])
	copy(a2[:], a2[1:15])
	copy(a2[14:15], tailA)
	b4 := make([]byte, 16)

	copy(b4[:], b[:])
	tailB := make([]byte, 2)
	copy(tailB, b4[0:2])
	copy(b4[:], b4[2:15])
	copy(b4[13:15], tailB)

	k := XorBytes(a2, b4)
	hash := randomOracle(k, nonce)
	cXorK := XorBytes(c, k)
	return XorBytes(cXorK, hash)
}

func Decrypt(a, b []byte, nonce uint32, c []byte) []byte {
	return Encrypt(a, b, nonce, c)
}

// Circular Correlation Robustness Hash
func CCRHash(length int, msg []byte) []byte {
	h, err := blake2b.New(length, nil)
	if err != nil {
		panic("error in CCRHash")
	}
	_, err = h.Write(msg)
	if err != nil {
		panic("error in CCRHash")
	}
	return h.Sum(nil)
}
