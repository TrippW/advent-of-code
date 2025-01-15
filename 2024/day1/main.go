package main

import (
	"github.com/TrippW/advent-of-code/utils"
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var rawInput string

//go:embed test.txt
var rawTest string

func findListDistance(left, right []int) int {
	result := 0
	for i := 0; i < len(left); i++ {
		result += utils.AbsDiff(left[i], right[i])
	}
	return result
}

func splitInputs(input []string) ([]int, []int) {
	lists := utils.SplitInputs(input, "   ", true)
	left, right := lists[0], lists[1]
	sortedLeft := utils.SortStrListAsIntList(left)
	sortedRight := utils.SortStrListAsIntList(right)
	return sortedLeft, sortedRight
}

func part1(sortedLeft, sortedRight []int) {
	result := findListDistance(sortedLeft, sortedRight)
	fmt.Printf("1.1 Answer: %d\n", result)
}

func part2(sortedLeft, sortedRight []int) {
	rightP := 0
	sum := 0
	count := 0;
	for leftP := 0; leftP < len(sortedLeft); leftP++ {
		for rightP < len(sortedRight) && sortedRight[rightP] <= sortedLeft[leftP] {
			if sortedRight[rightP] == sortedLeft[leftP] {
				count++
			}
			rightP++
		}
		sum += (count * sortedLeft[leftP])
		count = 0
	}

	fmt.Printf("1.2 Answer: %d\n", sum)
}

func Solve() {
	left, right := splitInputs(strings.Split(rawInput, "\n"))
	part1(left, right)
	part2(left, right)
}

func main() {
	Solve()
}
