package utils

import "crypto/sha256"

func ToSha256(data []byte) []byte {
	ret := sha256.Sum256(data)
	return ret[:]
}
