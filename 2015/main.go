package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	re "regexp"
	"slices"
	"strings"
)

func readFile(filename string) []string {
	f, err := os.ReadFile("inputs/" + filename)
	if err != nil {
		fmt.Println(err)
	}
	return strings.Split(string(f), "\n")
}


func calcWrappingPaper(l, w, h int) int {
	sides := []int{l * w, w * h, h * l}
	return 2 * sides[0] + 2 * sides[1] + 2 * sides[2] + min(sides...)
}

func min(nums ...int) int {
	val := nums[0]
	for _, n := range nums {
		if n < val {
			val = n
		}
	}
	return val
}

func calcRibbon(l, w, h int) int {
	sides := []int{l, w, h}
	slices.Sort(sides)
	return 2 * sides[0] + 2 * sides[1] + l * w * h
}

func calcTotalWrappingPaper(lines [][]int) int {
	sum := 0
	for _, line := range lines {
		sum += calcWrappingPaper(line[0], line[1], line[2])
	}
	return sum
}

func calcTotalRibbon(lines [][]int) int {
	sum := 0
	for _, line := range lines {
		sum += calcRibbon(line[0], line[1], line[2])
	}
	return sum
}

func solve_2_1() {
	input := readFile("2_input.txt")
	var lines [][]int
	for _, line := range input {
		var l []int
		for _, s := range strings.Split(line, "x") {
			i := 0
			fmt.Sscanf(s, "%d", &i)
			l = append(l, i)
		}
		lines = append(lines, l)
	}
	res := calcTotalWrappingPaper(lines)

	fmt.Println("2.1 Answer: ", res)

	res = calcTotalRibbon(lines)
	fmt.Println("2.2 Answer: ", res)
}

func solve_3_1() {
	input := readFile("3_input.txt")[0]
	visited := make(map[string]bool)
	visited["0,0"] = true
	sX, sY := 0, 0
	rX, rY := 0, 0
	cnt := 1
	for i, c := range input {
		var x, y int
		if i % 2 == 0 {
			x, y = rX, rY
		} else {
			x, y = sX, sY
		}

		switch c {
		case '^':
			y++
			break
		case 'v':
			y--
			break
		case '>':
			x++
			break
		case '<':
			x--
			break
		}

		if i % 2 == 0 {
			rX, rY = x, y
		} else {
			sX, sY = x, y
		}

		key := fmt.Sprintf("%d,%d", x, y)
		if _, ok := visited[fmt.Sprintf("%d,%d", x, y)]; !ok {
			visited[key] = true
			cnt++
		}
	}
	fmt.Println("3.1 Answer: ", cnt)
}

func solve_4() {
	input := "ckczppom"
	valid := false
	addition := 0
	for ; !valid; addition++ {
		test := md5.Sum([]byte(fmt.Sprintf("%s%d", input, addition)))
		val := hex.EncodeToString(test[:])
		valid = val[:6] == "000000"
		if valid {
			break
		}
	}

	fmt.Println("4.1 Answer: ", addition)
}

func isNice(s string) bool {
	vowels := 0
	hasRepeat := false
	badStrings := map[byte]byte {
		'a': 'b',
		'c': 'd',
		'p': 'q',
		'x': 'y',
	}
	vowelSet := []byte("aeiou")
	for i := 0; i < len(s); i++ {
		if i != len(s) - 1 {
			if s[i] == s[i+1] {
				hasRepeat = true
			}

			if next, ok := badStrings[s[i]]; ok && (s[i+1] == next) {
				return false
			}
		}

		if slices.Contains(vowelSet, s[i]) {
			vowels++
		}
	}

	return vowels >= 3 && hasRepeat
}

func containsPair(pair, s string) bool {
	if len(s) < 2 {
		return false
	}
	for i := 0; i <= len(s) - len(pair); i++ {
		if s[i:i+len(pair)] == pair {
			return true
		}
	}

	return false
}

func isNice2(s string) bool {
	containsMirror := false
	hasPair := false
	for i := 0; i < len(s) - 2; i++ {
		containsMirror = containsMirror || s[i] == s[i+2]
		hasPair = hasPair || ((i < len(s) - 2) && containsPair(s[i:i+2], s[i+2:]))
		if containsMirror && hasPair {
			return true
		}
	}

	return false
}
	

func solve_5() {
	input := readFile("5_input.txt")

	nice := 0
	nice2 := 0
	for _, line := range input {
		if isNice(line) {
			nice++
		}
		if isNice2(line) {
			nice2++
		}
	}
	fmt.Println("5.1 Answer:", nice)
	fmt.Println("5.2 Answer:", nice2)
}

type LightCommand string

const (
	TurnOn LightCommand = "turn on"
	TurnOff LightCommand = "turn off"
	Toggle LightCommand = "toggle"
)

type Command struct {
	cmd LightCommand
	startX, startY, endX, endY int
}

func parseCommand(s string) Command {
	exp := re.MustCompile(`(?P<cmd>turn on|turn off|toggle) (?P<x0>\d+),(?P<y0>\d+) through (?P<x1>\d+),(?P<y1>\d+)`)
	data := exp.FindStringSubmatch(s)

	var x, y, a, b int
	cmd := LightCommand(data[1])
	fmt.Sscanf(data[2], "%d", &x)
	fmt.Sscanf(data[3], "%d", &y)
	fmt.Sscanf(data[4], "%d", &a)
	fmt.Sscanf(data[5], "%d", &b)

	return Command{
		cmd: cmd,
		startX: x,
		startY: y,
		endX: a,
		endY: b,
	}
}
		
