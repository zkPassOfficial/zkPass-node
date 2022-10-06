package keystore

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/binary"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	u "zkpass-node/utils"
)

type Keystore struct {
	sync.Mutex

	// startTime|endTime|publicKey|signature for client to verify
	keyData []byte

	// key for each connection
	connKey *ecdsa.PrivateKey

	// used to sign connection keys
	masterKey *ecdsa.PrivateKey

	// master public key PEM format
	masterPublicKeyPEM []byte

	// key valid duration in second
	keyValidDuration int
}

func (k *Keystore) Run() {
	k.genMasterKey()
	go k.rotatingKeys()
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

// rotating connection keys and sign it with the master key during a period of time
func (k *Keystore) rotatingKeys() {
	k.keyValidDuration = 15 * 60 // hard code it to be 15 mins at current situation

	// init to 0, so that the connection will do the rotation at beginning
	nextKeyRotationTime := time.Unix(0, 0)

	for {
		time.Sleep(time.Second * 1)
		now := time.Now()

		// 4 mins means the connection is expired
		if nextKeyRotationTime.Sub(now) > time.Minute*2 {
			continue
		}

		// pick a random interval to avoid side-channel attacks
		randInterval := u.RandInt(k.keyValidDuration/2, k.keyValidDuration)
		nextKeyRotationTime = now.Add(time.Second * time.Duration(randInterval))

		// key valid start time
		startTime := make([]byte, 4)
		binary.BigEndian.PutUint32(startTime, uint32(now.Unix()))

		// key valid end time
		endTime := make([]byte, 4)
		interval := now.Add(time.Second * time.Duration(k.keyValidDuration))
		binary.BigEndian.PutUint32(endTime, uint32(interval.Unix()))

		rotateMasterKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			log.Fatalln("Could not create keys:", err)
		}

		publicKey := u.Concat([]byte{0x04}, u.To32Bytes(rotateMasterKey.PublicKey.X), u.To32Bytes(rotateMasterKey.PublicKey.Y))
		signature := u.ECDSASign(k.masterKey, startTime, endTime, publicKey)
		keyData := u.Concat(startTime, endTime, publicKey, signature)

		k.Lock()
		k.keyData = keyData
		k.masterKey = rotateMasterKey
		k.Unlock()
	}
}
