package day16

import (
	"fmt"

	"github.com/mbark/aoc2023/maps"
)

const testInput = `
.|...\....
|.-.\.....
.....|-...
........|.
..........
.........\
..../.\\..
.-.-/..|..
.|....-|.\
..//.|....
`

func Day16(input string) {
	if input == "" {
		input = testInput
	}

	m := maps.New(input, func(_, _ int, b byte) byte { return b })
	fmt.Printf("first: %d\n", first(m, Beam{Start: maps.Coordinate{}, Dir: maps.Right}))
	fmt.Printf("second: %d\n", second(m))
}

func first(m maps.Map[byte], start Beam) int {
	beams := map[Beam]bool{start: false}
	energized := make(map[maps.Coordinate][]maps.Direction)

	addBeam := func(at maps.Coordinate, dir maps.Direction) {
		b := Beam{Start: at, Dir: dir}
		if _, ok := beams[b]; ok {
			return
		}

		beams[b] = false
	}

	for {
		before := len(beams)
		for b, ok := range beams {
			at := b.Start
			dir := b.Dir

			for !ok {
				energized[at] = append(energized[at], dir)

				switch m.At(at) {
				case Empty:
				case RightMirror:
					switch dir {
					case maps.Up:
						dir = maps.Right
					case maps.Left:
						dir = maps.Down
					case maps.Right:
						dir = maps.Up
					case maps.Down:
						dir = maps.Left
					}
				case LeftMirror:
					switch dir {
					case maps.Up:
						dir = maps.Left
					case maps.Left:
						dir = maps.Up
					case maps.Right:
						dir = maps.Down
					case maps.Down:
						dir = maps.Right
					}
				case SplitterVert:
					switch dir {
					case maps.Left, maps.Right:
						addBeam(at, maps.Up)
						addBeam(at, maps.Down)
						ok = true
					}
				case SplitterHori:
					switch dir {
					case maps.Up, maps.Down:
						addBeam(at, maps.Left)
						addBeam(at, maps.Right)
						ok = true
					}
				}

				nextc := dir.Apply(at)
				if !m.Exists(nextc) {
					break
				}

				at = nextc
			}

			beams[b] = true
		}

		if len(beams) <= before {
			break
		}
	}

	return len(energized)
}

func second(m maps.Map[byte]) int {
	var mx int
	for x := 0; x < m.Columns; x++ {
		tiles := first(m, Beam{Start: maps.Coordinate{X: x, Y: 0}, Dir: maps.Down})
		if tiles > mx {
			mx = tiles
		}
		tiles = first(m, Beam{Start: maps.Coordinate{X: x, Y: m.Rows - 1}, Dir: maps.Up})
		if tiles > mx {
			mx = tiles
		}
	}
	for y := 0; y < m.Rows; y++ {
		tiles := first(m, Beam{Start: maps.Coordinate{X: 0, Y: y}, Dir: maps.Right})
		if tiles > mx {
			mx = tiles
		}
		tiles = first(m, Beam{Start: maps.Coordinate{X: m.Columns - 1, Y: y}, Dir: maps.Left})
		if tiles > mx {
			mx = tiles
		}
	}
	return mx
}

type Beam struct {
	Start maps.Coordinate
	Dir   maps.Direction
}

const (
	Empty        = '.'
	RightMirror  = '/'
	LeftMirror   = '\\'
	SplitterVert = '|'
	SplitterHori = '-'
)
