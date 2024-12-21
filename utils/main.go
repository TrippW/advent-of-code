package utils

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

var now = time.Now()

func TrackTimeFromNow(name string) {
	TrackTime(now, name)
}
func TrackTime(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func ReadFile(filename string) []string {
	f, err := os.ReadFile("inputs/" + filename)
	if err != nil {
		fmt.Println(err)
	}
	return strings.Split(string(f), "\n")
}

func AbsDiff(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

// if ColAsList is true, then each column will be a list of strings. All rows must have the same number of columns
// if ColAsList is false, then each row will be a list of strings. Each row can have a different number of columns
func SplitInputs(input []string, delim string, colAsList bool) ([][]string) {
	var lines [][]string
	for _, line := range input {
		lines = append(lines, strings.Split(line, delim))
	}

	var result [][]string
	if colAsList {
		result = make([][]string, len(lines[0]))
		for _, line := range lines {
			for j, s := range line {
				result[j] = append(result[j], strings.TrimSpace(s))
			}
		}
	} else {
		result = make([][]string, len(lines))
		for i, line := range lines {
			for _, s := range line {
				result[i] = append(result[i], strings.TrimSpace(s))
			}
		}
	}
	return result
}

func StrListToIntList(list []string) []int {
	var result []int
	for _, s := range list {
		i := 0
		fmt.Sscanf(s, "%d", &i)
		result = append(result, i)
	}
	return result
}

func StrListTo2DIntList(list []string) [][]int {
	var result [][]int
	for _, s := range list {
		var row []int
		for _, v := range strings.Split(s, "") {
			i := 0
			fmt.Sscanf(v, "%d", &i)
			row = append(row, i)
		}
		result = append(result, row)
	}
	return result
}

func SortStrListAsIntList(list []string) []int {
	var result sort.IntSlice
	result = StrListToIntList(list)
	result.Sort()
	return result
}

func RemoveIndex[T any](s []T, index int) []T {
	removed := make([]T, 0)
	removed = append(removed, s[:index]...)
	removed = append(removed, s[index+1:]...)
	return removed
}
