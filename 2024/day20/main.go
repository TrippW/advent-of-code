package main

import (
	"fmt"
	"sort"
	"strings"
	_ "embed"
	_ "github.com/TrippW/advent-of-code/utils"
	"github.com/TrippW/advent-of-code/graphs"
)

//go:embed input.txt
var rawInput string

//go:embed test.txt
var rawTest string

var allowedDirections = []graphs.Direction{graphs.North, graphs.South, graphs.East, graphs.West}

// memory is a map of coordinates to a map of cheats remaining to a map of picoSaved to count of cheats
func cheat(scores map[graphs.Coordinate2D]int, coord graphs.Coordinate2D, directions []graphs.Direction, cheatsRemaining int, memory map[graphs.Coordinate2D]map[int]map[int]int, seen map[graphs.Coordinate2D]bool) map[int]int { // map[picoSaved]count
	// if we have no cheats left or we've already seen this coordinate, return
	if cheatsRemaining < 2 || seen[coord] {
		return make(map[int]int, 0)
	}
	seen[coord] = true
	if _, ok := scores[coord]; !ok {
		return make(map[int]int, 0)
	}
	if v, ok := memory[coord][cheatsRemaining]; ok {
		return v
	}

	if cheatsRemaining == 2 {
		cheats := make([]graphs.Coordinate2D, 0)
		for _, d := range allowedDirections {
			cheats = append(cheats, coord.Add(d).Add(d))
		}
		saves := make(map[int]int, 0)
		for _, c := range cheats {
			if cheat_score, ok := scores[c]; ok {
				if cheat_score > (scores[coord] + 2) {
					diff := cheat_score - scores[coord] - 2
					if _, ok := saves[diff]; !ok {
						saves[diff] = 0
					}
					saves[diff]++
				}
			}
		}
		fmt.Println("At", coord, "with", cheatsRemaining, "cheats, we can save", saves)
		if _, ok := memory[coord]; !ok {
			memory[coord] = make(map[int]map[int]int, 0)
		}
		memory[coord][cheatsRemaining] = saves
		return saves
	} else {
		total_saves := make(map[int]int, 0)
		for _, d := range directions {
			next := coord.Add(d)
			saves := cheat(scores, next, directions, cheatsRemaining - 1, memory, seen)
			for k, v := range saves {
				if _, ok := total_saves[k]; !ok {
					total_saves[k] = 0
				}
				total_saves[k] += v
			}
		}
		return total_saves
	}
}

func main() {
	grid := make([][]byte, 0)
	start := graphs.Coordinate2D{X: -1, Y: -1}
	end := graphs.Coordinate2D{X: -1, Y: -1}
	for l, line := range strings.Split(rawTest, "\n") {
		grid = append(grid, []byte(line))
		if i := strings.Index(line, "S"); i != -1 {
			start = graphs.Coordinate2D{X: i, Y: l}
		}
		if i := strings.Index(line, "E"); i != -1 {
			end = graphs.Coordinate2D{X: i, Y: l}
		}
	}
	board := graphs.AStarBoard{
		Grid: grid,
		Start: start,
		End: end,
		ObstructedTiles: map[byte]bool{'#' : true},
	}

	resp := graphs.AStar(board, graphs.ManhattanDistance, allowedDirections)
	saves := make(map[int]int, 0)
	flat_memory := make(map[int]int, 0)
	memory := make(map[graphs.Coordinate2D]map[int]map[int]int, 0)
	for _, p := range resp.Path {
		seen := make(map[graphs.Coordinate2D]bool, 0)
		if _, ok := memory[p]; !ok {
			memory[p] = make(map[int]map[int]int, 0)
		}
		saves = cheat(resp.GScore, p, allowedDirections, 2, memory, seen)
		memory[p][19] = saves
		for k, v := range saves {
			if _, ok := flat_memory[k]; !ok {
				flat_memory[k] = 0
			}
			flat_memory[k] += v
		}
	}
	sort_keys := make([]int, 0)
	for k := range flat_memory {
		sort_keys = append(sort_keys, k)
	}
	sort.Ints(sort_keys)
	cnt := 0
	for _, k := range sort_keys {
		if k >= 50 {
			cnt += flat_memory[k]
			fmt.Println("There are", flat_memory[k], "cheats that save", k, "picoseconds")
		}
	}
	fmt.Println("There are", cnt, "cheats that save 100 picoseconds or more")
}
