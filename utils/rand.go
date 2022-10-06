package utils

import (
	"encoding/binary"
	"math/rand"
	"time"
)

func GenRandom(size int) []byte {
	randomBytes := make([]byte, size)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return randomBytes
}

func RandInt(min, max int) int {
	rand.Seed(int64(binary.BigEndian.Uint64(GenRandom(8))))
	return rand.Intn(max-min) + min
}

func RandString() string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
