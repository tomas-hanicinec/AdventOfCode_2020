package main

import (
	"fmt"
	"strings"
)

func main() {
	groupAnswers := readGroupAnswers()

	totalResultAny := 0
	totalResultAll := 0
	for _, groupAnswer := range groupAnswers {
		totalResultAny += groupAnswer.getAnswerCountAny()
		totalResultAll += groupAnswer.getAnswerCountAll()
	}

	fmt.Printf("Number of questions answered YES by anyone in group (sum across all groups): %d\n", totalResultAny)   // Part I.
	fmt.Printf("Number of questions answered YES by everyone in group (sum across all groups): %d\n", totalResultAll) // Part II.
}

type GroupAnswer struct {
	id             int
	size           int
	answerCountMap map[uint8]int
}

func readGroupAnswers() []*GroupAnswer {
	lines := ReadLines("inputs/day_06.txt")

	previousNewline := -1
	groupCounter := 0
	result := make([]*GroupAnswer, 0)
	for i, line := range lines {
		if line == "" {
			groupCounter++
			result = append(result, NewGroupAnswer(groupCounter, lines[previousNewline+1:i]))
			previousNewline = i
		}
	}

	return result
}

func NewGroupAnswer(groupId int, groupLines []string) *GroupAnswer {
	allAnswers := strings.Join(groupLines, "")
	answerCountMap := make(map[uint8]int, len(allAnswers))
	for i := range allAnswers {
		answerCountMap[allAnswers[i]]++
	}

	return &GroupAnswer{
		id:             groupId,
		size:           len(groupLines),
		answerCountMap: answerCountMap,
	}
}

func (ga GroupAnswer) getAnswerCountAny() int {
	result := 0
	for _, count := range ga.answerCountMap {
		if count > 0 {
			result++ // Al least one person from group answered yes to this question
		}
	}

	return result
}

func (ga GroupAnswer) getAnswerCountAll() int {
	result := 0
	for _, count := range ga.answerCountMap {
		if count == ga.size {
			result++ // All people from group answered yes to this question
		}
	}

	return result
}
