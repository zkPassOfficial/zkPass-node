package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"zkpass-node/connection"
	"zkpass-node/keystore"
	u "zkpass-node/utils"
)

var cm *connection.Manager
var ks *keystore.Keystore

func destroyOnPanic(c *connection.Connection) {
	r := recover()
	if r == nil {
		return
	}

	fmt.Println("caught a panic message: ", r)
	debug.PrintStack()

	c.SignalDisconnect()
}

func router(w http.ResponseWriter, req *http.Request) {
	cid := string(req.URL.RawQuery)
	path := req.URL.Path[1:]
	log.Println("got request ", path, " from ", req.RemoteAddr)

	var out []byte
	if path == "connect" {
		c := cm.Connect(cid)
		key, keyData := ks.GetConnKey()
		c.SetSigningKey(key)
		// keyData is sent to Client unencrypted
		out = append(out, keyData...)
	}

	c := cm.GetConnection(cid)
	defer destroyOnPanic(c)

	router := cm.GetRouter(path, cid)
	body := u.ReadBody(req)
	out = append(out, router(body)...)
	u.WriteResponse(out, w)

	if path == "done" { // last interaction
		c.SignalDisconnect()
	}
}

// heartbeat
func ping(w http.ResponseWriter, req *http.Request) {
	log.Println("in ping", req.RemoteAddr)
	u.WriteResponse(nil, w)
}

// sends node's master public key to the client
func masterPublicKey(w http.ResponseWriter, req *http.Request) {
	log.Println("in getPubKey", req.RemoteAddr)
	u.WriteResponse(ks.GetMasterPublicKeyPEM(), w)
}

func main() {
	//  Memory profile:
	// 	  install with: go get github.com/pkg/profile
	// 		then run: curl http://localhost:8080/debug/pprof/heap > heap
	// 		go tool pprof -png heap

	// defer profile.Start(profile.MemProfile).Stop()
	// go func() {
	//   http.ListenAndServe(":8080", nil)
	// }()

	host := flag.String("host", "0.0.0.0", "Server Host")
	port := flag.Int("port", 3333, "Server Port")
	flag.Parse()
	log.Println("server info", *host, *port)

	// start keystore
	ks = new(keystore.Keystore)
	ks.Run()

	// start receiving connection
	cm = new(connection.Manager)
	cm.Run()

	// TODO: need to start initialization of GC

	// heartbeat
	http.HandleFunc("/ping", ping)

	// Get master public key for each connection
	http.HandleFunc("/key", masterPublicKey)

	// process of 3 parties of TLS
	http.HandleFunc("/", router)

	addr := fmt.Sprintf("%s:%d", *host, *port)
	// listen to port
	http.ListenAndServe(addr, nil)
}
