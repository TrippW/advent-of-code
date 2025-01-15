package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"image/color"
	"image/png"
	"image"
	"github.com/TrippW/advent-of-code/utils"
)

func SolveDay7() {
	input := utils.ReadFile("input_7.txt")
	cnt := 0
	for _, line := range input {
		splitLine := strings.Split(line, ":")
		targetStr, valStr := splitLine[0], strings.Trim(splitLine[1], " ")
		target, err := strconv.Atoi(targetStr)
		if err != nil {
			panic(err)
		}

		values := make([]int, 0)

		for _, val := range strings.Split(valStr, " ") {
			if v, err := strconv.Atoi(val); err == nil {
				values = append(values, v)
			} else {
				panic(err)
			}
		}


		for i := 0; float64(i) < math.Pow(3, float64(len(values) - 1)); i++ {
			total := values[0]
			for j := 1; j < len(values); j++ {
				op := (i / int(math.Pow(3, float64(j - 1)))) % 3
				switch op {
				case 0:
					total += values[j]
				case 1:
					total *= values[j]
				case 2:
					str := fmt.Sprintf("%d%d", total, values[j])
					if v, err := strconv.Atoi(str); err == nil {
						total = v
					} else {
						panic(err)
					}
				default:
					panic("Invalid operation")
				}
			}
			if total == target {
				cnt += total
				break
			}
		}
	}

	fmt.Println("7.1:", cnt)
}

func GetValidAnitNodes(coord Coordinate2D, dx, dy int, maxMapCoord Coordinate2D) []Coordinate2D {
	antiNodes := make([]Coordinate2D, 0)
	for x, y := coord.X + dx, coord.Y + dy; x >= 0 && y >= 0 && x < maxMapCoord.X && y < maxMapCoord.Y; x, y = x + dx, y + dy {
		antiNodes = append(antiNodes, Coordinate2D{x, y})
	}
	return antiNodes
}

func SolveDay8() {
	input := utils.ReadFile("input_8.txt")
	reqNodeCoordList := make(map[rune][]Coordinate2D)
	seen := make(map[Coordinate2D]bool)
	for i, line := range input {
		for j, c := range line {
			if c != '.' {
				coord := Coordinate2D{Y:i, X:j}
				reqNodeCoordList[c] = append(reqNodeCoordList[c], coord)
			}
		}
	}

	cnt := 0
	mapBounds := Coordinate2D{len(input), len(input[0])}
	for _, v := range reqNodeCoordList {
		for i, coord := range v {
			for _, coord2 := range v[i+1:] {
				dx := coord2.X - coord.X
				dy := coord2.Y - coord.Y
				coord2AntiNodes := GetValidAnitNodes(coord2, dx, dy, mapBounds)
				for _, antiNode := range coord2AntiNodes {
					if _, ok := seen[antiNode]; !ok {
						cnt++
						seen[antiNode] = true
					}
				}
				if _, ok := seen[coord]; !ok {
					cnt++
					seen[coord] = true
				}
				if _, ok := seen[coord2]; !ok {
					cnt++
					seen[coord2] = true
				}

				coordAntiNodes := GetValidAnitNodes(coord, -dx, -dy, mapBounds)
				for _, antiNode := range coordAntiNodes {
					if _, ok := seen[antiNode]; !ok {
						cnt++
						seen[antiNode] = true
					}
				}
				if _, ok := seen[coord2]; !ok {
					cnt++
					seen[coord2] = true
				}
				if _, ok := seen[coord]; !ok {
					cnt++
					seen[coord] = true
				}
			}
		}
	}
	fmt.Println("8.2:", len(seen))
	if len(seen) <= 1258  {
		for _, line := range input {
			fmt.Println(line)
		}
		fmt.Println()
		for k := range seen {
			input[k.Y] = input[k.Y][:k.X] + "." + input[k.Y][k.X+1:]
		}
		for _, line := range input {
			fmt.Println(line)
		}
		panic("Number too low")
	}
}

func CreateArrayOfSizeAndConstValue(sizeStr string, value int) ([]int, int) {
	var fileBlock []int
	v, err := strconv.Atoi(sizeStr)

	if err == nil {
		fileBlock = make([]int, v)
		for j := 0; j < v; j++ {
			fileBlock[j] = value
		}
	} else {
		panic(err)
	}
	return fileBlock, v
}

