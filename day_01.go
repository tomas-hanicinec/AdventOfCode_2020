package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

const TargetValue = 2020

func main() {
	input, err := getInputIntegers()
	if err != nil {
		panic(fmt.Errorf("failed to get input: %w", err))
	}
	sort.Ints(input)

	for _, a := range input {
		for _, b := range input {
			if a+b > TargetValue {
				break // input array is sorted, no need to go further
			}
			for _, c := range input {
				if a+b+c > TargetValue {
					break // input array is sorted, no need to go further
				}
				if a+b+c == TargetValue {
					fmt.Printf("Solution found: %d + %d + %d = %d, %d * %d * %d = %d\n", a, b, c, a+b+c, a, b, c, a*b*c)
					return
				}
			}
		}
	}

	fmt.Println("No solution found")
}

func getInputIntegers() ([]int, error) {
	data, err := ioutil.ReadFile("inputs/day_01")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	stringsArray := strings.Split(string(data), "\n")
	result := make([]int, len(stringsArray))
	for i, stringVal := range stringsArray {
		intVal, err := strconv.Atoi(stringVal)
		if err != nil {
			return nil, fmt.Errorf("failed to convert item string %s to int: %w", stringVal, err)
		}
		result[i] = intVal
	}
	return result, nil
}
