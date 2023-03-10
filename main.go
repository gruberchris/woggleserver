package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddress  string
	listener       net.Listener
	quitChannel    chan struct{}
	messageChannel chan Message

	clientsMu sync.RWMutex
	clients   map[string]net.Conn
}

func NewServer(listenAddress string) *Server {
	return &Server{
		listenAddress:  listenAddress,
		quitChannel:    make(chan struct{}),
		messageChannel: make(chan Message, 10),
		clients:        make(map[string]net.Conn),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.listenAddress)

	if err != nil {
		return err
	}

	defer listener.Close()

	s.listener = listener

	go s.acceptConnections()

	<-s.quitChannel

	close(s.messageChannel)

	return nil
}

func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()

		s.clientsMu.Lock()
		s.clients[conn.RemoteAddr().String()] = conn
		s.clientsMu.Unlock()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		fmt.Println("Accepted connection from: ", conn.RemoteAddr().String())

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		s.clientsMu.Lock()
		delete(s.clients, conn.RemoteAddr().String())
		s.clientsMu.Unlock()

		_ = conn.Close()
		fmt.Printf("Connection from %s closed.\n", conn.RemoteAddr().String())
	}()

	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			// TODO: This is not a graceful way to exit the thread when the client disconnects
			//fmt.Println("Error reading from connection: ", err.Error())
			return
		}

		s.messageChannel <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}

		// Reverse the message
		reverseBuf := make([]byte, n+1)

		for i := 0; i < n; i++ {
			reverseBuf[i] = buf[n-i-1]
		}

		reverseBuf[n] = '\n'

		conn.Write(reverseBuf)
	}
}

func main() {
	server := NewServer(":3000")

	go func() {
		/*
			for {
				msg := <-server.messageChannel
				fmt.Printf("[%s]: %s", msg.from, string(msg.payload))
			}

		*/

		for msg := range server.messageChannel {
			fmt.Printf("[%s]: %s", msg.from, string(msg.payload))
		}
	}()

	log.Fatal(server.Start())
}
