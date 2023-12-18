package day18

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mbark/aoc2023/maps"
	"github.com/mbark/aoc2023/maths"
	"github.com/mbark/aoc2023/util"
)

const testInput = `
R 6 (#70c710)
D 5 (#0dc571)
L 2 (#5713f0)
D 2 (#d2c081)
R 2 (#59c680)
D 2 (#411b91)
L 5 (#8ceee2)
U 2 (#caa173)
L 1 (#1b58a2)
U 2 (#caa171)
R 2 (#7807d2)
U 3 (#a77fa3)
L 2 (#015232)
U 2 (#7a21e3)
`

func Day18(input string) {
	if input == "" {
		input = testInput
	}

	var plans []plan
	for _, s := range util.ReadInput(input, "\n") {
		split := strings.Split(s, " ")

		var dir maps.Direction
		switch split[0] {
		case "D":
			dir = maps.Down
		case "R":
			dir = maps.Right
		case "L":
			dir = maps.Left
		case "U":
			dir = maps.Up
		}

		plans = append(plans, plan{
			dir:    dir,
			meters: util.Str2Int(split[1]),
			color:  strings.TrimPrefix(strings.Trim(split[2], "()"), "#"),
		})
	}

	fmt.Printf("first: %d\n", first(plans))
	fmt.Printf("second: %d\n", second(plans))
}

type plan struct {
	dir    maps.Direction
	meters int
	color  string
}

func (p plan) String() string {
	return fmt.Sprintf("%s %d (#%s)", p.dir, p.meters, p.color)
}

func first(plans []plan) int {
	return shoelace(plans)
}

func second(incorrectPlans []plan) int {
	var plans []plan
	for _, p := range incorrectPlans {
		distance, _ := strconv.ParseInt(p.color[:5], 16, 64)

		var dir maps.Direction
		switch p.color[len(p.color)-1] {
		case '0':
			dir = maps.Right
		case '1':
			dir = maps.Down
		case '2':
			dir = maps.Left
		case '3':
			dir = maps.Up
		}

		plans = append(plans, plan{
			dir:    dir,
			meters: int(distance),
		})
	}

	return shoelace(plans)
}

func shoelace(plans []plan) int {
	var at maps.Coordinate
	var coordinates []maps.Coordinate

	for _, p := range plans {
		at = p.dir.ApplyN(at, p.meters)
		coordinates = append(coordinates, at)
	}

	var add int
	for i := range coordinates {
		x := coordinates[i].X
		y := coordinates[(i+1)%len(coordinates)].Y
		add += x * y
	}

	var remove int
	for i := range coordinates {
		x := coordinates[(i+1)%len(coordinates)].X
		y := coordinates[i].Y
		remove += y * x
	}

	var edges int
	for _, p := range plans {
		edges += p.meters
	}

	return maths.AbsInt(add-remove)/2 + edges/2 + 1
}
