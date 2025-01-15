package main

import (
	"fmt"
	"slices"
	"strings"
	_ "embed"
)

//go:embed test.txt
var testInput string

//go:embed input.txt
var rawInput string

type Node struct {
	Name string
	Edges []*Node
}

func (n *Node) String() string {
	return n.Name
}

func ParseInput(input []string) ([]*Node, map[string]*Node, map[string]bool) {
	nodes := make([]*Node, 0)
	graph := make(map[string]*Node)
	edges := make(map[string]bool)
	
	for _, line := range input {
		if line == "" {
			continue
		}
		names := strings.Split(line, "-")
		leftNodeName, rightNodeName := names[0], names[1]
		var leftNode, rightNode *Node

		if _, ok := graph[leftNodeName]; !ok {
			leftNode = &Node{Name: leftNodeName, Edges: make([]*Node, 0)}
			graph[leftNodeName] = leftNode
			nodes = append(nodes, leftNode)
		} else {
			leftNode = graph[leftNodeName]
		}

		if _, ok := graph[rightNodeName]; !ok {
			rightNode = &Node{Name: rightNodeName, Edges: make([]*Node, 0)}
			graph[rightNodeName] = rightNode
			nodes = append(nodes, rightNode)
		} else {
			rightNode = graph[rightNodeName]
		}
		slices.Sort(names)

		if _, ok := edges[strings.Join(names, "-")]; !ok {
			leftNode.Edges = append(leftNode.Edges, rightNode)
			rightNode.Edges = append(rightNode.Edges, leftNode)
		}
		edges[strings.Join(names, "-")] = true
	}
	
	return nodes, graph, edges
}

func countCycles(seen, edges map[string]bool, start *Node, depth int) int {
	if depth < 2 {
		return len(start.Edges)
	}
	if depth == 3 {
		count := 0
		for i, l := range start.Edges {
			for _, r := range start.Edges[i + 1:] {
				var connection string
				names := []string{l.Name, r.Name}
				slices.Sort(names)
				connection = strings.Join(names, "-")
				if _, ok := edges[connection]; ok {
					names = append(names, start.Name)
					slices.Sort(names)
					connection = strings.Join(names, "-")
					if _, ok := seen[connection]; !ok {
						count++
						seen[connection] = true
					}
				}
			}
		}
		
	}

	panic("Not implemented")
}

func main() {
	input := strings.Split(rawInput, "\n")
	nodes, _, edges := ParseInput(input)

	sum := 0
	seen := make(map[string]bool)
	for _, node := range nodes {
		if node.Name[0] == 't' {
			fmt.Println(node)
			cnt := countCycles(seen, edges, node, 3)
			fmt.Println(node.Name, "Count:", cnt)
			sum += cnt
		}
	}
	fmt.Println("Part 1:", sum)

	fmt.
}
