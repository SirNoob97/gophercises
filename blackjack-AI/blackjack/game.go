package blackjack

import (
	"errors"

	"github.com/SirNoob97/gophercises/deck"
)

// Status
const (
	//stateBet state = iota
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

var (
	errBust = errors.New("Hand score exceeded 21")
)

// Move ...
type Move func(*Game) error

// State ..
type state uint8

// Options game options
type Options struct {
	Deck            int
	Hands           int
	BlackJackPayout float64
}

type hand struct {
	cards []deck.Card
	bet   int
}

// Game ...
type Game struct {
	nDecks          int
	nHands          int
	blackjackPayout float64

	deck  []deck.Card
	state state

	player    []hand
	handIndx  int
	playerBet int
	balance    int

	dealer   []deck.Card
	dealerAI AI
}

// New ...
func New(opts Options) Game {
	g := Game{
		state:    statePlayerTurn,
		dealerAI: dealerAI{},
	}
	if opts.Deck == 0 {
		opts.Deck = 3
	}
	if opts.Hands == 0 {
		opts.Hands = 100
	}
	if opts.BlackJackPayout == 0.0 {
		opts.BlackJackPayout = 1.5
	}
	g.nDecks = opts.Deck
	g.nHands = opts.Hands
	g.blackjackPayout = opts.BlackJackPayout
	return g
}

// Play ...
func (g *Game) Play(ai AI) int {
	g.deck = nil
	min := 52 * g.nDecks / 3
	shuffle := false

	for i := 0; i < g.nHands; i++ {
		if len(g.deck) < min {
			g.deck = deck.New(deck.Deck(g.nDecks), deck.Shuffle)
			shuffle = true
		}
		bet(g, ai, shuffle)
		deal(g)

		if BlackJack(g.dealer...) {
			finalRound(g, ai)
			continue
		}

		for g.state == statePlayerTurn {
			hand := make([]deck.Card, len(*g.currentHand()))
			copy(hand, *g.currentHand())
			move := ai.Play(hand, g.dealer[0])
			err := move(g)

			switch err {
			case errBust:
				MoveStand(g)
			case nil:

			default:
				panic(err)
			}
		}

		for g.state == stateDealerTurn {
			hand := make([]deck.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAI.Play(hand, g.dealer[0])
			move(g)
		}
		finalRound(g, ai)
	}
	return g.balance
}

func bet(g *Game, ai AI, shuffled bool) {
	bet := ai.Bet(shuffled)
	if bet < 100 {
		panic("bet must be at least 100")
	}
	g.playerBet = bet
}

func (g *Game) currentHand() *[]deck.Card {
	switch g.state {
	case statePlayerTurn:
		return &g.player[g.handIndx].cards
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("It isn't currently any player turn")
	}
}

func deal(g *Game) {
	playerHand := make([]deck.Card, 0, 5)
	g.handIndx = 0
	g.dealer = make([]deck.Card, 0, 5)

	var card deck.Card
	for i := 0; i < 2; i++ {
		card, g.deck = draw(g.deck)
		playerHand = append(playerHand, card)
		card, g.deck = draw(g.deck)
		g.dealer = append(g.dealer, card)
	}
	g.player = []hand{
		{
			cards: playerHand,
			bet:   g.playerBet,
		},
	}
	g.state = statePlayerTurn
}

// MoveHit keep playing
func MoveHit(g *Game) error {
	hand := g.currentHand()
	var card deck.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)

	if Score(*hand...) > 21 {
		return errBust
	}
	return nil
}

// MoveDouble double the bet
func MoveDouble(g *Game) error {
	if len(*g.currentHand()) != 2 {
		return errors.New("Can only double on a hand with 2 cards")
	}
	g.playerBet *= 2
	MoveHit(g)
	return MoveStand(g)
}

// MoveStand stop playing
func MoveStand(g *Game) error {
	if g.state == stateDealerTurn {
		g.state++
		return nil
	}

	if g.state == statePlayerTurn {
		g.handIndx++
		if g.handIndx >= len(g.player) {
			g.state++
		}
		return nil
	}
	return errors.New("Invalid state")
}

// MoveSplit split the hand
func MoveSplit(g *Game) error {
	cards := g.currentHand()
	if len(*cards) != 2 {
		return errors.New("You can only split with two cards in your hand")
	}
	if (*cards)[0].Rank != (*cards)[1].Rank {
		return errors.New("Both cards must have the same rank to split")
	}

	g.player = append(g.player, hand{
		cards: []deck.Card{(*cards)[1]},
		bet:   g.player[g.handIndx].bet,
	})
	g.player[g.handIndx].cards = (*cards)[:1]
	return nil
}

func finalRound(g *Game, ai AI) {
	dScore := Score(g.dealer...)
	dBlaBlackJack := BlackJack(g.dealer...)
	allHands := make([][]deck.Card, len(g.player))

	for i, hand := range g.player {
		cards := hand.cards
		allHands[i] = cards
		pScore, pBlackJack := Score(cards...), BlackJack(cards...)
		winnings := hand.bet

		switch {
		case pBlackJack && dBlaBlackJack:
			winnings = 0
		case dBlaBlackJack:
			winnings *= -1
		case pBlackJack:
			winnings = int(float64(winnings) * g.blackjackPayout)
		case pScore > 21:
			winnings *= -1
		case dScore > 21:
			// win
		case pScore > dScore:
			// win
		case dScore > pScore:
			winnings *= -1
		case pScore == dScore:
			winnings = 0
		}
		g.balance += winnings
	}

	ai.Result(allHands, g.dealer)
	g.player = nil
	g.dealer = nil
}

// BlackJack verfify blackjack
func BlackJack(hand ...deck.Card) bool {
	return len(hand) == 2 && Score(hand...) == 21
}

// SoftScore return true if the score of a hand is a soft score
func SoftScore(h ...deck.Card) bool {
	minScore := minScore(h...)
	score := Score(h...)
	return minScore != score
}

// Score get the score of the hand
func Score(h ...deck.Card) int {
	minScore := minScore(h...)
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

func minScore(h ...deck.Card) int {
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

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}
