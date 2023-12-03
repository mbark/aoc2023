package day3

import (
	"fmt"
	"regexp"

	"github.com/mbark/aoc2023/maps"
	"github.com/mbark/aoc2023/util"
)

const testInput = `
467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..
`

func Day3(input string) {
	if input == "" {
		input = testInput
	}

	m := maps.New(input, func(x, y int, b byte) C { return C(b) })
	fmt.Printf("first: %d\n", first(m))
	fmt.Printf("second: %d\n", second(m))
}

func first(m maps.Map[C]) int {
	var sum int

	var number []C
	var isAdjacent bool
	for _, coord := range m.Coordinates() {
		at := m.At(coord)
		switch {
		case !at.IsNumber():
			if len(number) > 0 && isAdjacent {
				asInt := util.ParseInt[int](string(number))
				sum += asInt
			}

			number = nil
			isAdjacent = false

		case at.IsNumber():
			number = append(number, at)
			for _, scoord := range m.Surrounding(coord) {
				if m.At(scoord).IsSymbol() {
					isAdjacent = true
				}
			}
		}
	}

	return sum
}

func second(m maps.Map[C]) int {
	var number []C
	adjacentTo := make(map[maps.Coordinate]bool)

	var gears []maps.Coordinate
	numbers := make(map[maps.Coordinate][]int)

	for _, coord := range m.Coordinates() {
		at := m.At(coord)

		if at == '*' {
			gears = append(gears, coord)
		}

		switch {
		case !at.IsNumber():
			if len(number) > 0 && len(adjacentTo) > 0 {
				asInt := util.ParseInt[int](string(number))
				for c := range adjacentTo {
					numbers[c] = append(numbers[c], asInt)
				}
			}

			number = nil
			adjacentTo = make(map[maps.Coordinate]bool)

		case at.IsNumber():
			number = append(number, at)
			for _, scoord := range m.Surrounding(coord) {
				if m.At(scoord).IsSymbol() {
					adjacentTo[scoord] = true
				}
			}
		}
	}

	var sum int
	for _, gear := range gears {
		nrs := numbers[gear]
		if len(nrs) == 2 {
			ratio := 1
			for _, nr := range nrs {
				ratio *= nr
			}
			sum += ratio
		}
	}

	return sum
}

type C byte

var digit = regexp.MustCompile(`[0-9]`)

func (c C) IsNumber() bool {
	return digit.Match([]byte{byte(c)})
}

func (c C) IsEmpty() bool {
	return c == '.'
}

func (c C) IsSymbol() bool {
	return !(c.IsEmpty() || c.IsNumber())
}
