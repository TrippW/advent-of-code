package main

import (
	"fmt"
	_ "embed"
	"strings"
	"github.com/TrippW/advent-of-code/utils"
)

//go:embed input.txt
var rawInput string

func checkValid(x, y int, expectIncreasing, expectDecreasing bool) bool {
	diff := x - y
	if diff == 0 || diff < -3 || diff > 3 {
		return false
	}
	if x < y && expectDecreasing {
		return false
	}
	if x > y && expectIncreasing {
		return false
	}
	return true
}

func validateLine(line []int) bool{
	increasing := line[0] < line[1]
	for i := 0; i < len(line) - 1; i++ {
		if !checkValid(line[i], line[i+1], increasing, !increasing) {
			return false
		}
	}

	return true
}

func findValid(input [][]int, allowSingleError bool) int {
		count := 0
		for _, line := range input {
			if validateLine(line) {
				count++
			} else if allowSingleError {
				for i := 0; i < len(line); i++ {
					if validateLine(utils.RemoveIndex(line, i)) {
						count++
						break;
					}
				}
			}
		}
		return count
}

func getDay2InputLines() [][]int {
	input := strings.Split(rawInput, "\n")
	lines := utils.SplitInputs(input, " ", false)
	res := make([][]int, len(lines))
	for i, line := range lines {
		res[i] = utils.StrListToIntList(line)
	}
	return res
}

func solve_2_1() {
	res := findValid(getDay2InputLines(), false)
	fmt.Println("2.1 Answer:", res)
}

func solve_2_2() {
	res := findValid(getDay2InputLines(), true)
	fmt.Println("2.2 Answer:", res)
}

func SolveDay2() {
	solve_2_1()
	solve_2_2()
}

func main() {
	SolveDay2()
}
