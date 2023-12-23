package day20

import (
	"fmt"
	"strings"

	"github.com/mbark/aoc2023/fns"
	"github.com/mbark/aoc2023/maths"
	"github.com/mbark/aoc2023/util"
)

const testInput = `
broadcaster -> a
%a -> inv, con
&inv -> b
%b -> con
&con -> output
`

func Day20(input string) {
	if input == "" {
		input = testInput
	}

	mapping := make(map[string][]string)
	var modules1 []Module
	var modules2 []Module
	for _, s := range util.ReadInput(input, "\n") {
		split := strings.Split(s, " ")

		var toModules []string
		for _, to := range strings.Split(strings.Split(s, " -> ")[1], ", ") {
			toModules = append(toModules, to)
		}

		var module1, module2 Module
		switch {
		case split[0] == "broadcaster":
			module1 = Module{name: split[0], receiver: new(Broadcaster)}
			module2 = Module{name: split[0], receiver: new(Broadcaster)}
		case split[0][0] == '%':
			module1 = Module{name: split[0][1:], receiver: new(FlipFlop)}
			module2 = Module{name: split[0][1:], receiver: new(FlipFlop)}
		case split[0][0] == '&':
			module1 = Module{name: split[0][1:], receiver: new(Conjunction)}
			module2 = Module{name: split[0][1:], receiver: new(Conjunction)}
		}

		modules1 = append(modules1, module1)
		modules2 = append(modules2, module2)
		mapping[module1.name] = toModules
	}

	for i, module := range modules1 {
		c, ok := module.receiver.(*Conjunction)
		if !ok {
			continue
		}
		c.memory = make(map[string]Pulse)

		for j, other := range modules1 {
			if i == j {
				continue
			}
			_, ok := fns.Find(mapping[other.name], func(s string) bool { return s == module.name })
			if ok {
				c.memory[other.name] = Low
			}
		}
	}
	for i, module := range modules2 {
		c, ok := module.receiver.(*Conjunction)
		if !ok {
			continue
		}
		c.memory = make(map[string]Pulse)

		for j, other := range modules2 {
			if i == j {
				continue
			}
			_, ok := fns.Find(mapping[other.name], func(s string) bool { return s == module.name })
			if ok {
				c.memory[other.name] = Low
			}
		}
	}

	byName1 := fns.Associate(modules1, func(m Module) string { return m.name })
	byName2 := fns.Associate(modules2, func(m Module) string { return m.name })
	fmt.Printf("first: %d\n", first(byName1, mapping))
	fmt.Printf("second: %d\n", second(byName2, mapping))
}

func first(modules map[string]Module, mapping map[string][]string) int {
	var high, low int
	for i := 0; i < 1000; i++ {
		sent := []Send{{from: "button", pulse: Low, to: "broadcaster"}}
		for len(sent) > 0 {
			s := sent[0]
			sent = sent[1:]

			switch s.pulse {
			case Low:
				low++
			case High:
				high++
			}

			next, ok := modules[s.to]
			if !ok {
				continue
			}

			from := Send{pulse: s.pulse, from: s.from, to: s.to}
			p, ok := next.receiver.recv(from)
			if !ok {
				continue
			}

			for _, to := range mapping[s.to] {
				sent = append(sent, Send{from: s.to, pulse: p, to: to})
			}
		}
	}

	return high * low
}

func second(modules map[string]Module, mapping map[string][]string) int {
	for i := 0; i < 100000; i++ {
		sent := []Send{{from: "button", pulse: Low, to: "broadcaster"}}

		for len(sent) > 0 {
			s := sent[0]
			sent = sent[1:]

			next, ok := modules[s.to]
			if !ok {
				continue
			}

			from := Send{count: i, pulse: s.pulse, from: s.from, to: s.to}
			p, ok := next.receiver.recv(from)
			if !ok {
				continue
			}

			for _, to := range mapping[s.to] {
				sent = append(sent, Send{from: s.to, pulse: p, to: to})
			}
		}
	}

	at := make(map[string]int)
	for s, module := range modules {
		c, ok := module.receiver.(*Conjunction)
		if ok {
			at[s] = c.onAtNext - c.onAt
			fmt.Println(s, "at", c.onAtNext, c.onAt, c.onAtNext-c.onAt)
		}
	}

	fmt.Println(maths.LCM(at["nb"], at["vc"], at["ls"], at["vg"]))

	return 0
}

type (
	Pulse bool
	Send  struct {
		count int
		pulse Pulse
		to    string
		from  string
	}
)

func (f Send) String() string {
	return fmt.Sprintf("%s %s %s", f.from, f.pulse, f.to)
}

func (p Pulse) String() string {
	switch p {
	case true:
		return "-high->"
	default:
		return "-low->"
	}
}

const (
	Low  Pulse = false
	High Pulse = true
)

type Module struct {
	name     string
	receiver Receiver
}

func (m Module) String() string {
	return fmt.Sprintf("%s(%s)", m.name, m.receiver)
}

type Receiver interface {
	recv(pulse Send) (Pulse, bool)
	String() string
}

var (
	_ Receiver = new(FlipFlop)
	_ Receiver = new(Conjunction)
	_ Receiver = new(Broadcaster)
	_ Receiver = new(Output)
)

type (
	FlipFlop    struct{ on bool }
	Conjunction struct {
		onAt     int
		onAtNext int
		memory   map[string]Pulse
	}
	Broadcaster struct{}
	Output      struct{}
)

func (b *Broadcaster) recv(pulse Send) (Pulse, bool) {
	return pulse.pulse, true
}

func (b *Broadcaster) String() string {
	return ""
}

func (c *Conjunction) recv(pulse Send) (Pulse, bool) {
	c.memory[pulse.from] = pulse.pulse
	if fns.EveryMap(c.memory, func(key string, val Pulse) bool { return val == High }) {
		return Low, true
	}

	if c.onAt == 0 {
		c.onAt = pulse.count
	} else if c.onAtNext == 0 {
		c.onAtNext = pulse.count
	}

	fmt.Println(pulse.to, pulse.count)
	return High, true
}

func (c *Conjunction) String() string {
	var ss []string
	for s, pulse := range c.memory {
		ss = append(ss, fmt.Sprintf("%s=%s", s, pulse))
	}
	return strings.Join(ss, ",")
}

func (f *FlipFlop) recv(pulse Send) (Pulse, bool) {
	if pulse.pulse == High {
		return Low, false
	}

	switch f.on {
	case true:
		f.on = false
		return Low, true
	default:
		f.on = true
		return High, true
	}
}

func (f *FlipFlop) String() string {
	return "%off"
}

func (f *Output) recv(pulse Send) (Pulse, bool) {
	return Low, false
}

func (f *Output) String() string {
	return "-"
}
