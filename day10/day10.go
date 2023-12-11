package day10

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2023/fns"
	"github.com/mbark/aoc2023/maps"
)

const testInput = `
FF7FSF7F7F7F7F7F---7
L|LJ||||||||||||F--J
FL-7LJLJ||||||LJL-77
F--JF--7||LJLJ7F7FJ-
L---JF-JLJ.||-FJLJJ7
|F|F-JF---7F7-L7L|7|
|FFJF7L7F-JF7|JL---7
7-L-JL7||F7|L7F-7F7|
L.L7LFJ|||||FJL7||LJ
L7JLJL-JLJLJL--JLJ.L
`

func Day10(input string) {
	if input == "" {
		input = testInput
	}

	m := maps.New[Tile](input, func(_, _ int, b byte) Tile { return Tile(b) })
	fmt.Printf("first: %d\n", first(m))

	var doubledColumns []string
	for _, row := range strings.Split(input, "\n") {
		var doubledRow []byte
		for _, b := range row {
			switch Tile(b) {
			case Ground:
				doubledRow = append(doubledRow, '.', '#')
			case PipeVert:
				doubledRow = append(doubledRow, '|', '#')
			case PipeHori:
				doubledRow = append(doubledRow, '-', '-')
			case BendNE:
				doubledRow = append(doubledRow, 'L', '-')
			case BendNW:
				doubledRow = append(doubledRow, 'J', '#')
			case BendSW:
				doubledRow = append(doubledRow, '7', '#')
			case BendSE:
				doubledRow = append(doubledRow, 'F', '-')
			case Start:
				doubledRow = append(doubledRow, 'S', '#')
			}
		}

		for x, t := range doubledRow {
			if t == '#' {
				doubledRow[x] = '.'
			}
		}

		doubledColumns = append(doubledColumns, string(doubledRow))
	}

	var doubledRows []string
	for _, row := range doubledColumns {
		var nextRow []byte
		for _, b := range row {
			switch Tile(b) {
			case Ground:
				nextRow = append(nextRow, '#')
			case PipeVert:
				nextRow = append(nextRow, '|')
			case PipeHori:
				nextRow = append(nextRow, '#')
			case BendNE:
				nextRow = append(nextRow, '#')
			case BendNW:
				nextRow = append(nextRow, '#')
			case BendSW:
				nextRow = append(nextRow, '|')
			case BendSE:
				nextRow = append(nextRow, '|')
			case Start:
				nextRow = append(nextRow, '#')
			}
		}

		for x, t := range nextRow {
			if t == '#' {
				nextRow[x] = '.'
			}
		}
		doubledRows = append(doubledRows, row, string(nextRow))
	}

	m2 := maps.New[Tile](strings.Join(doubledRows, "\n"), func(_, _ int, b byte) Tile { return Tile(b) })
	fmt.Printf("second: %d\n", second(m2))

	fmt.Println(m)
	fmt.Println(m2)
}

func first(m maps.Map[Tile]) int {
	start, _ := fns.Find(m.Coordinates(), func(c maps.Coordinate) bool { return m.At(c) == Start })

	var at maps.Coordinate
	var direction maps.Direction
	for _, c := range m.Adjacent(start) {
		if m.At(c) != Ground {
			at = c
			direction = maps.Direction{
				X: c.X - start.X,
				Y: c.Y - start.Y,
			}
		}
	}

	steps := 1
	for ; ; steps++ {
		if at == start {
			break
		}

		switch m.At(at) {
		case PipeVert, PipeHori:
		case BendNE:
			switch direction {
			case maps.South:
				direction = maps.East
			case maps.West:
				direction = maps.North
			}

		case BendNW:
			switch direction {
			case maps.South:
				direction = maps.West
			case maps.East:
				direction = maps.North
			}

		case BendSW:
			switch direction {
			case maps.North:
				direction = maps.West
			case maps.East:
				direction = maps.South
			}

		case BendSE:
			switch direction {
			case maps.North:
				direction = maps.East
			case maps.West:
				direction = maps.South
			}
		}

		at = direction.Apply(at)
	}

	return steps / 2
}

func second(m maps.Map[Tile]) int {
	start, _ := fns.Find(m.Coordinates(), func(c maps.Coordinate) bool { return m.At(c) == Start })

	var at maps.Coordinate
	var direction maps.Direction
	for _, c := range m.Adjacent(start) {
		if m.At(c) != Ground {
			at = c
			direction = maps.Direction{
				X: c.X - start.X,
				Y: c.Y - start.Y,
			}
		}
	}

	loop := map[maps.Coordinate]bool{start: true}
	for {
		loop[at] = true
		if at == start {
			break
		}

		switch m.At(at) {
		case PipeVert, PipeHori:
		case BendNE:
			switch direction {
			case maps.South:
				direction = maps.East
			case maps.West:
				direction = maps.North
			}

		case BendNW:
			switch direction {
			case maps.South:
				direction = maps.West
			case maps.East:
				direction = maps.North
			}

		case BendSW:
			switch direction {
			case maps.North:
				direction = maps.West
			case maps.East:
				direction = maps.South
			}

		case BendSE:
			switch direction {
			case maps.North:
				direction = maps.East
			case maps.West:
				direction = maps.South
			}
		}

		next := direction.Apply(at)
		at = next
	}

	contains := make(map[maps.Coordinate]bool)
	for _, c := range m.Coordinates() {
		if (c.X)%2 == 1 || c.Y%2 == 1 {
			continue
		}
		if loop[c] {
			continue
		}

		if floodFill(m, loop, c) {
			contains[c] = true
		}
	}

	return len(contains)
}

func floodFill(m maps.Map[Tile], loop map[maps.Coordinate]bool, at maps.Coordinate) bool {
	queue := []maps.Coordinate{at}
	visited := map[maps.Coordinate]bool{at: true}

	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]

		adjacent := m.Adjacent(next)
		if len(adjacent) != 4 {
			return false
		}

		for _, c := range adjacent {
			if visited[c] {
				continue
			}
			if loop[c] {
				continue
			}

			queue = append(queue, c)
			visited[c] = true
		}
	}

	return true
}

const (
	PipeVert = Tile('|')
	PipeHori = Tile('-')
	BendNE   = Tile('L')
	BendNW   = Tile('J')
	BendSW   = Tile('7')
	BendSE   = Tile('F')
	Ground   = Tile('.')
	Start    = Tile('S')
)

type Tile byte

func (t Tile) String() string {
	return string(t)
}