func HardCompactData(start, end *FileData) *FileData{
	//VisualizeDay9FilePointers(start)
	for ep := end; ep.PrevFile != nil; ep = ep.PrevFile{
		if ep.IsFree || ep.Size == 0 || ep.Moved{
			continue
		}

		for sp := start; sp != nil; sp = sp.NextFile {
			if sp == ep {
				break
			}
			if !sp.IsFree {
				continue
			}
			if sp.Size >= ep.Size {
				// Move ep to sp, removing space from sp
				sp.Size -= ep.Size

				nextFile := &FileData{ep.Size, -1, true, ep.PrevFile, ep.NextFile, false}
				// Create new empty space at ep, this may require compaction since before and after may both be free
				ep.PrevFile.NextFile = nextFile

				// Move ep to before sp
				sp.PrevFile.NextFile = ep
				ep.PrevFile, ep.NextFile, sp.PrevFile = sp.PrevFile, sp, ep

				ep.Moved = true

				ep = nextFile

				if ep.PrevFile != nil && ep.PrevFile.IsFree {
					ep.PrevFile.NextFile = ep.NextFile
					ep.PrevFile.Size += ep.Size
					ep = ep.PrevFile
				}

				if ep.NextFile != nil && ep.NextFile.IsFree {
					ep.Size += ep.NextFile.Size
					ep.NextFile = ep.NextFile.NextFile
					if ep.NextFile != nil {
						ep.NextFile.PrevFile = ep
					}
				}

				break
			}
		}
	}
	//VisualizeDay9FilePointers(start)

	return start
}

func ChecksumFileData(data *FileData) int {
	checksum := 0
	pos := 0
	for d := data; d != nil; d = d.NextFile{
		for i := 0; i < d.Size; i++ {
			if !d.IsFree {
				checksum += (pos * d.FileNumber)
			}
			pos++
		}
	}
	return checksum
}

func CompactData(data []int) []int {
	i := 0
	j := len(data) - 1
	for i < j {
		if data[i] != -1 {
			i++
			continue
		}
		if data[j] == -1 {
			j--
			continue
		}
		data[i] = data[j]
		data[j] = -1
		i++
		j--
	}
	return data
}


func CalculateChecksum(data []int) int {
	checksum := 0
	for i, v := range data {
		if v > 0 {
			checksum += i * v
		}
	}
	return checksum
}

func VisualizeDay9Data(data []int) {
	for i := 0; i < len(data); i += 64 {
		for j := 0; j < 64; j++ {
			if data[i+j] == -1 {
				fmt.Print(".")
			} else {
				fmt.Print(data[i+j])
			}
		}
		fmt.Println()
	}
}

func VisualizeDay9FilePointers(first *FileData) {
	for f := first; f != nil; f = f.NextFile {
		VisualizeDay9FileData(*f)
	}
	fmt.Println()
}


func VisualizeDay9FileData(data FileData) {
		var char rune
		if data.IsFree {
			char = '.'
		} else {
			char = rune('0' + (data.FileNumber % 36))
			if char > '9' {
				char = 'A' + (char - '9' - 1)
			}
		}
		for j := 0; j < data.Size; j++ {
			fmt.Print(string(char))
		}
		if data.NextFile != nil && data.Size > 0 {
			fmt.Print("")
		}
}

func VisualizeDay9FileDataList(data []FileData) {
	for _, d := range data {
		VisualizeDay9FileData(d)
		fmt.Println()
	}
}

type FileData struct {
	Size int
	FileNumber int 
	IsFree bool
	PrevFile *FileData
	NextFile *FileData
	Moved bool
}


func SolveDay9() {
	input := utils.ReadFile("input_9.txt")[0]
	isFileBlock := true
	data := make([]int, 0)
	fileNumber := 0
	size := 0
	var firstFile *FileData
	var prevFile *FileData
	for _, c := range input {
		var value int
		if isFileBlock {
			value = fileNumber
		} else {
			value = -1
		}

		fileBlock, block_size := CreateArrayOfSizeAndConstValue(string(c), value)
		var f FileData

		if isFileBlock {
			f = FileData{block_size, fileNumber, false, prevFile, nil, false}
			fileNumber++
		} else {
			f = FileData{block_size, -1, true, prevFile, nil, false}
		}
		isFileBlock = !isFileBlock

		if block_size == 0 {
			continue
		}

		if prevFile != nil {
			prevFile.NextFile = &f
		} else {
			firstFile = &f
		}
		prevFile = &f


		size += block_size
		data = append(data, fileBlock...)
	}

	data = CompactData(data)

	fmt.Println("9.1:", CalculateChecksum(data))
	fmt.Println("9.2:", ChecksumFileData(HardCompactData(firstFile, prevFile)))
}

func SumTrailHeads(input [][]int, i, j, target int, seen map[Coordinate2D]int) int {
	if i < 0 || j < 0 || i >= len(input) || j >= len(input[0]) {
		return 0
	}
	if input[i][j] != target {
		return 0
	}
	if v, ok := seen[Coordinate2D{i, j}]; ok {
		seen[Coordinate2D{i, j}] = v + 1
		return 0
	}
	if input[i][j] == 9 {
		seen[Coordinate2D{i, j}] = 1
		return 1
	}
	target++
	return SumTrailHeads(input, i + 1, j, target, seen) + 
	SumTrailHeads(input, i, j + 1, target, seen) + 
	SumTrailHeads(input, i - 1, j, target, seen) + 
	SumTrailHeads(input, i, j - 1, target, seen) 
}

