package main

import (
	"fmt"

	"github.com/ayushman101/pokerGameGG/deck"
)

func main() {
	d := deck.New()

	d = d.Shuffle(3)

	fmt.Println(d)

}
