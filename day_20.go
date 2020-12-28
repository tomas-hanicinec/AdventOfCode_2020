package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const topIndex = 0
const rightIndex = 1
const bottomIndex = 2
const leftIndex = 3

var monsterPattern = []string{
	"..................#.",
	"#....##....##....###",
	".#..#..#..#..#..#...",
}

func main() {
	tiles := getTiles("inputs/day_20.txt")

	// Part I.
	image := NewImage(tiles)
	cornerTileIdProduct := image.tilePlacement[0][0].id *
		image.tilePlacement[0][image.imageSize-1].id *
		image.tilePlacement[image.imageSize-1][0].id *
		image.tilePlacement[image.imageSize-1][image.imageSize-1].id
	fmt.Printf("Product of the corner tiles IDs: %d\n", cornerTileIdProduct)

	// Part II.
	imageTile := image.getImageTile()
	monsterPositions := imageTile.findMonsters(monsterPattern)
	counterMonster := 0
	for i := range monsterPattern {
		counterMonster += strings.Count(monsterPattern[i], "#")
	}
	counterTotal := 0
	for i := range imageTile.content {
		counterTotal += strings.Count(string(imageTile.content[i]), "#")
	}
	fmt.Printf("Final water roughness: %d (%d monsters found)\n", counterTotal-len(monsterPositions)*counterMonster, len(monsterPositions))
}

// --------------------------------------- IMAGE

type Image struct {
	tilePlacement [][]*Tile
	imageSize     int
}

func NewImage(tiles Tiles) *Image {
	imageSize := int(math.Sqrt(float64(len(tiles))))
	tilePlacement := make([][]*Tile, imageSize)
	for i := range tilePlacement {
		tilePlacement[i] = make([]*Tile, imageSize)
	}

	// place tiles into the image
	usedIds := make(map[int]struct{})
	for i := 0; i < imageSize; i++ {
		for j := 0; j < imageSize; j++ {
			var tile *Tile
			if i == 0 && j == 0 {
				tile = tiles.findTopLeftCorner()
			} else {
				matchBorder := ""
				borderIndex := -1
				if j == 0 {
					matchBorder = tilePlacement[i-1][j].borders[bottomIndex]
					borderIndex = topIndex
				} else {
					matchBorder = tilePlacement[i][j-1].borders[rightIndex]
					borderIndex = leftIndex
				}

				tile = tiles.findNext(matchBorder, borderIndex, usedIds)
				if i > 0 {
					if tile.borders[topIndex] != tilePlacement[i-1][j].borders[bottomIndex] {
						panic(fmt.Errorf("tile id [%d] does not fit in the %d-th row and %d-th column (left border matches, but the top one does not)", tile.id, i, j))
					}
				}
			}

			tilePlacement[i][j] = tile
			usedIds[tile.id] = struct{}{}
		}
	}

	return &Image{
		tilePlacement: tilePlacement,
		imageSize:     imageSize,
	}
}

func (im *Image) getImageTile() *Tile {
	tileSize := im.tilePlacement[0][0].size
	imageTileSize := len(im.tilePlacement) * (tileSize - 2)
	content := make([][]byte, imageTileSize)

	for i := range im.tilePlacement {
		for j := range im.tilePlacement[i] {
			currentTile := im.tilePlacement[i][j]
			for k := 1; k < tileSize-1; k++ {
				// for each non-border tile row
				imageRow := i*(tileSize-2) + k - 1
				if content[imageRow] == nil {
					content[imageRow] = make([]byte, imageTileSize)
				}
				for l := 1; l < tileSize-1; l++ {
					imageColumn := j*(tileSize-2) + l - 1
					content[imageRow][imageColumn] = currentTile.content[k][l]
				}
			}
		}
	}

	imageTile := Tile{
		id:      0,
		borders: [4]string{},
		content: content,
		size:    imageTileSize,
	}
	imageTile.updateBorders()

	return &imageTile
}

// --------------------------------------- TILE ARRAY

type Tiles []*Tile

func getTiles(input string) Tiles {
	lines := ReadLines(input)
	tiles := make(Tiles, 0)
	currentTileStart := 0
	currentTileId := 0
	for i, line := range lines {
		if line == "" {
			tile := NewTile(currentTileId, lines[currentTileStart+1:i])
			tiles = append(tiles, &tile)
		} else if line[0:5] == "Tile " {
			id, err := strconv.Atoi(line[5 : len(line)-1])
			if err != nil {
				panic(fmt.Errorf("invalid tile header [%s]", line))
			}
			currentTileStart = i
			currentTileId = id
		}
	}

	return tiles
}

func (ta Tiles) findTopLeftCorner() *Tile {
	for i := 0; i < len(ta); i++ {
		freeBorders := make([]int, 0, 2)
		for j := range ta[i].borders {
			currentBorder := ta[i].borders[j]
			reversed := Reverse(currentBorder)
			matchFound := false
			for k := 0; k < len(ta) && !matchFound; k++ {
				if k == i {
					continue // do not match with itself
				}
				for l := range ta[k].borders {
					if currentBorder == ta[k].borders[l] || reversed == ta[k].borders[l] {
						matchFound = true
						break
					}
				}
			}
			if !matchFound {
				freeBorders = append(freeBorders, j)
			}
		}

		if len(freeBorders) == 2 {
			//This is a corner tile - rotate it the way it matches the TL corner (corners 3 and 0)
			if AbsInt(freeBorders[0]-freeBorders[1]) != 3 {
				//not rotated properly
				min := int(math.Min(float64(freeBorders[0]), float64(freeBorders[1])))
				rotateSteps := (3 - min + 4) % 4
				ta[i].rotateClockwise(rotateSteps)
			}

			return ta[i] // corner tile found and positioned correctly to TL orientation
		}
	}

	panic("top left corner tile not found")
}

