package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	forrest, err := NewForrest()
	if err != nil {
		panic(fmt.Errorf("failed to crete forrest: %w", err))
	}

	increments := [][]int{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}
	result := 1
	for _, increment := range increments {
		treeHits, err := forrest.traverse(increment[0], increment[1])
		if err != nil {
			panic(fmt.Errorf("error while traversing: %w", err))
		}

		fmt.Printf("Hitting %d trees on the way down with increments [%d, %d]\n", treeHits, increment[0], increment[1])
		result *= treeHits
	}

	fmt.Printf("Final result: %d", result)
}

type Forrest struct {
	width  int
	height int
	trees  []string
}

func NewForrest() (*Forrest, error) {
	data, err := ioutil.ReadFile("inputs/day_03")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	lines := strings.Split(string(data), "\n")

	return &Forrest{
		width:  len(lines[0]),
		height: len(lines),
		trees:  lines,
	}, nil
}

func (f Forrest) traverse(incrementX int, incrementY int) (int, error) {
	x, y := 0, 0
	stepCounter, treeCounter := 0, 0
	for !f.traverseFinished(y) {
		stepCounter++
		//fmt.Printf("Step #%d [%d, %d] ", stepCounter, x, y)
		x, y = f.getNextPosition(x, y, incrementX, incrementY)
		isTree, err := f.isTree(x, y)
		if err != nil {
			return 0, err
		}
		//fmt.Printf("-> [%d, %d], tree found: %v\n", x, y, isTree)
		if isTree {
			treeCounter++
		}
	}

	return treeCounter, nil
}

func (f Forrest) traverseFinished(y int) bool {
	return y >= f.height-1
}

func (f Forrest) getNextPosition(currentX int, currentY int, incrementX int, incrementY int) (int, int) {
	return (currentX + incrementX) % f.width, currentY + incrementY
}

func (f Forrest) isTree(x int, y int) (bool, error) {
	char := string(f.trees[y][x])
	if char == "#" {
		return true, nil
	}
	if char == "." {
		return false, nil
	}
	return false, fmt.Errorf("unsupported character [%s] on position %d, %d", char, x, y)
}
