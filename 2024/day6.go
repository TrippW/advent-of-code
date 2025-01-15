package main

import (
	"fmt"
	tm "github.com/buger/goterm"
	"time"
	"github.com/TrippW/advent-of-code/utils"
)

type GuardMapPosition struct {
	X int
	Y int
	IsObstacle bool
	IsVisited bool
	AdjacentPositions map[GuardDirection]*GuardMapPosition
	SeenDirection map[GuardDirection]bool
}

func (p *GuardMapPosition) Visit() {
	p.IsVisited = true
}

type GuardMap [][]*GuardMapPosition

type GuardDirection = int

const (
	GuardUp GuardDirection = iota
	GuardRight
	GuardDown
	GuardLeft
	Done
)

var guardCharacterDirections = map[string]GuardDirection{
	"^": GuardUp,
	">": GuardRight,
	"v": GuardDown,
	"<": GuardLeft,
}

type Guard struct {
	X int
	Y int
	Direction int
	CurrentPosition *GuardMapPosition
}

func (g Guard) Peek() GuardDirection {
	for offset := 0; offset < 4; offset++ {
		next := g.CurrentPosition.AdjacentPositions[(g.Direction + offset) % 4]
		if next != nil {
			if next.IsObstacle {
				continue
			}
			return (g.Direction + offset) % 4
		} else {
			return Done
		}
	}

	return Done
}

func (g *Guard) Move(d GuardDirection) {
	g.Direction = d
	g.CurrentPosition.Visit()
	g.CurrentPosition.SeenDirection[d] = true
	g.CurrentPosition = g.CurrentPosition.AdjacentPositions[d]
}

func toGuardMap(input []string) (Guard, GuardMap) {
	var guard Guard
	var guardMap GuardMap
	for i, line := range input {
		var row []*GuardMapPosition
		for j, c := range line {
			p := GuardMapPosition{
				X: j,
				Y: i,
				IsObstacle: false, 
				IsVisited: false, 
				AdjacentPositions: map[GuardDirection]*GuardMapPosition{
					GuardUp: nil,
					GuardRight: nil,
					GuardDown: nil,
					GuardLeft: nil,
				},
				SeenDirection: map[GuardDirection]bool{
					GuardUp: false,
					GuardRight: false,
					GuardDown: false,
					GuardLeft: false,
				},
			}
			if c == '#' {
				p.IsObstacle = true
			}
			if d, ok := guardCharacterDirections[string(c)]; ok {
				guard = Guard{i, j, d, &p}
			}
			if i > 0 {
				p.AdjacentPositions[GuardUp] = guardMap[i-1][j]
				guardMap[i-1][j].AdjacentPositions[GuardDown] = &p
			}
			if j > 0 {
				p.AdjacentPositions[GuardLeft] = row[j-1]
				row[j-1].AdjacentPositions[GuardRight] = &p
			}
			row = append(row, &p)
		}
		guardMap = append(guardMap, row)
	}

	return guard, guardMap
}