func SolveDay10() {
	input := utils.StrListTo2DIntList(utils.ReadFile("input_10.txt"))

	sum := 0
	sum2 := 0
	for i, line := range input {
		for j, c := range line {
			if c == 0 {
				m := make(map[Coordinate2D]int)
				sum += SumTrailHeads(input, i, j, 0, m)
				for _, v := range m {
					sum2 += v
				}
			}
		}
	}	

	fmt.Println("10.1:", sum)
	fmt.Println("10.2:", sum2)
}

type Stone struct {
	Value int
	Children []*Stone
}

func (s *Stone) Process() {
	numDigits := int(math.Log10(float64(s.Value))) + 1
	if s.Value == 0 {
		s.Value = 1
	} else if numDigits % 2 == 0  {
		digitsDenom := int(math.Pow10(numDigits / 2))
		RightValue := int(s.Value % digitsDenom)
		s.Value = int(s.Value / digitsDenom)
		newRight := Stone{RightValue, make([]*Stone, 0)}
		s.Children = append(s.Children, &newRight)
	} else {
		s.Value *= 2024
	}
}

func InputToStones(input []int) []*Stone {
	var stones []*Stone
	for _, v := range input {
		stone := Stone{v, make([]*Stone, 0)}
		stones = append(stones, &stone)
	}
	return stones
}

type StoneDirection int

func ProcessStones(stones []*Stone) {
	for _, stone := range stones {
		ProcessStones(stone.Children)
	}
	for _, stone := range stones {
		stone.Process()
	}
}

func CountStones(stones []*Stone) int {
	if stones == nil {
		return 0
	}
	cnt := len(stones)
	
	for _, stone := range stones {
		cnt += CountStones(stone.Children)
		fmt.Print(stone.Value, " ")
	}

	return cnt
}


func SolveDay11() {
	input := utils.ReadFile("test_11.txt")
	data := utils.StrListToIntList(utils.SplitInputs(input, " ", false)[0])
	stone := InputToStones(data)
//	for i := 0; i < 6; i++ {
//		ProcessStones(stone)
//	}
//	fmt.Println("\n11.a:", CountStones(stone))
//	for i := 0; i < 19; i++ {
//		ProcessStones(stone)
//	}
//	fmt.Println("11.b:", CountStones(stone))

	input = utils.ReadFile("input_11.txt")
	data = utils.StrListToIntList(utils.SplitInputs(input, " ", false)[0])
	stone = InputToStones(data)

	ProcessStones(stone)
	fmt.Println("\n11.1:", CountStones(stone))
	panic("temp stop")

	ProcessStones(stone)
	fmt.Println("11.2:", CountStones(stone))
}


type GardenPlant struct {
	Plant string
	Visited bool
	CornersProcessed bool
}

type Edges struct {
	Left, Right, Up, Down bool
}

type GardenPlot struct {
	Area int
	Perimeter int
	Edges Edges
	Sides int
}

func NewGardenPlot() *GardenPlot {
	return &GardenPlot{0, 0, Edges{false, false, false, false}, 0}
}

type GridDirection int
const (
	Up GridDirection = iota
	Down
	Left
	Right
	UpLeft
	UpRight
	DownLeft
	DownRight
	Self
)

func (g GridDirection) String() string {
	switch g {
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Left:
		return "Left"
	case Right:
		return "Right"
	case UpLeft:
		return "UpLeft"
	case UpRight:
		return "UpRight"
	case DownLeft:
		return "DownLeft"
	case DownRight:
		return "DownRight"
	case Self:
		return "Self"
	}
	return "Invalid"
}

func GetAdjacent2D[T any](grid [][]*T, i, j int, direction GridDirection) *T {
	y, x := direction.toYX(i, j)
	if x < 0 || y < 0 || y >= len(grid) || x >= len(grid[0]) {
		return nil
	}
	return grid[y][x]
}

func (d GridDirection) reverse() GridDirection {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	case UpLeft:
		return DownRight
	case UpRight:
		return DownLeft
	case DownLeft:
		return UpRight
	case DownRight:
		return UpLeft
	case Self:
		return Self
	}
	panic("Invalid direction")
}

