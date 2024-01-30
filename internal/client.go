package internal

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/SakBlues/golang-IM-System/pkg"
)

const (
	exitMode = iota
	publicChatMode
	privateChatMode
	updateNameMode
)

var modeMap = map[int]string{
	exitMode:        fmt.Sprintf("%d: exit", exitMode),
	publicChatMode:  fmt.Sprintf("%d: public chat", publicChatMode),
	privateChatMode: fmt.Sprintf("%d: private chat", privateChatMode),
	updateNameMode:  fmt.Sprintf("%d: update username", updateNameMode),
}

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn

	// client mode
	flag int

	// unified input processing to avoid read blocking
	// here use the stdout as output, so do not need WriteCh
	ReadCh chan string

	// everyone can use this chan to send a close request.
	toClose chan struct{}

	// real close, deal with the resources.
	CloseCh chan struct{}
}

func NewClient(serverIP string, serverPort int) *Client {
	// create client
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		flag:       -1, // make sure first time not exit loop
		toClose:    make(chan struct{}, 1),
		CloseCh:    make(chan struct{}),
		ReadCh:     make(chan string),
	}
	// connect socket
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn

	// goroutinues to deal something
	go client.waitToClose()
	go client.DealResponse()
	go client.DealInput()

	return client
}

// A mediator to close
func (this *Client) waitToClose() {
	<-this.toClose
	close(this.CloseCh)
	close(this.ReadCh)
}

// only a mediator, will just do one time in a client.
func (this *Client) tryClose() {
	select {
	case this.toClose <- struct{}{}:
	default:
	}
}

func (this *Client) IsClosed() bool {
	select {
	case <-this.CloseCh:
		return true
	default:
	}
	return false
}

func (this *Client) DealResponse() {
	// once this.conn has data, copy the data to stdout,
	// this will block and listen forever.
	io.Copy(os.Stdout, this.conn)

	// if conn close, close client
	this.tryClose()
}

func (this *Client) DealInput() {
	for {
		str, err := pkg.Readline()
		if err != nil {
			// TODO: maybe can do something...
			fmt.Println("DealInput::Readline err:", err)
			continue
		}
		select {
		case <-this.CloseCh:
			return
		case this.ReadCh <- str:
		}
	}
}

func (this *Client) showMenu() {
	for i := 0; i < len(modeMap); i++ {
		fmt.Println(modeMap[i])
	}
}

func (this *Client) checkMode() bool {
	var flag int
	flagStr, ok := <-this.ReadCh
	if !ok {
		return false
	}
	flag, err := strconv.Atoi(flagStr)
	if err != nil {
		fmt.Println(">>>>> not a num, please try again ... <<<<<")
		return false
	}
	if flag < 0 || flag > 3 {
		fmt.Println(">>>>> num is illegal, please try again ... <<<<<")
		return false
	}
	this.flag = flag
	return true
}

func (this *Client) Run() {
	var biz Biz
	for this.flag != 0 {
		if this.IsClosed() {
			return
		}

		this.showMenu()
		// select a mode
		if this.checkMode() != true {
			continue
		}

		switch this.flag {
		case publicChatMode:
			biz = &PublicChat{}
		case privateChatMode:
			biz = &PrivateChat{}
		case updateNameMode:
			biz = &UpdateName{}
		}
		biz.Do(this)
	}

}
