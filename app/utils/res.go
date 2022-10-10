package utils

import (
	"io/ioutil"
	"log"
	"net/http"
)

// Read the HTTP request body from client
func ReadBody(req *http.Request) []byte {
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic("can't read request body")
	}
	return body
}

// write HTTP response to client, add cross core origin
func WriteResponse(resp []byte, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(resp)
	log.Println("wrote response of size: ", len(resp))
}
