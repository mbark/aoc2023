package day23

import (
	"fmt"

	"github.com/mbark/aoc2023/fns"
	"github.com/mbark/aoc2023/maps"
	"github.com/mbark/aoc2023/maths"
)

const testInput = `
#.#####################
#.......#########...###
#######.#########.#.###
###.....#.>.>.###.#.###
###v#####.#v#.###.#.###
###.>...#.#.#.....#...#
###v###.#.#.#########.#
###...#.#.#.......#...#
#####.#.#.#######.#.###
#.....#.#.#.......#...#
#.#####.#.#.#########v#
#.#...#...#...###...>.#
#.#.#v#######v###.###v#
#...#.>.#...>.>.#.###.#
#####v#.#.###v#.#.###.#
#.....#...#...#.#.#...#
#.#########.###.#.#.###
#...###...#...#...#.###
###.###.#.###v#####v###
#...#...#.#.>.>.#.>.###
#.###.###.#.###.#.#v###
#.....###...###...#...#
#####################.#
`

func Day23(input string) {
	if input == "" {
		input = testInput
	}

	m := maps.New(input, func(_, _ int, b byte) byte { return b })
	m2 := maps.New(input, func(_, _ int, b byte) byte {
		if b == '.' || b == '#' {
			return b
		}
		return '.'
	})
	fmt.Printf("first: %d\n", first(m))
	fmt.Printf("second: %d\n", second(m2))
}

func first(m maps.Map[byte]) int {
	start := maps.Coordinate{X: 1, Y: 0}
	return dfs(m, start, map[maps.Coordinate]bool{start: true}, 0)
}

func second(m maps.Map[byte]) int {
	start := maps.Coordinate{X: 1, Y: 0}
	end := maps.Coordinate{X: m.Columns - 2, Y: m.Rows - 1}

	neighbors := bfs(m, start)
	visited := make(map[maps.Coordinate]bool)
	compacted := make(map[maps.Coordinate]map[maps.Coordinate]int)

	for _, c := range m.Coordinates() {
		compacted[c] = make(map[maps.Coordinate]int)
	}

	for _, c := range m.Coordinates() {
		visited[c] = true
		for nc := range neighbors[c] {
			if visited[nc] {
				continue
			}

			next := nc
			visited[next] = true
			ok := true
			steps := 1
			for ok && len(neighbors[next]) == 2 {
				visited[next] = true
				next, ok = fns.Find(fns.Keys(neighbors[next]), func(n maps.Coordinate) bool { return !visited[n] })
				steps++
			}
			if !ok {
				continue
			}

			compacted[c][next] = steps
			compacted[next][c] = steps
		}

		if len(compacted[c]) > 0 {
			fmt.Println("connected", c, "to", compacted[c])
		}
	}

	for c, n := range compacted {
		if len(n) == 0 {
			delete(compacted, c)
		}
	}

	return dfsGeneral(func(c maps.Coordinate) []withCost {
		return fns.Map(fns.Keys(compacted[c]), func(nc maps.Coordinate) withCost {
			return withCost{c: nc, cost: compacted[c][nc]}
		})
	}, start, end, map[maps.Coordinate]bool{start: true}, 0)
}

func dfs(m maps.Map[byte], at maps.Coordinate, visited map[maps.Coordinate]bool, steps int) int {
	end := maps.Coordinate{X: m.Columns - 2, Y: m.Rows - 1}
	if at == end {
		return steps
	}

	nextv := make(map[maps.Coordinate]bool, len(visited)+1)
	for c := range visited {
		nextv[c] = true
	}
	nextv[at] = true

	var dir maps.Direction
	switch m.At(at) {
	case '>':
		dir = maps.Right
	case '<':
		dir = maps.Left
	case 'v':
		dir = maps.Down
	case '^':
		dir = maps.Up
	}

	if dir != maps.NoDirection {
		to := dir.Apply(at)
		if nextv[to] {
			return 0
		}

		return dfs(m, to, nextv, steps+1)
	}

	adjacent := m.Adjacent(at)
	var nextSteps int
	for _, c := range adjacent {
		if _, ok := nextv[c]; ok {
			continue
		}
		if m.At(c) == '#' {
			continue
		}

		nextSteps = maths.MaxInt(nextSteps, dfs(m, c, nextv, steps+1))
	}

	return nextSteps
}

func bfs(m maps.Map[byte], start maps.Coordinate) map[maps.Coordinate]map[maps.Coordinate]bool {
	queue := []maps.Coordinate{start}
	visited := make(map[maps.Coordinate]bool)
	neighbors := make(map[maps.Coordinate]map[maps.Coordinate]bool)
	for _, c := range m.Coordinates() {
		neighbors[c] = make(map[maps.Coordinate]bool)
	}

	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]
		visited[next] = true

		for _, ac := range m.Adjacent(next) {
			if m.At(ac) == '#' {
				continue
			}
			if visited[ac] {
				continue
			}

			neighbors[next][ac] = true
			neighbors[ac][next] = true
			queue = append(queue, ac)
		}
	}

	return neighbors
}

type withCost struct {
	c    maps.Coordinate
	cost int
}

func dfsGeneral(adjacent func(c maps.Coordinate) []withCost, at maps.Coordinate, end maps.Coordinate, visited map[maps.Coordinate]bool, steps int) int {
	if at == end {
		return steps
	}

	nextv := make(map[maps.Coordinate]bool, len(visited)+1)
	for c := range visited {
		nextv[c] = true
	}
	nextv[at] = true

	var nextSteps int
	for _, c := range adjacent(at) {
		if _, ok := nextv[c.c]; ok {
			continue
		}

		nextSteps = maths.MaxInt(nextSteps, dfsGeneral(adjacent, c.c, end, nextv, steps+c.cost))
	}

	return nextSteps
}

func topologicalSort(m maps.Map[byte], at maps.Coordinate, visited map[maps.Coordinate]bool) []maps.Coordinate {
	if visited[at] {
		return nil
	}
	visited[at] = true

	var coordinates []maps.Coordinate
	for _, c := range m.Adjacent(at) {
		if m.At(c) == '#' {
			continue
		}

		coordinates = append(coordinates, topologicalSort(m, c, visited)...)
	}

	return append(coordinates, at)
}

func longestPath(m maps.Map[byte], at maps.Coordinate, sorting []maps.Coordinate) int {
	distances := map[maps.Coordinate]int{at: 0}
	for len(sorting) > 0 {
		next := sorting[len(sorting)-1]
		sorting = sorting[:len(sorting)-1]

		for _, c := range m.Adjacent(next) {
			if m.At(c) == '#' {
				continue
			}

			distances[c] = maths.MaxInt(distances[c], distances[next]+1)
		}
	}

	end := maps.Coordinate{X: m.Columns - 2, Y: m.Rows - 1}
	return distances[end]
}