func (d GridDirection) toCoordinate(c Coordinate2D) Coordinate2D {
	switch d {
	case Up:
		return Coordinate2D{c.X, c.Y - 1}
	case Down:
		return Coordinate2D{c.X, c.Y + 1}
	case Left:
		return Coordinate2D{c.X - 1, c.Y}
	case Right:
		return Coordinate2D{c.X + 1, c.Y}
	case UpLeft:
		return Coordinate2D{c.X - 1, c.Y - 1}
	case UpRight:
		return Coordinate2D{c.X + 1, c.Y - 1}
	case DownLeft:
		return Coordinate2D{c.X - 1, c.Y + 1}
	case DownRight:
		return Coordinate2D{c.X + 1, c.Y + 1}
	case Self:
		return c
	}
	panic("Invalid direction")
}

func (d GridDirection) toYX(y, x int) (int, int) {
	switch d {
	case Up:
		return y - 1, x
	case Down:
		return y + 1, x
	case Left:
		return y, x - 1
	case Right:
		return y, x + 1
	case UpLeft:
		return y - 1, x - 1
	case UpRight:
		return y - 1, x + 1
	case DownLeft:
		return y + 1, x - 1
	case DownRight:
		return y + 1, x + 1
	case Self:
		return y, x
	}
	panic("Invalid direction")
}

func GetGardenPlotArea(grid [][]*GardenPlant, i, j int, plant string) *GardenPlot {
	cur_plant := GetAdjacent2D(grid, i, j, Self)
	if cur_plant == nil || cur_plant.Plant != plant {
		return nil
	}
	if cur_plant.Visited {
		return NewGardenPlot()
	}
	cur_plant.Visited = true
	plot := NewGardenPlot()
	plot.Area = 1

	left, right, up, down := GetGardenPlotArea(grid, i, j - 1, plant), GetGardenPlotArea(grid, i, j + 1, plant), GetGardenPlotArea(grid, i - 1, j, plant), GetGardenPlotArea(grid, i + 1, j, plant)
	dirs := map[GridDirection]*GardenPlot{Left: left, Right: right, Up: up, Down: down}
	for k, v := range dirs {
		if v != nil {
			plot.Area += v.Area
			plot.Perimeter += v.Perimeter
			plot.Sides += v.Sides
		} else {
			plot.Perimeter++
			switch k {
			case Left:
				plot.Edges.Left = true
			case Right:
				plot.Edges.Right = true
			case Up:
				plot.Edges.Up = true
			case Down:
				plot.Edges.Down = true
			}
		}
	}

	// Count Corners
	mySides := 0
	if plot.Edges.Left && plot.Edges.Up {
		mySides++
	}
	if plot.Edges.Left && plot.Edges.Down {
		mySides++
	}
	if plot.Edges.Right && plot.Edges.Up {
		mySides++
	}
	if plot.Edges.Right && plot.Edges.Down {
		mySides++
	}
	plot.Sides += mySides
	if mySides == 4 {
		return plot
	}


	upLeft, upRight, downLeft, downRight := GetAdjacent2D(grid, i, j, UpLeft), GetAdjacent2D(grid, i, j, UpRight), GetAdjacent2D(grid, i, j, DownLeft), GetAdjacent2D(grid, i, j, DownRight)

	if (up != nil && right == nil) || (up == nil && right != nil) {
		if upRight != nil && upRight.Plant == plant && !upRight.CornersProcessed {
			plot.Sides++
		}
	}
	if (up != nil && left == nil) || (up == nil && left != nil) {
		if upLeft != nil && upLeft.Plant == plant && !upLeft.CornersProcessed {
			plot.Sides++
		}
	}
	if (down != nil && right == nil) || (down == nil && right != nil) {
		if downRight != nil && downRight.Plant == plant && !downRight.CornersProcessed {
			plot.Sides++
		}
	}
	if (down != nil && left == nil) || (down == nil && left != nil) {
		if downLeft != nil && downLeft.Plant == plant && !downLeft.CornersProcessed {
			plot.Sides++
		}
	}

	cur_plant.CornersProcessed = true

	return plot
}

func FormGardenPlots(grid [][]*GardenPlant) []*GardenPlot {
	plots := make([]*GardenPlot, 0)
	for i, row := range grid {
		for j, plant := range row {
			if plant.Visited {
				continue
			}
			plot := GetGardenPlotArea(grid, i, j, plant.Plant)
			plots = append(plots, plot)
		}
	}
	return plots
}

func SolveDay12() {
	input := utils.ReadFile("input_12.txt")
	grid := make([][]*GardenPlant, len(input))
	for i, line := range input {
		plants := strings.Split(line, "")
		grid[i] = make([]*GardenPlant, len(plants))
		for j, plant := range plants {
			grid[i][j] = &GardenPlant{plant, false, false}
		}
	}

	plots := FormGardenPlots(grid)

	cost := 0
	discount_cost := 0
	for _, plot := range plots {
		cost += plot.Area * plot.Perimeter
		discount_cost += plot.Area * plot.Sides
	}
	fmt.Printf("12.1: %v expects 1374934, %v\n", cost, cost == 1374934)
	fmt.Printf("12.2: %v expects 841078, %v\n", discount_cost, discount_cost == 841078)
}

