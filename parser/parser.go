package parser

import (
	"errors"
)

var InvalidRegex = errors.New("invalid regex")

var priority = map[rune]int{
	'*': 4,
	'+': 4,
	'?': 4,
	'.': 3,
	'|': 2,
	'(': 1,
}

func transform(re string) string {
	if len(re) == 0 {
		return re
	}

	runes := []rune(re)
	chars := make([]rune, 0, len(runes)*2)

	chars = append(chars, runes[0])
	for i, c := range runes[1:] {
		switch c {
		case '|', ')', '*', '+', '?':
			chars = append(chars, c)
		default: // '(' and litters
			switch runes[i] { // check before
			case '(', '|':
				chars = append(chars, c)
			default:
				chars = append(chars, '.', c)
			}
		}
	}

	return string(chars)
}

func re2post(re string) (post string, err error) {
	chars := make([]rune, 0, len(re))
	var operator Stack = make([]rune, 0, len(re))

	for _, c := range re {
		switch c {
		case '+', '*', '?', '|', '.':
			for {
				if op, ok := operator.Peek(); ok && priority[op] >= priority[c] {
					op, _ = operator.Pop()
					chars = append(chars, op)
				} else {
					break
				}
			}
			operator.Push(c)
		case '(':
			operator.Push(c)
		case ')':
			for {
				op, e := operator.Pop()
				if e != nil {
					err = InvalidRegex
					return
				}

				if op == '(' {
					break
				}

				chars = append(chars, op)
			}
		default:
			chars = append(chars, c)
		}
	}

	for !operator.Empty() {
		op, _ := operator.Pop()
		if op == '(' {
			err = InvalidRegex
			return
		}
		chars = append(chars, op)
	}

	return string(chars), nil
}
