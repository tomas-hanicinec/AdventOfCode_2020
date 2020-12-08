package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const myBagColor = "shiny gold"

func main() {
	bagRules := readBagRules()

	// Part I.
	result1 := make(map[BagColor]bool)
	traverseGraphUpwards(myBagColor, bagRules, result1)
	fmt.Printf("Number of bags that can contain my bag: %d\n", len(result1)-1) // Exclude the root (myBagColor) from the result

	// Part II.
	result2 := traverseGraphDownwards(myBagColor, bagRules, 1)
	fmt.Printf("Number of bags my bag has to contain: %d\n", result2-1) // Exclude the root (myBagColor) from the result
}

type BagColor string

type BagRule struct {
	children map[BagColor]int
	parents  map[BagColor]int
}

func (br *BagRule) addParent(color BagColor, count int) {
	br.parents[color] = count
}

func readBagRules() map[BagColor]*BagRule {
	lines := ReadLines("inputs/day_07.txt")
	bagRules := make(map[BagColor]*BagRule, len(lines))

	// Fill "children" links (parent bag can contain child bag)
	for _, line := range lines {
		color, children := parseLine(line)
		bagRules[color] = &BagRule{
			children: children,
			parents:  make(map[BagColor]int),
		}
	}

	// Fill "parent" links (child bag can be contained by parent bag)
	for color, bagRule := range bagRules {
		for childColor, count := range bagRule.children {
			bagRules[childColor].parents[color] = count
		}
	}

	return bagRules
}

func parseLine(line string) (BagColor, map[BagColor]int) {
	linePattern := regexp.MustCompile("^([a-z]+ [a-z]+) bags contain (.+)\\.$")
	bagPattern := regexp.MustCompile("^([0-9]+) ([a-z]+ [a-z]+) bags?$")
	noOtherBagsPattern := "no other bags"
	lineMatches := getMatches(line, linePattern)
	bagRuleItems := make(map[BagColor]int)
	if lineMatches[1] != noOtherBagsPattern {
		for _, containedBag := range strings.Split(lineMatches[1], ", ") {
			bagMatches := getMatches(containedBag, bagPattern)
			bagCount, err := strconv.Atoi(bagMatches[0])
			if err != nil {
				panic(fmt.Errorf("invalid bag count [%s]", bagMatches[0]))
			}
			bagRuleItems[BagColor(bagMatches[1])] = bagCount
		}
	}

	return BagColor(lineMatches[0]), bagRuleItems
}

func getMatches(input string, pattern *regexp.Regexp) []string {
	matches := pattern.FindStringSubmatch(input)
	if len(matches) < 1 {
		panic(fmt.Errorf("failed to match input [%s] with pattern[%s]", input, pattern))
	}

	return matches[1:]
}

func traverseGraphUpwards(color BagColor, graph map[BagColor]*BagRule, result map[BagColor]bool) {
	currentRule, exists := graph[color]
	if !exists {
		panic(fmt.Errorf("logic error: bag color [%s] not present in bag rules", color))
	}

	if _, exists = result[color]; exists {
		return // This color is already in the map, no need to traverse further (would lead to infinite loop anyway)
	}

	result[color] = true

	// Handle parents
	for parentColor := range currentRule.parents {
		traverseGraphUpwards(parentColor, graph, result)
	}
}

func traverseGraphDownwards(color BagColor, graph map[BagColor]*BagRule, currentDepth int) int {
	if currentDepth > len(graph) {
		panic(fmt.Errorf("infinite loop in input graph discovered"))
	}
	currentRule, exists := graph[color]
	if !exists {
		panic(fmt.Errorf("logic error: bag color [%s] not present in bag rules", color))
	}
	result := 1 // Include this bag as well

	// Handle children
	for childColor, count := range currentRule.children {
		val := traverseGraphDownwards(childColor, graph, currentDepth+1)
		result += count * val // Add total number of bags in "childColor" bag (including the childColorBag)
	}

	return result // Total number of bags in "color" bag (including itself)
}
