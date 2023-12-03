package day1

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mbark/aoc2023/util"
)

const testInput = `
1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet
`

func Day1(input string) {
	if input == "" {
		input = testInput
	}

	in := util.ReadInput(input, "\n")

	fmt.Printf("first: %d\n", part1(in))
	fmt.Printf("second: %d\n", part2(in))
}

var digit = regexp.MustCompile("[1-9]")

func part1(in []string) int {
	sum := 0
	for _, s := range in {
		matches := digit.FindAllString(s, -1)
		fmt.Println(matches)
		first := util.ParseInt[int](matches[0])
		last := util.ParseInt[int](matches[len(matches)-1])
		sum += 10*first + last
	}

	return sum
}

var digits = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

func part2(in []string) int {
	var fwd, bwd []string
	for s := range digits {
		fwd = append(fwd, s)
		bwd = append(bwd, util.Reverse(s))
	}

	digitFwd := regexp.MustCompile(fmt.Sprintf("([0-9]|%s)", strings.Join(fwd, "|")))
	digitBackwards := regexp.MustCompile(fmt.Sprintf("([0-9]|%s)", strings.Join(bwd, "|")))

	sum := 0
	for _, s := range in {
		forwardMatch := digitFwd.FindString(s)
		backwardMatch := digitBackwards.FindString(util.Reverse(s))

		var first, last int
		if f, ok := digits[forwardMatch]; ok {
			first = f
		} else {
			first = util.ParseInt[int](forwardMatch)
		}
		backwardMatch = util.Reverse(backwardMatch)
		if l, ok := digits[backwardMatch]; ok {
			last = l
		} else {
			last = util.ParseInt[int](backwardMatch)
		}

		sum += 10*first + last
	}

	return sum
}
