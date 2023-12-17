package day17

import (
	"container/heap"
	"fmt"

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
	fmt.Printf("first: %d\n", djikstra(m, crucibleDirections))
	fmt.Printf("second: %d\n", djikstra(m, ultraDirections))
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

func djikstra(m maps.Map[int], getDirections GetDirection) int {
	start := maps.Coordinate{}
	goal := maps.Coordinate{X: m.Columns - 1, Y: m.Rows - 1}
	distances := make(map[key]int)
	var goalItem item
	from := map[key]key{}

	var pq queue.PriorityQueue[item]
	heap.Init(&pq)
	heap.Push(&pq, &queue.Item[item]{Value: item{c: start, cost: 0}, Priority: 0})

	for len(pq) > 0 {
		n := heap.Pop(&pq).(*queue.Item[item])
		val := n.Value

		directions := getDirections(val)
		for _, dir := range directions {
			if val.direction != maps.NoDirection && dir == val.direction.Opposite() {
				continue
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

			if a == goal {
				goalItem = i
				pq = nil
			}

			func(prev, i item) {
				if c, ok := distances[i.key()]; !ok || i.cost < c {
					from[i.key()] = prev.key()
					distances[i.key()] = i.cost
					heap.Push(&pq, &queue.Item[item]{Value: i, Priority: i.cost})
				}
			}(val, i)
		}
	}

	return backtrack(m, start, goalItem, from)
}

type GetDirection func(at item) []maps.Direction

func crucibleDirections(at item) []maps.Direction {
	var directions []maps.Direction
	switch {
	case at.consecutive >= 3:
		directions = fns.Filter([]maps.Direction{maps.Up, maps.Left, maps.Right, maps.Down}, func(d maps.Direction) bool {
			return d != at.direction
		})

	default:
		directions = []maps.Direction{maps.Up, maps.Left, maps.Right, maps.Down}
	}

	return directions
}

func ultraDirections(at item) []maps.Direction {
	var directions []maps.Direction
	switch {
	case at.consecutive < 4:
		directions = []maps.Direction{at.direction}

	case at.consecutive >= 10:
		directions = fns.Filter([]maps.Direction{maps.Up, maps.Left, maps.Right, maps.Down}, func(d maps.Direction) bool {
			return d != at.direction
		})

	default:
		directions = []maps.Direction{maps.Up, maps.Left, maps.Right, maps.Down}
	}

	return directions
}

func backtrack(m maps.Map[int], start maps.Coordinate, goal item, from map[key]key) int {
	path := make(map[maps.Coordinate]maps.Direction)
	var cost int
	for at := goal.key(); at.c != start; {
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

	return cost
}
