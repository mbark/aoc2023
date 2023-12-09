package day4

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mbark/aoc2023/fns"
	"github.com/mbark/aoc2023/maths"
	"github.com/mbark/aoc2023/util"
)

const testInput = `
Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53
Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19
Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1
Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83
Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36
Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11
`

func Day4(input string) {
	if input == "" {
		input = testInput
	}

	in := util.ReadInput(input, "\n")
	var cards []Card
	for i, s := range in {
		split := strings.Split(s, ":")
		split = strings.Split(split[1], "|")

		w := strings.Split(strings.TrimSpace(split[0]), " ")
		w = fns.Filter(w, func(s string) bool { return s != "" })
		winning := fns.Map(w, func(num string) int { return util.ParseInt[int](num) })
		h := strings.Split(strings.TrimSpace(split[1]), " ")
		h = fns.Filter(h, func(s string) bool { return s != "" })
		have := fns.Map(h, func(num string) int { return util.ParseInt[int](num) })

		cards = append(cards, Card{
			Index:   i,
			Winning: winning,
			Have:    have,
		})
	}

	fmt.Printf("first: %d\n", first(cards))
	fmt.Printf("second: %d\n", second(cards))
}

func first(cards []Card) int {
	var sum int
	for _, card := range cards {
		count := card.WinningNumbers()
		if count > 0 {
			sum += maths.PowInt(2, count-1)
		}
	}

	return sum
}

func second(cards []Card) int {
	stack := cards
	var count int
	for len(stack) > 0 {
		count += 1
		next := stack[0]
		stack = stack[1:]

		won := next.WinningNumbers()
		for i := 1; i <= won; i++ {
			stack = append(stack, cards[next.Index+i])
		}
	}

	return count
}

type Card struct {
	Index   int
	Winning []int
	Have    []int
}

func (c Card) String() string {
	return fmt.Sprintf("Card %d: %s | %s", c.Index,
		strings.Join(fns.Map(c.Winning, func(i int) string { return strconv.Itoa(i) }), " "),
		strings.Join(fns.Map(c.Have, func(i int) string { return strconv.Itoa(i) }), " "),
	)
}

func (c Card) WinningNumbers() int {
	var count int
	for _, have := range c.Have {
		for _, win := range c.Winning {
			if have == win {
				count += 1
			}
		}
	}

	return count
}
