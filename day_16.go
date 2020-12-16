package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	fields, myTicket, tickets := readTicketInput()

	// Part I.
	scanningErrorRate := 0
	validTickets := make([]Ticket, 0)
	for _, ticket := range tickets {
		invalidValues := fields.getInvalidValues(ticket)
		for _, value := range invalidValues {
			scanningErrorRate += value
		}
		if len(invalidValues) == 0 {
			validTickets = append(validTickets, ticket)
		}
	}
	fmt.Printf("Ticket scanning error rate: %d\n", scanningErrorRate)

	// Part II.
	positionMap := fields.getFieldPositions(validTickets)
	result := 1
	for _, field := range fields {
		if field.isDepartureField() {
			result *= myTicket[positionMap[field.name]]
		}
	}
	fmt.Printf("My departure fields product: %d\n", result)
}

func readTicketInput() (TicketFields, Ticket, []Ticket) {
	lines := ReadLines("inputs/day_16.txt")

	i := 0
	validations := make(map[string]TicketField)
	for lines[i] != "" {
		field := NewTicketField(lines[i])
		validations[field.name] = field
		i++
	}
	myTicket := NewTicket(lines[i+2])
	tickets := make([]Ticket, 0)
	for i = i + 5; i < len(lines); i++ {
		tickets = append(tickets, NewTicket(lines[i]))
	}

	return validations, myTicket, tickets
}

// --------------------------- TICKET

type Ticket []int

func NewTicket(inputLine string) Ticket {
	return StringsToInts(strings.Split(inputLine, ","))
}

// --------------------------- TICKET ARRAY

type Tickets []Ticket

func (ts Tickets) isPossibleFieldPosition(position int, field TicketField) bool {
	for _, ticket := range ts {
		if !field.validateValue(ticket[position]) {
			return false
		}
	}
	return true
}

// --------------------------- TICKET FIELD

type TicketField struct {
	name      string
	intervals []Interval
}

func NewTicketField(inputLine string) TicketField {
	pattern := regexp.MustCompile("^([a-z ]+): (.*)$")
	matches := pattern.FindStringSubmatch(inputLine)
	if len(matches) != 3 {
		panic(fmt.Errorf("invalid validation input line [%s]", inputLine))
	}

	intervals := make([]Interval, 0)
	for _, intervalString := range strings.Split(matches[2], " or ") {
		intervals = append(intervals, NewInterval(intervalString))
	}

	return TicketField{
		name:      matches[1],
		intervals: intervals,
	}
}

func (f TicketField) validateValue(value int) bool {
	for _, interval := range f.intervals {
		if isInInterval(interval[0], interval[1], value) {
			return true
		}
	}

	return false
}

func (f TicketField) isDepartureField() bool {
	departure := "departure"
	if len(f.name) < len(departure) {
		return false
	}
	return f.name[0:len(departure)] == departure
}

// --------------------------- TICKET FIELD ARRAY

type TicketFields map[string]TicketField

func (fa TicketFields) getInvalidValues(ticket Ticket) []int {
	result := make([]int, 0)
	for _, value := range ticket {
		if !fa.valueIsValid(value) {
			result = append(result, value)
		}
	}

	return result
}

func (fa TicketFields) valueIsValid(value int) bool {
	for _, field := range fa {
		if field.validateValue(value) {
			return true // valid at leas one validation
		}
	}

	return false
}

func (fa TicketFields) getFieldPositions(tickets Tickets) map[string]int {
	remainingFields := fa.getCopy()
	remainingPositions := getPositionMap(len(fa))
	fieldPositions := make(map[string]int)
	iteration := 1
	for len(remainingFields) > 0 {

		var positionFound = false
		for _, field := range remainingFields {
			possiblePositionCounter := 0
			lastPossiblePosition := -1
			for position := range remainingPositions {
				if tickets.isPossibleFieldPosition(position, field) {
					lastPossiblePosition = position
					possiblePositionCounter++
				}
			}
			if possiblePositionCounter == 1 {
				delete(remainingFields, field.name)
				delete(remainingPositions, lastPossiblePosition)
				fieldPositions[field.name] = lastPossiblePosition
				positionFound = true
			}
		}

		if !positionFound {
			panic(fmt.Errorf("no unique position found for any of the remaining fields in iteration [%d]", iteration))
		}
		iteration++
	}

	return fieldPositions
}

func (fa TicketFields) getCopy() TicketFields {
	result := make(map[string]TicketField)
	for k, v := range fa {
		result[k] = v
	}

	return result
}

func getPositionMap(size int) map[int]struct{} {
	result := make(map[int]struct{})
	for i := 0; i < size; i++ {
		result[i] = struct{}{}
	}

	return result
}

// --------------------------- INTERVAL

type Interval [2]int

func NewInterval(intervalString string) Interval {
	parts := strings.Split(intervalString, "-")
	a, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(fmt.Errorf("invalid interval start [%s]", intervalString))
	}
	b, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(fmt.Errorf("invalid interval end [%s]", intervalString))
	}

	return Interval{a, b}
}
