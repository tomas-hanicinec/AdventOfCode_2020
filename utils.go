package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

func ReadLines(filePath string) []string {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %w", err))
	}

	return strings.Split(string(data), "\n")
}

func StringsToInts(strings []string) []int {
	result := make([]int, len(strings))
	for i, stringVal := range strings {
		intVal, err := strconv.Atoi(stringVal)
		if err != nil {
			panic(fmt.Errorf("failed to convert item string %s to int: %w", stringVal, err))
		}
		result[i] = intVal
	}

	return result
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

func StringsToBytes(strings []string) [][]byte {
	result := make([][]byte, len(strings))
	for i, stringVal := range strings {
		result[i] = []byte(stringVal)
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

func AbsInt(val int) int {
	return int(math.Abs(float64(val)))
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int64, integers ...int64) int64 {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}