func makeLightsGrid(x, y int) [][]bool {
	lights := make([][]bool, y)
	for i := 0; i < y; i++ {
		lights[i] = make([]bool, x)
	}
	return lights
}

func makeVariableLightsGrid(x, y int) [][]int {
	lights := make([][]int, y)
	for i := 0; i < y; i++ {
		lights[i] = make([]int, x)
	}
	return lights
}

func countLightsInState(lights [][]bool, state bool) int {
	cnt := 0
	for _, row := range lights {
		for _, light := range row {
			if light == state {
				cnt++
			}
		}
	}
	return cnt
}

func sumLightsState(lights [][]int) int {
	sum := 0
	for _, row := range lights {
		for _, light := range row {
			sum += light
		}
	}
	return sum
}

func processCommands(commands []Command, lights [][]bool) [][]bool {
	for _, command := range commands {
		for y := command.startY; y <= command.endY; y++ {
			for x := command.startX; x <= command.endX; x++ {
				if command.cmd == TurnOn {
					lights[y][x] = true
				} else if command.cmd == TurnOff {
					lights[y][x] = false
				} else if command.cmd == Toggle {
					lights[y][x] = !lights[y][x]
				}
			}
		}
	}

	return lights
}

func processCommandsVariable(commands []Command, lights [][]int) [][]int {
	for _, command := range commands {
		for y := command.startY; y <= command.endY; y++ {
			for x := command.startX; x <= command.endX; x++ {
				if command.cmd == TurnOn {
					lights[y][x] += 1
				} else if command.cmd == TurnOff {
					if lights[y][x] > 0 {
						lights[y][x]--
					}
				} else if command.cmd == Toggle {
					lights[y][x] += 2
				}
			}
		}
	}

	return lights
}

func solve_6() {
	input := readFile("6_input.txt")
	commands := make([]Command, len(input))
	for i, s := range input {
		commands[i] = parseCommand(s)
	}

	lights := makeLightsGrid(1000, 1000)

	lights = processCommands(commands, lights)
	fmt.Println("6.1 Answer:", countLightsInState(lights, true))

	variableLights := makeVariableLightsGrid(1000, 1000)
	fmt.Println("6.2 Answer:", sumLightsState(processCommandsVariable(commands, variableLights)))
}

func solve_8() {
	input := readFile("8_input.txt")

	exp := re.MustCompile(`\\(?:\"|\\|x[0-9a-f]{2})`)
	sum := 0
	memSum := 0
	for _, v := range input {
		sum += len(v)
		memSum += len(v)
		matches := exp.FindAllStringSubmatch(v, -1)
		matchesLen := 0
		for _, sub := range matches {
			matchesLen += len(sub[0]) - 1
		}
		memSum -= (matchesLen + 2)
	}

	fmt.Println("8.1:", sum - memSum)

	exp = re.MustCompile(`\\|"`)
	sum = 0
	memSum = 0
	for _, v := range input {
		sum += len(v)
		memSum += len(v) + 2
		matches := exp.FindAllStringSubmatch(v, -1)
		memSum += len(matches)
	}

	fmt.Println("8.2", memSum - sum)
}

type Node struct {
	Name string
	Visited bool
	Edges []Edge
}

type Edge struct {
	To *Node
	Distance int
}

func buildWeightedGraph(input []string) map[string]*Node {
	nodes := make(map[string]*Node)
	for _, line := range input {
		var from, to string
		var distance int
		fmt.Sscanf(line, "%s to %s = %d", &from, &to, &distance)
		fromNode, ok := nodes[from]
		if !ok {
			fromNode = &Node{
				Name: from,
				Edges: make([]Edge, 0),
			}
			nodes[from] = fromNode
		}
		toNode, ok := nodes[to]
		if !ok {
			toNode = &Node{
				Name: to,
				Edges: make([]Edge, 0),
			}
			nodes[to] = toNode
		}

		toEdge := Edge{
			To: toNode,
			Distance: distance,
		}
		fromEdge := Edge{
			To: fromNode,
			Distance: distance,
		}
		fromNode.Edges = append(fromNode.Edges, toEdge)
		toNode.Edges = append(toNode.Edges, fromEdge)
	}
	return nodes
}

func resetGraph(nodes map[string]*Node) {
	for _, node := range nodes {
		node.Visited = false
	}
}

func findShortestPath(nodes map[string]*Node) int {
	minDist := -1
	visitLocations := len(nodes)
	fmt.Println("Must visit all locations:", visitLocations)
	var shortestEdge *Edge
	for _, v := range nodes {
		if shortestEdge == nil {
			shortestEdge = &v.Edges[0]
		}

		for _, edge := range v.Edges {
			if edge.Distance < shortestEdge.Distance {
				shortestEdge = &edge
			}
		}
	}

	startNode := shortestEdge.To
	fmt.Println("Starting at:", startNode.Name)

	return minDist
}

func solve_9() {
	input := readFile("9_test.txt")
	nodes := buildWeightedGraph(input)
	dist :=findShortestPath(nodes)
	fmt.Println("9.1 Answer:", dist)
}

func main() {
	solve_7()
}

