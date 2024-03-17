package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/sirupsen/logrus"
)

type GameVariant uint8

const (
	TexasHoldem GameVariant = iota
	Other
)

func (g GameVariant) String() string {
	switch g {
	case TexasHoldem:
		return "TEXASHOLDEM"
	case Other:
		return "OTHER"
	default:
		return "UNKNOWN"
	}
}

type ServerConfig struct {
	ListenAddr  string // its port number
	GameVariant GameVariant
}

type Server struct {
	ServerConfig
	Handler
	transport *tcpTransport
	peers     map[net.Addr]*Peer
	addpeer   chan *Peer
	delpeer   chan *Peer
	msgs      chan *Message
}

func NewServer(config ServerConfig) *Server {
	s := &Server{
		ServerConfig: config,
		peers:        make(map[net.Addr]*Peer),
		addpeer:      make(chan *Peer),
		delpeer:      make(chan *Peer),
		msgs:         make(chan *Message),
		Handler:      &DefaultHandler{},
	}

	tr := NewTcpTransport(s.ListenAddr)

	s.transport = tr

	tr.AddPeer = s.addpeer
	tr.DelPeer = s.delpeer

	return s
}

func (s *Server) Start() {

	go s.loop()

	fmt.Println("The server has started at ", s.ListenAddr)

	logrus.WithFields(logrus.Fields{
		"port":        s.transport.ListenAddr,
		"GameVariant": s.GameVariant,
	}).Info("Game has Started")

	if err := s.transport.ListenAndAccept(); err != nil {
		panic(err)
	}

}

func (s *Server) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return err
	}

	peer := &Peer{
		conn: conn,
	}

	s.addpeer <- peer

	return peer.SendHandshake(s.GameVariant)
}

func (s *Server) handShake(p *Peer) error {
	hs := &HandShake{}

	buf := make([]byte, 1024)

	if _, err := p.conn.Read(buf); err != nil {
		return err
	}

	b := bytes.NewBuffer(buf)

	decoder := gob.NewDecoder(b)

	fmt.Println("Decoder created")

	if err := decoder.Decode(hs); err != nil {
		if err == io.EOF {
			logrus.Info("Handshake: Reached end of input data stream")
		} else if err == io.ErrUnexpectedEOF {
			fmt.Println("Handshake : Unexpected EOF")
			fmt.Println(hs)
		} else {
			fmt.Println("Error in Decoding hs")
			return err

		}

	}

	if hs.GameVariant != s.GameVariant {
		logrus.Error(" Handshake failed")
		return fmt.Errorf("player gameVariant incompatible with server game variant")
	}

	logrus.WithFields(logrus.Fields{
		"peer":    p.conn.RemoteAddr(),
		"Version": hs.Version,
		"Variant": hs.GameVariant,
	}).Info(" Handshake Received")

	return nil
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addpeer:

			if err := s.handShake(peer); err != nil {
				logrus.Info("Handshake failed with player: ", err)
				continue
			}

			s.peers[peer.conn.RemoteAddr()] = peer

			fmt.Printf("New player connected from %s\n", peer.conn.RemoteAddr())

			go peer.ReadLoop(s.msgs)

		case peer := <-s.delpeer:
			delete(s.peers, peer.conn.RemoteAddr())
			fmt.Printf("Player Disconnectd | Address: %s\n", peer.conn.RemoteAddr())

		case msg := <-s.msgs:
			if err := s.Handler.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}
