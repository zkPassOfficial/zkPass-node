package connection

import (
	"log"
	"sync"
	"time"
)

type Router func([]byte) []byte

type ConnItem struct {
	conn *Connection
	// methodLookup is a map used to look up the session's method by its name
	router   map[string]Router
	activeAt int64
	createAt int64
}

type Manager struct {
	conns       map[string]*ConnItem
	destroyChan chan string
	sync.Mutex
}

func (cm *Manager) Run() {
	cm.conns = make(map[string]*ConnItem)

	// auto disconnect the inactive connection
	go cm.autoDisconnect()

	cm.destroyChan = make(chan string)

	// drop the connection with signal
	go cm.signalDisconnect()
}

// addSession creates a new session and sets its creation time
func (cm *Manager) Connect(cid string) *Connection {
	if _, ok := cm.conns[cid]; ok {
		log.Println("Error: connection already exists ", cid)
		return nil // TODO: for connection with new rotated of key, needs to recreate it
	}
	c := new(Connection)
	c.cid = cid
	c.destroyChan = cm.destroyChan
	now := int64(time.Now().UnixNano() / 1e9)

	router := map[string]Router{
		"test": c.test,

		// TODO: each function to deal with

	}

	cm.Lock()
	defer cm.Unlock()
	cm.conns[cid] = &ConnItem{c, router, now, now}
	return c
}

func (cm *Manager) GetConnection(cid string) *Connection {
	item, ok := cm.conns[cid]
	if !ok {
		log.Println("Error: connection not exist ", cid)
		return nil
	}

	item.activeAt = int64(time.Now().UnixNano() / 1e9)

	return item.conn
}

func (cm *Manager) GetRouter(path string, cid string) Router {
	item, ok := cm.conns[cid]

	if !ok {
		log.Println("Error: connection not exist", cid)
		panic("Error: connection not exist")
	}

	f, ok := item.router[path]
	if !ok {
		log.Println("Error: router not exist ", cid)
		panic("Error: router not exist")
	}
	return f
}

func (cm *Manager) disconnect(cid string) {
	cm.Lock()
	defer cm.Unlock()
	delete(cm.conns, cid)
}

// auto remove connection which is timeout
func (cm *Manager) autoDisconnect() {
	const activeTimeout int64 = 180 // 3 mins may enough to gen zkp
	const createTimeout int64 = 360 // 6 mins

	for {
		time.Sleep(time.Second)
		now := int64(time.Now().UnixNano() / 1e9)
		for cid, t := range cm.conns {
			if now-t.activeAt > activeTimeout || now-t.createAt > createTimeout {
				log.Println("auto remove connection ", cid)
				cm.disconnect(cid)
			}
		}
	}
}

// drop for a signal from a connection
func (cm *Manager) signalDisconnect() {
	for {
		cid := <-cm.destroyChan
		log.Println("auto drop cid: ", cid)
		cm.disconnect(cid)
	}
}
