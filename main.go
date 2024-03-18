package main

import (
	"log"
	"time"

	// "github.com/ayushman101/pokerGameGG/deck"
	"github.com/ayushman101/pokerGameGG/p2p"
)

func main() {

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

	time.Sleep(1 * time.Second)

	OtherCfg := p2p.ServerConfig{
		ListenAddr:  ":4001",
		GameVariant: p2p.TexasHoldem,
	}

	OtherServer := p2p.NewServer(OtherCfg)

	go OtherServer.Start()

	if err := OtherServer.Connect(":3000"); err != nil {
		log.Fatal(err)
	}

	select {}

}
