package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	rules, messages := readSatelliteMessages()
	v1Counter, v2Counter := 0, 0

	// The way the input data are specified, any matching message must be a 8 11. That expanded means 42 42 31 (this is without recursion in Part II)
	// With recursion it grows to 42 42 42 31, 42 42 42 31 31 and so on so it's always a sequence of 42 "sub-matches" followed by a (smaller) sequence of 31 "sub-matches"
	// Moreover, rules 42 and 31 both have the same amount of possible combinations (128) and length of the possible combinations (8) so we can match the message by bytes
	combinations42, combinations31 := rules.getCombinationsForRule(42), rules.getCombinationsForRule(31)
	sort.Strings(combinations42) // faster search in messageMatchesRule()
	sort.Strings(combinations31)
	for _, message := range messages {
		if len(message)%8 > 0 {
			continue // no need to check messages with length not divisible by 8
		}

		matchPattern := ""
		for byteNumber := 0; byteNumber < len(message)/8; byteNumber++ {
			// Encode the "shape" of byte matches in letters to create an easy-to-check pattern for this message
			messageByte := message[0+8*byteNumber : 8+8*byteNumber]
			if messageMatchesRule(messageByte, combinations42) {
				matchPattern += "A"
			} else if messageMatchesRule(messageByte, combinations31) {
				matchPattern += "B"
			} else {
				matchPattern = "X" // not matching at all
				break
			}
		}

		if regexp.MustCompile("^AA+B+$").MatchString(matchPattern) {
			//AAB, AAAB, AAABB allowed; AABA, BBA not...
			if strings.Count(matchPattern, "A") > strings.Count(matchPattern, "B") {
				// AAB, AAAB, AAABB allowed; AB AABB, ABBB not...
				v2Counter++
				if matchPattern == "AAB" {
					v1Counter++ // "basic" version without loops
				}
			}
		}
	}
	fmt.Printf("Number of satellite message matching the rule set (without loops): %d\n", v1Counter)                        // Part I.
	fmt.Printf("Number of satellite message matching the rule set (with rules 8 and 11 containing loops): %d\n", v2Counter) // Part II.
}

func readSatelliteMessages() (Rules, []string) {
	rules := make(Rules)
	messages := make([]string, 0)

	isMessage := false
	for _, line := range ReadLines("inputs/day_19.txt") {
		if line == "" {
			isMessage = true
			continue
		}

		if isMessage {
			messages = append(messages, line)
		} else {
			rule := NewRule(line)
			rules[rule.index] = rule
		}
	}

	return rules, messages
}

func messageMatchesRule(message string, ruleCombinations []string) bool {
	i := sort.SearchStrings(ruleCombinations, message)
	return i < len(ruleCombinations) && ruleCombinations[i] == message
}

type Rule struct {
	index       int
	char        string
	ruleOptions []RuleSequence
}

type RuleSequence []int

func NewRule(line string) Rule {
	parts := strings.Split(line, ": ")
	index, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(fmt.Errorf("invalid format for rule [%s]: invalid rule number", line))
	}

	pattern := regexp.MustCompile("^\"([a-z])\"$")
	matches := pattern.FindStringSubmatch(parts[1])
	if len(matches) == 2 {
		return Rule{
			index:       index,
			char:        matches[1],
			ruleOptions: nil,
		}
	}

	ruleOptionsStrings := strings.Split(parts[1], " | ")
	ruleOptions := make([]RuleSequence, len(ruleOptionsStrings))
	for i, ruleSequenceString := range ruleOptionsStrings {
		ruleIndexes := strings.Split(ruleSequenceString, " ")
		ruleOptions[i] = make([]int, len(ruleIndexes))
		for j, ruleIndex := range ruleIndexes {
			ruleOptions[i][j], err = strconv.Atoi(ruleIndex)
			if err != nil {
				panic(fmt.Errorf("invalid format for rule [%s]: invalid rune index", line))
			}
		}

	}

	return Rule{
		index:       index,
		char:        "",
		ruleOptions: ruleOptions,
	}
}

type Rules map[int]Rule

func (r *Rules) getCombinationsForRule(ruleIndex int) []string {
	rule, ok := (*r)[ruleIndex]
	if !ok {
		panic(fmt.Errorf("rule index [%d] not in rule set", ruleIndex))
	}

	if rule.char != "" {
		return []string{rule.char} // terminal rule, just return
	}

	combinations := make([]string, 0)
	for _, ruleOption := range rule.ruleOptions {
		// ruleOption = 3 1 4
		optionCombinations := []string{""}
		for _, subRuleIndex := range ruleOption {
			// subRuleIndex = 3
			subCombinations := r.getCombinationsForRule(subRuleIndex)
			newOptionCombinations := make([]string, 0)
			for _, newComb := range subCombinations {
				for _, currentComb := range optionCombinations {
					newOptionCombinations = append(newOptionCombinations, currentComb+newComb)
				}
			}
			optionCombinations = newOptionCombinations
		}
		combinations = append(combinations, optionCombinations...)
	}

	return combinations
}
