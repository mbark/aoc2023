package day8

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2023/fns"
	"github.com/mbark/aoc2023/maths"
	"github.com/mbark/aoc2023/util"
)

const testInput = `
LR

11A = (11B, XXX)
11B = (XXX, 11Z)
11Z = (11B, XXX)
22A = (22B, XXX)
22B = (22C, 22C)
22C = (22Z, 22Z)
22Z = (22B, 22B)
XXX = (XXX, XXX)
`

func Day8(input string) {
	if input == "" {
		input = testInput
	}
	in := util.ReadInput(input, "\n\n")

	directions := in[0]

	var nodes []Node
	for _, s := range strings.Split(in[1], "\n") {
		split := strings.Split(s, " = ")
		name := split[0]
		split = strings.Split(split[1], ", ")

		nodes = append(nodes, Node{
			Name:  name,
			Left:  strings.TrimPrefix(split[0], "("),
			Right: strings.TrimSuffix(split[1], ")"),
		})
	}

	fmt.Printf("first: %d\n", first(directions, nodes))
	fmt.Printf("second: %d\n", second(directions, nodes))
}

func first(directions string, nodes []Node) int {
	m := fns.Associate(nodes, func(n Node) string { return n.Name })
	at := m["AAA"]
	end := m["ZZZ"]

	steps := 0
	for ; ; steps++ {
		if at == end {
			break
		}

		switch directions[steps%len(directions)] {
		case Left:
			at = m[at.Left]
		case Right:
			at = m[at.Right]
		}
	}

	return steps
}

func second(directions string, nodes []Node) int {
	m := fns.Associate(nodes, func(n Node) string { return n.Name })
	at := fns.Filter(nodes, func(n Node) bool { return strings.HasSuffix(n.Name, "A") })
	ends := fns.AsMap(
		fns.Filter(nodes, func(t Node) bool {
			return strings.HasSuffix(t.Name, "Z")
		}),
		func(t Node) (string, bool) {
			return t.Name, true
		})

	lcms := make(map[string]int)

	steps := 0
	for ; ; steps++ {
		for _, node := range at {
			if !ends[node.Name] {
				continue
			}
			if _, ok := lcms[node.Name]; ok {
				continue
			}

			lcms[node.Name] = steps
		}

		if len(lcms) == len(at) {
			break
		}

		switch directions[steps%len(directions)] {
		case Left:
			for i, n := range at {
				at[i] = m[n.Left]
			}
		case Right:
			for i, n := range at {
				at[i] = m[n.Right]
			}
		}
	}

	ints := fns.Values(lcms)
	return maths.LCM(ints[0], ints[1], ints[2:]...)
}

const (
	Left  = 'L'
	Right = 'R'
)

type Node struct {
	Name  string
	Left  string
	Right string
}

func (n Node) String() string {
	return fmt.Sprintf("%s = (%s,%s)", n.Name, n.Left, n.Right)
}
