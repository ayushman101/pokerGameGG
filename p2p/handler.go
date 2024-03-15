package p2p

import (
	"fmt"
	"io"
)

type Handler interface {
	HandleMessage(*Message) error
}

type DefaultHandler struct{}

func (h *DefaultHandler) HandleMessage(msg *Message) error {

	b, err := io.ReadAll(msg.Payload)

	if err != nil {
		return err
	}

	fmt.Printf("handling Message from %s: %s\n", msg.ListenAddr, string(b))
	return nil
}
