package day19

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2023/fns"
	"github.com/mbark/aoc2023/util"
)

const testInput = `
px{a<2006:qkq,m>2090:A,rfg}
pv{a>1716:R,A}
lnx{m>1548:A,A}
rfg{s<537:gd,x>2440:R,A}
qs{s>3448:A,lnx}
qkq{x<1416:A,crn}
crn{x>2662:A,R}
in{s<1351:px,qqz}
qqz{s>2770:qs,m<1801:hdj,R}
gd{a>3333:R,R}
hdj{m>838:A,pv}

{x=787,m=2655,a=1222,s=2876}
{x=1679,m=44,a=2067,s=496}
{x=2036,m=264,a=79,s=2244}
{x=2461,m=1339,a=466,s=291}
{x=2127,m=1623,a=2188,s=1013}
`

func Day19(input string) {
	if input == "" {
		input = testInput
	}

	var workflows []workflow
	in := util.ReadInput(input, "\n\n")

	for _, s := range strings.Split(in[0], "\n") {
		s = strings.TrimSuffix(s, "}")
		split := strings.Split(s, "{")
		name := split[0]

		var rules []rule
		splits := strings.Split(split[1], ",")
		fallback := splits[len(splits)-1]

		for _, s := range splits[:len(splits)-1] {
			rating := string(s[0])
			var fn Comparison
			comparer := s[1]

			s = s[1:]
			split := strings.Split(s[1:], ":")
			b := util.Str2Int(split[0])
			switch comparer {
			case '<':
				fn = lessThan{b: b}
			case '>':
				fn = moreThan{b: b}
			}

			rules = append(rules, rule{
				to:     split[1],
				rating: rating,
				fn:     fn,
				b:      b,
			})
		}

		workflows = append(workflows, workflow{
			name:     name,
			fallback: fallback,
			rules:    rules,
		})
	}

	var parts []part
	for _, s := range strings.Split(in[1], "\n") {
		p := part{ratings: make(map[string]int)}
		s = strings.Trim(s, "{}")
		for _, s := range strings.Split(s, ",") {
			split := strings.Split(s, "=")
			p.ratings[split[0]] = util.Str2Int(split[1])
		}

		parts = append(parts, p)
	}

	fmt.Printf("first: %d\n", first(workflows, parts))
	fmt.Printf("second: %d\n", second(workflows))
}

func first(workflows []workflow, parts []part) int {
	byName := fns.Associate(workflows, func(w workflow) string {
		return w.name
	})
	byName["A"] = workflow{name: "A"}
	byName["R"] = workflow{name: "R"}

	var sum int
	for _, p := range parts {
		accepted := resolve(byName, byName["in"], p)
		if accepted {
			for _, i := range p.ratings {
				sum += i
			}
		}
	}
	return sum
}

func resolve(workflows map[string]workflow, wf workflow, p part) bool {
	if wf.name == "A" {
		return true
	}
	if wf.name == "R" {
		return false
	}

	for _, r := range wf.rules {
		rating := p.ratings[r.rating]
		if r.fn.Cmp(rating) {
			return resolve(workflows, workflows[r.to], p)
		}
	}

	return resolve(workflows, workflows[wf.fallback], p)
}

func second(workflows []workflow) int {
	byName := fns.Associate(workflows, func(w workflow) string {
		return w.name
	})
	byName["A"] = workflow{name: "A"}
	byName["R"] = workflow{name: "R"}

	ranges := resolve2(byName, byName["in"], []PartRange{
		{ranges: map[string]Range{
			"x": {lower: 1, upper: 4000},
			"m": {lower: 1, upper: 4000},
			"a": {lower: 1, upper: 4000},
			"s": {lower: 1, upper: 4000},
		}},
	})

	var sum int
	for _, pr := range ranges {
		parts := 1
		for _, r := range pr.ranges {
			parts *= r.upper - r.lower + 1
		}

		sum += parts
	}

	return sum
}

