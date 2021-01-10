package parser

import (
	"fmt"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func Draw(s *State, path string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = InvalidRegex
		}
	}()

	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		return err
	}
	defer func() {
		_ = graph.Close()
		_ = g.Close()
	}()
	graph.SetRankDir(cgraph.LRRank)

	visited := make(map[*State]*cgraph.Node)
	var build func(node *cgraph.Node, s *State, char rune) error

	idx := 0
	build = func(node *cgraph.Node, s *State, char rune) error {
		if s == nil {
			return nil
		}
		newNode, seen := visited[s]
		if !seen {
			n, err := graph.CreateNode(fmt.Sprintf("S%d", idx))
			if err != nil {
				return err
			}
			visited[s] = n
			idx++
			newNode = n

			if s.char == CharMatch {
				n.SetShape(cgraph.DoubleCircleShape)
			} else {
				n.SetShape(cgraph.CircleShape)
			}
		}

		edge, err := graph.CreateEdge("e", node, newNode)
		if err != nil {
			return err
		}
		if char >= 0 {
			edge.SetLabel(string(char))
		}

		if !seen {
			if err := build(newNode, s.out, s.char); err != nil {
				return err
			}
			if err := build(newNode, s.out2, s.char); err != nil {
				return err
			}
		}

		return nil
	}

	start, err := graph.CreateNode("")
	if err != nil {
		return err
	}
	start.SetStyle("invis")
	if err := build(start, s, -1); err != nil {
		return err
	}

	if err := g.RenderFilename(graph, graphviz.PNG, path); err != nil {
		return err
	}
	return nil
}
