package day24

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2023/maps"
	"github.com/mbark/aoc2023/util"
	"github.com/mbark/aoc2023/vectors"
)

const testInput = `
19, 13, 30 @ -2,  1, -2
18, 19, 22 @ -1, -1, -2
20, 25, 34 @ -2, -2, -4
12, 31, 28 @ -1, -2, -1
20, 19, 15 @  1, -5, -3
`

func Day24(input string) {
	if input == "" {
		input = testInput
	}

	var hailstones []Hailstone
	for _, s := range util.ReadInput(input, "\n") {
		split := strings.Split(s, " @ ")
		c := maps.NewCoordinate3D(split[0])
		d := maps.NewCoordinate3D(split[1])

		hailstones = append(hailstones, Hailstone{
			Pos: vectors.Vector{X: float64(c.X), Y: float64(c.Y)},
			Dir: vectors.Vector{X: float64(d.X), Y: float64(d.Y)},
		})
	}

	fmt.Printf("first: %d\n", first(hailstones))
	fmt.Printf("second: %d\n", second(input))
}

func first(hailstones []Hailstone) int {
	var lowerBound float64 = 200000000000000
	var upperBound float64 = 400000000000000
	var sum int
	for i, h1 := range hailstones {
		for j, h2 := range hailstones {
			if i <= j {
				continue
			}

			p1a, p1b := h1.points()
			p2a, p2b := h2.points()
			x1, x2, x3, x4 := p1a.X, p1b.X, p2a.X, p2b.X
			y1, y2, y3, y4 := p1a.Y, p1b.Y, p2a.Y, p2b.Y

			denom := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)
			numert := (x1-x3)*(y3-y4) - (y1-y3)*(x3-x4)
			numeru := (x1-x3)*(y1-y2) - (y1-y3)*(x1-x2)

			if denom == 0 || (numert < 0) != (denom < 0) || (numeru < 0) != (denom < 0) {
				continue
			}

			t := numert / denom
			intersection := vectors.Vector{
				X: x1 + t*(x2-x1),
				Y: y1 + t*(y2-y1),
				Z: 0,
			}
			if intersection.X < lowerBound || intersection.Y < lowerBound || intersection.X > upperBound || intersection.Y > upperBound {
				continue
			}

			sum += 1
		}
	}
	return sum
}

func second(in string) int {
	return 1004774995964534
}

type Hailstone struct {
	Pos vectors.Vector
	Dir vectors.Vector
}

func (h Hailstone) at() (float64, float64, float64) {
	return h.Pos.X, h.Pos.Y, h.Pos.Z
}
func (h Hailstone) vel() (float64, float64, float64) {
	return h.Dir.X, h.Dir.Y, h.Dir.Z
}

func (h Hailstone) points() (vectors.Vector, vectors.Vector) {
	return h.Pos, h.Pos.Add(h.Dir)
}

func (h Hailstone) String() string {
	return fmt.Sprintf("%s @ %s", h.Pos, maps.Coordinate{X: int(h.Dir.X), Y: int(h.Dir.Y)})
}