func resolve2(workflows map[string]workflow, wf workflow, ranges []PartRange) []PartRange {
	if wf.name == "A" {
		return ranges
	}
	if wf.name == "R" {
		return nil
	}

	var added []PartRange
	for _, rule := range wf.rules {
		var successes []PartRange
		var failed []PartRange
		for _, r := range ranges {
			truthy, falsy := rule.fn.Split(rule.rating, r)
			if truthy != nil {
				successes = append(successes, *truthy)
			}
			if falsy != nil {
				failed = append(failed, *falsy)
			}
		}

		added = append(added, resolve2(workflows, workflows[rule.to], successes)...)
		ranges = failed
	}

	return append(added, resolve2(workflows, workflows[wf.fallback], ranges)...)
}

type workflow struct {
	name     string
	fallback string
	rules    []rule
}

func (w workflow) String() string {
	return fmt.Sprintf("%s{%s,%s}", w.name, strings.Join(
		fns.Map(w.rules, func(r rule) string {
			return r.String()
		}), ","),
		w.fallback,
	)
}

type rule struct {
	to     string
	rating string
	fn     Comparison
	b      int
}

func (r rule) String() string {
	var fn string
	switch r.fn.(type) {
	case lessThan:
		fn = "<"
	case moreThan:
		fn = ">"
	}

	return fmt.Sprintf("%s%s%d:%s", r.rating, fn, r.b, r.to)
}

type PartRange struct {
	ranges map[string]Range
}

func (pr PartRange) String() string {
	var ss []string
	for s, r := range pr.ranges {
		ss = append(ss, fmt.Sprintf("%s:%d-%d", s, r.lower, r.upper))
	}

	return fmt.Sprintf("{%s}", strings.Join(ss, ","))
}

func (pr PartRange) copyWith(rating string, r Range) *PartRange {
	if !r.isValid() {
		return nil
	}

	ranges := make(map[string]Range, len(pr.ranges))
	for k, v := range pr.ranges {
		ranges[k] = v
	}
	ranges[rating] = r

	return &PartRange{ranges: ranges}
}

type Range struct {
	lower int
	upper int
}

func (r *Range) isValid() bool {
	return r != nil && r.upper >= r.lower
}

type Comparison interface {
	Cmp(a int) bool
	Split(rating string, r PartRange) (*PartRange, *PartRange) // true, false
}

var (
	_ Comparison = lessThan{}
	_ Comparison = moreThan{}
)

type (
	lessThan struct{ b int }
	moreThan struct{ b int }
)

func (l lessThan) Cmp(a int) bool {
	return a < l.b
}

func (l lessThan) Split(rating string, pr PartRange) (*PartRange, *PartRange) {
	r := pr.ranges[rating]
	switch {
	case r.upper <= l.b:
		return &pr, nil
	case r.lower >= l.b:
		return nil, &pr
	default:
		return pr.copyWith(rating, Range{lower: r.lower, upper: l.b - 1}),
			pr.copyWith(rating, Range{lower: l.b, upper: r.upper})
	}
}

func (m moreThan) Cmp(a int) bool {
	return a > m.b
}

func (m moreThan) Split(rating string, pr PartRange) (*PartRange, *PartRange) {
	r := pr.ranges[rating]
	switch {
	case r.lower >= m.b:
		return &pr, nil
	case r.upper <= m.b:
		return nil, &pr
	default:
		return pr.copyWith(rating, Range{lower: m.b + 1, upper: r.upper}),
			pr.copyWith(rating, Range{lower: r.lower, upper: m.b})
	}
}

type part struct {
	ratings map[string]int
}

func (p part) String() string {
	var ss []string
	for s, i := range p.ratings {
		ss = append(ss, fmt.Sprintf("%s=%d", s, i))
	}

	return fmt.Sprintf("{%s}", strings.Join(ss, ","))
}
