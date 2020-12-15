package main

import "fmt"

var StartingNumbers = []int{5, 1, 9, 18, 13, 8, 0}

func main() {
	fmt.Printf("%dth number spoken: %d\n", 2020, getNthNumber(2020, StartingNumbers))         // Part I.
	fmt.Printf("%dth number spoken: %d\n", 30000000, getNthNumber(30000000, StartingNumbers)) // Part II.
}

func getNthNumber(n int, startingNumbers []int) int {
	history := make(NumberHistory)
	spokenNumber := 0
	for i := 0; i < n-1; i++ {
		if i < len(StartingNumbers) {
			if spokenNumber != 0 {
				panic(fmt.Errorf("invalid input list - number [%d] is repeated", StartingNumbers[i-1]))
			}
			spokenNumber = startingNumbers[i] // start by feeding the starting numbers into the history
		}
		spokenNumber = history.add(spokenNumber, i) // add the current number (and get the next one)
	}

	return spokenNumber
}

type NumberHistory map[int]int

func (m NumberHistory) add(number int, iteration int) int {
	lastOccurrence, exists := m[number]
	if exists {
		m[number] = iteration
		return iteration - lastOccurrence
	}

	m[number] = iteration
	return 0
}