type Point2D struct {
	X, Y int
}

func getClawCost(a_move, b_move Point2D, a_cost, b_cost int, target Point2D) int {
	x_intersect := (b_move.X * target.Y - b_move.Y * target.X) / (a_move.Y * b_move.X - a_move.X * b_move.Y)
	y_intersect := (a_move.X * target.Y - a_move.Y * target.X) / (a_move.X * b_move.Y - a_move.Y * b_move.X)

	final_spot := Point2D{x_intersect * a_move.X + y_intersect * b_move.X, x_intersect * a_move.Y + y_intersect * b_move.Y}
	if final_spot == target {
		return int(x_intersect) * a_cost + int(y_intersect) * b_cost
	}
	return 0
}

type ClawInput struct {
	A, B Point2D
	Target Point2D
}

func SolveDay13() {
	start := time.Now()
	input := utils.ReadFile("input_13.txt")
	clawInput := make([]ClawInput, 0)
	claw2Input := make([]ClawInput, 0)
	buffer := make([]string, 3)
	for i, line := range input {
		if i % 4 == 3 {
			continue
		}
		buffer[i % 4] = line
		if i % 4 == 2 {
			var a_x, a_y, b_x, b_y, t_x, t_y int
			fmt.Sscanf(buffer[0], "Button A: X+%d, Y+%d", &a_x, &a_y)
			fmt.Sscanf(buffer[1], "Button B: X+%d, Y+%d", &b_x, &b_y)
			fmt.Sscanf(buffer[2], "Prize: X=%d, Y=%d", &t_x, &t_y)
			c := ClawInput{Point2D{a_x, a_y}, Point2D{b_x, b_y}, Point2D{t_x , t_y}}
			clawInput = append(clawInput, c)
			c2 := ClawInput{Point2D{a_x, a_y}, Point2D{b_x, b_y}, Point2D{t_x + 10000000000000, t_y + 10000000000000}}
			claw2Input = append(claw2Input, c2)
		}
	}
	trackTime(start, "13 Parsing")
	
	start = time.Now()
	total_cost_1 := 0;
	for _, claw := range clawInput {
		tokens := getClawCost(claw.A, claw.B, 3, 1, claw.Target)
		total_cost_1 += tokens
	}
	trackTime(start, "13.1 Calculating")

	start = time.Now()
	total_cost_2 := 0
	for _, claw := range claw2Input {
		tokens := getClawCost(claw.A, claw.B, 3, 1, claw.Target)
		total_cost_2 += tokens
	}
	trackTime(start, "13.2 Calculating")
	fmt.Println("13.1:", total_cost_1)
	fmt.Println("13.2:", total_cost_2)
}

type BathroomSecurity struct {
	Width, Height int
}

type BathroomRobot struct {
	Position Point2D
	Velocity Point2D
}

func moveRobots(robots []BathroomRobot, bathroom BathroomSecurity, steps int) []BathroomRobot {
	for j, robot := range robots {
		dx, dy := robot.Velocity.X, robot.Velocity.Y
		if dx < 0 {
			dx = bathroom.Width + dx
		}
		if dy < 0 {
			dy = bathroom.Height + dy
		}
		robots[j].Position.X = (robot.Position.X + dx * steps) % bathroom.Width
		robots[j].Position.Y = (robot.Position.Y + dy * steps) % bathroom.Height
	}
	return robots
}

type Quadrant int 

const (
	QuadrantTopLeft = iota
	QuadrantTopRight
	QuadrantBottomLeft
	QuadrantBottomRight
)


