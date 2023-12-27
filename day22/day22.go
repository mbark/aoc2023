package day22

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2023/fns"
	"github.com/mbark/aoc2023/maps"
	"github.com/mbark/aoc2023/util"
)

const testInput = `
1,0,1~1,2,1
0,0,2~2,0,2
0,2,3~2,2,3
0,0,4~0,2,4
2,0,5~2,2,5
0,1,6~2,1,6
1,1,8~1,1,9
`

func Day22(input string) {
	if input == "" {
		input = testInput
	}

	var cubes []maps.Cuboid
	for _, s := range util.ReadInput(input, "\n") {
		split := strings.Split(s, "~")
		c := maps.Cuboid{
			From: maps.NewCoordinate3D(split[0]),
			To:   maps.NewCoordinate3D(split[1]),
		}
		cubes = append(cubes, c)
	}

	f, supports, supportedBy := first(cubes)
	fmt.Printf("first: %d\n", f)
	fmt.Printf("second: %d\n", second(supports, supportedBy))
}

func first(cubes []maps.Cuboid) (int, map[maps.Cuboid][]maps.Cuboid, map[maps.Cuboid][]maps.Cuboid) {
	var stable, falling []maps.Cuboid
	for _, cube := range cubes {
		if cube.From.Z == 1 || cube.To.Z == 1 {
			stable = append(stable, cube)
		} else {
			falling = append(falling, cube)
		}
	}

	for len(falling) > 0 {
		var next []maps.Cuboid

		var collided bool
		for i := 0; !collided && i < len(falling); i++ {
			moved := falling[i].Move(maps.ZUp)

			_, ok := fns.Find(stable, func(c maps.Cuboid) bool { return c.IsOverlapping(moved) })
			if ok {
				stable = append(stable, falling[i])
				collided = true
				next = remove(falling, i)
			} else if moved.To.Z == 0 || moved.From.Z == 0 {
				stable = append(stable, falling[i])
				collided = true
				next = remove(falling, i)
			} else {
				next = append(next, moved)
			}
		}

		falling = next
	}

	supports := make(map[maps.Cuboid][]maps.Cuboid)
	supportedBy := make(map[maps.Cuboid][]maps.Cuboid)
	for _, c := range stable {
		supports[c] = nil
		supportedBy[c] = nil
	}

	for i, c := range stable {
		above := c.Move(maps.ZDown)
		for j, other := range stable {
			if i == j {
				continue
			}

			if other.IsOverlapping(above) {
				supports[c] = append(supports[c], other)
				supportedBy[other] = append(supportedBy[other], c)
			}
		}
	}

	var count int
	for _, cube := range stable {
		supported := supports[cube]
		hasOtherSupport := fns.Every(supported, func(other maps.Cuboid) bool { return len(supportedBy[other]) > 1 })

		if hasOtherSupport {
			count++
		}
	}

	return count, supports, supportedBy
}

func second(supports map[maps.Cuboid][]maps.Cuboid, supportedBy map[maps.Cuboid][]maps.Cuboid) int {
	var sum int
	for cube := range supports {
		supported := supports[cube]
		hasOtherSupport := fns.Every(supported, func(other maps.Cuboid) bool { return len(supportedBy[other]) > 1 })

		if hasOtherSupport {
			continue
		}

		falling := map[maps.Cuboid]bool{cube: true}
		for {
			atStart := len(falling)
			// check all bricks
			for other, otherSupport := range supportedBy {
				// if this brick is already falling, ignore it
				if falling[other] || len(otherSupport) == 0 {
					continue
				}

				// check all that this rock is supported by, are they all falling?
				allFalling := fns.Every(otherSupport, func(o maps.Cuboid) bool {
					return falling[o]
				})

				if allFalling {
					falling[other] = true
					break
				}
			}

			if len(falling) == atStart {
				break
			}
		}

		sum += len(falling) - 1
	}

	return sum
}

func remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
