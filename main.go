package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ayushman101/pokerGameGG/deck"
	"github.com/ayushman101/pokerGameGG/p2p"
)

func main() {
	d := deck.New()

	d = d.Shuffle(3)

	fmt.Println(d)

	cfg := p2p.ServerConfig{
		ListenAddr:  ":3000",
		GameVariant: p2p.TexasHoldem,
	}

	s := p2p.NewServer(cfg)

	go s.Start()

	time.Sleep(1 * time.Second)

	remoteCfg := p2p.ServerConfig{
		ListenAddr:  ":4000",
		GameVariant: p2p.TexasHoldem,
	}

	remoteServer := p2p.NewServer(remoteCfg)

	go remoteServer.Start()

	if err := remoteServer.Connect(":3000"); err != nil {
		log.Fatal(err)
	}

	select {}

}
