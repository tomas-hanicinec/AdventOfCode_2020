package main

import (
	"fmt"
	"math"
	"strconv"
)

const preambleSize = 25

func main() {
	xmasCode := NewXmasCode()

	// Part I.
	weaknessSumIndex, weaknessSum := xmasCode.findWeaknessSum(preambleSize)
	fmt.Printf("Weakness sum found on index [%d]: %d\n", weaknessSumIndex, weaknessSum)

	// Part II.
	i, j, weakness, found := xmasCode.findWeakness(weaknessSum)
	if !found {
		panic(fmt.Errorf("no weakness found"))
	}
	fmt.Printf("Weakness found between indexes [%d, %d], final number: %d\n", i, j, weakness)
}

type XmasCode []int64

func NewXmasCode() XmasCode {
	return StringsToLongints(ReadLines("inputs/day_09.txt"))
}

func (xc XmasCode) findWeaknessSum(preambleSize int) (int, int64) {
	for i := preambleSize; i < len(xc); i++ {
		_, _, found := xc.findSum(xc[i], xc[i-preambleSize:i])
		if !found {
			return i, xc[i]
		}

	}

	panic(fmt.Errorf("no weakness sum found in XmasCode"))
}

func (xc XmasCode) findWeakness(weaknessSum int64) (int, int, int64, bool) {
	startIndex, stopIndex, sum := 0, 1, xc[0]
	for stopIndex < len(xc) {
		if sum == weaknessSum {
			// Solution found
			minValue, maxValue := GetMinMax(xc[startIndex:stopIndex])
			return startIndex, stopIndex, minValue + maxValue, true
		}

		if sum > weaknessSum {
			// Sum too big, remove the first number
			sum -= xc[startIndex]
			startIndex++
		} else {
			// Sum not big enough, add next number
			sum += xc[stopIndex]
			stopIndex++
		}
	}

	panic(fmt.Errorf("no contiguous set for XmasCode weakness found"))
}

func (xc XmasCode) findSum(sum int64, values []int64) (int64, int64, bool) {
	// O(n) complexity (2 passes)
	numberMap := make(map[int64]bool, len(values))
	remainders := make([]int64, len(values))
	for i, value := range values {
		numberMap[value] = true
		remainders[i] = sum - value
	}

	for _, remainder := range remainders {
		if _, exists := numberMap[remainder]; exists {
			return remainder, sum - remainder, true
		}
	}

	return 0, 0, false
}

func StringsToLongints(strings []string) []int64 {
	result := make([]int64, len(strings))
	for i, stringVal := range strings {
		intVal, err := strconv.ParseInt(stringVal, 10, 64)
		if err != nil {
			panic(fmt.Errorf("failed to convert item string %s to int64: %w", stringVal, err))
		}
		result[i] = intVal
	}

	return result
}

func GetMinMax(values []int64) (int64, int64) {
	if len(values) < 1 {
		panic(fmt.Errorf("cannot get Min, Max from empty array"))
	}
	minVal := int64(math.MaxInt64)
	maxVal := int64(math.MinInt64)
	for _, val := range values {
		minVal = int64(math.Min(float64(minVal), float64(val)))
		maxVal = int64(math.Max(float64(maxVal), float64(val)))
	}

	return minVal, maxVal
}
