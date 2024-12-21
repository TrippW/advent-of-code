package graphs

import (
	"fmt"
	"github.com/TrippW/advent-of-code/utils"
)

type Coordinate2D struct {
	X, Y int
}

func (c Coordinate2D) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

type AStarResponse struct {
	Path []Coordinate2D
	CameFrom map[Coordinate2D]Coordinate2D
	FScore map[Coordinate2D]int
	GScore map[Coordinate2D]int
}

type Direction int

const (
	None Direction = iota
	North
	East
	South
	West
	NorthEast
	NorthWest
	SouthEast
	SouthWest
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
	case NorthEast:
		return Coordinate2D{c.X + 1, c.Y - 1}
	case NorthWest:
		return Coordinate2D{c.X - 1, c.Y - 1}
	case SouthEast:
		return Coordinate2D{c.X + 1, c.Y + 1}
	case SouthWest:
		return Coordinate2D{c.X - 1, c.Y + 1}
	}
	panic("Invalid direction")
}

type Runner struct {
	Start Coordinate2D
	Direction Direction
}

type AStarBoard struct {
	Grid [][]byte
	Start Coordinate2D
	End Coordinate2D
	ObstructedTiles map[byte]bool
}

func AStar(board AStarBoard, heuristic func(Coordinate2D, Coordinate2D)int, ValidDirections []Direction) *AStarResponse {
	start := board.Start
	end := board.End

	openSet := make(map[Coordinate2D]bool)
	openSet[start] = true

	cameFrom := make(map[Coordinate2D]Coordinate2D)

	gScore := make(map[Coordinate2D]int)
	gScore[start] = 0

	fScore := make(map[Coordinate2D]int)
	fScore[start] = heuristic(start, end)

	for len(openSet) != 0 {
		current := getLowestFScore(openSet, fScore)
		if current == end {
			resp := reconstructPath(cameFrom, current, start)
			resp.FScore = fScore
			resp.GScore = gScore
			return resp
		}

		delete(openSet, current)

		for _, neighbor := range getNeighbors(board.Grid, current, board.ObstructedTiles) {
			tentativeGScore := gScore[current] + 1
			if score, ok := gScore[neighbor] ; !ok || tentativeGScore < score {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentativeGScore
				fScore[neighbor] = gScore[neighbor] + heuristic(neighbor, end)
				openSet[neighbor] = true
			}
		}
	}

	return nil
}

func getLowestFScore(openSet map[Coordinate2D]bool, fScore map[Coordinate2D]int) Coordinate2D {
	lowest := ^1
	var lowestCoord *Coordinate2D
	for coord := range openSet {
		if fScore[coord] < lowest || lowestCoord == nil{
			lowest = fScore[coord]
			lowestCoord = &coord
		}
	}

	if lowestCoord == nil {
		panic("No lowest coord found")
	}
	return *lowestCoord
}

func getNeighbors(grid [][]byte, coord Coordinate2D, obstructedTileSymbols map[byte]bool) []Coordinate2D {
	neighbors := make([]Coordinate2D, 0)
	if coord.X > 0 && !obstructedTileSymbols[grid[coord.Y][coord.X - 1]] {
		neighbors = append(neighbors, Coordinate2D{coord.X - 1, coord.Y})
	}
	if coord.X < len(grid[0]) - 1 && !obstructedTileSymbols[grid[coord.Y][coord.X + 1]] {
		neighbors = append(neighbors, Coordinate2D{coord.X + 1, coord.Y})
	}
	if coord.Y > 0 && !obstructedTileSymbols[grid[coord.Y - 1][coord.X]] {
		neighbors = append(neighbors, Coordinate2D{coord.X, coord.Y - 1})
	}
	if coord.Y < len(grid) - 1 && !obstructedTileSymbols[grid[coord.Y + 1][coord.X]] {
		neighbors = append(neighbors, Coordinate2D{coord.X, coord.Y + 1})
	}
	return neighbors
}

func reconstructPath(cameFrom map[Coordinate2D]Coordinate2D, current, start Coordinate2D) *AStarResponse {
	path := make([]Coordinate2D, 0)
	for current != start {
		path = append(path, current)
		current = cameFrom[current]
	}
	path = append(path, start)

	return &AStarResponse{
		Path: path,
		CameFrom: cameFrom,
	}
}

func ManhattanDistance(a, b Coordinate2D) int {
	return utils.AbsDiff(a.X,b.X) + utils.AbsDiff(a.Y,b.Y)
}

func NewSquareGrid(size int) [][]byte {
	grid := make([][]byte, size + 1)
	for i := range grid {
		grid[i] = make([]byte, size + 1)
	}
	return grid
}

func PlaceAllTiles(grid [][]byte, coords []Coordinate2D, symbol byte) {
	for _, coord := range coords {
		grid[coord.Y][coord.X] = symbol
	}
}

func PlaceTiles(grid [][]byte, coords []Coordinate2D, maxCoords int, symbol byte) {
	if maxCoords > len(coords) {
		PlaceAllTiles(grid, coords, symbol)
	} else {
		PlaceAllTiles(grid, coords[:maxCoords], symbol)
	}
}

func NewCoordinateMap(path []Coordinate2D) map[Coordinate2D]bool {
	pathMap := make(map[Coordinate2D]bool)
	for _, coord := range path {
		pathMap[coord] = true
	}
	return pathMap
}

func PrintGrid(grid [][]byte, path []Coordinate2D, wall byte) {
	pathMap := NewCoordinateMap(path)
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if pathMap[Coordinate2D{x, y}] {
				if grid[y][x] == wall {
					fmt.Print("X")
				} else {
					fmt.Print("O")
				}
			} else {
				if grid[y][x] == wall {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
		}
		fmt.Println()
	}
}

