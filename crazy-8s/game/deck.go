package game

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

const deckSize = 52

type Deck interface {
	Draw() (*Card, error)
	Shuffle()
}

type StandardDeck struct {
	cards []*Card
}

func NewStandardDeck() *StandardDeck {
	cards := make([]*Card, 0, deckSize)

	for _, color := range Colors {
		for i := 1; i < 14; i++ {
			cards = append(cards, NewCard(color, i))
		}
	}

	deck := &StandardDeck{
		cards: cards,
	}

	deck.Shuffle()

	return deck
}

func (deck *StandardDeck) Draw() (*Card, error) {
	if len(deck.cards) == 0 {
		return nil, errors.New("deck is empty, unable to draw")
	}

	lastIdx := len(deck.cards) - 1
	top := deck.cards[lastIdx:][0]
	deck.cards = deck.cards[:len(deck.cards) - 1]

	return top, nil
}

func (deck *StandardDeck) Shuffle() {
	rand.Shuffle(len(deck.cards), func(i, j int) { deck.cards[i], deck.cards[j] = deck.cards[j], deck.cards[i]})
}

func (deck *StandardDeck) String() string {
	var sb strings.Builder
	for _, card := range deck.cards {
		sb.WriteString(fmt.Sprint(card))
		sb.WriteString(fmt.Sprint(","))
	}
	return fmt.Sprintf("StandardDeck[cards=%v]", sb.String())
}

