package main;

import (
	"fmt"
	"sort"
	"strings"
	_ "embed"
)

//go:embed input.txt
var rawInput string

//go:embed test.txt
var rawTest string


func parseFile(input string) ([]string, []string) {
	lines := strings.Split(input, "\n")
	buildingPatterns := true
	patterns := make([]string, 0)
	designs := make([]string, 0)
	for _, line := range lines {
		if line == "" {
			buildingPatterns = false
			continue
		}
		if buildingPatterns {
			patterns = append(patterns, strings.Split(line, ", ")...)
		} else {
			designs = append(designs, line)
		}
	}

	return patterns, designs
}

type PatternHolder struct {
	sizes []int
	patterns map[int]map[rune][]string
	isValid map[string]bool
}

func (p *PatternHolder) addPattern(pattern string) {
	size := len(pattern)
	if _, ok := p.patterns[size]; !ok {
		p.sizes = append(p.sizes, size)
		p.patterns[size] = make(map[rune][]string)
		sort.Ints(p.sizes)
	}
	firstChar := rune(pattern[0])
	p.patterns[size][firstChar] = append(p.patterns[size][firstChar], pattern)
}

func (ph *PatternHolder) validate(pattern string) bool {
	if isValid, ok := ph.isValid[pattern]; ok {
		return isValid
	}

	if pattern == "" {
		return true
	}

	for i := len(ph.sizes) - 1; i >= 0; i-- {
		size := ph.sizes[i]
		if size > len(pattern) {
			continue
		}
		firstChar := rune(pattern[0])
		if _, ok := ph.patterns[size][firstChar]; !ok {
			continue
		}
		for _, p := range ph.patterns[size][firstChar] {
			if pattern[:size] == p {
				if ph.validate(pattern[size:]) {
					ph.isValid[pattern] = true
					return true
				}
			}
		}
	}
	ph.isValid[pattern] = false
	return false
}

func (p *PatternHolder) addPatterns(patterns []string) {
	for _, pattern := range patterns {
		p.addPattern(pattern)
	}
}

func main() {
	p, d := parseFile(rawInput)
	ph := PatternHolder{patterns: map[int]map[rune][]string{}, sizes: []int{}, isValid: map[string]bool{}}
	ph.addPatterns(p)
	matches := 0
	for _, d := range d {
		v := ph.validate(d)
		if v {
			matches++
		}
	}
	fmt.Println("Matches:", matches)
}
