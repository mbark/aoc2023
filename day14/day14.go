package day14

import (
	"fmt"
	"math/rand"

	"github.com/mbark/aoc2023/maps"
)

const testInput = `
O....#....
O.OO#....#
.....##...
OO.#O....O
.O.....O#.
O.#..O.#.#
..O..#O..O
.......O..
#....###..
#OO..#....
`

func Day14(input string) {
	if input == "" {
		input = testInput
	}

	m := maps.New(input, func(_, _ int, b byte) byte { return b })

	fmt.Printf("first: %d\n", first(m))
	fmt.Printf("second: %d\n", second(m))
}

func first(m maps.Map[byte]) int {
	coords := m.Coordinates()
	rand.Shuffle(len(coords), func(i, j int) {
		coords[i], coords[j] = coords[j], coords[i]
	})

	for rolled := true; rolled; {
		rolled = false
		for _, c := range coords {
			if m.At(c) != Round {
				continue
			}

			nc := roll(m, c, maps.North)
			m.Set(c, Empty)
			m.Set(nc, Round)
			if c != nc {
				rolled = true
			}
		}
	}

	var sum int
	for _, c := range m.Coordinates() {
		if m.At(c) == Round {
			sum += m.Rows - c.Y
		}
	}

	return sum
}

func second(m maps.Map[byte]) int {
	steps := []string{m.String()}
	cycle := map[string]int{m.String(): 0}

	var cycleLength int
	var start int
	for i := 1; ; i++ {
		for _, direction := range []maps.Direction{maps.North, maps.West, maps.South, maps.East} {
			for rolled := true; rolled; {
				rolled = false
				for _, c := range m.Coordinates() {
					if m.At(c) != Round {
						continue
					}

					nc := roll(m, c, direction)
					m.Set(c, Empty)
					m.Set(nc, Round)
					if nc != c {
						rolled = true
					}
				}
			}
		}

		s := m.String()
		steps = append(steps, s)
		if idx, ok := cycle[s]; ok {
			start = idx
			cycleLength = i - idx
			break
		}

		cycle[s] = i
	}

	idx := (1000000000 - start) % cycleLength
	idx += start
	m = maps.New(steps[idx], func(_, _ int, b byte) byte { return b })
	var sum int
	for _, c := range m.Coordinates() {
		if m.At(c) == Round {
			sum += m.Rows - c.Y
		}
	}

	return sum
}

func roll(m maps.Map[byte], nc maps.Coordinate, direction maps.Direction) maps.Coordinate {
	for {
		nextc := direction.Apply(nc)
		if !m.Exists(nextc) || m.At(nextc) != Empty {
			return nc
		}

		nc = nextc
	}
}

const (
	Round = 'O'
	Cube  = '#'
	Empty = '.'
)
