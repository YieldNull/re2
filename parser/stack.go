package parser

import "errors"

type Stack []rune

func (s *Stack) Push(c rune) {
	*s = append(*s, c)
}

func (s *Stack) Pop() (c rune, err error) {
	if s.Empty() {
		err = errors.New("empty stack")
		return
	}

	c = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	return
}

func (s *Stack) Peek() (c rune, ok bool) {
	if s.Empty() {
		return
	}

	return (*s)[len(*s)-1], true
}

func (s *Stack) Empty() bool {
	return len(*s) == 0
}
