package day12

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mbark/aoc2023/fns"
	"github.com/mbark/aoc2023/maths"
	"github.com/mbark/aoc2023/util"
)

const testInput = `
???.### 1,1,3
.??..??...?##. 1,1,3
?#?#?#?#?#?#?#? 1,3,1,6
????.#...#... 4,1,1
????.######..#####. 1,6,5
?###???????? 3,2,1
`

func Day12(input string) {
	if input == "" {
		input = testInput
	}

	var records1, records2 []Record
	for _, s := range util.ReadInput(input, "\n") {
		split := strings.Split(s, " ")

		r1 := Record{
			Springs: split[0],
			Groups:  fns.Map(strings.Split(split[1], ","), func(s string) int { return util.ParseInt[int](s) }),
		}
		records1 = append(records1, r1)

		springs2 := strings.Repeat(split[0]+"?", 5)
		springs2 = springs2[:len(springs2)-1]

		var groups []int
		for i := 0; i < 5; i++ {
			groups = append(groups, r1.Groups...)
		}

		records2 = append(records2, Record{Springs: springs2, Groups: groups})
	}

	fmt.Printf("first: %d\n", first(records1))
	fmt.Printf("second: %d\n", first(records2))
}

func first(records []Record) int {
	var sum int
	for _, record := range records {
		memo = map[memoKey]int{}
		count := resolve(record.Springs, record.Groups)
		sum += count
	}
	return sum
}

type memoKey struct {
	springs string
	groups  int
}

func key(springs string, groups []int) memoKey {
	var g int
	for i, group := range groups {
		g += maths.PowInt(10, i) * group
	}

	return memoKey{springs: springs, groups: g}
}

var memo = map[memoKey]int{}

func resolve(springs string, groups []int) (c int) {
	k := key(springs, groups)
	if i, ok := memo[k]; ok {
		return i
	}
	defer func() { memo[k] = c }()

	// no groups left
	if len(groups) == 0 {
		if strings.ContainsRune(springs, '#') {
			return 0
		}

		return 1
	}

	// groups but no springs left
	if len(springs) == 0 {
		return 0
	}

	group := groups[0]

	handleDamaged := func() int {
		if len(springs) < group {
			return 0
		}

		damaged := springs[:group]
		damaged = strings.ReplaceAll(damaged, "?", "#")
		// make sure the group can fit
		if damaged != strings.Repeat("#", group) {
			return 0
		}

		if len(springs) == group {
			if len(groups) == 1 {
				return 1
			}

			return 0
		}

		// we need a ? or . to space them
		if springs[group] == Damaged {
			return 0
		}

		return resolve(springs[group+1:], groups[1:])
	}

	handleOperational := func() int {
		return resolve(springs[1:], groups)
	}

	switch springs[0] {
	case Damaged:
		return handleDamaged()
	case Operational:
		return handleOperational()
	default:
		return handleDamaged() + handleOperational()
	}
}

type Record struct {
	Springs string
	Groups  []int
}

func (r Record) String() string {
	return fmt.Sprintf("%s %s",
		r.Springs,
		strings.Join(fns.Map(r.Groups, func(i int) string { return strconv.Itoa(i) }), ","))
}

const (
	Operational = '.'
	Damaged     = '#'
	Unknown     = '?'
)
