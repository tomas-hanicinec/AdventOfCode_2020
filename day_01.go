package main

import (
	"fmt"
	"sort"
)

const TargetValue = 2020

func main() {
	input := StringsToInts(ReadLines("inputs/day_01.txt"))
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
