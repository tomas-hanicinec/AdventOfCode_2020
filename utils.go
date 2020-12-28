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

func AbsInt(val int) int {
	return int(math.Abs(float64(val)))
}

func isInInterval(min, max, index int) bool {
	return min <= index && index <= max
}
