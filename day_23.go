package main

import (
	"fmt"
	"strconv"
)

func main() {
	// Part I.
	gameV1 := NewCupGame([]int{4, 6, 7, 5, 2, 8, 1, 9, 3}, 9)
	for round := 1; round <= 100; round++ {
		gameV1.playRound()
	}
	cupOrder := ""
	startCup := gameV1.cups[1]
	for currentCup := startCup.next; currentCup != startCup; currentCup = currentCup.next {
		cupOrder += strconv.Itoa(currentCup.value)
	}
	fmt.Printf("Cup order after 100 rounds of V1 game: %s\n", cupOrder)

	// Part II.
	gameV2 := NewCupGame([]int{4, 6, 7, 5, 2, 8, 1, 9, 3}, 1000000)
	for round := 1; round <= 10000000; round++ {
		gameV2.playRound()
	}
	cup1, _ := gameV2.cups[1]
	fmt.Printf("Cups immediately following cup 1 after 10M rounds of V2 game: [%d, %d], result: %d\n", cup1.next.value, cup1.next.next.value, cup1.next.value*cup1.next.next.value)
}

type CupGame struct {
	cups        map[int]*Cup
	currentCup  *Cup
	roundNumber int
}

type Cup struct {
	value int
	next  *Cup
}

func NewCupGame(cups []int, gameSize int) CupGame {
	cupList := make(map[int]*Cup, gameSize)
	var previousCup *Cup
	lastValue := 0
	for i := 0; i < gameSize; i++ {
		value := i + 1 // values start from 1...
		if i < len(cups) {
			value = cups[i]
		}
		cupList[value] = &Cup{
			value: value,
			next:  nil,
		}
		lastValue = value
		if previousCup != nil {
			previousCup.next = cupList[value] // link with previous cup
		}
		previousCup = cupList[value]
	}

	cupList[lastValue].next = cupList[cups[0]] // close the circle

	return CupGame{
		cups:        cupList,
		currentCup:  cupList[cups[0]],
		roundNumber: 0,
	}
}

func (g *CupGame) playRound() {
	g.roundNumber++

	//get destination cup
	destinationCup := g.getDestinationCup()

	//move the 3 selected cups to their new position
	g.moveCups(destinationCup)

	//get the new current cup
	g.currentCup = g.currentCup.next
}

func (g *CupGame) getDestinationCup() *Cup {

	selectedValues := make(map[int]struct{})
	for selected := g.currentCup.next; selected != g.getLastSelected().next; selected = selected.next {
		selectedValues[selected.value] = struct{}{}
	}
	destinationCupValue := g.currentCup.value - 1
	for i := 0; i < len(selectedValues)+1; i++ {
		if destinationCupValue == 0 {
			destinationCupValue = len(g.cups) // reached zero, we want the highest remaining number so reset back to the game size
		}
		if _, ok := selectedValues[destinationCupValue]; !ok {
			return g.cups[destinationCupValue]
		}
		destinationCupValue--
	}

	panic(fmt.Errorf("destination cup not found (logic or input error, should not happen)"))
}

func (g *CupGame) moveCups(destinationCup *Cup) {
	first := g.currentCup.next
	last := g.getLastSelected()
	g.currentCup.next = last.next   // remove selected block from list
	last.next = destinationCup.next // link selected block after the destination cup
	destinationCup.next = first
}

func (g *CupGame) getLastSelected() *Cup {
	return g.currentCup.next.next.next
}
