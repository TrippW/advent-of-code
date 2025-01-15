package main

import (
	"fmt"
	re "regexp"
	"github.com/TrippW/advent-of-code/utils"
)

func solve_3_1() {
	regex := re.MustCompile(`(mul\(\d{1,3},\d{1,3}\))`)
	sum := solve_3(regex)
	fmt.Println("3.1 Answer:", sum)
}

func solve_3_2() {
	regex := re.MustCompile(`(do\(\)|don't\(\)|mul\(\d\d?\d?,\d\d?\d?\))`)
	sum := solve_3(regex)
	fmt.Println("3.2 Answer:", sum)
}

func solve_3(exp *re.Regexp) int {
	input := utils.ReadFile("input_3.txt")
	enabled := true
	sum := 0
	for _, line := range input {
		matches := exp.FindAllString(line, -1)
		for _, match := range matches {
			if match == "do()" {
				enabled = true
				continue
			}
			if match == "don't()" {
				enabled = false
				continue
			}
			if !enabled {
				continue
			}
			i, j := 0, 0
			fmt.Sscanf(match, "mul(%d,%d)", &i, &j)
			sum += (i * j)
		}
	}
	return sum
}

func SolveDay3 () {
	solve_3_1()
	solve_3_2()
}

