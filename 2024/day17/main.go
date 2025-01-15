package main

import (
	"fmt"
	"strings"
	"time"
	_ "embed"
	"github.com/TrippW/advent-of-code/utils"
)

//go:embed input.txt
var rawInput string

//go:embed test.txt
var rawTest string

type Register string

const (
	A Register = "A"
	B Register = "B"
	C Register = "C"
)

func (r Register) String() string {
	return string(r)
}

type Command int

const (
	ADV Command = iota
	BXL 
	BST 
	JNZ
	BXC
	OUT
	BDV
	CDV
)

func (c Command) String() string {
	switch c {
	case ADV:
		return "A DIV"
	case BXL:
		return "B XOR"
	case BST:
		return "Modulo"
	case JNZ:
		return "JNZ"
	case BXC:
		return "B XOR C"
	case OUT:
		return "OUT"
	case BDV:
		return "B DIV"
	case CDV:
		return "C DIV"
	}
	panic("Invalid Command")
}

type ComboCode int

const (
	CC_Zero ComboCode = iota
	CC_One
	CC_Two
	CC_Three
	CC_ComboA
	CC_ComboB
	CC_ComboC
	COMBO_RESERVED
)

func (c ComboCode) String() string {
	switch c {
	case CC_Zero:
		return "Zero"
	case CC_One:
		return "One"
	case CC_Two:
		return "Two"
	case CC_Three:
		return "Three"
	case CC_ComboA:
		return "ComboA"
	case CC_ComboB:
		return "ComboB"
	case CC_ComboC:
		return "ComboC"
	case COMBO_RESERVED:
		return "Reserved"
	}
	panic("Invalid ComboCode")
}

func (c ComboCode) validate() bool {
	if(c >= CC_Zero && c < COMBO_RESERVED) {
		return true
	}
	panic("Invalid ComboCode")
}

func (c ComboCode) ToRegister() Register {
	if c >= CC_ComboA && c <= CC_ComboC {
		return Register(rune(c - CC_ComboA) + rune(A[0]))
	}
	panic("Invalid ComboCode")
}

func (c *Computer) value(cc ComboCode) int {
	cc.validate()	
	if cc >= CC_Zero && cc <= CC_Three {
		return int(cc)
	}
	return c.Registers[cc.ToRegister()]
}

func (c *Computer) div_register(r Register, cc ComboCode) {
	c.Registers[r] = c.Registers[A] >> c.value(cc)
}

func (c *Computer) adv(code ComboCode) {
	c.div_register(A, code)
}

func (c *Computer) bxl(code ComboCode) {
	c.Registers[B] ^= int(code)
}

func (c *Computer) bst(code ComboCode) {
	c.Registers[B] = c.value(code) % 8 
}

func (c *Computer) jnz(code ComboCode) {
	if c.Registers[A] != 0 {
		c.Index = int(code)
	}
}

func (c *Computer) bxc(_ ComboCode) {
	c.Registers[B] = c.Registers[B] ^ c.Registers[C]
}

func (c *Computer) out(code ComboCode) {
	c.output = append(c.output, c.value(code) % 8)
}

func (c *Computer) bdv(code ComboCode) {
	c.div_register(B, code)
}

func (c *Computer) cdv(code ComboCode) {
	c.div_register(C, code)
}

func NewComputer(registers map[Register]int) *Computer {
	c := &Computer{
		Registers: registers,
		Commands: make(map[Command]func(ComboCode)),
		Index: 0,
		output: []int{},
	}
	c.Commands[ADV] = c.adv
	c.Commands[BXL] = c.bxl
	c.Commands[BST] = c.bst
	c.Commands[JNZ] = c.jnz
	c.Commands[BXC] = c.bxc
	c.Commands[OUT] = c.out
	c.Commands[BDV] = c.bdv
	c.Commands[CDV] = c.cdv
	return c
}

type Computer struct {
	Registers map[Register]int
	Commands map[Command]func(ComboCode)
	Index int
	output []int
}

func (c *Computer) Print() {
	data := []string{}
	for _, v := range c.output {
		data = append(data, fmt.Sprintf("%d", v))
	}
	fmt.Println(strings.Join(data, ","))
}

func (c *Computer) Execute(program []int) {
	c.Index = 0
	for c.Index < len(program) {
		command := Command(program[c.Index])
		code := ComboCode(program[c.Index + 1])
		c.Index += 2
		c.Commands[command](code)
	}
}

func (c *Computer) Reset(i int) {
	c.output = []int{}
	c.Registers[A] = i
	c.Registers[B] = 0
	c.Registers[C] = 0
}

func (c *Computer) GenerateTargetOuput(program, output []int) ([]int) {
	size := len(output)
	if size  == 0 {
		panic("Invalid output, must have at least one element")
	}
	if size == 1 {
		target := output[0]

		c.Reset(0)
		potentials := []int{}
		for i := 0; i < 8; i++ {
			c.Reset(i)
			c.Execute(program)
			if c.output[0] == target {
				potentials = append(potentials, i)
			}
		}
		return potentials
	} else {
		potentials := []int{}
		for _, out := range c.GenerateTargetOuput(program, output[1:]) {
			r := out << 3
			for i := 0; i < 8; i++ {
				c.Reset(r | i)
				c.Execute(program)
				if c.output[0] == output[0] {
					potentials = append(potentials, r | i)
				}
			}
		}
		return potentials
	}
}

func (c *Computer) GenerateInitialRegister(program []int) int{
	return utils.MinOf(c.GenerateTargetOuput(program, program)...)
}

func trackTime(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func ParseInput(input []string) ([]int, map[Register]int) {
	registers := map[Register]int{}
	var program []int
	for _, line := range input {
		if line == "" {
			continue
		}
		if strings.Contains(line, "Register") {
			var regChar rune
			var value int
			_, err := fmt.Sscanf(line, "Register %c: %d", &regChar, &value)
			if err != nil {
				fmt.Println(err)
			}
			registers[Register(regChar)] = value
		} else if strings.Contains(line, "Program") {
			var rawProgram string
			fmt.Sscanf(line, "Program: %s", &rawProgram)
			commands := strings.Split(rawProgram, ",")
			program = utils.StrListToIntList(commands)
		}
	}
	return program, registers
}

func Part1(c *Computer, program []int) {
	defer trackTime(time.Now(), "part 1")
	fmt.Println("Part 1")
	c.Execute(program)
	c.Print()
}

func Part2(c *Computer, program []int) {
	defer trackTime(time.Now(), "part 2")
	fmt.Println("Part 2")
	initialRegister := c.GenerateInitialRegister(program)
	fmt.Println("Initial Register:", initialRegister)
	fmt.Println("Validating output...")
	c.Reset(initialRegister)
	c.Execute(program)
	c.Print()
	fmt.Println(program)
}

func main() {
	defer trackTime(time.Now(), "main")
	fmt.Println("Run day 17")
	input := strings.Split(rawInput, "\n")
	program, registers := ParseInput(input)
	c := NewComputer(registers)
	Part1(c, program)
	Part2(c, program)

	fmt.Println("End day 17")
}
