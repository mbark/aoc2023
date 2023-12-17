package day17

import (
	"container/heap"
	"fmt"
	"math"
	"strconv"

	"github.com/mbark/aoc2023/fns"
	"github.com/mbark/aoc2023/maps"
	"github.com/mbark/aoc2023/queue"
	"github.com/mbark/aoc2023/util"
)

const testInput = `
2413432311323
3215453535623
3255245654254
3446585845452
4546657867536
1438598798454
4457876987766
3637877979653
4654967986887
4564679986453
1224686865563
2546548887735
4322674655533
`

func Day17(input string) {
	if input == "" {
		input = testInput
	}

	m := maps.New(input, func(_, _ int, b byte) int { return util.ParseInt[int](string(b)) })
	fmt.Printf("first: %d\n", first(m))
	fmt.Printf("second: %d\n", second(m))
}

type key struct {
	c           maps.Coordinate
	consecutive int
	direction   maps.Direction
}

func (k key) String() string {
	return fmt.Sprintf("%s %s (%d)", k.c, k.direction, k.consecutive)
}

type item struct {
	c           maps.Coordinate
	cost        int
	direction   maps.Direction
	consecutive int
}

func (i item) String() string {
	return fmt.Sprintf("{c: %s, cost: %d, dir: %s, consecutive: %d}",
		i.c, i.cost, i.direction, i.consecutive)
}

func (i item) key() key {
	return key{
		c:           i.c,
		consecutive: i.consecutive,
		direction:   i.direction,
	}
}

func first(m maps.Map[int]) int {
	start := maps.Coordinate{}
	goal := maps.Coordinate{X: m.Columns - 1, Y: m.Rows - 1}

	distances := make(map[key]int)
	var pq queue.PriorityQueue[item]
	heap.Init(&pq)
	heap.Push(&pq, &queue.Item[item]{Value: item{c: start, cost: 0}, Priority: 0})
	doneAt := math.MaxInt
	var goalItem item

	from := map[key]key{}

	fastest := djikstra(m)

	pushIf := func(prev, i item) {
		f := fastest[i.c]
		if i.cost+f > doneAt {
			return
		}

		if c, ok := distances[i.key()]; !ok || i.cost < c {
			from[i.key()] = prev.key()
			distances[i.key()] = i.cost
			heap.Push(&pq, &queue.Item[item]{Value: i, Priority: 0})
		}
	}

	for len(pq) > 0 {
		n := heap.Pop(&pq).(*queue.Item[item])
		val := n.Value

		for _, dir := range []maps.Direction{maps.Up, maps.Left, maps.Right, maps.Down} {
			switch val.direction {
			case maps.Up:
				if dir == maps.Down {
					continue
				}
			case maps.Down:
				if dir == maps.Up {
					continue
				}
			case maps.Left:
				if dir == maps.Right {
					continue
				}
			case maps.Right:
				if dir == maps.Left {
					continue
				}
			}

			a := dir.Apply(val.c)
			if !m.Exists(a) {
				continue
			}

			i := item{c: a, cost: val.cost + m.At(a), direction: dir}
			switch dir {
			case val.direction:
				i.consecutive = val.consecutive + 1
			default:
				i.consecutive = 1
			}
			if i.consecutive > 3 {
				continue
			}

			if a == goal && i.cost < doneAt {
				doneAt = i.cost
				goalItem = i
			}

			pushIf(val, i)
		}
	}

	path := make(map[maps.Coordinate]maps.Direction)
	var cost int
	for at := goalItem.key(); at.c != start; {
		cost += m.At(at.c)
		if _, ok := path[at.c]; ok {
			break
		}

		prev := from[at]
		path[at.c] = maps.Direction{
			X: at.c.X - prev.c.X,
			Y: at.c.Y - prev.c.Y,
		}
		at = prev
	}

	fmt.Println(m.Stringf(func(c maps.Coordinate, val int) string {
		if dir, ok := path[c]; ok {
			return dir.String()
		}

		return strconv.Itoa(val)
	}))

	return cost
}

