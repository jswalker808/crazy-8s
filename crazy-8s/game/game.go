package game

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const MaxPlayers = 4

type Status string

const (
	PENDING Status = "pending"
	IN_PROGRESS Status = "in_progress"
)

func ParseStatus(str string) (Status, error) {
	switch strings.ToLower(str) {
    case "pending":
        return PENDING, nil
    case "in_progress":
        return IN_PROGRESS, nil
    default:
        return "", fmt.Errorf("unable to parse \"%v\" to a status", str)
    }
}

type Game struct {
	id string
	ownerId string
	maxPoints int
	players []*Player
	deck Deck
	discardPile []*Card
	status Status
	currentTurn string
}

func NewGame(ownerId string, ownerName string) *Game {
	owner := NewPlayer(ownerId, ownerName)
	players := make([]*Player, 0, MaxPlayers)
	players = append(players, owner)

	return &Game{
		id: uuid.NewString(),
		ownerId: owner.id,
		maxPoints: 100,
		players: players,
		deck: NewStandardDeck(),
		discardPile: make([]*Card, 0),
		status: PENDING,
		currentTurn: "",
	}
}

func NewGameFromExisting(
	id string,
	ownerId string,
	maxPoints int,
	players []*Player,
	deck []*Card ,
	discardPile []*Card,
	status Status,
	currentTurn string,
) *Game {
	return &Game{
		id: id,
		ownerId: ownerId,
		maxPoints: maxPoints,
		players: players,
		deck: NewStandardDeckFromExisting(deck),
		discardPile: discardPile,
		status: status,
		currentTurn: currentTurn,
	}
}

func (game *Game) String() string {
	return fmt.Sprintf("Game[id=%v, ownerId=%v, maxPoints=%v, players=%v]", game.id, game.ownerId, game.maxPoints, len(game.players))
}

func (game *Game) GetId() string {
	return game.id
}

func (game *Game) GetOwnerId() string {
	return game.ownerId
}

func (game *Game) GetMaxPoints() int {
	return game.maxPoints
}

func (game *Game) GetPlayers() []*Player {
	return game.players
} 

func (game *Game) GetDeck() Deck {
	return game.deck
}

func (game *Game) GetDiscardPile() []*Card {
	return game.discardPile
}

func (game *Game) GetStatus() Status {
	return game.status
}

func (game *Game) GetCurrentTurn() string {
	return game.currentTurn
}

func (game *Game) GetOwner() *Player {
	for _, player := range game.players {
		if player.id == game.ownerId {
			return player
		}
	}
	panic("Invalid state: Unable to find owner for game")
}

func (game *Game) GetPlayer(id string) *Player {
	for _, player := range game.players {
		if player.id == id {
			return player
		}
	}
	return nil
}

func (game *Game) GetOpponents(id string) []*Player {
	opponents := make([]*Player, 0)
	for _, player := range game.players {
		if player.id != id {
			opponents = append(opponents, player)
		}
	}
	return opponents
}

func (game *Game) AddPlayer(player *Player) error {
	if (len(game.players) == MaxPlayers) {
		return fmt.Errorf("reached max number of players: %v", MaxPlayers)
	}
	game.players = append(game.players, player)
	return nil
}

func (game *Game) RemovePlayer(id string) {
	for i, player := range game.players {
		if player.id == id {
			game.deck.ReturnCards(player.hand)
			game.players = append(game.players[:i], game.players[i+1:]...)
		}
	}
}