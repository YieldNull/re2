package main

import (
	"fmt"
	"os"

	"github.com/yieldnull/re2/parser"
)

func main() {
	s, err := parser.Compile(os.Args[1])
	if err != nil {
		panic(err)
	}

	if err := parser.Draw(s, "nfa.png"); err != nil {
		panic(err)
	}

	for _, p := range os.Args[2:] {
		if s.Match(p) {
			fmt.Println(p)
		} else {
			_, _ = fmt.Fprintln(os.Stderr, "Not Match")
		}
	}
}
