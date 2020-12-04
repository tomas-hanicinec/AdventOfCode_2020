package main

import (
	"fmt"
	"io/ioutil"
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
