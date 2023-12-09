package main

import (
	"flag"

	"github.com/mbark/aoc2023/day1"
	"github.com/mbark/aoc2023/day2"
	"github.com/mbark/aoc2023/day3"
	"github.com/mbark/aoc2023/day4"
	"github.com/mbark/aoc2023/day5"
	"github.com/mbark/aoc2023/util"
)

func main() {
	var flagDay = flag.Int("day", 0, "use test input")
	var flagTest = flag.Bool("test", false, "use test input")
	flag.Parse()

	var input string
	if !*flagTest {
		input = util.GetInput(*flagDay)
	}

	switch *flagDay {
	case 1:
		day1.Day1(input)
	case 2:
		day2.Day2(input)
	case 3:
		day3.Day3(input)
	case 4:
		day4.Day4(input)
	case 5:
		day5.Day5(input)
	}
}
