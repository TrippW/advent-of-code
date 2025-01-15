package main

import (
	_ "embed"
	"fmt"
	"math"
	"slices"
	"sort"
	"strings"

	"github.com/TrippW/advent-of-code/utils"
	_ "github.com/TrippW/advent-of-code/utils"
)

//go:embed test.txt
var rawTest string

//go:embed input.txt
var rawInput string

const (
	Wall rune = '#'
	Empty rune = '.'
	Start rune = 'S'
	End rune = 'E'
)

type Coordinate2D struct {
	X, Y int
}

type Direction int

const (
	North Direction = iota
	East
	South
	West
	None
)

func (c Coordinate2D) Add(d Direction) Coordinate2D {
	switch d {
	case North:
		return Coordinate2D{c.X, c.Y - 1}
	case East:
		return Coordinate2D{c.X + 1, c.Y}
	case South:
		return Coordinate2D{c.X, c.Y + 1}
	case West:
		return Coordinate2D{c.X - 1, c.Y}
	}
	panic("Invalid direction")
}

type Runner struct {
	Start *Node
	Direction 
}

type Node struct {
	position Coordinate2D
	connections map[Direction]*Node
	f, g, h float64
	ApproachDirection Direction
	parent *Node
	next *Node
}

func (n Node) String() string {
	return fmt.Sprintf("Node{Position: %v, f: %d, g: %d, h: %d, ApproachDirection: %v}", n.position, int(n.f), int(n.g), int(n.h), n.ApproachDirection)
}

func (r Runner) RunAStar(end *Node) (int, [][]Coordinate2D) {
	// Swap with heap
	open := []*Node{r.Start}

	// Swap with map[Coordinate2D]bool
	closed := []*Node{}

	paths := [][]Coordinate2D{}

	bestScore := math.Inf(0)
	//Create Materializedd Map
	for len(open) > 0 {
		var current *Node
		sort.Slice(open, func(i, j int) bool {
			return open[i].f + 0 < open[j].f
		})
		current, open = open[0], open[1:]
		for _, conn := range current.connections {
			if conn == nil {
				fmt.Println("Skipping nil connection")
				continue
			}
			if conn == end {
				conn.parent = current
				current.next = conn
				conn.h = 0
				if current.position.Add(current.ApproachDirection) == conn.position {
					conn.g = current.g + 1
				} else {
					conn.g = current.g + 1001
				}
				conn.f = conn.g + conn.h
				if conn.f <= bestScore {
					if conn.f < bestScore {
						paths = [][]Coordinate2D{}
					}

					bestScore = conn.f
					path := []Coordinate2D{}
					for parent := conn.parent; parent != nil; parent = parent.parent {
						path = append(path, parent.position)
					}
					paths = append(paths, path)
				}
				fmt.Println("Found end", conn.position, "with score", conn.f)
				continue
			}

			var g, h, f float64
			var a Direction
			if current.position.Add(current.ApproachDirection) == conn.position {
				g = current.g + 1
				a = current.ApproachDirection
			} else {
				g = current.g + 1001
				if current.position.Add(North) == conn.position {
					a = North
				} else if current.position.Add(East) == conn.position {
					a = East
				} else if current.position.Add(South) == conn.position {
					a = South
				} else if current.position.Add(West) == conn.position {
					a = West
				} else {
					panic("Invalid direction")
				}
			}
			h = float64(FindDistance(conn, end))
			f = g + h

//			if i := slices.Index(open, conn); i != -1 && open[i].f < f {
//				continue
//			}
			if i := slices.Index(closed, conn); i != -1 && closed[i].f < f {
				if closed[i].next == nil || closed[i].next.f >= f {
					fmt.Println("Skipping", conn.position, "because it's already in closed")
					fmt.Println("Should add path", conn.position, "to open")
				} else {
					fmt.Println("Skipping", conn.position, "because it's already in closed", closed[i].f, f)
					continue
				}
			}
			conn.f = f
			conn.g = g
			conn.h = h
			conn.ApproachDirection = a
			conn.parent = current
			current.next = conn
			open = append(open, conn)
		}
		closed = append(closed, current)
	}
	return int(bestScore), paths
}

func FindDistance(start, end *Node) int {
	return utils.AbsDiff(start.position.X, end.position.X) + utils.AbsDiff(start.position.Y, end.position.Y)  
}

func NewMaze(maze []string) (*Node, *Node) {
	mazeMap := map[Coordinate2D]*Node{}
	var start, end *Node
	for y, row := range maze {
		for x, elem := range row {
			if elem == Wall {
				continue
			}
			pos := Coordinate2D{x, y}
			thisNode := &Node{pos, map[Direction]*Node{}, math.Inf(0), math.Inf(0), math.Inf(0), None, nil, nil}
			mazeMap[pos] = thisNode
			if conn, ok := mazeMap[pos.Add(North)]; ok {
				mazeMap[pos].connections[North] = conn
				conn.connections[South] = thisNode
			}
			if conn, ok := mazeMap[pos.Add(West)]; ok {
				mazeMap[pos].connections[West] = conn
				conn.connections[East] = thisNode
			}

			if elem == Start {
				start = thisNode
				thisNode.f = 0
				thisNode.g = 0
				thisNode.ApproachDirection = East
			}
			if elem == End {
				end = thisNode
			}
		}
	}
	start.h = float64(FindDistance(start, end))
	return start, end
}

func Visualize(maze []string, end *Node) {
	nodes := make( map[Coordinate2D]*Node )
	for parent := end.parent; parent != nil; parent = parent.parent {
		nodes[parent.position] = parent
		fmt.Println(parent)
	}
	for y, row := range maze {
		for x, elem := range row {
			pos := Coordinate2D{x, y}
			if node, ok := nodes[pos]; ok {
				switch node.ApproachDirection {
				case North:
					elem = '^'
				case East:
					elem = '>'
				case South:
					elem = 'v'
				case West:
					elem = '<'
				}
			}
			fmt.Print(string(elem))
		}
		fmt.Println()
	}
}

func SolveDay16() {
	fmt.Println("Day 16")
	input := rawTest
	maze := strings.Split(input, "\n")
	start, end := NewMaze(maze)
	runner := Runner{start, East}
	score, paths := runner.RunAStar(end)
	Visualize(maze, end)
	fmt.Println(score)

	seen_seat := map[Coordinate2D]bool{}
	seat_count := 0
	rawMap := [][]rune{}
	for _, row := range maze {
		rawMap = append(rawMap, []rune(row))
	}
	for _, path := range paths {
		fmt.Println(path)
	}
	for _, path := range paths {
		for _, pos := range path {
			if _, ok := seen_seat[pos]; ok {
				continue
			}
			seen_seat[pos] = true
			if rawMap[pos.Y][pos.X] != Wall {
				rawMap[pos.Y][pos.X] = '0'
				seat_count++
			}
		}
	}
	fmt.Println(seat_count)
	for _, row := range rawMap {
		fmt.Println(string(row))
	}
}

func main() {
	SolveDay16()
}
