package main

import (
	"fmt"
	"sort"
)

func main() {
	boardingTickets := readBoardingTickets()

	maxIdIndex := -1
	maxId := 0
	seatIds := make([]int, len(boardingTickets))
	for i, ticket := range boardingTickets {
		currentId := ticket.getSeatId()
		if currentId > maxId {
			maxId = currentId
			maxIdIndex = i
		}
		seatIds[i] = currentId
	}
	fmt.Printf("Max seat ID: %d (boardingPass #%d)\n", maxId, maxIdIndex) // Part I.

	sort.Ints(seatIds)
	previousSeatId := seatIds[0]
	for i := 1; i < len(seatIds); i++ {
		if seatIds[i]-previousSeatId == 2 {
			fmt.Printf("My seat ID: %d\n", seatIds[i]-1) // Part II.
			return
		}
		previousSeatId = seatIds[i]
	}

	panic("no free seat found\n")
}

type BoardingTicket struct {
	number     int
	binaryCode string
}

func readBoardingTickets() []*BoardingTicket {
	lines := ReadLines("inputs/day_05.txt")

	result := make([]*BoardingTicket, len(lines))
	for i, line := range lines {
		result[i] = newBoardingTicket(line)
	}

	return result
}

func newBoardingTicket(binaryCode string) *BoardingTicket {
	return &BoardingTicket{
		binaryCode: binaryCode,
	}
}

func (bt BoardingTicket) getRow() int {
	rowInterval := interval{
		a: 0,
		b: 127,
	}
	return parseBinaryInterval(bt.binaryCode[0:7], rowInterval, 'F', 'B')
}

func (bt BoardingTicket) getColumn() int {
	columnInterval := interval{
		a: 0,
		b: 7,
	}
	return parseBinaryInterval(bt.binaryCode[7:10], columnInterval, 'L', 'R')
}

func parseBinaryInterval(binaryCode string, fullInterval interval, codeLow uint8, codeUp uint8) int {
	currentInterval := fullInterval
	//fmt.Printf("interval [%d, %d]\n", currentInterval.a, currentInterval.b)
	for i := 0; i < len(binaryCode); i++ {
		if binaryCode[i] == codeLow {
			currentInterval = currentInterval.getLowerHalf()
		} else if binaryCode[i] == codeUp {
			currentInterval = currentInterval.getUpperHalf()
		} else {
			panic(fmt.Errorf("unknown input letter %s in binary code %s", string(binaryCode[i]), binaryCode))
		}
		//fmt.Printf("interval [%d, %d] after [%s]\n", currentInterval.a, currentInterval.b, string(binaryCode[i]))
	}
	if currentInterval.a != currentInterval.b {
		panic(fmt.Errorf("invalid interval [%d, %d] in the end of parsing for binary code %s", currentInterval.a, currentInterval.b, binaryCode))
	}
	return currentInterval.a
}

func (bt BoardingTicket) getSeatId() int {
	return bt.getRow()*8 + bt.getColumn()
}

type interval struct {
	a int
	b int
}

func (i interval) getLowerHalf() interval {
	return interval{
		a: i.a,
		b: i.b - (i.b-i.a)/2 - 1, // 7 / 2 = 3 in Go...
	}
}

func (i interval) getUpperHalf() interval {
	return interval{
		a: i.a + (i.b-i.a)/2 + 1,
		b: i.b,
	}
}
