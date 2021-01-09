package parser

import (
	"errors"
)

type State struct {
	char rune
	out  *State
	out2 *State
}

const (
	TypeSplit rune = -1
	TypeMatch rune = -2
)

func Char(char rune, out *State) *State {
	return &State{char: char, out: out}
}

func Split(out *State, out2 *State) *State {
	return &State{char: TypeSplit, out: out, out2: out2}
}

func Match() *State {
	return &State{char: TypeMatch}
}

type Frag struct {
	start *State
	outs  []**State
}

type fragStack []*Frag

func (s *fragStack) Push(c *Frag) {
	*s = append(*s, c)
}

func (s *fragStack) Pop() (c *Frag, err error) {
	if len(*s) == 0 {
		err = errors.New("empty stack")
		return
	}

	c = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	return
}

func post2nfa(re string) (_ *State, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = InvalidRegex
		}
	}()

	var stack fragStack = make([]*Frag, 0, 2)

	for _, char := range re {
		switch char {
		default:
			//     a
			// O ---->
			s := Char(char, nil)
			stack.Push(&Frag{
				start: s,
				outs:  []**State{&s.out},
			})
		case '.':
			//      +----------+      +----------+
			// ---> |---Frag---| ---> |---Frag---| --->
			//      +----------+      +----------+
			f2, _ := stack.Pop()
			f1, _ := stack.Pop()
			for _, o := range f1.outs {
				*o = f2.start
			}

			stack.Push(&Frag{
				start: f1.start,
				outs:  f2.outs,
			})
		case '?':
			//            +----------+
			// --> O ---> |---Frag---| --->
			//     |      +----------+
			//     |
			//     | --------------------->
			//
			f, _ := stack.Pop()
			split := Split(f.start, nil)

			stack.Push(&Frag{
				start: split,
				outs:  append(f.outs, &split.out2),
			})
		case '+':
			//      +----------+
			// ---> |---Frag---| ---> O ------>
			//      +----------+      |
			//            ↑           |
			//            |-----------|
			f, _ := stack.Pop()
			s := Split(f.start, nil)
			for _, o := range f.outs {
				*o = s
			}

			stack.Push(&Frag{
				start: f.start,
				outs:  []**State{&s.out2},
			})
		case '*':
			//     |--------------------- |
			//     |                      |
			//     V      +----------+    |
			// --> O ---> |---Frag---| ---|
			//     |      +----------+
			//     |
			//     |---------------------->
			f, _ := stack.Pop()
			s := Split(f.start, nil)
			for _, o := range f.outs {
				*o = s
			}

			stack.Push(&Frag{
				start: s,
				outs:  []**State{&s.out2},
			})
		case '|':
			//            +----------+
			//     |----> |---Frag---| ---->
			//     |      +----------+
			// --> O
			//     |      +----------+
			//     |----> |---Frag---| ---->
			//            +----------+
			f1, _ := stack.Pop()
			f2, _ := stack.Pop()
			s := Split(f1.start, f2.start)

			stack.Push(&Frag{
				start: s,
				outs:  append(f1.outs, f2.outs...),
			})
		}
	}

	f, _ := stack.Pop()
	m := Match()
	for _, o := range f.outs {
		*o = m
	}
	return f.start, nil
}
