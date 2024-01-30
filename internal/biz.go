package internal

import "fmt"

type Biz interface {
	Do(*Client)
}

type PublicChat struct {
}

func (b *PublicChat) Do(c *Client) {
	for {
		if c.IsClosed() {
			return
		}

		fmt.Println(">>>>> please enter chat content. enter \"exit\" to quit <<<<<")
		chatMsg := <-c.ReadCh
		if chatMsg == "" {
			continue
		}
		if chatMsg == "exit" {
			return
		}
		sendMsg := chatMsg + "\n"
		if _, err := c.conn.Write([]byte(sendMsg)); err != nil {
			fmt.Println("conn.Write err:", err)
			break
		}
	}
}

type PrivateChat struct {
}

func (b *PrivateChat) Do(c *Client) {
	for {
		if c.IsClosed() {
			return
		}

		b := &SearchOnlineUsers{}
		b.Do(c)

		fmt.Println(">>>>> please enter a user to chat, enter \"exit\" to quit <<<<<")
		remoteUser := <-c.ReadCh
		if remoteUser == "" {
			continue
		}
		if remoteUser == "exit" {
			return
		}

		// chat with remoteUser
		// notice: use break, not return to break the second loop.
		for {
			if c.IsClosed() {
				return
			}

			fmt.Println(">>>>> please enter chat content. enter \"exit\" to quit <<<<<")
			chatMsg := <-c.ReadCh
			if chatMsg == "" {
				continue
			}
			if chatMsg == "exit" {
				break
			}
			sendMsg := "to|" + remoteUser + "|" + chatMsg + "\n"
			if _, err := c.conn.Write([]byte(sendMsg)); err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}
	}
}

type UpdateName struct {
}

func (b *UpdateName) Do(c *Client) {
	fmt.Println(">>>>> please enter username <<<<<")
	name := <-c.ReadCh
	c.Name = name
	sendMsg := "rename|" + c.Name + "\n"
	if _, err := c.conn.Write([]byte(sendMsg)); err != nil {
		fmt.Println("Client::conn.Write err:", err)
	}
}

type SearchOnlineUsers struct {
}

func (b *SearchOnlineUsers) Do(c *Client) {
	sendMsg := "who\n"
	if _, err := c.conn.Write([]byte(sendMsg)); err != nil {
		fmt.Println("conn.Write err:", err)
	}
}
