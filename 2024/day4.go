package main

import (
	"fmt"
	"github.com/TrippW/advent-of-code/utils"
)

func findWordMatches(input []string, word string, x, y int, dir []int) int {
	if word == "" {
		return 1
	}
	if y >= len(input) || y < 0 || x < 0 || x >= len(input[y]) {
		return 0
	}
	if input[y][x] != word[0] {
		return 0
	}
	changes := [][]int {}

	if len(dir) == 0 {
		sets := []int {-1, 0, 1}
		for _, dx := range sets {
			for _, dy := range sets {
				changes = append(changes, []int{dx, dy})
			}
		}
	} else {
		changes = append(changes, dir)
	}

	nextWord := word[1:]
	cnt := 0
	for _, change := range changes {
		cnt += findWordMatches(input, nextWord, x + change[0], y + change[1], change)
	}

	return cnt
}

func findAllWordMatches(input []string, word string) int {
	cnt := 0
	for y := 0; y < len(input); y++ {
		for x := 0; x < len(input[y]); x++ {
			cnt += findWordMatches(input, word, x, y, []int{})
		}
	}
	return cnt
}

func isMS(a, b byte) bool {
	return (a == 'M' && b == 'S') || (a == 'S' && b == 'M')
}

func findMSCross(input []string, x, y int) bool {
	topLeft := input[y-1][x-1]
	topRight := input[y-1][x+1]
	botLeft := input[y+1][x-1]
	botRight := input[y+1][x+1]

	return isMS(topLeft, botRight) && isMS(topRight, botLeft)
}

func findMasCrosses(input []string) int {
	cnt := 0
	for y := 1; y < len(input)-1; y++ {
		for x := 1; x < len(input[y])-1; x++ {
			if input[y][x] == 'A' && findMSCross(input, x, y) {
				cnt++
			}
		}
	}

	return cnt
}

func solve_day_4_test() {
	input := []string{
		"MMMSXXMASM",
		"MSAMXMSMSA",
		"AMXSXMAAMM",
		"MSAMASMSMX",
		"XMASAMXAMM",
		"XXAMMXXAMA",
		"SMSMSASXSS",
		"SAXAMASAAA",
		"MAMMMXMMMM",
		"MXMXAXMASX",
	}

	fmt.Println("Day 4 Test: expects 18 == ", findAllWordMatches(input, "XMAS"))
	fmt.Println("Day 4 Test: expects 9 == ", findMasCrosses(input))
}

func solve_day_4_1() {
	input := utils.ReadFile("4_input.txt")
	fmt.Println("Day 4 pt1:", findAllWordMatches(input, "XMAS"))
}

func solve_day_4_2() {
	input := utils.ReadFile("4_input.txt")
	fmt.Println("Day 4 pt2:", findMasCrosses(input))
}

func SolveDay4 () {
	solve_day_4_test()
	solve_day_4_1()
	solve_day_4_2()
}
