package main

import (
	"fmt"
	"log"

	"github.com/ayushman101/pokerGameGG/deck"
)

func main() {
	c, err := deck.NewCard(deck.Hearts, 1)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hello Poker World")
	fmt.Println(c)

}
