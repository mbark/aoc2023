package day13

import (
	"fmt"

	"github.com/mbark/aoc2023/maps"
	"github.com/mbark/aoc2023/util"
)

const testInput = `
#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.

#...##..#
#....#..#
..##..###
#####.##.
#####.##.
..##..###
#....#..#
`

func Day13(input string) {
	if input == "" {
		input = testInput
	}

	var ms []maps.Map[byte]
	for _, in := range util.ReadInput(input, "\n\n") {
		m := maps.New(in, func(_, _ int, b byte) byte { return b })
		ms = append(ms, m)
	}

	fmt.Printf("first: %d\n", first(ms))
	fmt.Printf("second: %d\n", second(ms))
}

func first(ms []maps.Map[byte]) int {
	var sum int
	for _, m := range ms {
		for x := 0; x < m.Columns; x++ {
			if reflectionAtX(m, x) == nil {
				sum += x + 1
			}
		}

		for y := 0; y < m.Rows; y++ {
			if reflectionAtY(m, y) == nil {
				sum += 100 * (y + 1)
			}
		}
	}

	return sum
}

func second(ms []maps.Map[byte]) int {
	var sum int
	for _, m := range ms {
		for x := 0; x < m.Columns; x++ {
			if len(reflectionAtX(m, x)) == 1 {
				sum += x + 1
			}
		}

		for y := 0; y < m.Rows; y++ {
			if len(reflectionAtY(m, y)) == 1 {
				sum += 100 * (y + 1)
			}
		}
	}

	return sum
}

func reflectionAtY(m maps.Map[byte], y int) []maps.Coordinate {
	var diffs []maps.Coordinate
	for i := 0; ; i++ {
		row := 10*i + 5
		yUp := ((10*y + 5) - row) / 10
		yDown := ((10*y + 5) + row) / 10
		if yUp < 0 || yDown >= m.Rows {
			break
		}

		for x := 0; x < m.Columns; x++ {
			up := m.At(maps.Coordinate{X: x, Y: yUp})
			down := m.At(maps.Coordinate{X: x, Y: yDown})
			if up != down {
				diffs = append(diffs, maps.Coordinate{X: x, Y: yUp})
			}
		}
	}

	return diffs
}

func reflectionAtX(m maps.Map[byte], x int) []maps.Coordinate {
	var diffs []maps.Coordinate

	for i := 0; ; i++ {
		row := 10*i + 5
		xUp := ((10*x + 5) - row) / 10
		xDown := ((10*x + 5) + row) / 10
		if xUp < 0 || xDown >= m.Columns {
			break
		}

		for y := 0; y < m.Rows; y++ {
			up := m.At(maps.Coordinate{X: xUp, Y: y})
			down := m.At(maps.Coordinate{X: xDown, Y: y})
			if up != down {
				diffs = append(diffs, maps.Coordinate{X: xUp, Y: y})
			}
		}
	}

	return diffs
}

const (
	Ash   = '.'
	Rocks = '#'
)
