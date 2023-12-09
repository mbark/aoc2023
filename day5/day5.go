package day5

import (
	"fmt"
	"math"
	"strings"

	"github.com/mbark/aoc2023/fns"
	"github.com/mbark/aoc2023/maths"
	"github.com/mbark/aoc2023/util"
)

const testInput = `
seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4
`

func Day5(input string) {
	if input == "" {
		input = testInput
	}

	in := util.ReadInput(input, "\n\n")
	split := strings.Split(in[0], " ")

	seeds := fns.Map(split[1:], func(s string) int { return util.ParseInt[int](s) })

	var maps []Map
	for _, m := range in[1:] {
		split := strings.Split(m, "\n")
		names := strings.Split(strings.Split(split[0], " ")[0], "-to-")

		mp := Map{
			From: names[0],
			To:   names[1],
		}

		for _, s := range split[1:] {
			sp := strings.Split(s, " ")
			dest := util.ParseInt[int](sp[0])
			src := util.ParseInt[int](sp[1])
			length := util.ParseInt[int](sp[2])
			mp.Ranges = append(mp.Ranges, RangeMapping{
				Dest:   NewRange(dest, length),
				Source: NewRange(src, length),
			})
		}

		maps = append(maps, mp)
	}

	fmt.Printf("first: %d\n", first(seeds, maps))
	fmt.Printf("second: %d\n", second(seeds, maps))
}

func first(seeds []int, maps []Map) int {
	bySource := fns.Associate(maps, func(m Map) string { return m.From })

	minLocation := math.MaxInt
	for _, seed := range seeds {
		curr := bySource["seed"]
		n := seed
		for {
			r, ok := fns.Find(curr.Ranges, func(r RangeMapping) bool { return r.Source.Has(n) })
			if ok {
				n = r.Map(n)
			}

			curr, ok = bySource[curr.To]
			if !ok {
				break
			}
		}

		if n < minLocation {
			minLocation = n
		}
	}

	return minLocation
}

func second(seeds []int, maps []Map) int {
	var seedRanges []Range
	for i := 0; i < len(seeds)-1; i += 2 {
		seedRanges = append(seedRanges, NewRange(seeds[i], seeds[i+1]))
	}

	bySource := fns.Associate(maps, func(m Map) string { return m.From })

	minLocation := math.MaxInt
	for _, seed := range seedRanges {
		curr := bySource["seed"]
		ranges := []Range{seed}

		for {
			var added []Range
			for _, r := range ranges {
				unmapped := []Range{r}

				for _, rm := range curr.Ranges {
					var next []Range
					for _, ur := range unmapped {
						mapped, un := rm.MapRange(ur)
						next = append(next, un...)

						if mapped.Length() > 0 {
							added = append(added, mapped)
						}
					}

					unmapped = next
				}

				added = append(added, unmapped...)
			}

			ranges = added

			next, ok := bySource[curr.To]
			if !ok {
				break
			}

			curr = next
		}

		for _, r := range ranges {
			if r.Start < minLocation {
				minLocation = r.Start
			}
		}
	}

	return minLocation
}

type Map struct {
	From   string
	To     string
	Ranges []RangeMapping
}

func (m Map) String() string {
	ranges := fns.Map(m.Ranges, func(r RangeMapping) string { return r.String() })
	return fmt.Sprintf("%s-to-%s map:\n%s", m.From, m.To, strings.Join(ranges, "\n"))
}

type Range struct {
	Start int
	End   int
}

func NewRange(start, length int) Range {
	return Range{Start: start, End: start + length}
}

func (r Range) String() string {
	if r.Length() == 0 {
		return "<empty>"
	}
	return fmt.Sprintf("%d %d", r.Start, r.End-r.Start)
}

func (r Range) Has(n int) bool {
	return r.Start <= n && n < r.End
}

func (r Range) Length() int {
	return r.End - r.Start
}

type RangeMapping struct {
	Source Range
	Dest   Range
}

func (r RangeMapping) Map(n int) int {
	diff := r.Dest.Start - r.Source.Start
	return n + diff
}

func (r RangeMapping) MapRange(rn Range) (Range, []Range) {
	var mapped Range
	var unmapped []Range
	r1 := r.Source

	// before
	if rn.Start < r1.Start {
		unmapped = append(unmapped, Range{Start: rn.Start, End: maths.MinInt(rn.End, r1.Start)})
	}

	// overlap
	if rn.End >= r1.Start && rn.Start < r1.End {
		start := maths.MaxInt(r1.Start, rn.Start)
		offset := start - r1.Start
		length := maths.MinInt(rn.End, r1.End) - start
		mapped = NewRange(r.Dest.Start+offset, length)
	}

	// after
	if rn.End >= r1.End {
		unmapped = append(unmapped, Range{Start: maths.MaxInt(r1.End, rn.Start), End: rn.End})
	}

	unmapped = fns.Filter(unmapped, func(r Range) bool { return r.Length() > 0 })
	return mapped, unmapped
}

func (r RangeMapping) String() string {
	return fmt.Sprintf("%d %d %d", r.Source, r.Dest, r.Dest.Length())
}
