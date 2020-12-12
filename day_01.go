package main

import (
	"fmt"
	"strings"
)

const TargetValue = 2020

func main() {
	inputMap := NewInputMap()

	// Part I.
	sumGroupV1 := inputMap.findSumGroup(TargetValue, 2)
	fmt.Println(sumGroupV1.format("V1"))

	// Part II.
	sumGroupV2 := inputMap.findSumGroup(TargetValue, 3)
	fmt.Println(sumGroupV2.format("V2"))
}

type InputMap map[int]int

func NewInputMap() *InputMap {
	input := StringsToInts(ReadLines("inputs/day_01.txt"))
	inputMap := make(InputMap, len(input))
	for _, val := range input {
		inputMap[val] = val
	}

	return &inputMap
}

func (im InputMap) findSumGroup(groupSum int, groupSize int) SumGroupMap {
	for val := range im {
		if groupSize == 1 {
			if _, exists := im[groupSum]; exists {
				return SumGroupMap{groupSum: groupSum}
			} else {
				return nil
			}
		}
		// process remainder of the group
		if tail := im.findSumGroup(groupSum-val, groupSize-1); tail != nil {
			if _, exists := tail[val]; !exists {
				// we do not allow using one value multiple times
				tail[val] = val
				return tail
			}
		}
	}

	return nil // tried all possibilities, found nothing
}

type SumGroupMap map[int]int

func (sgm SumGroupMap) format(sumGroupName string) string {
	if sgm == nil {
		return fmt.Sprintf("%s solution not found", sumGroupName)
	}
	index, sum, prod := 0, 0, 1
	sumGroupStrings := make([]string, len(sgm))
	for val := range sgm {
		sum += val
		prod *= val
		sumGroupStrings[index] = fmt.Sprintf("%d", val)
		index++
	}
	return fmt.Sprintf(
		"%s solution found: %s = %d, %s = %d",
		sumGroupName,
		strings.Join(sumGroupStrings, " + "), sum,
		strings.Join(sumGroupStrings, " * "), prod,
	)
}
