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

func (p *Peer) ReadLoop(msgc chan *Message) {
	buf := make([]byte, 1024)

	for {
		n, err := p.conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				fmt.Printf("Connection terminated from client %s\n", p.conn.RemoteAddr())
			}
			break
		}

		msgc <- &Message{
			ListenAddr: p.conn.RemoteAddr(),
			Payload:    bytes.NewReader(buf[:n]),
		}

	}

	//TODO: unregister this peer
	p.conn.Close()

}

type Message struct {
	ListenAddr net.Addr
	Payload    io.Reader
}

type tcpTransport struct {
	ListenAddr string
	listener   net.Listener
	AddPeer    chan *Peer
	DelPeer    chan *Peer
}

func (t *tcpTransport) ListenAndAccept() error {

	ln, err := net.Listen("tcp", t.ListenAddr)

	if err != nil {
		return err
	}

	t.listener = ln

	for {
		conn, err := t.listener.Accept()

		if err != nil {
			panic(err)
		}

		peer := &Peer{
			conn: conn,
		}

		t.AddPeer <- peer

	}

}

func NewTcpTransport(addr string) *tcpTransport {
	return &tcpTransport{
		ListenAddr: addr,
	}
}
