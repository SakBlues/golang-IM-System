package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/SakBlues/golang-IM-System/internal"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	server := internal.NewServer("127.0.0.1", 8888)
	server.Start()
}
