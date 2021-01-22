package main

import (
	"fmt"
	"strings"

	"github.com/SirNoob97/gophercises/deck"
)

// Hand player and dealer hand
type Hand []deck.Card

func (h Hand) String() string {
	str := make([]string, len(h))
	for i := range h {
		str[i] = h[i].String()
	}
	return strings.Join(str, ", ")
}

// DealerString dealer needs to hide the second card
func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

// Score get the score of the hand
func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}

	for _, c := range h {
		if c.Rank == deck.Ace {
			//ace is currently worth 1, and we are changing to be worth 11
			//11 - 1 = 10
			return minScore + 10
		}
	}
	return minScore
}

// MinScore get the minimun score of the hand
func (h Hand) MinScore() int {
	score := 0
	for _, c := range h {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Shuffle ..
func Shuffle(gs GameState) GameState {
	ret := clone(gs)
	ret.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return ret
}

// Deal ...
func Deal(gs GameState) GameState {
	ret := clone(gs)
	ret.Player = make(Hand, 0, 5)
	ret.Dealer = make(Hand, 0, 5)

	var card deck.Card
	for i := 0; i < 2; i++ {
		card, ret.Deck = draw(ret.Deck)
		ret.Player = append(ret.Player, card)
		card, ret.Deck = draw(ret.Deck)
		ret.Dealer = append(ret.Dealer, card)
	}
	ret.State = StatePlayerTurn
	return ret
}

// Hit ...
func Hit(gs GameState) GameState {
	ret := clone(gs)
	hand := ret.CurrentHand()
	var card deck.Card
	card, ret.Deck = draw(ret.Deck)
	*hand = append(*hand, card)

	if hand.Score() > 21 {
		return Stand(ret)
	}

	return ret
}

// Stand ...
func Stand(gs GameState) GameState {
	ret := clone(gs)
	ret.State++
	return ret
}

// FinalHand ...
func FinalHand(gs GameState) GameState {
	ret := clone(gs)
	pScore, dScore := ret.Player.Score(), ret.Dealer.Score()
	fmt.Println("***FINAL HAND***")
	fmt.Println("Player", ret.Player, "\nScore:", pScore)
	fmt.Println("Dealer", ret.Dealer, "\nScore:", dScore)

	switch {
	case pScore > 21:
		fmt.Println("You Busted")
	case dScore > 21:
		fmt.Println("Dealer Busted")
	case pScore > dScore:
		fmt.Println("You Win!")
	case dScore > pScore:
		fmt.Println("You Lose!")
	case pScore == dScore:
		fmt.Println("DRAW")
	}

	ret.Player = nil
	ret.Dealer = nil

	return ret
}

func main() {
	var gs GameState
	gs = Shuffle(gs)

	for i := 0; i < 10; i++ {
		gs = Deal(gs)

		var input string
		for gs.State == StatePlayerTurn {
			fmt.Println("Player", gs.Player)
			fmt.Println("Dealer", gs.Dealer.DealerString())
			fmt.Println("What will you do? (h)it, (s)tand")
			fmt.Scanf("%s\n", &input)

			switch input {
			case "h":
				gs = Hit(gs)
			case "s":
				gs = Stand(gs)
			default:
				fmt.Println("Invalid Option:", input)
			}
		}

		for gs.State == StateDealerTurn {
			if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
				gs = Hit(gs)
			} else {
				gs = Stand(gs)
			}
		}

		gs = FinalHand(gs)
	}
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

// State ..
type State uint8

// Status
const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

// GameState ...
type GameState struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

// CurrentHand get the current player
func (g *GameState) CurrentHand() *Hand {
	switch g.State {
	case StatePlayerTurn:
		return &g.Player
	case StateDealerTurn:
		return &g.Dealer
	default:
		panic("It isn't currently any player turn")
	}
}

func clone(gs GameState) GameState {
	ret := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		State:  gs.State,
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
	}
	copy(ret.Deck, gs.Deck)
	copy(ret.Player, gs.Player)
	copy(ret.Dealer, gs.Dealer)

	return ret
}
