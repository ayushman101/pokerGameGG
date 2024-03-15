package p2p

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

type Peer struct {
	conn net.Conn
}

func (p *Peer) Send(b []byte) error {
	_, err := p.conn.Write(b)

	return err
}

type ServerConfig struct {
	ListenAddr string // its port number
}

type Message struct {
	ListenAddr net.Addr
	Payload    io.Reader
}

type Server struct {
	ServerConfig
	listner net.Listener
	Handler
	peers   map[net.Addr]*Peer
	addpeer chan *Peer
	delpeer chan *Peer
	msgs    chan *Message
}

func NewServer(config ServerConfig) *Server {
	return &Server{
		ServerConfig: config,
		peers:        make(map[net.Addr]*Peer),
		addpeer:      make(chan *Peer),
		delpeer:      make(chan *Peer),
		msgs:         make(chan *Message),
		Handler:      &DefaultHandler{},
	}
}

func (s *Server) Start() {

	go s.loop()

	if err := s.Listen(); err != nil {
		panic(err)
	}

	fmt.Println("The server has started at ", s.ListenAddr)

	s.acceptLoop()
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listner.Accept()

		if err != nil {
			panic(err)
		}

		peer := &Peer{
			conn: conn,
		}

		s.addpeer <- peer

		msg := "Welcome to Poker GameGG Version 1.01\n"

		peer.Send([]byte(msg))

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				fmt.Printf("Connection terminated from client %s\n", conn.RemoteAddr())
			}
			break
		}

		s.msgs <- &Message{
			ListenAddr: conn.RemoteAddr(),
			Payload:    bytes.NewReader(buf[:n]),
		}

	}

	s.delpeer <- s.peers[conn.RemoteAddr()]

}

func (s *Server) Listen() error {

	ln, err := net.Listen("tcp", s.ListenAddr)

	if err != nil {
		return err
	}

	s.listner = ln
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addpeer:
			s.peers[peer.conn.RemoteAddr()] = peer

			fmt.Printf("New player connected from %s\n", peer.conn.RemoteAddr())

		case peer := <-s.delpeer:
			delete(s.peers, peer.conn.LocalAddr())
			fmt.Printf("Player Disconnectd | Address: %s\n", peer.conn.RemoteAddr())

		case msg := <-s.msgs:
			if err := s.Handler.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}