func VisualizeRobots(robots []BathroomRobot, bathroom BathroomSecurity) {
	grid := make([][]int, bathroom.Height)
	for i := 0; i < bathroom.Height; i++ {
		grid[i] = make([]int, bathroom.Width)
	}
	for _, robot := range robots {
		grid[robot.Position.Y][robot.Position.X] = grid[robot.Position.Y][robot.Position.X] + 1
	}

	for _, row := range grid {
		for _, v := range row {
			if v > 0 {
				fmt.Print(v)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func robotsToQuadrantCounts(robots []BathroomRobot, bathroom BathroomSecurity) map[Quadrant]int {
	res := make(map[Quadrant]int)

	for _, robot := range robots {
		if robot.Position.X < bathroom.Width / 2 {
			if robot.Position.Y < bathroom.Height / 2 {
				res[QuadrantTopLeft]++
			} else if robot.Position.Y > bathroom.Height / 2{
				res[QuadrantBottomLeft]++
			}
		} else if robot.Position.X > bathroom.Width / 2 {
			if robot.Position.Y < bathroom.Height / 2 {
				res[QuadrantTopRight]++
			} else if robot.Position.Y > bathroom.Height / 2{
				res[QuadrantBottomRight]++
			}
		}
	}
	return res
}

func CalcRobotSafetyProduct(robots []BathroomRobot, bathroom BathroomSecurity, steps int) int {
	robots = moveRobots(robots, bathroom, steps)
	robotsGrid := robotsToQuadrantCounts(robots, bathroom)

	product := robotsGrid[QuadrantTopLeft]
	product *= robotsGrid[QuadrantTopRight]
	product *= robotsGrid[QuadrantBottomLeft]
	product *= robotsGrid[QuadrantBottomRight]
	return product
}

func inputToRobots(input []string) []BathroomRobot {
	robots := make([]BathroomRobot, 0)
	for _, line := range input {
		if line == "" || line[0] != 'p' {
			continue
		}
		x, y, dx, dy := 0, 0, 0, 0
		fmt.Sscanf(line, "p=%d,%d v=%d,%d", &x, &y, &dx, &dy)
		robots = append(robots, BathroomRobot{Point2D{x, y}, Point2D{dx, dy}})
	}
	return robots
}

func RobotsToImage(robots []BathroomRobot, img *image.RGBA) *image.RGBA {
	for _, robot := range robots {
		img.Set(robot.Position.X, robot.Position.Y, color.Black)
	}
	return img
}

func SolveDay14() {
	input := utils.ReadFile("test_14.txt")
	testBathroom := BathroomSecurity{11, 7}

	robots := inputToRobots(input)
	product := CalcRobotSafetyProduct(robots, testBathroom, 100)
	fmt.Println("14.test:", product)

	input = utils.ReadFile("input_14.txt")
	realBathroom := BathroomSecurity{101, 103}
	robots = inputToRobots(input)
	if _, err := os.ReadDir("outputs"); err != nil {
		os.Mkdir("outputs", 0755)
	}
	if _, err := os.ReadDir("outputs/day14"); err != nil {
		os.Mkdir("outputs/day14", 0755)
	}

	for i := 0; i < 10000; i++ {
		f, err := os.Create(fmt.Sprintf("outputs/day14/%d.png", i))
		if err != nil {
			panic(err)
		}
		robot_image := image.NewRGBA(image.Rect(0, 0, realBathroom.Width, realBathroom.Height))
		robot_image = RobotsToImage(robots, robot_image)
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		encoder.Encode(f, robot_image)
		robots = moveRobots(robots, realBathroom, 1)
		f.Close()
	}
	
	product = CalcRobotSafetyProduct(robots, realBathroom, 100)
	fmt.Println("14.1:", product)
}

func isFullyTransparentPng(img image.Image) bool {
    for x := img.Bounds().Min.X; x < img.Bounds().Dx(); x++ {
        for y := img.Bounds().Min.Y; y < img.Bounds().Dy(); y++ {
            _, _, _, alpha := img.At(x, y).RGBA()
            if alpha != 0 {
                return false
            }
        }
    }
    return true
}

type Warehouse struct {
	Width, Height int
	Grid [][]WarehouseObjects
	Robot Coordinate2D
}

type WarehouseObjects rune

const (
	WarehouseWall = '#'
	WarehouseOpen = '.'
	WarehouseRobot = '@'
	WarehouseCrate = 'O'
	WarehouseCrateLeft = '['
	WarehouseCrateRight = ']'
)

func (w WarehouseObjects) String() string {
	switch w {
	case WarehouseWall:
		return "Wall"
	case WarehouseOpen:
		return "Open"
	case WarehouseRobot:
		return "Robot"
	case WarehouseCrate:
		return "Crate"
	case WarehouseCrateLeft:
		return "CrateLeft"
	case WarehouseCrateRight:
		return "CrateRight"
	}
	return string(w)
}

func (w Warehouse) getNextOpenCoordinate(move GridDirection, simple bool) (Coordinate2D, error) {
	checkPositions := []Coordinate2D{w.Robot}
	seen := make(map[Coordinate2D]bool)
	steps := 0

	y, x := w.Robot.Y, w.Robot.X
	for len(checkPositions) > 0 {
		steps++
		var coord Coordinate2D
		coord, checkPositions = checkPositions[0], checkPositions[1:]
		if seen[coord] {
			continue
		}
		seen[coord] = true
		y, x = move.toYX(coord.Y, coord.X)

		if y < 0 || x < 0 || y >= w.Height || x >= w.Width {
			return Coordinate2D{}, errors.New("Out of bounds")
		}
		switch w.Grid[y][x] {
		case WarehouseWall:
			return Coordinate2D{}, errors.New("Wall")
		case WarehouseOpen:
			if simple {
				return Coordinate2D{x, y}, nil
			}
		case WarehouseCrate:
			checkPositions = append(checkPositions, Coordinate2D{x, y})
		case WarehouseCrateLeft:
			checkPositions = append(checkPositions, Coordinate2D{x, y})
			checkPositions = append(checkPositions, Coordinate2D{x + 1, y})
		case WarehouseCrateRight:
			checkPositions = append(checkPositions, Coordinate2D{x, y})
			checkPositions = append(checkPositions, Coordinate2D{x - 1, y})
		case WarehouseRobot:
			fmt.Println(w)
			fmt.Println("Robot", move, w.Robot, x, y, steps)
			panic("Should be unreachable")
		}
	}

	return Coordinate2D{x, y}, nil
}

func (w Warehouse) getSortedMoves(move GridDirection) []Coordinate2D {
	ToMove := []Coordinate2D{w.Robot}
	ToCheck := []Coordinate2D{w.Robot}
	Seen := make(map[Coordinate2D]bool)
	for len(ToCheck) > 0 {
		var check Coordinate2D
		check, ToCheck = ToCheck[0], ToCheck[1:]
		if Seen[check] {
			continue
		}
		Seen[check] = true
		y, x := move.toYX(check.Y, check.X)

		if y < 0 || x < 0 || y >= w.Height || x >= w.Width {
			panic("Out of bounds")
		}
		coord := Coordinate2D{x, y}

		switch w.Grid[y][x] {
		case WarehouseWall:
			return make([]Coordinate2D, 0)
		case WarehouseOpen:
			ToMove = append(ToMove, check)
		case WarehouseCrate:
			ToCheck = append(ToCheck, coord)
			ToMove = append(ToMove, check)
		case WarehouseCrateLeft:
			ToCheck = append(ToCheck, coord)
			ToCheck = append(ToCheck, Right.toCoordinate(coord))

			ToMove = append(ToMove, coord)
			ToMove = append(ToMove, Right.toCoordinate(coord))
		case WarehouseCrateRight:
			ToCheck = append(ToCheck, coord)
			ToCheck = append(ToCheck, Left.toCoordinate(coord))

			ToMove = append(ToMove, coord)
			ToMove = append(ToMove, Left.toCoordinate(coord))
		case WarehouseRobot:
			fmt.Println("Robot", y, x)
			fmt.Println(w)
		}
	}

	ToMove = SortMoves(ToMove, move)
	return ToMove
}

func (w Warehouse) moveRobots(moves []GridDirection, useSimpleMovement bool) Warehouse {
	for _, move := range moves {
		coord, err := w.getNextOpenCoordinate(move, useSimpleMovement)
		if err != nil {
			continue
		}

		// Cheat and shuffle crates
		if useSimpleMovement {
			nextY, nextX := move.toYX(w.Robot.Y, w.Robot.X)
			if w.Grid[nextY][nextX] == WarehouseCrate {
				w.Grid[coord.Y][coord.X] = WarehouseCrate
			}

			// Move robot
			w.Grid[w.Robot.Y][w.Robot.X] = WarehouseOpen
			w.Grid[nextY][nextX] = WarehouseRobot
			w.Robot = Coordinate2D{nextX, nextY}
		} else {
			seen := make(map[Coordinate2D]bool)
			ToMove := w.getSortedMoves(move)

			for i := len(ToMove) - 1; i >= 0; i-- {
				if seen[ToMove[i]] {
					continue
				}
				seen[ToMove[i]] = true
				moveThis := ToMove[i]
				intoThis := move.toCoordinate(moveThis)
				w.Grid[intoThis.Y][intoThis.X], w.Grid[moveThis.Y][moveThis.X] = w.Grid[moveThis.Y][moveThis.X], w.Grid[intoThis.Y][intoThis.X]
			}
			// This doesn't seem to track perfectly, full scan
			// w.Robot = move.toCoordinate(ToMove[0])

			foundRobot := false
			for y, row := range w.Grid {
				for x, c := range row {
					if c == WarehouseRobot {
						w.Robot = Coordinate2D{x, y}
						foundRobot = true
						break
					}
				}
				if foundRobot {
					break
				}
			}
		}
	}
	return w
}

func SortMoves(moves []Coordinate2D, move GridDirection) []Coordinate2D {
	var sortFunc func(i, j int) bool
	switch move {
	case Up:
		sortFunc = func(i, j int) bool {
			return moves[i].Y > moves[j].Y
		}
	case Down:
		sortFunc = func(i, j int) bool {
			return moves[i].Y < moves[j].Y
		}

	case Left:
		sortFunc = func(i, j int) bool {
			return moves[i].X > moves[j].X
		}
	case Right:
		sortFunc = func(i, j int) bool {
			return moves[i].X < moves[j].X
		}
	}
	sort.Slice(moves, sortFunc)
	return moves
}

func (w Warehouse) String() string {
	var b strings.Builder
	for _, row := range w.Grid {
		b.WriteString(string(row))
		b.WriteString("\n")
	}
	return b.String()
}

func parseDay15Input(input []string) (Warehouse, []GridDirection) {
	var warehouse Warehouse
	var moves []GridDirection

	buildWarehouse := true
	for i, line := range input {
		if buildWarehouse {
			if line == "" {
				buildWarehouse = false
				continue
			}
			if i == 0 {
				warehouse.Width = len(line)
			}
			warehouse.Height++
			warehouse.Grid = append(warehouse.Grid, []WarehouseObjects(line))
			if strings.ContainsRune(line, WarehouseRobot) {
				warehouse.Robot = Coordinate2D{strings.IndexRune(line, WarehouseRobot), i}
			}
		} else {
			for _, c := range line {
				switch c {
				case '^':
					moves = append(moves, Up)
				case 'v':
					moves = append(moves, Down)
				case '<':
					moves = append(moves, Left)
				case '>':
					moves = append(moves, Right)
				default:
					panic(fmt.Sprintf("Invalid move: %v,\n %v\n", c, line))
				}
			}
		}
	}

	return warehouse, moves
}

func parseDay15InputAsWide(input []string) (Warehouse, []GridDirection) {
	var warehouse Warehouse
	var moves []GridDirection

	buildWarehouse := true
	for i, line := range input {
		if buildWarehouse {
			if line == "" {
				buildWarehouse = false
				fmt.Println("End of warehouse")
				fmt.Println(warehouse.Robot)
				fmt.Println(warehouse.Width, warehouse.Height)
				fmt.Println(warehouse)
				if warehouse.Grid[warehouse.Robot.Y][warehouse.Robot.X] != WarehouseRobot {
					fmt.Println(warehouse.Grid[warehouse.Robot.Y][warehouse.Robot.X])
					for i := warehouse.Robot.Y - 1; i <= warehouse.Robot.Y + 1; i++ {
						for j := warehouse.Robot.X - 1; j <= warehouse.Robot.X + 1; j++ {
							fmt.Printf("%v", warehouse.Grid[i][j])
						}
						fmt.Println()
					}

					panic("Robot not found")
				}
				continue
			}
			if i == 0 {
				warehouse.Width = len(line) * 2
			}
			warehouse.Height++
			row := make([]WarehouseObjects, 0)
			for j, c := range line {
				switch c {
				case WarehouseWall:
					row = append(row, WarehouseWall, WarehouseWall)
				case WarehouseOpen:
					row = append(row, WarehouseOpen, WarehouseOpen)
				case WarehouseCrate:
					row = append(row, WarehouseCrateLeft, WarehouseCrateRight)
				case WarehouseRobot:
					row = append(row, WarehouseRobot, WarehouseOpen)
					warehouse.Robot = Coordinate2D{j * 2, i}
				}
			}
			warehouse.Grid = append(warehouse.Grid, row)
		} else {
			for _, c := range line {
				switch c {
				case '^':
					moves = append(moves, Up)
				case 'v':
					moves = append(moves, Down)
				case '<':
					moves = append(moves, Left)
				case '>':
					moves = append(moves, Right)
				default:
					panic(fmt.Sprintf("Invalid move: %v,\n %v\n", c, line))
				}
			}
		}
	}

	return warehouse, moves
}

func (w Warehouse) SumCoordinates() int {
	sum := 0
	for y, row := range w.Grid {
		for x, c := range row {
			if c == WarehouseCrate || c == WarehouseCrateLeft{
				sum += y * 100 + x
			}
		}
	}
	return sum
}

func SolveDay15() {
	input := utils.ReadFile("input_15.txt")
	warehouse, moves := parseDay15Input(input)
	warehouse = warehouse.moveRobots(moves, false)
	fmt.Println(warehouse)
	fmt.Println("15.1:", warehouse.SumCoordinates())
	warehouse, moves = parseDay15InputAsWide(input)
	warehouse = warehouse.moveRobots(moves, false)
	fmt.Println(warehouse)
	fmt.Println("15.2:", warehouse.SumCoordinates())
}

func main() {
	defer trackTime(time.Now(), "main")
	funcs := []func(){
		day1.Solve,
		SolveDay2,
		SolveDay3,
		SolveDay4,
		SolveDay5,
		SolveDay6,
		SolveDay7,
		SolveDay8,
		SolveDay9,
		SolveDay10,
		SolveDay11,
		SolveDay12,
		SolveDay13,
		SolveDay14,
		SolveDay15,
	}
	n := len(funcs)
	if len(os.Args) > 1 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil {
			if v > 0 && v <= n {
				n = v
			}
		}
	}
	funcs[n-1]()
}
