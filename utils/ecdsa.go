package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func ECDSASign(key *ecdsa.PrivateKey, items ...[]byte) []byte {
	var concatAll []byte
	for _, item := range items {
		concatAll = append(concatAll, item...)
	}
	digest_to_be_signed := ToSha256(concatAll)
	r, s, err := ecdsa.Sign(rand.Reader, key, digest_to_be_signed)
	if err != nil {
		panic("ecdsa.Sign")
	}
	signature := append(To32Bytes(r), To32Bytes(s)...)
	return signature
}

func ECDSAPubkeyToPEM(key *ecdsa.PublicKey) []byte {
	derBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		fmt.Println(err)
		panic("x509.MarshalPKIXPublicKey")
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derBytes,
	}
	pubKeyPEM := pem.EncodeToMemory(block)
	return pubKeyPEM
}
