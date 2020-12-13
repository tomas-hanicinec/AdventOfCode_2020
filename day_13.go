package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {
	input := ReadLines("inputs/day_13.txt")
	arrivalTime, err := strconv.ParseInt(input[0], 10, 64)
	if err != nil {
		panic(fmt.Errorf("invalid arrval time [%s]", input[0]))
	}
	schedule := NewBusSchedule(input[1])

	// Part I.
	busPeriod, waitTime := schedule.getBestWaitTime(arrivalTime)
	fmt.Printf("Best bus number is %d, will have to wait %d minutes, result: %d\n", busPeriod, waitTime, busPeriod*waitTime)

	// Part II.
	normalizedSchedule := schedule.getNormalized()
	solution := normalizedSchedule.getSolution()
	fmt.Printf("Winning departure time: %d\n", solution)
}

type BusSchedule []int64

func NewBusSchedule(inputLine string) BusSchedule {
	scheduleItems := strings.Split(inputLine, ",")
	schedule := make([]int64, len(scheduleItems))
	for i := range scheduleItems {
		period := int64(0)
		if scheduleItems[i] != "x" {
			var err error
			period, err = strconv.ParseInt(scheduleItems[i], 10, 64)
			if err != nil {
				panic(fmt.Errorf("invalid bus number [%s]", scheduleItems[i]))
			}
		}
		schedule[i] = period
	}
	return schedule
}

func (bs BusSchedule) getBestWaitTime(arrivalTime int64) (int64, int64) {
	bestBusPeriod := int64(0)
	bestWaitTime := int64(math.MaxInt64)
	for _, busPeriod := range bs {
		if busPeriod == 0 {
			continue
		}
		waitTime := busPeriod - arrivalTime%busPeriod
		if waitTime < bestWaitTime {
			bestBusPeriod = busPeriod
			bestWaitTime = waitTime
		}
	}

	return bestBusPeriod, bestWaitTime
}

// tries to shorten the schedule by joining various lines with intersecting periods together
func (bs BusSchedule) getNormalized() BusSchedule {
	// bs is immutable
	bsCopy := make(BusSchedule, len(bs))
	copy(bsCopy, bs)

	// tries to move the specified item to the new specified position in the schedule
	moveItem := func(currentIndex int, newIndex int) {
		if newIndex < 0 || newIndex >= len(bsCopy) {
			return // out of bounds, cannot move
		}
		if bsCopy[newIndex] == 0 {
			bsCopy[newIndex] = bsCopy[currentIndex] // there is nothing on this position yet, move the item here
		} else {
			bsCopy[newIndex] = LCM(bsCopy[currentIndex], bsCopy[newIndex]) // there already is another item (bus period), replace with LCM of both ( = smallest period both busses will be in this point of time together again)
		}
		bsCopy[currentIndex] = 0 // free the current item (it was moved elsewhere)

	}

	for i := range bsCopy {
		moveItem(i, i+int(bsCopy[i])) // try moving one period to the right
	}
	for i := len(bsCopy) - 1; i >= 0; i-- {
		moveItem(i, i-int(bsCopy[i])) // try moving one period to the left
	}

	return bsCopy
}

func (bs BusSchedule) getSolution() int64 {
	// it is more efficient to loop through the maximum period (less iterations)
	maxPeriod := int64(0)
	maxPeriodIndex := -1
	for i, period := range bs {
		if period > maxPeriod {
			maxPeriod = period
			maxPeriodIndex = i
		}
	}

	possibleSolutionTime := maxPeriod - int64(maxPeriodIndex) // solution time counts form the first bus arriving, not the max period bus
	for !bs.isSolution(possibleSolutionTime) {
		possibleSolutionTime += maxPeriod // try next period
	}

	return possibleSolutionTime // found!
}

func (bs BusSchedule) isSolution(solutionTime int64) bool {
	for index, period := range bs {
		if period == 0 {
			continue
		}
		shouldArrive := solutionTime + int64(index)
		if shouldArrive%period != 0 {
			return false // should arrive but does not -> this is not a solution
		}
	}

	return true // all bus lines arrived in the times they should
}
