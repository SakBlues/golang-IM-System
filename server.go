package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	// base.
	IP   string
	Port int

	// record online user list.
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// a channel to broadcast message.
	BroadcastChan chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:            ip,
		Port:          port,
		OnlineMap:     make(map[string]*User),
		BroadcastChan: make(chan string),
	}
	return server
}

// listen to user.Message channel, broadcast to all online user
func (this *Server) Broadcast() {
	for {
		msg := <-this.BroadcastChan
		this.mapLock.RLock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.RUnlock()
	}
}

func (this *Server) SendMsgToBroadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.BroadcastChan <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	user := NewUser(conn, this)
	user.Online()

	isAlive := make(chan struct{})

	// receive user Message
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			// remove '\n'
			msg := string(buf[:n-1])
			user.DoMessage(msg)
			// send any message, stands for user is alive.
			isAlive <- struct{}{}
		}
	}()

	for {
		select {
		case <-isAlive: // user is alive, go to next loop
		case <-time.After(5 * time.Minute): // timeout
			user.SendMsg("you have been kicked offline")
			// although close channel will cause goroutine above offline, and delete OnlineMap here.
			// however, maybe another goroutine broadcast and send to the close channel cause panic.
			// so need to delete the user here.
			this.mapLock.Lock()
			delete(this.OnlineMap, user.Name)
			this.mapLock.Unlock()

			// recycle resources and exit.
			close(user.C)
			user.conn.Close()
			return
		}
	}
}

func (this *Server) Start() {
	// TODO: Use concurrency model

	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.IP, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()

	// start a goroutine to broadcast some user messages.
	go this.Broadcast()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		// do business
		go this.Handler(conn)
	}
}
