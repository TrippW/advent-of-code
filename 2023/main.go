package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
	"unicode"
	"regexp"
)

func FindFirstAndLastNumbers(s string) (int, int) {
	firstIndex := len(s)
	firstNumber := 0
	lastIndex := 0 
	lastNumber := 0

	numberstrings := []string{ "zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine" }


	for i, r := range s {
		if unicode.IsDigit(r) {
			if i < firstIndex {
				firstIndex = i
				firstNumber = int(r - '0')
				fmt.Println("Setting first number to ", firstNumber, " at index ", i, " with string ", s)
			}
			if i > lastIndex {
				lastIndex = i
				lastNumber = int(r - '0')
				fmt.Println("Setting last number to ", lastNumber, " at index ", i, " with string ", s)
				continue
			}
		}

		for j, numberstring := range numberstrings {
			if i+len(numberstring) > len(s) {
				continue
			}
			if i+len(numberstring) <= len(s) && s[i:i+len(numberstring)] == numberstring {
				fmt.Println("Checking ", numberstring, " against ", s[i:i+len(numberstring)])
				if i < firstIndex {
					fmt.Println("Setting first number to ", j, " at index ", i, " with string ", s)
					firstIndex = i
					firstNumber = j
				}
				if i > lastIndex {
					fmt.Println("Setting last number to ", j, " at index ", i, " with string ", s)
					lastIndex = i
					lastNumber = j
					continue
				}
			}
		}
	}

	fmt.Println("Returning first number ", firstNumber, " and last number ", lastNumber)

	return firstNumber, lastNumber
}

func FindNumber(s string) int {
	f, l := FindFirstAndLastNumbers(s)
	return f*10 + l
}

func PuzzleOne() int {
	inputFile, _ := os.Open("input1.txt")
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)
	sum := 0
	var err error
	for str := ""; err == nil; str, err = reader.ReadString('\n') {
		sum += FindNumber(str)
	}
	return sum
}

type CubeGame struct {
	Number int
	Pulls []Pull
}

type Pull struct {
	Cubes []Cube
}

type Cube struct {
	Number int
	Color Color
}

type Color string

const Red Color = "red"
const Blue Color = "blue"
const Green Color = "green"

func NewPulls(s string) []Pull {
	pulls := []Pull{}
	for _, pullStr := range strings.Split(s, ";") {
		pull := Pull{}
		for _, cubeStr := range strings.Split(pullStr, ",") {
			cube := Cube{}
			cubeNumberEndIndex := strings.Index(cubeStr[1:], " ") + 1
			cube.Number, _ = strconv.Atoi(cubeStr[1:cubeNumberEndIndex])
			switch string(cubeStr[cubeNumberEndIndex+1]) {
			case "r":
				cube.Color = Red
			case "b":
				cube.Color = Blue
			case "g":
				cube.Color = Green
			default:
				panic("Invalid color")
			}
			pull.Cubes = append(pull.Cubes, cube)
		}
		pulls = append(pulls, pull)
	}
	return pulls
}

var CubeLimits = map[Color]int{
	Red: 12,
	Green: 13,
	Blue: 14,
}

func CubeConundrum(filename string) {
	inputFile, _ := os.Open(filename)
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)
	partOneSum := 0
	partTwoSum := 0
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		gameStr := str[5:5+strings.Index(str[5:], ":")]
		gameNumber, _ := strconv.Atoi(gameStr)

		Game := CubeGame{
			Number: gameNumber,
			Pulls: NewPulls(str[5 + strings.Index(str[5:], ":") + 1:]),
		}

		// Part 1
		valid := true
		for _, pull := range Game.Pulls {
			for _, cube := range pull.Cubes {
				if cube.Number > CubeLimits[cube.Color] {
					fmt.Println("Game ", Game.Number, " is invalid")
					valid = false
					break
				}
			}
			if !valid {
				break
			}
		}
		if valid {
			partOneSum += Game.Number
		}

		// Part 2
		MinCube := map[Color]int{
			Red: 0,
			Green: 0,
			Blue: 0,
		}
		for _, pull := range Game.Pulls {
			for _, cube := range pull.Cubes {
				MinCube[cube.Color] = max(cube.Number, MinCube[cube.Color])
			}
		}
		power := 1
		for _, minCube := range MinCube {
			power *= minCube
		}
		partTwoSum += power
	}
	fmt.Println("Puzzle Two Part 1: ", partOneSum)
	fmt.Println("Puzzle Two Part 2: ", partTwoSum)
}

