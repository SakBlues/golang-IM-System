package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// start a goroutine to listen broadcast message.
	go user.ListenBroadCast()
	return user
}

func (this *User) Online() {
	// add user to OnlineMap
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// broadcast user online message
	this.server.SendMsgToBroadcast(this, "online")
}

func (this *User) Offline() {
	// delete u	this.server.mapLock.Lock()
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// broadcast user offline message
	this.server.SendMsgToBroadcast(this, "offline")
}

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// msg format: who
func (this *User) searchOnlineUsers() {
	this.server.mapLock.RLock()
	for _, user := range this.server.OnlineMap {
		onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "online...\n"
		this.SendMsg(onlineMsg)
	}
	this.server.mapLock.RUnlock()
}

// msg format: rename|name
func (this *User) rename(msg string) {
	sendMsg := "username has been used\n"
	newName := strings.Split(msg, "|")[1]
	this.server.mapLock.Lock()
	if _, ok := this.server.OnlineMap[newName]; !ok {
		delete(this.server.OnlineMap, this.Name)
		this.Name = newName
		this.server.OnlineMap[this.Name] = this
		sendMsg = fmt.Sprintf("update username successfully: %s\n", this.Name)
	}
	this.server.mapLock.Unlock()
	this.SendMsg(sendMsg)
}

// msg format: to|name|xxx
func (this *User) privateChat(msg string) {
	contents := strings.Split(msg, "|")
	name := contents[1]
	if name == "" {
		this.SendMsg("msg format incorrect, please use \"to|name|xxx\" format.\n")
		return
	}
	this.server.mapLock.RLock()
	recUser, ok := this.server.OnlineMap[name]
	this.server.mapLock.RUnlock()
	if !ok {
		this.SendMsg("user is not exist\n")
		return
	}
	content := contents[2]
	if content == "" {
		this.SendMsg("no message content, please retry.\n")
		return
	}
	recUser.SendMsg("[private]" + this.Name + ":" + content + "\n")
}

func (this *User) DoMessage(msg string) {
	// search online users
	if msg == "who" {
		this.searchOnlineUsers()
		return
	}

	if len(msg) > 7 && msg[:7] == "rename|" {
		this.rename(msg)
		return
	}

	// msg format:
	if len(msg) > 4 && msg[:3] == "to|" {
		this.privateChat(msg)
		return
	}

	this.server.SendMsgToBroadcast(this, msg)
}

// listen to user.C, send once a message is received
func (this *User) ListenBroadCast() {
	// use for range instead of for loop,
	// so that the method can exit when this.C close,
	// and the outlayer goroutine can exit.
	for msg := range this.C {
		this.conn.Write([]byte(msg + "\n"))
	}
}
