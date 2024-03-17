package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	id          uint64
	msgQueue    chan Message
	incomeQueue chan *string
	conn        net.Conn
	leftSignal  chan bool
}

func (client *Client) Run() {
	writeStop := make(chan bool, 1)
	go client.Read(&writeStop)
	go client.Write(&writeStop)
}

func (client *Client) Read(writeStop *chan bool) {
	buf := make([]byte, 256)
	for {
		len, err := client.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("A client left")
			} else {
				fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err)
			}
			break
		}
		client.msgQueue <- Message{sender: client, body: string(buf[:len])}
	}
	client.leftSignal <- true
	*writeStop <- true
	client.conn.Close()
}

func (client *Client) Write(stop *chan bool) {
	for {
		select {
		case <-(*stop):
			return

		case msg := <-client.incomeQueue:
			client.conn.Write([]byte(*msg))

		default:
			// do nothing
		}
	}
}
