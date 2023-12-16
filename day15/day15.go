package day15

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2023/fns"
	"github.com/mbark/aoc2023/util"
)

const testInput = `rn=1,cm-,qp=3,cm=2,qp-,pc=4,ot=9,ab=5,pc-,pc=6,ot=7`

func Day15(input string) {
	if input == "" {
		input = testInput
	}

	in := util.ReadInput(input, ",")
	fmt.Printf("first: %d\n", first(in))
	fmt.Printf("second: %d\n", second(in))
}

func first(seq []string) int {
	var sum int32
	for _, s := range seq {
		var val int32
		for _, b := range s {
			val += b
			val *= 17
			val %= 256
		}

		sum += val
	}

	return int(sum)
}

func second(seq []string) int {
	boxes := make([][]Lens, 256)

	for _, s := range seq {
		var key string
		var val int
		var del bool
		if strings.Contains(s, "=") {
			s := strings.Split(s, "=")
			key = s[0]
			val = util.ParseInt[int](s[1])
		} else {
			key = strings.TrimSuffix(s, "-")
			del = true
		}

		var hash int32
		for _, b := range key {
			hash += b
			hash *= 17
			hash %= 256
		}

		switch del {
		case true:
			boxes[hash] = fns.Filter(boxes[hash], func(l Lens) bool {
				return l.Key != key
			})
		default:
			var found bool
			for i := range boxes[hash] {
				if boxes[hash][i].Key == key {
					found = true
					boxes[hash][i].Val = val
				}
			}

			if !found {
				boxes[hash] = append(boxes[hash], Lens{Key: key, Val: val})
			}
		}
	}

	var sum int
	for i, lenses := range boxes {
		for j, lens := range lenses {
			sum += (1 + i) * (j + 1) * lens.Val
		}
	}

	return sum
}

type Lens struct {
	Key string
	Val int
}
