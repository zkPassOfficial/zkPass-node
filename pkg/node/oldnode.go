package node

import (
	"flag"
	"zkpass-node/app/keystore"
)

// var cm *session.Manager
var ks *keystore.Keystore

// func destroyOnPanic(c *session.Session) {
// 	r := recover()
// 	if r == nil {
// 		return
// 	}
// 	// config

// 	fmt.Println("caught a panic message: ", r)
// 	debug.PrintStack()

// 	c.SignalDisconnect()
// }

// func router(w http.ResponseWriter, req *http.Request) {
// 	cid := string(req.URL.RawQuery)
// 	path := req.URL.Path[1:]
// 	log.Println("got request ", path, " from ", req.RemoteAddr)

// 	var out []byte
// 	if path == "connect" {
// 		c := cm.Connect(cid)
// 		key, keyData := ks.GetConnKey()
// 		c.SetSigningKey(key)
// 		// keyData is sent to Client unencrypted
// 		out = append(out, keyData...)
// 	}

// 	c := cm.GetConnection(cid)
// 	defer destroyOnPanic(c)

// 	router := cm.GetRouter(path, cid)
// 	body := u.ReadBody(req)
// 	out = append(out, router(body)...)
// 	u.WriteResponse(out, w)

// 	if path == "done" { // last interaction
// 		c.SignalDisconnect()
// 	}
// }

// // sends node's master public key to the client
// func masterPublicKey(w http.ResponseWriter, req *http.Request) {
// 	log.Println("in getPubKey", req.RemoteAddr)
// 	u.WriteResponse(ks.GetMasterPublicKeyPEM(), w)
// }

func mainxxx() {
	//  Memory profile:
	// 	  install with: go get github.com/pkg/profile
	// 		then run: curl http://localhost:8080/debug/pprof/heap > heap
	// 		go tool pprof -png heap

	// defer profile.Start(profile.MemProfile).Stop()
	// go func() {
	//   http.ListenAndServe(":8080", nil)
	// }()

	flag.Parse()

	// start keystore
	ks = new(keystore.Keystore)
	ks.Run()

	// start receiving connection
	// cm = new(session.Manager)
	// cm.Run()

	// Get master public key for each connection
	// http.HandleFunc("/key", masterPublicKey)

	// process of 3 parties of TLS
	// http.HandleFunc("/", router)
}
