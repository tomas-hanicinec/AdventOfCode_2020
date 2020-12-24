package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	flipInstructions := getFlipInstructions()
	floor := make(Floor)

	// Part I.
	for _, instruction := range flipInstructions {
		row, column := instruction.getTargetCoordinates()
		key := getKey(row, column)
		if _, ok := floor[key]; !ok {
			floor[key] = false // init this tile, in the beginning it's white
		}
		floor[key] = !floor[key] // flip tile
	}
	floor = floor.reduceToBlack()
	fmt.Printf("Number of black floor tiles in the beginning: %d\n", len(floor))

	// Part II.
	for i := 0; i < 100; i++ {
		floor = floor.expand()
		floor = floor.execIteration()
		floor = floor.reduceToBlack()
	}
	fmt.Printf("Number of black floor tiles after 100 iterations: %d\n", len(floor))
}

type FlipInstruction string

func getFlipInstructions() []FlipInstruction {
	lines := ReadLines("inputs/day_24.txt")
	result := make([]FlipInstruction, len(lines))
	for i := range lines {
		result[i] = FlipInstruction(lines[i])
	}
	return result
}

func (i FlipInstruction) getTargetCoordinates() (int, int) {
	row, column := 0, 0
	for _, step := range i.parse() {
		increment, ok := stepIncrements[step]
		if !ok {
			panic(fmt.Errorf("invalid step [%s]", step))
		}
		row += increment[0]
		column += increment[1]
	}

	return row, column
}

func (i FlipInstruction) parse() []string {
	result := make([]string, 0)
	position := 0
	for position < len(i) {
		end := position + 1 // default is one letter ("e" or "w")
		if i[position] == 'n' || i[position] == 's' {
			end++ // 2-letter direction
		}
		result = append(result, string(i[position:end]))
		position = end
	}

	return result
}

type Floor map[string]bool

func (f Floor) execIteration() Floor {
	newFloor := make(Floor)
	for key, isBlack := range f {
		blackNeighboursCount := f.getBlackNeighbourCount(parseKey(key))
		willBeBlack := isBlack
		if isBlack {
			if blackNeighboursCount == 0 || blackNeighboursCount > 2 {
				willBeBlack = false
			}
		} else {
			if blackNeighboursCount == 2 {
				willBeBlack = true
			}
		}

		newFloor[key] = willBeBlack
	}

	return newFloor
}

func (f Floor) expand() Floor {
	newFloor := make(Floor)
	for key, value := range f {
		newFloor[key] = value
		for _, neighbour := range getNeighbours(parseKey(key)) {
			neighbourKey := getKey(neighbour[0], neighbour[1])
			if _, ok := f[neighbourKey]; !ok {
				newFloor[neighbourKey] = false // add the white tile
			}
		}
	}

	return newFloor
}

func (f Floor) reduceToBlack() Floor {
	newFloor := make(Floor)
	for key, value := range f {
		if value {
			newFloor[key] = value
		}
	}

	return newFloor
}

func (f Floor) getBlackNeighbourCount(row, column int) int {
	blackCounter := 0
	for _, neighbour := range getNeighbours(row, column) {
		if value, ok := f[getKey(neighbour[0], neighbour[1])]; ok && value {
			blackCounter++
		}
	}

	return blackCounter
}

var stepIncrements = map[string][2]int{
	"e":  {0, 2},
	"w":  {0, -2},
	"se": {1, 1},
	"sw": {1, -1},
	"ne": {-1, 1},
	"nw": {-1, -1},
}

func getKey(row, column int) string {
	return strconv.Itoa(row) + ":" + strconv.Itoa(column)
}

func parseKey(key string) (int, int) {
	parts := strings.Split(key, ":")
	row, errRow := strconv.Atoi(parts[0])
	column, errColumn := strconv.Atoi(parts[1])
	if errRow != nil || errColumn != nil {
		panic(fmt.Errorf("invalid key [%s]", key))
	}

	return row, column
}

func getNeighbours(row, column int) [6][2]int {
	var result [6][2]int
	i := 0
	for _, increment := range stepIncrements {
		result[i] = [2]int{row + increment[0], column + increment[1]}
		i++
	}
	return result
}
