package p2p

import (
	"fmt"
)

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

type Peers []string

type Handler interface {
	HandleMessage(*Message) error
}

type DefaultHandler struct{}

func (h *DefaultHandler) HandleMessage(msg *Message) error {

	switch v := msg.Payload.(type) {
	case []string:
		fmt.Printf("handling MessagePeerList from %s: %+v\n", msg.ListenAddr, v)
	default:
		fmt.Printf("Other Message: %s\n", v)

	}
	return nil
}
