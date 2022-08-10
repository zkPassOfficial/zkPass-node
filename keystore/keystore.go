package keystore

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"os"
	"path/filepath"
	"sync"

	u "zkpass-node/utils"
)

type Keystore struct {
	sync.Mutex

	// validFrom|validUntil|pubkey|signature for client to verify
	keyData []byte

	// key for each connection
	connKey *ecdsa.PrivateKey

	// used to sign connection keys
	masterKey *ecdsa.PrivateKey

	//master public key PEM format
	masterPublicKeyPEM []byte
}

func (k *Keystore) Run() {
	k.genMasterKey()
	k.rotatingKeys()
}

func (k *Keystore) GetConnKey() (ecdsa.PrivateKey, []byte) {
	k.Lock()
	keyData := make([]byte, len(k.keyData))
	copy(keyData, k.keyData)
	key := *k.connKey
	k.Unlock()
	return key, keyData
}

func (k *Keystore) GetMasterPublicKeyPEM() []byte {
	return k.masterPublicKeyPEM
}

func (k *Keystore) genMasterKey() {
	var err error

	k.masterKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalln("Could not create keys:", err)
		panic("ecdsa.GenerateKey")
	}

	k.masterPublicKeyPEM = u.ECDSAPubkeyToPEM(&k.masterKey.PublicKey)

	curDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filepath.Join(curDir, "public.key"), k.masterPublicKeyPEM, 0644)
	if err != nil {
		panic(err)
	}
}

func (k *Keystore) rotatingKeys() {
	// TODO: rotating connection keys and sign it with the master key during a period of time
}
