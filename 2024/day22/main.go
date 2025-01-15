package main

import (
	"fmt"
	"strings"
	_ "embed"
	"github.com/TrippW/advent-of-code/utils"
)

//go:embed input.txt
var rawInput string

//go:embed test.txt
var rawTest string

type Buyer struct {
	secrets []int
	priceChanges []int // Bound -9 -> 9, always n - 1 size where n is len(secrets)
	sequences map[string]Sequence
}

type Sequence struct {
	sequence []int
	value int
}

func createSequenceKey(k int) rune {
	return rune('a' + k + 10)
}

func genBuyers(seed, times int) Buyer {
	buyer	:= Buyer{
		secrets: make([]int, times + 1),
		priceChanges: make([]int, times),
		sequences: make(map[string]Sequence),
	}
	buyer.secrets[0] = seed

	var a, b, c rune
	for i := 1; i <= times; i++ {
		seed = ((seed << 6) ^ seed) % 16777216
		seed = ((seed >> 5) ^ seed) % 16777216
		seed = ((seed << 11) ^ seed) % 16777216
		buyer.secrets[i] = seed
		buyer.priceChanges[i-1] = (buyer.secrets[i] % 10) - (buyer.secrets[i-1] % 10)
		curKey := createSequenceKey(buyer.priceChanges[i-1])
		if i >= 4 {
			sequence := string([]rune{a, b, c, curKey})
			if _, ok := buyer.sequences[sequence]; !ok {
				buyer.sequences[sequence] = Sequence{
					sequence: buyer.priceChanges[i - 4:i],
					value: buyer.secrets[i] % 10,
				}
			}
		}
		a, b, c = b, c, curKey
	}

	return buyer
}

func main() {
	defer utils.TrackTimeFromNow("day 22")
	inputFile := rawInput

	secrets := utils.StrListToIntList(strings.Split(inputFile, "\n"))
	fmt.Println("Number of Buyers: ", len(secrets))
	results := make([]int, len(secrets))
	buyers := make([]Buyer, len(secrets))
	sum := 0
	sequences := make(map[string]Sequence)
	bestSequence := Sequence{
		value: 0,
	}
	for i, secret := range secrets {
		buyers[i] = genBuyers(secret, 2000)
		results[i] = buyers[i].secrets[2000]
		for k, v := range buyers[i].sequences {
			if s, ok := sequences[k]; !ok {
				sequences[k] = v
			} else {
				s.value += v.value
				sequences[k] = s
			}
			if sequences[k].value > bestSequence.value {
				bestSequence = sequences[k]
			}
		}
		sum += results[i]
	}
	fmt.Println("Sum of secret numbers part 1:", sum)

	fmt.Println("Best sequence:", bestSequence)
}
