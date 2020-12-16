package main

import (
	"fmt"
	"strings"
)

const activeCube = "#"
const inactiveCube = "."
const iterationCount = 6

func main() {
	// Part I.
	pd3d := NewPocketDimension3D()
	for i := 0; i < iterationCount; i++ {
		pd3d = pd3d.execBootCycle()
	}
	fmt.Printf("Total number of active cubes in 3D: %d\n", pd3d.getActiveCount())

	// Part II.
	pd4d := NewPocketDimension4D()
	for i := 0; i < iterationCount; i++ {
		pd4d = pd4d.execBootCycle()
	}
	fmt.Printf("Total number of active cubes in 4D: %d\n", pd4d.getActiveCount())
}

// -------------------------------------- 4D

type Space4D []Space3D

func NewPocketDimension4D() Space4D {
	result := make(Space4D, 1)
	result[0] = NewPocketDimension3D()

	return result
}

func (s4d *Space4D) execBootCycle() Space4D {
	newSpace := getEmptySpace4D(s4d.getWidth()+2, s4d.getHeight()+2, s4d.getDepth()+2, s4d.getSize()+2)
	for i := range newSpace {
		for j := range newSpace[i] {
			for k := range newSpace[i][j] {
				newLine := ""
				for l := range newSpace[i][j][k] {
					newLine += s4d.getNewState(i-1, j-1, k-1, l-1)
				}
				newSpace[i][j][k] = newLine
			}
		}
	}

	return newSpace
}

func (s4d *Space4D) getNewState(i, j, k, l int) string {
	return getNewState(s4d.getAdjacentActiveCount(i, j, k, l), s4d.isActive(i, j, k, l))
}

func (s4d *Space4D) isActive(i, j, k, l int) bool {
	if !isInInterval(0, s4d.getSize()-1, i) {
		return false
	}

	return (*s4d)[i].isActive(j, k, l)
}

func (s4d *Space4D) getActiveCount() int {
	count := 0
	for _, s3d := range *s4d {
		count += s3d.getActiveCount()
	}

	return count
}

func (s4d *Space4D) getAdjacentActiveCount(i, j, k, l int) int {
	offsets := [3]int{-1, 0, 1}

	counter := 0
	for _, iOffset := range offsets {
		for _, jOffset := range offsets {
			for _, kOffset := range offsets {
				for _, lOffset := range offsets {
					if s4d.isActive(i+iOffset, j+jOffset, k+kOffset, l+lOffset) {
						counter++
					}
				}
			}
		}
	}

	return counter
}

func (s4d *Space4D) getSize() int {
	return len(*s4d)
}

func (s4d *Space4D) getDepth() int {
	return (*s4d)[0].getDepth()
}

func (s4d *Space4D) getHeight() int {
	return (*s4d)[0].getHeight()
}

func (s4d *Space4D) getWidth() int {
	return (*s4d)[0].getWidth()
}

func getEmptySpace4D(width int, height int, depth int, x int) Space4D {
	result := make(Space4D, x)
	for i := range result {
		result[i] = getEmptySpace3D(width, height, depth)
	}
	return result
}

// -------------------------------------- 3D

type Space3D []Space2D

func NewPocketDimension3D() Space3D {
	lines := ReadLines("inputs/day_17.txt")

	result := make(Space3D, 1)
	result[0] = lines

	return result
}

func (s3d *Space3D) execBootCycle() Space3D {
	newSpace := getEmptySpace3D(s3d.getWidth()+2, s3d.getHeight()+2, s3d.getDepth()+2)
	for i := range newSpace {
		for j := range newSpace[i] {
			newLine := ""
			for k := range newSpace[i][j] {
				newLine += s3d.getNewState(i-1, j-1, k-1)
			}
			newSpace[i][j] = newLine
		}
	}

	return newSpace
}

func (s3d *Space3D) getNewState(i, j, k int) string {
	return getNewState(s3d.getAdjacentActiveCount(i, j, k), s3d.isActive(i, j, k))
}

func (s3d *Space3D) isActive(i, j, k int) bool {
	if !isInInterval(0, s3d.getDepth()-1, i) || !isInInterval(0, s3d.getHeight()-1, j) || !isInInterval(0, s3d.getWidth()-1, k) {
		return false
	}

	return (*s3d)[i][j][k] == activeCube[0]
}

func (s3d *Space3D) getActiveCount() int {
	count := 0
	for _, s2d := range *s3d {
		for _, line := range s2d {
			count += strings.Count(line, activeCube)
		}
	}

	return count
}

func (s3d *Space3D) getAdjacentActiveCount(i, j, k int) int {
	offsets := [3]int{-1, 0, 1}

	counter := 0
	for _, iOffset := range offsets {
		for _, jOffset := range offsets {
			for _, kOffset := range offsets {
				if s3d.isActive(i+iOffset, j+jOffset, k+kOffset) {
					counter++
				}
			}
		}
	}

	return counter
}

func (s3d *Space3D) getDepth() int {
	return len(*s3d)
}

func (s3d *Space3D) getHeight() int {
	return len((*s3d)[0])
}

func (s3d *Space3D) getWidth() int {
	return len((*s3d)[0][0])
}

func getEmptySpace3D(width int, height int, depth int) Space3D {
	result := make(Space3D, depth)
	for i := range result {
		result[i] = getEmptySpace2D(width, height)
	}
	return result
}

// -------------------------------------- 2D

type Space2D []string

func getEmptySpace2D(width int, height int) Space2D {
	result := make(Space2D, height)
	for i := range result {
		result[i] = strings.Repeat(inactiveCube, width)
	}
	return result
}

// -------------------------------------- HELPERS

func getNewState(adjacentActiveCount int, isActive bool) string {
	if isActive {
		adjacentActiveCount-- // getAdjacentActiveCount calculates all the active fields including the one in the middle
		if adjacentActiveCount == 2 || adjacentActiveCount == 3 {
			return activeCube
		}
		return inactiveCube
	} else {
		if adjacentActiveCount == 3 {
			return activeCube
		}
		return inactiveCube
	}
}
