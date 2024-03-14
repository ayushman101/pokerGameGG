package main

import (
	"fmt"

	"github.com/ayushman101/pokerGameGG/deck"
	"github.com/ayushman101/pokerGameGG/p2p"
)

func main() {
	d := deck.New()

	d = d.Shuffle(3)

	fmt.Println(d)

	cfg := p2p.ServerConfig{
		ListenAddr: ":3000",
	}

	s := p2p.NewServer(cfg)

	s.Start()

}
