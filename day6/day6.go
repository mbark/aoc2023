package day6

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2023/fns"
	"github.com/mbark/aoc2023/util"
)

const testInput = `
Time:      7  15   30
Distance:  9  40  200
`

func Day6(input string) {
	if input == "" {
		input = testInput
	}

	in := util.ReadInput(input, "\n")
	times := fns.Filter(strings.Split(in[0], " ")[1:], func(s string) bool {
		return s != ""
	})
	distances := fns.Filter(strings.Split(in[1], " ")[1:], func(s string) bool {
		return s != ""
	})

	var races []Race
	for i := range times {
		races = append(races, Race{
			Time:     util.ParseInt[int](times[i]),
			Distance: util.ParseInt[int](distances[i]),
		})
	}

	secondRace := Race{
		Time:     util.ParseInt[int](strings.Join(times, "")),
		Distance: util.ParseInt[int](strings.Join(distances, "")),
	}

	fmt.Printf("first: %d\n", first(races))
	fmt.Printf("second: %d\n", second(secondRace))
}

func first(races []Race) int {
	sum := 1
	for _, race := range races {
		var start, end int

		for hold := 1; hold < race.Time; hold++ {
			distance := (race.Time - hold) * hold
			if distance > race.Distance {
				start = hold
				break
			}
		}

		for hold := race.Time - 1; hold > 0; hold-- {
			distance := (race.Time - hold) * hold
			if distance > race.Distance {
				end = hold
				break
			}
		}

		sum *= end - start + 1
	}

	return sum
}

func second(race Race) int {
	var start, end int

	for hold := 1; hold < race.Time; hold++ {
		distance := (race.Time - hold) * hold
		if distance > race.Distance {
			start = hold
			break
		}
	}

	for hold := race.Time - 1; hold > 0; hold-- {
		distance := (race.Time - hold) * hold
		if distance > race.Distance {
			end = hold
			break
		}
	}

	return end - start + 1
}

type Race struct {
	Time     int
	Distance int
}
