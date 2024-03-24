package main

import (
	"log"
	"time"

	// "github.com/ayushman101/pokerGameGG/deck"
	"github.com/ayushman101/pokerGameGG/p2p"
)

func makeServerAndStart(addr string) *p2p.Server {

	cfg := p2p.ServerConfig{
		ListenAddr:  addr,
		GameVariant: p2p.TexasHoldem,
	}

	s := p2p.NewServer(cfg)

	go s.Start()

	return s

}

func main() {

	player1 := makeServerAndStart(":3000")

	player2 := makeServerAndStart(":4000")

	player3 := makeServerAndStart(":6000")

	player4 := makeServerAndStart(":5000")

	time.Sleep(1 * time.Second)

	if err := player1.Connect(player2.ListenAddr); err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	if err := player3.Connect(player2.ListenAddr); err != nil {
		log.Fatal(err)
	}

	if err := player4.Connect(player1.ListenAddr); err != nil {
		log.Fatal(err)
	}

	select {}

}
