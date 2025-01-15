package main

import (
	"fmt"
	"strings"
	_ "embed"
	"github.com/TrippW/advent-of-code/utils"
	"github.com/TrippW/advent-of-code/graphs"
)

//go:embed input.txt
var rawInput string

//go:embed test.txt
var testInput string

type spot byte

type Coordinate2D = graphs.Coordinate2D

const (
	Empty spot = iota
	Corrupted 
)

func parseCorrupted(corrupted []string) []Coordinate2D {
	coords := make([]Coordinate2D, 0)
	for _, line := range corrupted {
		if line == "" {
			continue
		}
		var x, y int
		fmt.Sscanf(line, "%d,%d", &x, &y)
		coords = append(coords, Coordinate2D{X: x, Y: y})
	}
	return coords
}


func GetTestSetup() ([][]byte, []Coordinate2D, int, Coordinate2D) {
	grid := graphs.NewSquareGrid(6)
	corrupted := parseCorrupted(strings.Split(testInput, "\n"))
	corruptedCount := 12
	target := Coordinate2D{X: 6, Y: 6}
	return grid, corrupted, corruptedCount, target
}

func GetInputSetup() ([][]byte, []Coordinate2D, int, Coordinate2D) {
	grid := graphs.NewSquareGrid(70)
	corrupted := parseCorrupted(strings.Split(rawInput, "\n"))
	corruptedCount := 1024
	target := Coordinate2D{X: 70, Y: 70}
	return grid, corrupted, corruptedCount, target
}

func main() {
	defer utils.TrackTimeFromNow("main")
	grid, corrupted, corruptedCount, target := GetInputSetup()
	
	graphs.PlaceTiles(grid, corrupted, corruptedCount, byte(Corrupted))
	board := graphs.AStarBoard{
		Grid: grid, 
		Start: Coordinate2D{X: 0, Y: 0}, 
		End: target,
		ObstructedTiles: map[byte]bool{byte(Corrupted): true},
	}

	resp := graphs.AStar(board, graphs.ManhattanDistance, []graphs.Direction{graphs.North, graphs.East, graphs.South, graphs.West})
	path := resp.Path
	graphs.PrintGrid(grid, path, byte(Corrupted))
	fmt.Println("Part 1:", len(path) - 1)
	defer utils.TrackTimeFromNow("part2")

	var Corruptions []Coordinate2D
	lastPath := path

	for ; resp != nil; resp = graphs.AStar(board, graphs.ManhattanDistance, []graphs.Direction{graphs.North, graphs.East, graphs.South, graphs.West}) {
		pathMap := graphs.NewCoordinateMap(resp.Path)
		lastPath = resp.Path
		i := 1
		for corruption := corrupted[i + corruptedCount]; ; corruption = corrupted[i + corruptedCount] {
			if _, ok := pathMap[Coordinate2D{X: corruption.X, Y: corruption.Y}]; ok {
				fmt.Println("Blocking due to", corruption)
				Corruptions = append(Corruptions, corruption)
				break
			}
			i++
			if i + corruptedCount >= len(corrupted) {
				break
			}
		}
		
		graphs.PlaceTiles(grid, corrupted[corruptedCount+1:], i, byte(Corrupted))
		corruptedCount += i
	}

	graphs.PrintGrid(grid, lastPath, byte(Corrupted))

	fmt.Printf("Part 2: %v\n", Corruptions[len(Corruptions) - 1])
}
