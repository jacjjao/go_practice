package main

import (
	"fmt"
	"net"
	"os"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func listenPort(port string) *net.TCPListener {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", port)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	return listener
}

func acceptConnection(server *Server, listener *net.TCPListener) {
	var count uint64 = 0
	defer listener.Close()
	for {
		con, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when accepting connection: %s\n", err)
			continue
		}
		client := Client{
			id:          count,
			msgQueue:    make(chan Message, 10),
			incomeQueue: make(chan *string, 10),
			conn:        con,
			leftSignal:  make(chan bool, 1),
		}
		count += 1
		fmt.Println("A client joined: ", con.RemoteAddr().String())
		server.joinQueue <- &client
	}
}

func main() {
	server := Server{
		joinQueue: make(chan *Client, 10),
		clients:   make([]*Client, 0, 10),
	}
	go acceptConnection(&server, listenPort(":7777"))
	server.Run()
}
