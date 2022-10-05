package utils

import (
	"math/big"
)

// hex string to big.Int
func H2Bi(a string) *big.Int {
	res, ok := new(big.Int).SetString(a, 16)
	if !ok {
		panic("in h2bi")
	}
	return res
}
