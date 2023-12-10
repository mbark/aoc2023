package day7

import (
	"fmt"
	"slices"
	"strings"

	"github.com/mbark/aoc2023/fns"
	"github.com/mbark/aoc2023/util"
)

const testInput = `
32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483
`

func Day7(input string) {
	if input == "" {
		input = testInput
	}

	var hands []Hand
	for i, s := range util.ReadInput(input, "\n") {
		split := strings.Split(s, " ")
		bid := util.ParseInt[int](split[1])

		cards := make([]Card, 5)
		for i, c := range split[0] {
			cards[i] = Card(c)
		}

		hands = append(hands, Hand{Index: i, Cards: cards, Bid: bid})
	}

	fmt.Printf("first: %d\n", first(hands))
	fmt.Printf("second: %d\n", second(hands))
}

func first(hands []Hand) int {
	handTypes := make(map[int]int)
	for _, hand := range hands {
		counts := make(map[Card]int)
		for _, card := range hand.Cards {
			counts[card] += 1
		}

		values := fns.Values(counts)
		slices.SortFunc(values, func(a, b int) int { return b - a })

		handTypes[hand.Index] = getHandType(hand, values)
	}

	slices.SortFunc(hands, func(a, b Hand) int {
		d := handTypes[b.Index] - handTypes[a.Index]
		if d != 0 {
			return d
		}

		for i := 0; i < len(b.Cards); i++ {
			d := b.Cards[i].Value() - a.Cards[i].Value()
			if d != 0 {
				return d
			}
		}

		fmt.Println("same hand?", a, b)
		return 0
	})

	slices.Reverse(hands)
	var sum int
	for i, hand := range hands {
		sum += (i + 1) * hand.Bid
	}

	return sum
}

func second(hands []Hand) int {
	handTypes := make(map[int]int)
	for _, hand := range hands {
		counts := make(map[Card]int)
		var jokers int
		for _, card := range hand.Cards {
			if card == 'J' {
				jokers += 1
				continue
			}

			counts[card] += 1
		}

		values := fns.Values(counts)
		slices.SortFunc(values, func(a, b int) int { return b - a })

		handType := getHandType(hand, values)
		switch handType {
		case FiveOfAKind:
		case FourOfAKind:
			switch jokers {
			case 1:
				handType = FiveOfAKind
			}

		case FullHouse:
		case ThreeOfAKind:
			switch jokers {
			case 1:
				handType = FourOfAKind
			case 2:
				handType = FiveOfAKind
			}

		case TwoPair:
			switch jokers {
			case 1:
				handType = FullHouse
			}

		case OnePair:
			switch jokers {
			case 1:
				handType = ThreeOfAKind
			case 2:
				handType = FourOfAKind
			case 3:
				handType = FiveOfAKind
			}

		case HighCard:
			switch jokers {
			case 1:
				handType = OnePair
			case 2:
				handType = ThreeOfAKind
			case 3:
				handType = FourOfAKind
			case 4:
				handType = FiveOfAKind
			}
		}

		handTypes[hand.Index] = handType
	}

	slices.SortFunc(hands, func(a, b Hand) int {
		d := handTypes[b.Index] - handTypes[a.Index]
		if d != 0 {
			return d
		}

		for i := 0; i < len(b.Cards); i++ {
			d := b.Cards[i].ValueJoker() - a.Cards[i].ValueJoker()
			if d != 0 {
				return d
			}
		}

		return 0
	})

	slices.Reverse(hands)
	for _, h := range hands {
		fmt.Printf("%s %d\n", h, handTypes[h.Index])
	}
	var sum int
	for i, hand := range hands {
		sum += (i + 1) * hand.Bid
	}

	return sum
}

func getHandType(hand Hand, values []int) int {
	if len(values) == 0 {
		return FiveOfAKind
	}

	switch {
	case values[0] == 5:
		return FiveOfAKind
	case values[0] == 4:
		return FourOfAKind
	case len(values) > 1 && values[0] == 3 && values[1] == 2:
		return FullHouse
	case values[0] == 3:
		return ThreeOfAKind
	case len(values) > 1 && values[0] == 2 && values[1] == 2:
		return TwoPair
	case values[0] == 2:
		return OnePair
	case values[0] == 1:
		return HighCard
	default:
		panic(fmt.Sprintf("can't determine type for hand %s: %+v", hand.Cards, values))
	}
}

const (
	HighCard int = iota
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

type Hand struct {
	Index int
	Cards []Card
	Bid   int
}

func (h Hand) String() string {
	return fmt.Sprintf("%s %d",
		strings.Join(fns.Map(h.Cards, func(c Card) string {
			return c.String()
		}), ""),
		h.Bid,
	)
}

type Card byte

func (c Card) String() string {
	return string([]byte{byte(c)})
}

func (c Card) Value() int {
	switch c {
	case '2', '3', '4', '5', '6', '7', '8', '9':
		return util.ParseInt[int](string([]byte{byte(c)}))

	case 'T':
		return 10
	case 'J':
		return 11
	case 'Q':
		return 12
	case 'K':
		return 13
	case 'A':
		return 14

	default:
		panic(fmt.Sprintf("unknown card %b", c))
	}
}

func (c Card) ValueJoker() int {
	switch c {
	case '2', '3', '4', '5', '6', '7', '8', '9':
		return util.ParseInt[int](string([]byte{byte(c)}))

	case 'T':
		return 10
	case 'J':
		return 1
	case 'Q':
		return 12
	case 'K':
		return 13
	case 'A':
		return 14

	default:
		panic(fmt.Sprintf("unknown card %b", c))
	}
}
