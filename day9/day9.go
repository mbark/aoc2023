package day9

import (
	"fmt"
	"slices"
	"strings"

	"github.com/mbark/aoc2023/fns"
	"github.com/mbark/aoc2023/util"
)

const testInput = `
0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45
`

func Day9(input string) {
	if input == "" {
		input = testInput
	}

	var readings [][]int

	for _, s := range util.ReadInput(input, "\n") {
		split := strings.Split(s, " ")
		readings = append(readings, fns.Map(split, func(s string) int { return util.ParseInt[int](s) }))
	}

	fmt.Printf("first: %d\n", first(readings))
	fmt.Printf("second: %d\n", second(readings))
}

func first(readings [][]int) int {
	var sum int
	for _, reading := range readings {
		var diffs = [][]int{reading}
		current := reading
		for {
			var diff []int
			for i := 0; i < len(current)-1; i++ {
				diff = append(diff, current[i+1]-current[i])
			}

			diffs = append(diffs, diff)
			current = diff

			if fns.Every(diff, func(t int) bool { return t == 0 }) {
				break
			}
		}

		slices.Reverse(diffs)
		for i := 1; i < len(diffs); i++ {
			diff := diffs[i]
			prev := diffs[i-1]
			diffs[i] = append(diff, diff[len(diff)-1]+prev[len(prev)-1])
		}

		last := diffs[len(diffs)-1]
		sum += last[len(last)-1]
	}

	return sum
}

func second(readings [][]int) int {
	var sum int
	readings = fns.Map(readings, func(r []int) []int {
		slices.Reverse(r)
		return r
	})

	for _, reading := range readings {
		var diffs = [][]int{reading}
		current := reading
		for {
			var diff []int
			for i := 0; i < len(current)-1; i++ {
				diff = append(diff, current[i+1]-current[i])
			}

			diffs = append(diffs, diff)
			current = diff

			if fns.Every(diff, func(t int) bool { return t == 0 }) {
				break
			}
		}

		slices.Reverse(diffs)
		for i := 1; i < len(diffs); i++ {
			diff := diffs[i]
			prev := diffs[i-1]
			diffs[i] = append(diff, diff[len(diff)-1]+prev[len(prev)-1])
		}

		last := diffs[len(diffs)-1]
		sum += last[len(last)-1]
	}

	return sum
}
