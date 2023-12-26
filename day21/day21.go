package day21

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2023/maps"
	"github.com/mbark/aoc2023/maths"
	"github.com/mbark/aoc2023/util"
)

const testInput = `
...........
.....###.#.
.###.##..#.
..#.#...#..
....#.#....
.##..S####.
.##..#...#.
.......##..
.##.#.####.
.##..##.##.
...........
`

func Day21(input string) {
	if input == "" {
		input = testInput
	}

	var start maps.Coordinate
	m := maps.New(input, func(x, y int, b byte) byte {
		if b == 'S' {
			start = maps.Coordinate{X: x, Y: y}
			return '.'
		}

		return b
	})

	var tbt []string
	for _, line := range util.ReadInput(input, "\n") {
		line = strings.Join([]string{line, line, line, line, line}, "")
		tbt = append(tbt, line)
	}

	oneByFive := strings.Join(tbt, "\n")
	fiveByFive := strings.Join([]string{oneByFive, oneByFive, oneByFive, oneByFive, oneByFive}, "\n")

	m2 := maps.New(fiveByFive, func(x, y int, b byte) byte {
		if b == 'S' {
			return '.'
		}

		return b
	})

	fmt.Printf("first: %d\n", first(start, m))
	fmt.Printf("second: %d\n", second(maps.Coordinate{
		X: start.X + 2*m.Columns,
		Y: start.Y + 2*m.Rows,
	}, m2))
}

func first(start maps.Coordinate, m maps.Map[byte]) int {
	possible := map[maps.Coordinate]bool{start: true}
	for i := 0; i < 64; i++ {
		next := make(map[maps.Coordinate]bool)
		for c := range possible {
			for _, ac := range m.Adjacent(c) {
				if m.At(ac) == '#' {
					continue
				}
				next[ac] = true
			}
		}

		possible = next
	}

	return len(possible)
}

func second(start maps.Coordinate, m maps.Map[byte]) int {
	possible := map[maps.Coordinate]bool{start: true}
	to := 26501365 // 202300*131+65
	cycles := to / (m.Rows / 5)
	extra := to % (m.Rows / 5)

	for steps := 0; steps < 2*(m.Rows/5)+extra; steps++ {
		next := make(map[maps.Coordinate]bool)
		for c := range possible {
			for _, ac := range m.Adjacent(c) {
				if m.At(ac) == '#' {
					continue
				}

				next[ac] = true
			}
		}

		possible = next
	}

	// segment into 5x5
	// x 0 1 2 3 4
	// 0 * * * * *
	// 1 * * * * *
	// 2 * * S * *
	// 3 * * * * *
	// 4 * * * * *
	// S starts at 2,2

	segments := make([][]int, 5)
	for i := 0; i < 5; i++ {
		segments[i] = make([]int, 5)
	}
	for c := range possible {
		column := c.X / (m.Columns / 5)
		row := c.Y / (m.Rows / 5)

		segments[row][column] += 1
	}
	for row, r := range segments {
		for col, count := range r {
			fmt.Printf("[row=%d,col=%d]: %d\n", row, col, count)
		}
	}

	var sum int
	sum += segments[0][1] * cycles
	sum += segments[0][2]
	sum += segments[0][3] * cycles
	sum += segments[1][1] * (cycles - 1)
	sum += segments[1][2] * maths.PowInt(cycles, 2)
	sum += segments[1][3] * (cycles - 1)
	sum += segments[2][0]
	sum += segments[2][2] * maths.PowInt(cycles-1, 2)
	sum += segments[2][4]
	sum += segments[3][0] * cycles
	sum += segments[3][1] * (cycles - 1)
	sum += segments[3][3] * (cycles - 1)
	sum += segments[3][4] * cycles
	sum += segments[4][2]

	return sum
}
