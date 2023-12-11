package day11

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2023/maps"
	"github.com/mbark/aoc2023/maths"
)

const testInput = `
...#......
.......#..
#.........
..........
......#...
.#........
.........#
..........
.......#..
#...#.....
`

func Day11(input string) {
	if input == "" {
		input = testInput
	}
	m := maps.New(input, func(_, _ int, b byte) Space { return Space(b) })

	fmt.Printf("first: %d\n", first(m))
	fmt.Printf("second: %d\n", second(m, 1000000))
}

func first(m maps.Map[Space]) int {
	emptyRows := make(map[int]bool)
	for y := 0; y < m.Rows; y++ {
		isEmpty := true
		for x := 0; x < m.Columns; x++ {
			if m.At(maps.Coordinate{X: x, Y: y}) == Galaxy {
				isEmpty = false
				break
			}
		}

		if isEmpty {
			emptyRows[y] = true
		}
	}

	emptyColumns := make(map[int]bool)
	for x := 0; x < m.Columns; x++ {
		isEmpty := true
		for y := 0; y < m.Columns; y++ {
			if m.At(maps.Coordinate{X: x, Y: y}) == Galaxy {
				isEmpty = false
				break
			}
		}

		if isEmpty {
			emptyColumns[x] = true
		}
	}

	var expanded []string
	for y := 0; y < m.Rows; y++ {
		var row []byte
		for x := 0; x < m.Columns; x++ {
			row = append(row, byte(m.At(maps.Coordinate{X: x, Y: y})))
			if emptyColumns[x] {
				row = append(row, byte(Empty))
			}
		}

		expanded = append(expanded, string(row))
		if emptyRows[y] {
			expanded = append(expanded, string(row))
		}
	}

	expandedMap := maps.New(strings.Join(expanded, "\n"), func(_, _ int, b byte) Space {
		return Space(b)
	})

	var galaxies []maps.Coordinate
	for _, c := range expandedMap.Coordinates() {
		if expandedMap.At(c) == Galaxy {
			galaxies = append(galaxies, c)
		}
	}

	var distance int
	for i := range galaxies {
		for j := range galaxies {
			if i >= j {
				continue
			}

			distance += galaxies[i].ManhattanDistance(galaxies[j])
		}
	}

	return distance
}

func second(m maps.Map[Space], increaseBy int) int {
	emptyRows := make(map[int]bool)
	for y := 0; y < m.Rows; y++ {
		isEmpty := true
		for x := 0; x < m.Columns; x++ {
			if m.At(maps.Coordinate{X: x, Y: y}) == Galaxy {
				isEmpty = false
				break
			}
		}

		if isEmpty {
			emptyRows[y] = true
		}
	}

	emptyColumns := make(map[int]bool)
	for x := 0; x < m.Columns; x++ {
		isEmpty := true
		for y := 0; y < m.Columns; y++ {
			if m.At(maps.Coordinate{X: x, Y: y}) == Galaxy {
				isEmpty = false
				break
			}
		}

		if isEmpty {
			emptyColumns[x] = true
		}
	}

	var galaxies []maps.Coordinate
	for _, c := range m.Coordinates() {
		if m.At(c) == Galaxy {
			galaxies = append(galaxies, c)
		}
	}

	var distance int
	for i := range galaxies {
		for j := range galaxies {
			if i >= j {
				continue
			}

			g1 := galaxies[i]
			g2 := galaxies[j]

			steps := 0
			for x := maths.MinInt(g1.X, g2.X); x < maths.MaxInt(g1.X, g2.X); x++ {
				if emptyColumns[x] {
					steps += increaseBy
				} else {
					steps += 1
				}
			}
			for y := maths.MinInt(g1.Y, g2.Y); y < maths.MaxInt(g1.Y, g2.Y); y++ {
				if emptyRows[y] {
					steps += increaseBy
				} else {
					steps += 1
				}
			}

			distance += steps
		}
	}

	return distance
}

type Space byte

const (
	Empty  Space = Space('.')
	Galaxy Space = Space('#')
)

func (s Space) String() string {
	return string(s)
}
