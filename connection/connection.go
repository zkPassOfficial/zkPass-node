package connection

import (
	"crypto/ecdsa"
	"log"
)

type Connection struct {
	// sign the connection
	signingKey ecdsa.PrivateKey
	// Used to signal to manager whether it can be dropped
	cid string
	// Used to send cid when the connection needs to be dropped
	destroyChan chan string
}

func (c *Connection) SetSigningKey(key ecdsa.PrivateKey) {
	c.signingKey = key
}

func (c *Connection) SignalDisconnect() {
	c.destroyChan <- c.cid
}

func (c *Connection) test(body []byte) []byte {
	log.Println("Connection test")

	return []byte{1, 2, 3, 4, 5, 6}
}
