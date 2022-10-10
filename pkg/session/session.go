package session

import (
	"crypto/ecdsa"
	"time"
)

type Session struct {
	Id          string
	SigningKey  ecdsa.PrivateKey
	DestroyChan chan string
}

func (c *Session) SetSigningKey(key ecdsa.PrivateKey) {
	c.SigningKey = key
}

func (c *Session) SignalDisconnect() {
	c.DestroyChan <- c.Id
}

type TimeoutSession struct {
	session   *Session
	updatedAt int64
	createdAt int64
}

func (ts *TimeoutSession) Alive(timeout, life int64) bool {
	now := int64(time.Now().UnixNano() / 1e9)
	return now-ts.updatedAt > timeout || now-ts.createdAt > life
}

func (ts *TimeoutSession) Update() {
	now := int64(time.Now().UnixNano() / 1e9)
	ts.updatedAt = now
}
