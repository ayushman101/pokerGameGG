package p2p

// "fmt"

// "github.com/sirupsen/logrus"

type Message struct {
	ListenAddr string
	Payload    any
}

func NewMessage(from string, payload any) *Message {
	return &Message{
		Payload:    payload,
		ListenAddr: from,
	}
}

type MessagePeerList struct {
	Peers []string
}

type Handler interface {
	HandleMessage(*Message) error
}

type DefaultHandler struct{}
