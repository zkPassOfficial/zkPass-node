package utils

import (
	"log"
	"math/big"
	"testing"
)

func TestGenRandom(t *testing.T) {

	for i := 0; i < 5; i++ {
		r := GenRandom(16)
		big_r := new(big.Int)
		big_r.SetBytes(r)
		// int_r := binary.BigEndian.Uint64(r)
		log.Println(i, r, big_r)
	}
}
