package main

import (
	"fmt"
	"sort"
)

func main() {
	adapterChain := NewAdapters()

	// Part I.
	d1, d2, d3 := adapterChain.getGapDistribution()
	fmt.Printf("Adapter gap distribution: [%d, %d, %d] -> result: %d\n", d1, d2, d3, d1*d3)

	// Part II.
	result := adapterChain.getCombinations()
	fmt.Printf("Number of possible chain combinations: %d\n", result)
}

type AdapterChain []int

func NewAdapters() AdapterChain {
	adapters := StringsToInts(ReadLines("inputs/day_10.txt"))
	sort.Ints(adapters)
	adapters = append(adapters, adapters[len(adapters)-1]+3) // Add the internal device adapter
	return append([]int{0}, adapters...)                     // Add the built-in charging socket
}

func (a AdapterChain) getGapDistribution() (int, int, int) {
	gapCounter := make([]int, 3)
	previousValue := a[0] // built-in charging socket
	for i := 1; i < len(a); i++ {
		diff := a[i] - previousValue
		if diff < 1 || diff > 3 {
			panic(fmt.Errorf("difference between [%d] and [%d] is [%d], cannot chain all the adapters", previousValue, a[i], diff))
		}
		gapCounter[diff-1]++
		previousValue = a[i]
	}

	return gapCounter[0], gapCounter[1], gapCounter[2]
}

func (a AdapterChain) getCombinations() int64 {
	sequenceStartIndex := 1
	totalCombinations := int64(1)
	for i := 1; i < len(a)-1; i++ {
		// split the array to sub-sequences divided by numbers differing by 3 (such sub-sequences are independent to each other in terms of chain combinations)
		if a[i]-a[i-1] > 2 {
			// end of sequence -  process the sequence separately
			// each sequence has at least one solution (leave as-is, do not skip any items), more solutions / possible skips multiply the number of all solutions
			totalCombinations *= a.getSequenceLinkCombinations(sequenceStartIndex, sequenceStartIndex-1, i+1) // sequence contains the last number, with the difference of 3
			sequenceStartIndex = i + 1
		}
	}

	return totalCombinations
}

func (a AdapterChain) getSequenceLinkCombinations(currentIndex int, previousIndex int, sequenceEndIndex int) int64 {
	if a[currentIndex]-a[previousIndex] > 3 {
		return 0 // this gap is too big, such a chain is not possible
	}

	if currentIndex == sequenceEndIndex {
		return 1 // successful chain constructed (we're at the last item, which is internal device adapter)
	}

	// we are in the middle of chain, process there are 2 possible paths from here
	withCurrentValue := a.getSequenceLinkCombinations(currentIndex+1, currentIndex, sequenceEndIndex)
	withoutCurrentValue := a.getSequenceLinkCombinations(currentIndex+1, previousIndex, sequenceEndIndex)

	return withCurrentValue + withoutCurrentValue
}
