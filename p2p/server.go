package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	ListenAddr  string // its port number
	GameVariant GameVariant
	GameStatus  GameStatus
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
		addpeer:      make(chan *Peer, 100),
		delpeer:      make(chan *Peer),
		msgs:         make(chan *Message),
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
		"GameStatus":  s.GameStatus,
	}).Info("Game has Started")

	if err := s.transport.ListenAndAccept(); err != nil {
		panic(err)
	}

}

func (s *Server) SendHandshake(p *Peer, Variant GameVariant, gameStatus GameStatus) error {

	defer fmt.Println("Handshake sent successfully by player : ", s.ListenAddr)

	hs := &HandShake{
		GameVariant: Variant,
		GameStatus:  gameStatus,
	}

	buf := new(bytes.Buffer)

	if err := gob.NewEncoder(buf).Encode(hs); err != nil {
		fmt.Println("There is an error in encoding")
		return err
	}

	return p.Send(buf.Bytes())

}

func (s *Server) Connect(addr string) error {

	fmt.Printf("%s is connecting to %s\n", s.ListenAddr, addr)

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return err
	}

	fmt.Println("addr: ", addr)

	peer := &Peer{
		conn:     conn,
		outbound: true,
	}

	s.addpeer <- peer
	fmt.Printf("%s connected with %s\n", s.ListenAddr, peer.conn.RemoteAddr().String())

	err = s.SendHandshake(peer, s.GameVariant, s.GameStatus)

	return err
}

func (s *Server) SendPeerList(p *Peer) error {

	peerList := make([]string, len(s.peers))

	it := 0
	for addr := range s.peers {
		peerList[it] = addr.String()
		it++
	}

	buf := new(bytes.Buffer)

	msg := NewMessage(s.ListenAddr, peerList)

	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	p.Send(buf.Bytes())

	return nil

}

func (s *Server) handShake(p *Peer) error {
	hs := &HandShake{}

	buf := make([]byte, 1024)

	fmt.Printf("Handling handshake in %s from %s\n", s.ListenAddr, p.conn.RemoteAddr())

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

	if hs.GameStatus != s.GameStatus {
		logrus.Error(" Handshake failed")
		return fmt.Errorf("player gameStatus incompatible with server gameStatus")
	}

	logrus.WithFields(logrus.Fields{
		"server":     s.ListenAddr,
		"peer":       p.conn.RemoteAddr(),
		"Version":    hs.Version,
		"Variant":    hs.GameVariant,
		"GameStatus": hs.GameStatus,
	}).Info(" Handshake Received")

	return nil
}

func (s *Server) handlePeerList(l *MessagePeerList) error {

	fmt.Printf("Handling Peer List in %s : PeerList: %v\n", s.ListenAddr, l.Peers)

	for _, str := range l.Peers {
		// str = ":4000"
		if err := s.Connect(str); err != nil {
			return fmt.Errorf("connection with peerlist element failed: %w", err)
		}

		fmt.Printf("Inside handlePeerList loop of %s\n", s.ListenAddr)
	}

	fmt.Printf("HandlePeerList ended: %s\n", s.ListenAddr)
	return nil
}

func (s *Server) HandleMessage(msg *Message) error {

	logrus.WithFields(logrus.Fields{
		"msg from: ": msg.ListenAddr,
	}).Info("received message")

	switch v := msg.Payload.(type) {
	case []string:
		fmt.Printf("handling MessagePeerList from %s: %+v\n", msg.ListenAddr, v)
		err := s.handlePeerList(&MessagePeerList{Peers: v})
		logrus.Errorf(" Handle Peer List failed: %+v", err)
	default:
		fmt.Printf("Other Message: %s\n", v)

	}
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addpeer:

			fmt.Printf("%s server adding peer with addr: %s\n", s.ListenAddr, peer.conn.RemoteAddr())

			if err := s.handShake(peer); err != nil {
				logrus.Error("Handshake failed with player: ", err)
				peer.conn.Close()
				continue
			}

			if !peer.outbound {
				if err := s.SendHandshake(peer, s.GameVariant, s.GameStatus); err != nil {
					logrus.Errorf("%s: handshake failed to be sent :%s", s.ListenAddr, err)
					peer.conn.Close()
					continue
				}

				fmt.Printf("%s has sent peerlist to %s\n", s.ListenAddr, peer.conn.RemoteAddr())

				if err := s.SendPeerList(peer); err != nil {
					logrus.Errorf("%s: peerList failed to be sent :%s", s.ListenAddr, err)
					peer.conn.Close()
					continue
				}
			}

			s.peers[peer.conn.RemoteAddr()] = peer

			fmt.Printf("New player connected from %s\n", peer.conn.RemoteAddr())

			go peer.ReadLoop(s.msgs)

		case peer := <-s.delpeer:
			delete(s.peers, peer.conn.RemoteAddr())
			fmt.Printf("Player Disconnectd | Address: %s\n", peer.conn.RemoteAddr())

		case msg := <-s.msgs:
			if err := s.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}
