package utils

import (
	"math/big"
)

func Mul(a, b *big.Int) *big.Int {
	res := new(big.Int)
	res.Mul(a, b)
	return res
}

func Mod(a, b *big.Int) *big.Int {
	res := new(big.Int)
	res.Mod(a, b)
	return res
}

func Sub(a, b *big.Int) *big.Int {
	res := new(big.Int)
	res.Sub(a, b)
	return res
}

func Add(a, b *big.Int) *big.Int {
	res := new(big.Int)
	res.Add(a, b)
	return res
}

func Exp(a, b, c *big.Int) *big.Int {
	res := new(big.Int)
	res.Exp(a, b, c)
	return res
}
