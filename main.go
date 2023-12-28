package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mbark/aoc2023/day1"
	"github.com/mbark/aoc2023/day10"
	"github.com/mbark/aoc2023/day11"
	"github.com/mbark/aoc2023/day12"
	"github.com/mbark/aoc2023/day13"
	"github.com/mbark/aoc2023/day14"
	"github.com/mbark/aoc2023/day15"
	"github.com/mbark/aoc2023/day16"
	"github.com/mbark/aoc2023/day17"
	"github.com/mbark/aoc2023/day18"
	"github.com/mbark/aoc2023/day19"
	"github.com/mbark/aoc2023/day2"
	"github.com/mbark/aoc2023/day20"
	"github.com/mbark/aoc2023/day21"
	"github.com/mbark/aoc2023/day22"
	"github.com/mbark/aoc2023/day23"
	"github.com/mbark/aoc2023/day24"
	"github.com/mbark/aoc2023/day3"
	"github.com/mbark/aoc2023/day4"
	"github.com/mbark/aoc2023/day5"
	"github.com/mbark/aoc2023/day6"
	"github.com/mbark/aoc2023/day7"
	"github.com/mbark/aoc2023/day8"
	"github.com/mbark/aoc2023/day9"
	"github.com/mbark/aoc2023/util"
)

func main() {
	var flagDay = flag.Int("day", 0, "use test input")
	var flagTest = flag.Bool("test", false, "use test input")
	var cpuprofile = flag.Bool("profile", false, "write cpu profile to file")
	flag.Parse()

	if *cpuprofile {
		fmt.Println("using cpu profile")
		fn := util.WithProfiling()
		defer fn()
	}

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
	case 6:
		day6.Day6(input)
	case 7:
		day7.Day7(input)
	case 8:
		day8.Day8(input)
	case 9:
		day9.Day9(input)
	case 10:
		day10.Day10(input)
	case 11:
		day11.Day11(input)
	case 12:
		day12.Day12(input)
	case 13:
		day13.Day13(input)
	case 14:
		day14.Day14(input)
	case 15:
		day15.Day15(input)
	case 16:
		day16.Day16(input)
	case 17:
		day17.Day17(input)
	case 18:
		day18.Day18(input)
	case 19:
		day19.Day19(input)
	case 20:
		day20.Day20(input)
	case 21:
		day21.Day21(input)
	case 22:
		day22.Day22(input)
	case 23:
		day23.Day23(input)
	case 24:
		day24.Day24(input)
	default:
		fmt.Println("not implemented")
		os.Exit(1)
	}
}
