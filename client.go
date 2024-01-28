package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

const (
	exit = iota
	publicChat
	privateChat
	updateName
)

var modeMap = map[int]string{
	exit:        fmt.Sprintf("%d: exit", exit),
	publicChat:  fmt.Sprintf("%d: public chat", publicChat),
	privateChat: fmt.Sprintf("%d: private chat", privateChat),
	updateName:  fmt.Sprintf("%d: update username", updateName),
}

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn

	// client mode
	flag int
}

func NewClient(serverIP string, serverPort int) *Client {
	// create client
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		flag:       -1, // make sure first time not exit loop
	}
	// connect socket
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn
	return client
}

func (this *Client) DealResponse() {
	// once this.conn has data, copy the data to stdout,
	// this will block and listen forever.
	io.Copy(os.Stdout, this.conn)
}

func (this *Client) showMenu() {
	for i := 0; i < len(modeMap); i++ {
		fmt.Println(modeMap[i])
	}
}

func (this *Client) chooseMode() bool {
	this.showMenu()

	var flag int
	flagStr, err := Readline()
	if err != nil {
		fmt.Println(PrintDesc("please enter a legal num"))
		return false
	}
	flag, err = strconv.Atoi(flagStr)
	if err != nil {
		fmt.Println(PrintDesc("please enter a legal num"))
		return false
	}
	if flag < 0 || flag > 3 {
		fmt.Println(PrintDesc("please enter a legal num"))
		return false
	}
	this.flag = flag
	return true
}

func (this *Client) updateName() bool {
	fmt.Println(PrintDesc("please enter username"))

	// fmt.Scanln(&this.Name)
	name, err := Readline()
	if err != nil {
		fmt.Println("Client::Readline err:", err)
		return false
	}
	this.Name = name
	sendMsg := "rename|" + this.Name + "\n"
	if _, err := this.conn.Write([]byte(sendMsg)); err != nil {
		fmt.Println("Client::conn.Write err:", err)
		return false
	}
	return true
}

func (this *Client) publicChat() {
	for {
		fmt.Println(PrintDesc("please enter chat content. enter \"exit\" to quit"))
		chatMsg, err := Readline()
		if err != nil {
			fmt.Println("input.Readline err:", err)
			return
		}
		if chatMsg == "" {
			continue
		}
		if chatMsg == "exit" {
			return
		}
		sendMsg := chatMsg + "\n"
		if _, err := this.conn.Write([]byte(sendMsg)); err != nil {
			fmt.Println("conn.Write err:", err)
			break
		}
	}
}

func (this *Client) searchOnlneUsers() {
	sendMsg := "who\n"
	if _, err := this.conn.Write([]byte(sendMsg)); err != nil {
		fmt.Println("conn.Write err:", err)
	}
}

func (this *Client) privateChat() {
	for {
		this.searchOnlneUsers()

		fmt.Println(PrintDesc("please enter a user to chat, enter \"exit\" to quit"))
		remoteUser, err := Readline()
		if err != nil {
			fmt.Println("privateChat::Readline err:", err)
			return
		}
		if remoteUser == "" {
			continue
		}
		if remoteUser == "exit" {
			return
		}

		// chat with remoteUser
		// notice: use break, not return to break the second loop.
		for {
			fmt.Println(PrintDesc("please enter chat content. enter \"exit\" to quit"))
			chatMsg, err := Readline()
			if err != nil {
				fmt.Println("privateChat::Readline err:", err)
				break
			}
			if chatMsg == "" {
				continue
			}
			if chatMsg == "exit" {
				break
			}
			sendMsg := "to|" + remoteUser + "|" + chatMsg + "\n"
			if _, err := this.conn.Write([]byte(sendMsg)); err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}
	}
}

func (this *Client) Run() {
	for this.flag != 0 {
		// select a mode
		for this.chooseMode() != true {
		}
		switch this.flag {
		case publicChat:
			this.publicChat()
		case privateChat:
			this.privateChat()
		case updateName:
			this.updateName()
		}
	}
}

var serverIP string
var serverPort int

func Init() {
	// ./client -h
	// use cast: ./client -ip 127.0.0.1 -port 8888
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "set server IP")
	flag.IntVar(&serverPort, "port", 8888, "set server port")
}

func main() {
	Init()

	// command parse
	flag.Parse()

	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println(PrintDesc("connnect server failed..."))
		return
	}
	fmt.Println(PrintDesc("connnect server successfully..."))

	// create a goroutinue to deal with server response
	go client.DealResponse()

	client.Run()
}
