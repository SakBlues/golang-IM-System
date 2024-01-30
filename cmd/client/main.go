package main

import (
	"flag"
	"fmt"

	"github.com/SakBlues/golang-IM-System/internal"
)

var ServerIP string
var ServerPort int

func Init() {
	// ./client -h
	// use cast: ./client -ip 127.0.0.1 -port 8888
	flag.StringVar(&ServerIP, "ip", "127.0.0.1", "set server IP")
	flag.IntVar(&ServerPort, "port", 8888, "set server port")
}

func main() {
	Init()

	// command parse
	flag.Parse()

	client := internal.NewClient(ServerIP, ServerPort)
	if client == nil {
		fmt.Println("<<<<< connnect server failed...")
		return
	}
	fmt.Println("<<<<< connnect server successfully...")

	client.Run()
}
