package blackjack

import (
	"fmt"

	"github.com/SirNoob97/gophercises/deck"
)

// AI ...
type AI interface {
	Bet(shuffled bool) int
	Play(hand []deck.Card, dealer deck.Card) Move
	Result(hands [][]deck.Card, dealer []deck.Card)
}

type dealerAI struct{}

// Play ...
func (ai dealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	dealerScore := Score(hand...)
	if dealerScore <= 16 || (dealerScore == 17 && SoftScore(hand...)) {
		return MoveHit
	}
	return MoveStand
}

// Bet ...
func (ai dealerAI) Bet(shuffled bool) int {
	return 1
}

// Result print the result of the match
func (ai dealerAI) Result(hands [][]deck.Card, dealer []deck.Card) {
}

// HumanAI ...
func HumanAI() AI {
	return humanAI{}
}

type humanAI struct{}

// Play player input
func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	for {
		fmt.Println("Player", hand)
		fmt.Println("Dealer", dealer)
		fmt.Println("What will you do? (h)it, (s)tand, (d)double")

		var input string
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		case "d":
			return MoveDouble
		case "p":
			return MoveSplit
		default:
			fmt.Println("Invalid Option:", input)
		}
	}
}

// Bet ...
func (ai humanAI) Bet(shuffled bool) int {
	if shuffled {
		fmt.Println("The deck was just shuffled")
	}
	fmt.Println("What would you like to bet")
	var bet int
	fmt.Scanf("%d\n", &bet)
	return bet
}

// Result print the result of the match
func (ai humanAI) Result(hands [][]deck.Card, dealer []deck.Card) {
	fmt.Println("***FINAL HAND***")
	fmt.Println("Player:")
	for _, h := range hands {
		fmt.Println(" ", h)
	}
	fmt.Println("Dealer:", dealer)
}
