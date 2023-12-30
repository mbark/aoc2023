package day25

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/mbark/aoc2023/util"
)

const testInput = `
jqt: rhn xhk nvd
rsh: frs pzl lsr
xhk: hfx
cmg: qnr nvd lhk bvb
rhn: xhk bvb hfx
bvb: xhk hfx
pzl: lsr hfx nvd
qnr: nvd
ntq: jqt hfx bvb xhk
nvd: lhk
lsr: lhk
rzs: qnr cmg lsr rsh
frs: qnr lhk lsr
`

func Day25(input string) {
	if input == "" {
		input = testInput
	}

	graph := make(Graph)
	addNode := func(name string) {
		if _, ok := graph[name]; !ok {
			graph[name] = make(map[string]int)
		}
	}

	var ss []string
	for _, row := range util.ReadInput(input, "\n") {
		split := strings.Split(row, ":")
		name := split[0]
		addNode(name)

		var names []string
		for _, s := range strings.Split(strings.TrimSpace(split[1]), " ") {
			names = append(names, s)
		}
		for _, s := range names {
			addNode(s)
			graph[name][s] = 0
			graph[s][name] = 0
		}

		ss = append(ss, fmt.Sprintf("%s -- %s", name, strings.Join(names, ", ")))
	}

	// dot -Tsvg -Ksfdp day25.dot > day25.svg
	fmt.Println(strings.Join(ss, "\n"))

	fmt.Printf("first: %d\n", first(graph))
}

type Graph map[string]map[string]int

func first(g Graph) int {
	remove := map[string]string{"krx": "lmg", "tnr": "vzb", "tqn": "tvf"}
	for s, s2 := range remove {
		remove[s2] = s
	}

	for s, neighbors := range g {
		other, ok := remove[s]
		if !ok {
			continue
		}

		delete(neighbors, other)
		fmt.Println("removing", other, "from", s)
		delete(remove, s)
	}

	fmt.Println(remove)

	start1, start2 := "krx", "lng"
	fmt.Println(bfs(g, start1), "*", bfs(g, start2))
	return bfs(g, start1) * bfs(g, start2)
}

func bfs(g Graph, start string) int {
	visited := map[string]bool{start: true}
	queue := []string{start}

	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]

		for s := range g[next] {
			if visited[s] {
				continue
			}

			queue = append(queue, s)
			visited[s] = true
		}
	}

	return len(visited)
}

func karger(g Graph) int {
	graph := make(Graph)
	for name, neighbors := range g {
		graph[name] = make(map[string]int)
		for n := range neighbors {
			graph[name][n] = 0
		}
	}

	for len(graph) > 2 {
		var from, to string
		i := rand.Intn(len(graph))
		for name, neighbors := range graph {
			from = name
			j := rand.Intn(len(graph))
			for n := range neighbors {
				to = n
				if j == 0 {
					break
				}
				j--
			}

			if i == 0 {
				break
			}
			i--
		}

		var neighbors []string
		for s := range graph[from] {
			if s == from || s == to {
				continue
			}
			neighbors = append(neighbors, s)
		}
		for s := range graph[to] {
			if s == from || s == to {
				continue
			}
			neighbors = append(neighbors, s)
		}

		for s := range graph {
			if s == from || s == to {
				continue
			}

			delete(graph[s], from)
			delete(graph[s], to)
		}

		delete(graph, from)
		delete(graph, to)
		graph[from+to] = make(map[string]int)
		for _, n := range neighbors {
			graph[from+to][n] = 0
			graph[n][from+to] = 0
		}
	}

	var v1, v2 string
	for s1, neighbors := range graph {
		v1 = s1
		for s2 := range neighbors {
			v2 = s2
			break
		}

		break
	}

	fmt.Println(v1, v2, len(v1)/3*len(v2)/3)
	return len(v1) / 3 * len(v2) / 3
}
