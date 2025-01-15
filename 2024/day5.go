 package main

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/TrippW/advent-of-code/utils"
)

func SolveDay5() {
	SolveDay5_1()
	SolveDay5_2()
}

func Day5CreateRule(line string, rules map[string][]string) {
	rule := strings.Split(line, "|")
	if list, ok := rules[rule[1]]; ok {
		rules[rule[1]] = append(list, rule[0])
	} else {
		rules[rule[1]] = []string{rule[0]}
	}
}

func SolveDay5_1() {
	input := utils.ReadFile("5_input.txt")
	rules := map[string][]string{}
	creatingRules := true
	correctionErrorsMidSum := 0
	for _, line := range input {
		if line == "" {
			creatingRules = false
			continue
		}
		if creatingRules {
			Day5CreateRule(line, rules)
		} else {
			corrections := strings.Split(line, ",")
			forbidden := map[string]bool{}
			correct := true
			for _, correction := range corrections {
				if _, ok := forbidden[correction]; ok {
					correct = false
					break
				}
				
				for _, rule := range rules[correction] {
					forbidden[rule] = true
				}
			}

			if correct {
				if v, err := strconv.Atoi(corrections[len(corrections)/2]); err == nil {
					correctionErrorsMidSum += v
				}
			}
		}
	}

	fmt.Println("5.1:", correctionErrorsMidSum)
}

func SolveDay5_2() {
	input := utils.ReadFile("5_input.txt")
	rules := map[string][]string{}
	creatingRules := true
	correctionErrorsMidSum := 0
	for _, line := range input {
		if line == "" {
			creatingRules = false
			continue
		}
		if creatingRules {
			Day5CreateRule(line, rules)
		} else {
			corrections := strings.Split(line, ",")
			record := false
			correct := false
			for !correct {
				correct = true
				forbidden := map[string]int{}
				for i, correction := range corrections {
					if wrongIndex, ok := forbidden[correction]; ok {
						correct = false
						record = true
						corrections[i], corrections[wrongIndex] = corrections[wrongIndex], corrections[i]
						break
					}

					for _, rule := range rules[correction] {
						forbidden[rule] = i
					}
				}

				if !record {
					correct = true
				}
			}

			if record {
				if v, err := strconv.Atoi(corrections[len(corrections)/2]); err == nil {
					correctionErrorsMidSum += v
				}
			}
		}
	}

	fmt.Println("5.2:", correctionErrorsMidSum)
}