type GearType string

const (
	Number GearType = "number"
	Symbol GearType = "symbol"
	Blank GearType = "blank"
)

type Gear struct {
	Used bool
	IsPartNumber bool
	Value int
}

type Coordinate struct {
	X int
	Y int
}

func (c Coordinate) Add(other Coordinate) Coordinate {
	return Coordinate{c.X + other.X, c.Y + other.Y}
}

func (c Coordinate) Surrounding() []Coordinate {
	return []Coordinate{
		c.Add(Coordinate{0, 1}),
		c.Add(Coordinate{0, -1}),
		c.Add(Coordinate{1, 0}),
		c.Add(Coordinate{-1, 0}),
		c.Add(Coordinate{1, 1}),
		c.Add(Coordinate{1, -1}),
		c.Add(Coordinate{-1, 1}),
		c.Add(Coordinate{-1, -1}),
	}
}

func GearRatios(filename string) {
	inputFile, _ := os.Open(filename)
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)
	partOneSum := 0
	partTwoSum := 0
	schematic := map[Coordinate]*Gear{}
	symbols := []Coordinate{}
	partNumbers := []Gear{}
	// Read in the schematic
	for row := 0; ; row++{
		str, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		currentGear := &Gear{Used: false, IsPartNumber: false, Value: 0}
		for x, char := range str {
			if unicode.IsDigit(char) {
				currentGear.Value = (currentGear.Value * 10) + int(char - '0')
				schematic[Coordinate{x, row}] = currentGear
			} else {
				currentGear = &Gear{Used: false, IsPartNumber: false, Value: 0}
				if char != '.' {
					symbol := Coordinate{x, row}
					symbols = append(symbols, symbol)
				}
			}
		}
	}

	// Find the gears that are part numbers
	for _, coord := range symbols {
		for _, surrounding := range coord.Surrounding() {
			if gear, ok := schematic[surrounding]; ok {
				fmt.Println("Found gear at ", surrounding)
				gear.IsPartNumber = true

				if !gear.Used {
					fmt.Println("Gear is not used, adding to part numbers: ", gear.Value)

					gear.Used = true
					partOneSum += gear.Value
					partNumbers = append(partNumbers, *gear)
				} else {
					fmt.Printf("Gear %d, %d is used, skipping [%d]\n", surrounding.X, surrounding.Y, gear.Value)
				}
			}
		}
	}
	fmt.Println("Puzzle Three Part 1: ", partOneSum)
	fmt.Println("Puzzle Three Part 2: ", partTwoSum)
}

var (
	onlyNumbers = regexp.MustCompile(`\d+`)
	onlySymbols = regexp.MustCompile(`[^.0-9]`)
)

func main() {
	GearRatios("input3_full.txt")
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run main.go input.txt")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid file:", err)
		os.Exit(1)
	}
	defer f.Close()

	numberLocations := map[BCoordinate]**int{}
	symbolLocations := []BCoordinate{}
	var row, sum int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		data := scanner.Bytes()

		numberIndexes := onlyNumbers.FindAllIndex(data, -1)
		for _, location := range numberIndexes {
			index := location[0]
			size := location[1]
			ParseNumber(numberLocations, data, row, index, size)
		}

		symbolIndexes := onlySymbols.FindAllIndex(data, -1)
		for _, location := range symbolIndexes {
			symbolLocations = append(symbolLocations, BCoordinate{Row: row, Col: location[0]})
		}
		row++
	}

	direction := []bool{false, true, false, false, true, true, false, false}
	inc := []int{1, 1, -1, -1, -1, -1, 1, 1}
	for _, coord := range symbolLocations {
		for i, incRow := range direction {
			if incRow {
				coord.Row += inc[i]
			} else {
				coord.Col += inc[i]
			}
			if v, ok := numberLocations[coord]; ok && *v != nil {
				sum += **v
				*v = nil
			}
		}
	}

	fmt.Fprintln(os.Stdout, "The sum of all of the part numbers in the engine schematic is:", sum)
}

func ParseNumber(locations map[BCoordinate]**int, data []byte, row, startIdx, size int) {
	value := new(int)
	for i, m := size-1, 1; i >= startIdx; i, m = i-1, m*10 {
		*value += int(data[i]-'0') * m
		locations[BCoordinate{Row: row, Col: i}] = &value
	}
}

type BCoordinate struct {
	Col int
	Row int
}

func old_main() {
	//CubeConundrum("input2.txt")
}
