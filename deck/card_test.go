package deck

import (
	"fmt"
	"math/rand"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Rank: Eight, Suit: Diamond})
	fmt.Println(Card{Rank: Queen, Suit: Spade})
	fmt.Println(Card{Rank: Two, Suit: Club})
	fmt.Println(Card{Suit: Joker})

	//Output:
	//Ace of Hearts
	//Eight of Diamonds
	//Queen of Spades
	//Two of Clubs
	//Joker
}

func TestNew(t *testing.T) {
	cards := New()

	if len(cards) != 13*4 {
		t.Error("Wrong Number of cards in a new deck.")
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	exp := Card{Suit: Spade, Rank: Ace}

	if cards[0] != exp {
		t.Error("Expected Ace of Spades as a first card, Received:", cards[0])
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))
	exp := Card{Suit: Spade, Rank: Ace}

	if cards[0] != exp {
		t.Error("Expected Ace of Spades as a first card, Received:", cards[0])
	}
}

func TestJokers(t *testing.T) {
	exp := 5
	cards := New(Jokers(exp))
	count := 0

	for _, c := range cards {
		if c.Suit == Joker {
			count++
		}
	}

	if count != exp {
		t.Errorf("Expected %d Jokers, Received: %d", exp, count)
	}
}

func TestFilter(t *testing.T) {
	filter := func(card Card) bool {
		return card.Rank == Two || card.Rank == Three
	}

	cards := New(Filter(filter))

	for _, c := range cards {
		if c.Rank == Two || c.Rank == Three {
			t.Error("Expected all twos and threes to be filteres out.")
		}
	}
}

func TestDeck(t *testing.T) {
	cards := New(Deck(3))
	exp := 13 * 4 * 3

	if len(cards) != exp {
		t.Errorf("Expected %d cards, Received %d cards.", exp, len(cards))
	}
}

func TestShuffle(t *testing.T) {
	// make shuffleRand deterministic
	// Fisrt call to shuffleRand.Perm(52) should be:
	// [40, 35 ...]
	shuffleRand = rand.New(rand.NewSource(0))

	orig := New()
	first := orig[40]
	second := orig[35]
	cards := New(Shuffle)

	if cards[0] != first {
		t.Errorf("Expected the first card to be %s, Received %s.", first, cards[0])
	}

	if cards[1] != second {
		t.Errorf("Expected the second card to be %s, Received %s.", second, cards[1])
	}
}
