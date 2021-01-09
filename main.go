package main

import (
	"os"

	"github.com/yieldnull/re2/parser"
)

func main() {
	_, err := parser.Compile(os.Args[1])
	if err != nil {
		panic(err)
	}
}
