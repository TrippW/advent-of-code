package main

import (
	"fmt"
	"strconv"
	"strings"
	"sort"
)

type CircuitOp = string

const (
	opAND CircuitOp = "AND"
	opOR CircuitOp = "OR"
	opXOR CircuitOp = "XOR"
	opRSHIFT CircuitOp = "RSHIFT"
	opLSHIFT CircuitOp = "LSHIFT"
	opNOT CircuitOp = "NOT"
	opConst CircuitOp = "opConst"
	opNone CircuitOp = "opNone"
	opDefault CircuitOp = "Unassigned"
)

type CircuitCommand struct {
	Id string
	Connections []*CircuitCommand
	Op CircuitOp
	Assigned bool
	Value uint
}

func (c *CircuitCommand) And(v *CircuitCommand) uint {
	return c.Calc() & v.Calc()
}

func (c *CircuitCommand) Or(v *CircuitCommand) uint {
	return c.Calc() | v.Calc()
}

func (c *CircuitCommand) LShift(v *CircuitCommand) uint {
	return c.Calc() << v.Calc()
}

func (c *CircuitCommand) RShift(v *CircuitCommand) uint {
	return c.Calc() >> v.Calc()
}

func (c *CircuitCommand) Not() uint {
	return ^c.Calc()
}

func (c *CircuitCommand) Xor(v *CircuitCommand) uint {
	return c.Calc() ^ v.Calc()
}

func (c *CircuitCommand) Calc() uint {
	if c.Op == opDefault {
		panic(fmt.Sprintf("Tried to use unvalued circuit %v", c.Id))
	}
	if c.Assigned || c.Op == opConst {
		return c.Value
	}

	var v uint
	switch c.Op {
		case opAND: v = c.Connections[0].And(c.Connections[1])
		case opOR: v = c.Connections[0].Or(c.Connections[1])
		case opNone: v = c.Connections[0].Calc()
		case opXOR: v = c.Connections[0].Xor(c.Connections[1])
		case opLSHIFT: v = c.Connections[0].LShift(c.Connections[1])
		case opRSHIFT: v = c.Connections[0].RShift(c.Connections[1])
		case opNOT: v = c.Connections[0].Not()
	}

	c.Assigned = true
	c.Value = v
	
	return v
}

func (c *CircuitCommand) Print(verbose bool) string {
	if !c.Assigned {
		c.Calc()
	}

	s := "Unknown"

	switch c.Op {
		case opConst: s = fmt.Sprintf("%s = %s(%d)", c.Id, c.Id, c.Value)
		case opXOR: s = fmt.Sprintf("%s(%d) ^ %s(%d) = %s(%d)", c.Connections[0].Id, c.Connections[0].Value, c.Connections[1].Id, c.Connections[1].Value, c.Id, c.Value)
		case opAND: s = fmt.Sprintf("%s & %s = %s(%d)", c.Connections[0].Id, c.Connections[1].Id, c.Id, c.Value)
		case opOR: s = fmt.Sprintf("%s | %s = %s(%d)", c.Connections[0].Id, c.Connections[1].Id, c.Id, c.Value)
		case opLSHIFT: s = fmt.Sprintf("%s << %s = %s(%d)", c.Connections[0].Id, c.Connections[1].Id, c.Id, c.Value)
		case opRSHIFT: s = fmt.Sprintf("%s >> %s = %s(%d)", c.Connections[0].Id, c.Connections[1].Id, c.Id, c.Value)
		case opNOT: s = fmt.Sprintf("~%s = %s(%d)", c.Connections[0].Id, c.Id, c.Value)
		case opNone: s = fmt.Sprintf("%s = %s(%d)", c.Connections[0].Id, c.Id, c.Value)
	}

	if (verbose) {
		fmt.Println(s)
	}

	return s
}

type CircuitCommandMap = map[string]*CircuitCommand

type CircuitCommandParser struct {
	Map CircuitCommandMap
}

func (cm CircuitCommandParser) FindOrDefault(k string) *CircuitCommand {
	var cmd *CircuitCommand
	if c, ok := cm.Map[k]; !ok {
		cmd = &CircuitCommand{
			Id: k,
			Op: opDefault,
			Connections: make([]*CircuitCommand, 0),
		}
		if intVal, err := strconv.Atoi(k); err == nil {
			cmd.Op = opConst
			cmd.Value = uint(intVal)
		}
		cm.Map[k] = cmd
	} else {
		cmd = c
	}

	return cmd
}

func (cm CircuitCommandParser) parseCircuitCommands(input []string) {
	for _, line := range input {
		if len(line) == 0 || line == ""{
			continue
		}
		splitLine := strings.Split(line, "->")
		fmt.Println(splitLine)
		provider, id := strings.Trim(splitLine[0], " "), strings.Trim(splitLine[1], " ")
		cmd := cm.FindOrDefault(id)
		if cmd.Op != opDefault {
			fmt.Println("Skipping", id, cmd.Op)
			continue
		}
		ops := strings.Split(provider, " ")
		switch len(ops) {
		case 1:
			val := strings.Trim(ops[0], " ")
			valCmd := cm.FindOrDefault(val)
			cmd.Op = opNone
			cmd.Connections = append(cmd.Connections, valCmd)
		case 2:
			if ops[0] != opNOT {
				panic(fmt.Sprintf("some other unary that isn't NOT %v in %v", ops[0], id))
			}
			notCmd := cm.FindOrDefault(ops[1])
			cmd.Op = opNOT
			cmd.Connections = append(cmd.Connections, notCmd)
		case 3:
			leftCmd := cm.FindOrDefault(ops[0])
			rightCmd := cm.FindOrDefault(ops[2])
			cmd.Connections = append(cmd.Connections, leftCmd)
			cmd.Connections = append(cmd.Connections, rightCmd)
			cmd.Op = ops[1]
		default:
			panic(fmt.Sprintf("Unable to handle %d ops\n", len(ops)))
		}
	}
}

func solve_7() {
	//input := readFile("7_input.txt")
	input := readFile("2024_24_input.txt")
	parser := CircuitCommandParser{
		Map: make(CircuitCommandMap),
	}
	parser.parseCircuitCommands(input)
	keys := make([]string, 0, len(parser.Map))
	for k := range parser.Map {
		if k[0] == 'z' {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	var value uint = 0
	for i := len(keys) - 1; i >= 0; i-- {
		k := keys[i]
		value <<= 1
		newValue := parser.Map[k].Calc()
		value += newValue
		fmt.Println(k, parser.Map[k].Print(false), value, newValue, fmt.Sprintf("%b", value))
	}

	fmt.Printf("7.1 Answer: %v\n", value)
	
	//fmt.Printf("7.1 Answer: %v\n", parser.Map["a"].Print(false))

	parser.Map = CircuitCommandMap{
		"b": &CircuitCommand{
			Op: opConst,
			Value: 956,
			Id: "b",
			Assigned: true,
		},
	}
}