func VisualizeGuardMap(guardMap GuardMap, guard Guard) {
	for _, row := range guardMap {
		for _, p := range row {
			if guard.X == p.X && guard.Y == p.Y {
				switch guard.Direction {
				case GuardUp:
					fmt.Print("^")
				case GuardRight:
					fmt.Print(">")
				case GuardDown:
					fmt.Print("v")
				case GuardLeft:
					fmt.Print("<")
				}
			} else if p.IsObstacle {
				fmt.Print("#")
			} else if p.IsVisited {
				num := 0
				if p.SeenDirection[GuardUp] {
					num += 1 << GuardUp
				}
				if p.SeenDirection[GuardRight] {
					num += 1 << GuardRight
				}
				if p.SeenDirection[GuardDown] {
					num += 1 << GuardDown
				}
				if p.SeenDirection[GuardLeft] {
					num += 1 << GuardLeft
				}
				fmt.Printf("%x", num)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

type VisualizeConfig struct {
	Width int
	Height int
	ShowTrail bool
	Clamp bool
}

func VisualizeGuardMapSquare(guardMap GuardMap, guard Guard, topTxt string, config VisualizeConfig) {
	x0 := 0
	x1 := len(guardMap[0]) - 1
	y0 := 0
	y1 := len(guardMap) - 1

	if !config.Clamp {
		if x1 > tm.Width() {
			config.Width = (tm.Width() - 1) / 2
		} else {
			config.Width = (x1 - 1) / 2
		}
		if y1 > tm.Height() {
			config.Height = (tm.Height() - 1) / 2 - 1
		} else {
			config.Height = (y1 - 1) / 2
		}
		config.Clamp = true
	}

	if config.Clamp {
		x0 = guard.CurrentPosition.X - config.Width
		x1 = guard.CurrentPosition.X + config.Width
		y0 = guard.CurrentPosition.Y - config.Height
		y1 = guard.CurrentPosition.Y + config.Height

		if y0 < 0 {
			y1 += utils.AbsDiff(y0, 0)
			if y1 >= len(guardMap) {
				y1 = len(guardMap) - 1
			}
			y0 = 0
		}
		if y1 >= len(guardMap) {
			y0 -= utils.AbsDiff(y1, len(guardMap) - 1)
			if y0 < 0 {
				y0 = 0
			}
			y1 = len(guardMap) - 1
		}
		if x0 < 0 {
			x1 += utils.AbsDiff(x0, 0)
			if x1 >= len(guardMap[0]) {
				x1 = len(guardMap[0]) - 1
			}
			x0 = 0
		}
		if x1 >= len(guardMap[0]) {
			x0 -= utils.AbsDiff(x1, len(guardMap[0]) - 1)
			if x0 < 0 {
				x0 = 0
			}
			x1 = len(guardMap[0]) - 1
		}
	}

	tm.MoveCursor(1, 1)
	tm.Printf("%s\n", topTxt)
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			p := guardMap[y][x]
			if p == guard.CurrentPosition {
				switch guard.Direction {
				case GuardUp:
					tm.Print("^")
				case GuardRight:
					tm.Print(">")
				case GuardDown:
					tm.Print("v")
				case GuardLeft:
					tm.Print("<")
				}
			} else if p.IsObstacle {
				tm.Print(" ")
			} else if p.IsVisited  && config.ShowTrail{
				num := 0
				if p.SeenDirection[GuardUp] {
					num += 1 << GuardUp
				}
				if p.SeenDirection[GuardRight] {
					num += 1 << GuardRight
				}
				if p.SeenDirection[GuardDown] {
					num += 1 << GuardDown
				}
				if p.SeenDirection[GuardLeft] {
					num += 1 << GuardLeft
				}
				tm.Printf("%x", num)
			} else {
				tm.Print(" ")
			}
		}
		tm.Println()
	}
	tm.Flush()
	time.Sleep(10 * time.Millisecond)
}

type Coordinate2D struct {
	X int
	Y int
}

func SolveDay6() {
	input := utils.ReadFile("day6.txt")
	guard, _ := toGuardMap(input)
	cnt := 0
	visitedCoords := make(map[int]map[int]bool)
	for d := guard.Peek(); d != Done; d = guard.Peek() {
		if _, ok := visitedCoords[guard.CurrentPosition.X]; !ok {
			visitedCoords[guard.CurrentPosition.X] = make(map[int]bool)
		}
		visitedCoords[guard.CurrentPosition.X][guard.CurrentPosition.Y] = true
		if !guard.CurrentPosition.IsVisited {
			cnt++
		}
		guard.Move(d)
		if _, ok := visitedCoords[guard.CurrentPosition.X]; !ok {
			visitedCoords[guard.CurrentPosition.X] = make(map[int]bool)
		}
		visitedCoords[guard.CurrentPosition.X][guard.CurrentPosition.Y] = true
	}

	if !guard.CurrentPosition.IsVisited {
		cnt++
		guard.CurrentPosition.Visit()
		if _, ok := visitedCoords[guard.CurrentPosition.X]; !ok {
			visitedCoords[guard.CurrentPosition.X] = make(map[int]bool)
		}
	}
	fmt.Println("Moves to exit:", cnt, "Expected 41")
	cnt = 0
	for _, row := range visitedCoords {
		for range row {
			cnt++
		}
	}
	fmt.Println("Visited Coords cnt:", cnt, "Expected 41")

	cnt = 0
	origin_guard, origin_map := toGuardMap(input)
	guard = Guard{origin_guard.X, origin_guard.Y, origin_guard.Direction, origin_guard.CurrentPosition}
	fmt.Println("Size", len(input), len(input[0]))
	for j, row := range visitedCoords {
		for i := range row {
			if origin_map[i][j].IsObstacle || (origin_guard.X == i && origin_guard.Y == j) {
				continue
			}
			tm.Clear()
			visited := []*GuardMapPosition{
				origin_guard.CurrentPosition,
			}
			origin_map[i][j].IsObstacle = true
			moves := 0
			for d := guard.Peek(); d != Done; d = guard.Peek() {
				visited = append(visited, guard.CurrentPosition)
				if guard.CurrentPosition.IsVisited {
					if seen, ok := guard.CurrentPosition.SeenDirection[d]; ok && seen {
						cnt++
						break
					}
				}
				guard.Move(d)
				moves++
				VisualizeGuardMapSquare(
					origin_map, 
					guard, 
					fmt.Sprintf("Added %d, %d; Move: %d; Pos: %03d, %03d", i, j, moves, guard.CurrentPosition.X, guard.CurrentPosition.Y),
					VisualizeConfig{Width: 20, Height: 20, ShowTrail: true, Clamp: false},
				)	
			}

			// Reset the guard and map
			for k := range visited {
				visited[k].IsVisited = false
				visited[k].SeenDirection = map[GuardDirection]bool{
					GuardUp: false,
					GuardRight: false,
					GuardDown: false,
					GuardLeft: false,
				}
			}
			origin_map[i][j].IsObstacle = false
			guard = Guard{origin_guard.X, origin_guard.Y, origin_guard.Direction, origin_guard.CurrentPosition}
		}
	}

	fmt.Println("Unique Loops:", cnt, "16083 is too high")
}

