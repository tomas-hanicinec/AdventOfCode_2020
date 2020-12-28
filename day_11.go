package main

import (
	"bytes"
	"fmt"
)

const maxIterations = 1000 // just to be sure we stop somewhere

func main() {
	println("Part I.")
	seatPlan := NewSeatPlan()
	iterationCount := seatPlan.iterateUntilStable(seatPlan.iterationTransformerV1, maxIterations)
	fmt.Printf("Seat plan stable state reached after [%d] iterations, number of occupied seats: %d\n", iterationCount, seatPlan.countOccupiedSeats())

	println("\nPart II.")
	seatPlan = NewSeatPlan()
	iterationCount = seatPlan.iterateUntilStable(seatPlan.iterationTransformerV2, maxIterations)
	fmt.Printf("Seat plan stable state reached after [%d] iterations, number of occupied seats: %d\n", iterationCount, seatPlan.countOccupiedSeats())
}

const occupiedSeat = "#"
const emptySeat = "L"
const floor = "."

type SeatPlan struct {
	plan   [][]byte
	width  int
	height int
}

type SeatTransformer func(i, j int) byte

func NewSeatPlan() *SeatPlan {
	lines := ReadLines("inputs/day_11.txt")
	return &SeatPlan{
		plan:   StringsToBytes(lines),
		width:  len(lines[0]),
		height: len(lines),
	}
}

func (sp *SeatPlan) iterateUntilStable(transformFunc SeatTransformer, maxIterations int) int {
	changed := true
	counter := 0
	for changed {
		if counter > maxIterations {
			panic(fmt.Errorf("max number of iterations [%d] reached and seat plan is still changing", maxIterations))
		}
		changed = sp.runIteration(transformFunc)
		counter++

		//fmt.Printf("\nAFTER ITERATION %d\n\n", counter)
		//seatPlan.print()
	}

	return counter
}

func (sp *SeatPlan) runIteration(transformFunc SeatTransformer) bool {
	planChanged := false
	newPlan := make([][]byte, sp.height)
	for i := range sp.plan {
		newPlan[i] = bytes.Repeat([]byte(floor), sp.width) // init with floor signs
		for j := range sp.plan[i] {
			newPlan[i][j] = transformFunc(i, j)
			planChanged = planChanged || sp.plan[i][j] != newPlan[i][j]
		}
	}
	sp.plan = newPlan

	return planChanged
}

func (sp *SeatPlan) iterationTransformerV1(i, j int) byte {
	if sp.plan[i][j] == emptySeat[0] {
		if sp.countOccupiedAdjacentSeats(i, j) == 0 {
			return occupiedSeat[0] // seat just became occupied
		}
	} else if sp.plan[i][j] == occupiedSeat[0] {
		if sp.countOccupiedAdjacentSeats(i, j) > 3 {
			return emptySeat[0] // seat just became empty
		}
	}

	return sp.plan[i][j] // no change otherwise
}

func (sp *SeatPlan) iterationTransformerV2(i, j int) byte {
	if sp.plan[i][j] == emptySeat[0] {
		if sp.countOccupiedVisibleSeats(i, j) == 0 {
			return occupiedSeat[0] // seat just became occupied
		}
	} else if sp.plan[i][j] == occupiedSeat[0] {
		if sp.countOccupiedVisibleSeats(i, j) > 4 {
			return emptySeat[0] // seat just became empty
		}
	}

	return sp.plan[i][j] // no change otherwise
}

func (sp *SeatPlan) countOccupiedAdjacentSeats(i, j int) int {
	counter := 0
	offsets := []int{-1, 0, 1}
	for _, iOffset := range offsets {
		for _, jOffset := range offsets {
			if iOffset == 0 && jOffset == 0 {
				continue // this is the main seat in the middle
			}
			newI := i + iOffset
			newJ := j + jOffset
			if sp.isOut(newI, newJ) {
				continue // out of bounds
			}

			if sp.plan[newI][newJ] == occupiedSeat[0] {
				counter++
			}
		}
	}

	return counter
}

func (sp *SeatPlan) countOccupiedVisibleSeats(i, j int) int {
	counter := 0
	offsets := []int{-1, 0, 1}
	for _, iOffset := range offsets {
		for _, jOffset := range offsets {
			if iOffset == 0 && jOffset == 0 {
				continue // this is the main seat in the middle
			}
			if sp.getVisibleSeat(i, j, iOffset, jOffset) == occupiedSeat[0] {
				counter++
			}
		}
	}

	return counter
}

func (sp *SeatPlan) getVisibleSeat(i, j int, iOffset, jOffset int) byte {
	i += iOffset
	j += jOffset
	for !sp.isOut(i, j) {
		if sp.plan[i][j] != floor[0] {
			return sp.plan[i][j] // reached the first visible seat
		}
		i += iOffset
		j += jOffset
	}

	return floor[0] // no seat in sight (just floor all the way to the edge)
}

func (sp *SeatPlan) isOut(i, j int) bool {
	return i < 0 || i >= sp.height || j < 0 || j >= sp.width
}

func (sp *SeatPlan) countOccupiedSeats() int {
	result := 0
	for i := range sp.plan {
		result += bytes.Count(sp.plan[i], []byte(occupiedSeat))
	}

	return result
}

func (sp *SeatPlan) print() {
	for i := range sp.plan {
		fmt.Println(string(sp.plan[i]))
	}
}

func StringsToBytes(strings []string) [][]byte {
	result := make([][]byte, len(strings))
	for i, stringVal := range strings {
		result[i] = []byte(stringVal)
	}

	return result
}