func (ta Tiles) findNext(borderToMatch string, wantBorderIndex int, exclude map[int]struct{}) *Tile {
	for _, tile := range ta {
		if _, ok := exclude[tile.id]; ok {
			continue // this tile is already used
		}
		for haveBorderIndex, border := range tile.borders {
			if border == borderToMatch || Reverse(border) == borderToMatch {
				tile.rotateClockwise((wantBorderIndex - haveBorderIndex + 4) % 4)
				if tile.borders[wantBorderIndex] != borderToMatch {
					if wantBorderIndex%2 == 0 {
						tile.flipHorizontal() // the matched (reversed) border is top or bottom -> flip horizontally
					} else {
						tile.flipVertical()
					}
				}
				return tile
			}
		}
	}

	panic(fmt.Errorf("no suitable tile found matching border [%s] on side [%d]", borderToMatch, wantBorderIndex))
}

// --------------------------------------- TILE

type Tile struct {
	id      int
	borders [4]string // top, right, bottom, left
	content [][]byte
	size    int
}

func NewTile(id int, lines []string) Tile {
	tile := Tile{
		id:      id,
		size:    len(lines),
		borders: [4]string{},
		content: make([][]byte, len(lines)),
	}
	for i := 0; i < tile.size; i++ {
		tile.content[i] = []byte(lines[i])
	}
	tile.updateBorders()

	return tile
}

func (t *Tile) updateBorders() {
	t.borders = [4]string{
		string(t.content[0]),
		"",
		string(t.content[t.size-1]),
		"",
	}
	for _, row := range t.content {
		t.borders[1] += string(row[t.size-1]) // right border
		t.borders[3] += string(row[0])        // left border
	}
}

func (t *Tile) flipVertical() {
	for i, j := 0, t.size-1; i < j; i, j = i+1, j-1 {
		t.content[i], t.content[j] = t.content[j], t.content[i]
	}
	t.updateBorders()
}

func (t *Tile) flipHorizontal() {
	for row := range t.content {
		for i, j := 0, t.size-1; i < j; i, j = i+1, j-1 {
			t.content[row][i], t.content[row][j] = t.content[row][j], t.content[row][i]
		}
	}
	t.updateBorders()
}

func (t *Tile) rotateClockwise(steps int) {
	steps = steps % 4 // trim

	for step := 0; step < steps; step++ {
		// rotate 90 degrees clockwise
		for i := 0; i < t.size/2; i++ {
			for j := i; j < t.size-i-1; j++ {
				temp := t.content[i][j]
				t.content[i][j] = t.content[t.size-1-j][i]
				t.content[t.size-1-j][i] = t.content[t.size-1-i][t.size-1-j]
				t.content[t.size-1-i][t.size-1-j] = t.content[j][t.size-1-i]
				t.content[j][t.size-1-i] = temp
			}
		}
	}
	t.updateBorders()
}

func (t *Tile) findMonsters(monsterPattern []string) [][2]int {
	// try rotating and flipping the tile until at least one monster found
	monsterFinderFunc := func() [][2]int {
		for i := 0; i < 4; i++ {
			t.rotateClockwise(i)
			monsterPositions := t.findPattern(monsterPattern)
			if len(monsterPositions) > 0 {
				return monsterPositions
			}
		}
		return nil // no monsters found
	}

	monsterPositions := monsterFinderFunc()
	if len(monsterPositions) == 0 {
		t.flipVertical()
		monsterPositions = monsterFinderFunc()
	}

	return monsterPositions
}

func (t *Tile) findPattern(pattern []string) [][2]int {
	regexpPattern := make([]*regexp.Regexp, len(pattern))
	for i := range pattern {
		regexpPattern[i] = regexp.MustCompile(pattern[i])
	}

	matchCorners := make([][2]int, 0) // coordinates of TL corner of the pattern in the image
	for i := range t.content {
		lastMatch := -1
		for true {
			if lastMatch+1+len(pattern[0]) > len(t.content[i])-1 {
				break // match no longer possible (pattern longer than the rest of the line)
			}
			firstLineMatch := regexpPattern[0].FindStringIndex(string(t.content[i][lastMatch+1:]))
			if firstLineMatch == nil {
				break // no match on this line anymore
			}
			matchStart := firstLineMatch[0] + lastMatch + 1 // we were matching a sub-slice of the row, must recalculate match bounds to absolute values
			matchEnd := firstLineMatch[1] + lastMatch + 1
			// try matching the rest of the pattern lines
			for j := 1; j < len(pattern) && i+j < len(t.content); j++ {
				lineMatch := regexpPattern[j].FindStringIndex(string(t.content[i+j][matchStart:matchEnd]))
				if lineMatch == nil {
					break // no match here - try another matching the first line again
				} else {
					if j == len(pattern)-1 {
						matchCorners = append(matchCorners, [2]int{i, matchStart}) // full match, add the TL corner to the result
					}
				}
			}
			lastMatch = matchStart
		}

	}

	return matchCorners
}

func (t *Tile) print() {
	for i := range t.content {
		println(string(t.content[i]))
	}
	println()
}

func Reverse(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, r := range s {
		n--
		runes[n] = r
	}
	return string(runes[n:])
}
