package utils

import (
	"crypto/rand"
	"log"

	paillier "github.com/roasbeef/go-go-gadget-paillier"
)

func PaillierGenerateKey() *paillier.PrivateKey {
	for {
		// double-check that n has 1536 bits
		paillierPrivKey, _ := paillier.GenerateKey(rand.Reader, 1536)
		if len(paillierPrivKey.PublicKey.N.Bytes()) == 192 {
			return paillierPrivKey
		}
		log.Println("n is not 1536 bits")
	}
}

func PaillierEncrypt(privKey *paillier.PrivateKey, payload []byte) []byte {
	res, err := paillier.Encrypt(&privKey.PublicKey, payload)
	if err != nil {
		panic(err)
	}
	return res
}

func PaillierDecrypt(privKey *paillier.PrivateKey, payload []byte) []byte {
	res, err := paillier.Decrypt(privKey, payload)
	if err != nil {
		panic(err)
	}
	return res
}
