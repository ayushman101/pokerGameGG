package deck

import (
	"errors"
	"fmt"
	"strconv"
)

type Suit int

const (
	Spades Suit = iota
	Hearts
	Clubs
	Diamonds
)

func (s Suit) String() string {

	switch s {
	case Spades:
		return "SPADES"
	case Hearts:
		return "HEARTS"
	case Clubs:
		return "CLUBS"
	case Diamonds:
		return "DIAMONDS"
	default:
		panic("Deck Pkg: suit String: Invalid suit type")
	}
}

func (s Suit) suitToUnicode() string {
	switch s {
	case Spades:
		return "♤"
	case Hearts:
		return "♡"
	case Clubs:
		return "♧"
	case Diamonds:
		return "♢"
	default:
		panic("deck pkg: suitToUnicode error: Invalid suit type")
	}

}

type Card struct {
	suit  Suit
	value int // min 1 and max 13
}

func (c Card) String() string {
	var s string

	switch c.value {
	case 1:
		s = "Ace"
	case 11:
		s = "Jack"
	case 12:
		s = "Queen"
	case 13:
		s = "King"
	default:
		s = strconv.Itoa(c.value)
	}
	return fmt.Sprintf("%s of %s %s", s, c.suit, c.suit.suitToUnicode())
}

func NewCard(s Suit, val int) (Card, error) {

	if val < 1 || val > 13 {
		return Card{}, errors.New("invalid card value")
	}

	if s != Spades && s != Clubs && s != Diamonds && s != Hearts {
		return Card{}, errors.New("invalid suit type")
	}

	return Card{
		suit:  s,
		value: val,
	}, nil
}
