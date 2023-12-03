package day2

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2023/util"
)

const testInput = `
Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green
`

func Day2(input string) {
	if input == "" {
		input = testInput
	}

	in := util.ReadInput(input, "\n")
	var games []Game
	for _, s := range in {
		split := strings.Split(s, ":")

		id := util.ParseInt[int](strings.Split(split[0], " ")[1])
		rs := strings.Split(split[1], ";")

		game := Game{id: id}
		for _, r := range rs {
			cs := strings.Split(r, ", ")

			colors := make(map[string]int)
			for _, c := range cs {
				c = strings.TrimSpace(c)
				sp := strings.Split(c, " ")
				colors[sp[1]] = util.ParseInt[int](sp[0])
			}

			game.reveals = append(game.reveals, colors)
		}

		games = append(games, game)
	}

	fmt.Printf("first: %d\n", first(games))
	fmt.Printf("second: %d\n", second(games))
}

func first(games []Game) int {
	noMoreThan := map[string]int{"red": 12, "green": 13, "blue": 14}
	var sum int
	for _, game := range games {
		isPossible := true
		for color, count := range noMoreThan {
			for _, reveal := range game.reveals {
				if reveal[color] > count {
					isPossible = false
					break
				}
			}

			if !isPossible {
				break
			}
		}

		if isPossible {
			sum += game.id
		}
	}

	return sum
}

func second(games []Game) int {
	var sum int
	for _, game := range games {
		fewest := map[string]int{"red": 0, "green": 0, "blue": 0}
		for _, reveal := range game.reveals {
			for color, count := range fewest {
				if c := reveal[color]; c > count {
					fewest[color] = reveal[color]
				}
			}
		}

		gs := 1
		for _, c := range fewest {
			gs *= c
		}

		sum += gs
	}

	return sum
}

type Reveal map[string]int

func (r Reveal) String() string {
	var s []string
	for color, count := range r {
		s = append(s, fmt.Sprintf("%d %s", count, color))
	}

	return strings.Join(s, ", ")
}

type Game struct {
	id      int
	reveals []Reveal
}

func (g Game) String() string {
	var rs []string
	for _, r := range g.reveals {
		rs = append(rs, r.String())
	}

	return fmt.Sprintf("Game %d: %s", g.id, strings.Join(rs, "; "))
}