func second(m maps.Map[int]) int {
	start := maps.Coordinate{}
	goal := maps.Coordinate{X: m.Columns - 1, Y: m.Rows - 1}

	distances := make(map[key]int)
	var pq queue.PriorityQueue[item]
	heap.Init(&pq)
	heap.Push(&pq, &queue.Item[item]{Value: item{c: start, cost: 0}, Priority: 0})
	doneAt := math.MaxInt
	var goalItem item

	from := map[key]key{}

	fastest := djikstra(m)

	pushIf := func(prev, i item) {
		f := fastest[i.c]
		if i.cost+f > doneAt {
			return
		}

		if c, ok := distances[i.key()]; !ok || i.cost < c {
			from[i.key()] = prev.key()
			distances[i.key()] = i.cost
			heap.Push(&pq, &queue.Item[item]{Value: i, Priority: 0})
		}
	}

	for len(pq) > 0 {
		n := heap.Pop(&pq).(*queue.Item[item])
		val := n.Value

		var directions []maps.Direction
		switch {
		case val.consecutive < 4:
			directions = []maps.Direction{val.direction}

		case val.consecutive >= 10:
			directions = fns.Filter([]maps.Direction{maps.Up, maps.Left, maps.Right, maps.Down}, func(d maps.Direction) bool {
				return d != val.direction
			})

		default:
			directions = []maps.Direction{maps.Up, maps.Left, maps.Right, maps.Down}
		}

		for _, dir := range directions {
			switch val.direction {
			case maps.Up:
				if dir == maps.Down {
					continue
				}
			case maps.Down:
				if dir == maps.Up {
					continue
				}
			case maps.Left:
				if dir == maps.Right {
					continue
				}
			case maps.Right:
				if dir == maps.Left {
					continue
				}
			}
			a := dir.Apply(val.c)
			if !m.Exists(a) {
				continue
			}

			i := item{c: a, cost: val.cost + m.At(a), direction: dir}
			switch dir {
			case val.direction:
				i.consecutive = val.consecutive + 1
			default:
				i.consecutive = 1
			}

			if i.consecutive >= 4 && a == goal && i.cost < doneAt {
				doneAt = i.cost
				goalItem = i
			}

			pushIf(val, i)
		}
	}

	path := make(map[maps.Coordinate]maps.Direction)
	var cost int
	for at := goalItem.key(); at.c != start; {
		cost += m.At(at.c)
		if _, ok := path[at.c]; ok {
			break
		}

		prev := from[at]
		path[at.c] = maps.Direction{
			X: at.c.X - prev.c.X,
			Y: at.c.Y - prev.c.Y,
		}
		at = prev
	}

	fmt.Println(m.Stringf(func(c maps.Coordinate, val int) string {
		if dir, ok := path[c]; ok {
			return dir.String()
		}

		return strconv.Itoa(val)
	}))

	return cost
}

func djikstra(m maps.Map[int]) map[maps.Coordinate]int {
	start := maps.Coordinate{X: m.Columns - 1, Y: m.Rows - 1}
	distances := map[maps.Coordinate]int{start: m.At(start)}
	visited := map[maps.Coordinate]bool{start: true}

	pq := make(queue.PriorityQueue[maps.Coordinate], len(m.Coordinates()))
	for i, c := range m.Coordinates() {
		prio := math.MaxInt
		if c == start {
			prio = 0
		}

		distances[c] = prio
		pq[i] = &queue.Item[maps.Coordinate]{Value: c, Priority: prio, Index: i}
	}
	heap.Init(&pq)

	for len(pq) > 0 {
		n := heap.Pop(&pq).(*queue.Item[maps.Coordinate])
		c := n.Value
		visited[c] = true

		for _, a := range m.Adjacent(c) {
			if _, ok := visited[a]; ok {
				continue
			}

			alt := distances[c] + m.At(a)
			if d, ok := distances[a]; !ok || alt < d {
				distances[a] = alt
				heap.Push(&pq, &queue.Item[maps.Coordinate]{Value: a, Priority: alt})
			}
		}
	}

	return distances
}
