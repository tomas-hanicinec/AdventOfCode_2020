package main

import (
	"fmt"
)

func main() {
	forrest := NewForrest()

	increments := [][]int{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}
	result := 1
	for _, increment := range increments {
		treeHits := forrest.traverse(increment[0], increment[1])
		fmt.Printf("Hitting %d trees on the way down with increments [%d, %d]\n", treeHits, increment[0], increment[1])
		result *= treeHits
	}

	fmt.Printf("Final result: %d\n", result)
}

type Forrest struct {
	width  int
	height int
	trees  []string
}

func NewForrest() *Forrest {
	lines := ReadLines("inputs/day_03.txt")

	return &Forrest{
		width:  len(lines[0]),
		height: len(lines),
		trees:  lines,
	}
}

func (f Forrest) traverse(incrementX int, incrementY int) int {
	x, y := 0, 0
	stepCounter, treeCounter := 0, 0
	for !f.traverseFinished(y) {
		stepCounter++
		//fmt.Printf("Step #%d [%d, %d] ", stepCounter, x, y)
		x, y = f.getNextPosition(x, y, incrementX, incrementY)
		if f.isTree(x, y) {
			treeCounter++
		}
	}

	return treeCounter
}

func (f Forrest) traverseFinished(y int) bool {
	return y >= f.height-1
}

func (f Forrest) getNextPosition(currentX int, currentY int, incrementX int, incrementY int) (int, int) {
	return (currentX + incrementX) % f.width, currentY + incrementY
}

func (f Forrest) isTree(x int, y int) bool {
	char := string(f.trees[y][x])
	if char == "#" {
		return true
	}
	if char == "." {
		return false
	}
	panic(fmt.Errorf("unsupported character [%s] on position %d, %d", char, x, y))
}
