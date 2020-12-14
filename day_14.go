package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	memoryV1 := make(Memory)
	memoryV2 := make(Memory)

	currentMask := BitMask("")
	for _, line := range ReadLines("inputs/day_14.txt") {
		if isBitMask(line) {
			currentMask = NewBitMask(line) // just set the mask
			continue
		}
		index, value := parseWriteInstruction(line)
		if currentMask == "" {
			panic(fmt.Errorf("mask not initialized befor the first [write] operation"))
		}
		memoryV1.add(index, currentMask.applyToMemoryValue(value))
		for _, indexWithMask := range currentMask.applyToMemoryIndex(index) {
			memoryV2.add(indexWithMask, value)
		}
	}

	fmt.Printf("Summary of memory value for V1: %d\n", memoryV1.sum()) // Part I.
	fmt.Printf("Summary of memory value for V2: %d\n", memoryV2.sum()) // Part II.
}

type Memory map[int64]int64

func (m Memory) add(index int64, value int64) {
	m[index] = value
}

func (m Memory) sum() int64 {
	sum := int64(0)
	for _, val := range m {
		sum += val
	}
	return sum
}

type BitMask string

func isBitMask(inputLine string) bool {
	return inputLine[0:7] == "mask = "
}

func NewBitMask(inputLine string) BitMask {
	return BitMask(strings.Split(inputLine, " = ")[1])
}

func (m BitMask) getBitNumber(bitIndex int) int {
	return len(m) - 1 - bitIndex // bits are in reverse order (0-th bit is last in mask)
}

func (m BitMask) applyToMemoryValue(value int64) int64 {
	result := value
	for i := range m {
		if m[i] == 'X' {
			continue // no change
		}
		result = setBitToValue(result, m.getBitNumber(i), m[i] == '1') // force the value of this bit
	}
	return result
}

func (m BitMask) applyToMemoryIndex(index int64) []int64 {
	// first get the modified memory index without floating bits
	baseMemoryIndex := index
	floatingBitNumbers := make([]int, 0)
	for i := range m {
		if m[i] == '1' {
			baseMemoryIndex = setBitToValue(baseMemoryIndex, m.getBitNumber(i), true)
		} else if m[i] == 'X' {
			floatingBitNumbers = append(floatingBitNumbers, m.getBitNumber(i))
		}
	}

	// start with the base index and resolve all the floating bits
	result := []int64{baseMemoryIndex}
	for _, floatingBitNumber := range floatingBitNumbers {
		for _, resultValue := range result {
			newBitValue := !getBitFromValue(resultValue, floatingBitNumber) // add the opposite bit value from the one already in this result item
			result = append(result, setBitToValue(resultValue, floatingBitNumber, newBitValue))
		}
	}

	return result
}

func parseWriteInstruction(inputLine string) (int64, int64) {
	pattern := regexp.MustCompile("^mem\\[([0-9]+)] = ([0-9]+)$")
	matches := pattern.FindStringSubmatch(inputLine)
	validationError := fmt.Errorf("intput line [%s] not a valid memory write instruction", inputLine)
	if len(matches) < 3 {
		panic(validationError)
	}

	index, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		panic(validationError)
	}

	value, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		panic(validationError)
	}

	return index, value
}
