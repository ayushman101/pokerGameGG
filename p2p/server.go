package p2p

import (
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

type Server struct {
	ServerConfig
	listner net.Listener

	peers   map[net.Addr]*Peer
	addpear chan *Peer
}

func NewServer(config ServerConfig) *Server {
	return &Server{
		ServerConfig: config,
		peers:        make(map[net.Addr]*Peer),

		addpear: make(chan *Peer),
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

		s.addpear <- peer

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
				fmt.Printf("The connection  %s has been closed\n", conn.RemoteAddr())
			}
			break
		}

		fmt.Println(string(buf[:n]))

	}

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
		case peer := <-s.addpear:
			s.peers[peer.conn.RemoteAddr()] = peer

			fmt.Printf("New player connected from %s", peer.conn.RemoteAddr())
		}
	}
}
