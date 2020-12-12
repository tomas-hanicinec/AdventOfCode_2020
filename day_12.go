package main

import (
	"fmt"
	"strconv"
)

const North = 0
const East = 90
const South = 180
const West = 270

func main() {
	instructionSet := NewNavInstructionSet()

	// Part I.
	ship := NewDefaultShip()
	ship.move(instructionSet, 1)
	fmt.Printf("V1 final ship position [N: %d, W: %d], manhattan distance: %d\n", ship.position.north, ship.position.west, ship.position.getManhattanDistance())

	// Part II.
	ship = NewDefaultShip()
	ship.move(instructionSet, 2)
	fmt.Printf("V2 final ship position [N: %d, W: %d], manhattan distance: %d\n", ship.position.north, ship.position.west, ship.position.getManhattanDistance())
}

type Position struct {
	north int
	west  int
}

func (sp Position) move(direction int, distance int) Position {
	switch direction {
	case North:
		return Position{sp.north + distance, sp.west}
	case East:
		return Position{sp.north, sp.west - distance}
	case South:
		return Position{sp.north - distance, sp.west}
	case West:
		return Position{sp.north, sp.west + distance}
	default:
		panic(fmt.Errorf("unsupported move direction value [%d]", direction))
	}
}

func (sp Position) moveToWaypoint(waypoint Position, multiplicator int) Position {
	return Position{
		north: sp.north + multiplicator*waypoint.north,
		west:  sp.west + multiplicator*waypoint.west,
	}
}

func (sp Position) rotateClockwise(degrees int) Position {
	switch degrees {
	case 0:
		return Position{sp.north, sp.west}
	case 90:
		return Position{sp.west, -sp.north}
	case 180:
		return Position{-sp.north, -sp.west}
	case 270:
		return Position{-sp.west, sp.north}
	default:
		panic(fmt.Errorf("unsupported rotation degrees value [%d]", degrees))
	}
}

func (sp Position) getManhattanDistance() int {
	return AbsInt(sp.north) + AbsInt(sp.west)
}

type Ship struct {
	direction int
	position  Position
	waypoint  Position
}

func NewDefaultShip() Ship {
	return Ship{
		direction: East,
		position:  Position{0, 0},
		waypoint:  Position{1, -10},
	}
}

func (s *Ship) move(is NavInstructionSet, version int) {
	for _, instruction := range is.instructions {
		s.position, s.waypoint, s.direction = is.transform[version-1][instruction.action](s.position, s.waypoint, s.direction, instruction.value)
	}
}

type NavInstructionSet struct {
	instructions []*NavInstruction
	transform    []map[string]NavInstructionTransformer
}

func NewNavInstructionSet() NavInstructionSet {
	lines := ReadLines("inputs/day_12.txt")
	instructions := make([]*NavInstruction, len(lines))
	for i, line := range lines {
		instructions[i] = NewNavInstruction(line)
	}

	return NavInstructionSet{
		instructions: instructions,
		transform: []map[string]NavInstructionTransformer{
			getInstructionTransformersV1(),
			getInstructionTransformersV2(),
		},
	}
}

type NavInstruction struct {
	action string
	value  int
}

func NewNavInstruction(line string) *NavInstruction {
	value, err := strconv.Atoi(line[1:])
	if err != nil {
		panic(fmt.Errorf("invalid instruction [%s], argument not numeric", line))
	}
	if line[0] == 'L' || line[0] == 'R' {
		if value%90 != 0 {
			// this would mean trigonometry hell, let's say it's not allowed
			panic(fmt.Errorf("invalid instruction [%s], cannot turn uneven angles", line))
		}
	}
	return &NavInstruction{
		action: line[0:1],
		value:  value % 360, // normalize rotation to 0-360
	}
}

type NavInstructionTransformer func(Position, Position, int, int) (Position, Position, int)

func getInstructionTransformersV1() map[string]NavInstructionTransformer {
	return map[string]NavInstructionTransformer{
		"N": func(p Position, w Position, d int, instructionValue int) (Position, Position, int) {
			return p.move(North, instructionValue), w, d
		},
		"S": func(p Position, w Position, d int, instructionValue int) (Position, Position, int) {
			return p.move(South, instructionValue), w, d
		},
		"E": func(p Position, w Position, d int, instructionValue int) (Position, Position, int) {
			return p.move(East, instructionValue), w, d
		},
		"W": func(p Position, w Position, d int, instructionValue int) (Position, Position, int) {
			return p.move(West, instructionValue), w, d
		},
		"L": func(p Position, w Position, d int, instructionValue int) (Position, Position, int) {
			return p, w, (d + 360 - instructionValue) % 360
		},
		"R": func(p Position, w Position, d int, instructionValue int) (Position, Position, int) {
			return p, w, (d + instructionValue) % 360
		},
		"F": func(p Position, w Position, d int, instructionValue int) (Position, Position, int) {
			return p.move(d, instructionValue), w, d
		},
	}
}

func getInstructionTransformersV2() map[string]NavInstructionTransformer {
	return map[string]NavInstructionTransformer{
		"N": func(p Position, w Position, d int, instructionValue int) (Position, Position, int) {
			return p, w.move(North, instructionValue), d
		},
		"S": func(p Position, w Position, d int, instructionValue int) (Position, Position, int) {
			return p, w.move(South, instructionValue), d
		},
		"E": func(p Position, w Position, d int, instructionValue int) (Position, Position, int) {
			return p, w.move(East, instructionValue), d
		},
		"W": func(p Position, w Position, d int, instructionValue int) (Position, Position, int) {
			return p, w.move(West, instructionValue), d
		},
		"L": func(p Position, w Position, d int, instructionValue int) (Position, Position, int) {
			return p, w.rotateClockwise(360 - instructionValue), d
		},
		"R": func(p Position, w Position, d int, instructionValue int) (Position, Position, int) {
			return p, w.rotateClockwise(instructionValue), d
		},
		"F": func(p Position, w Position, d int, instructionValue int) (Position, Position, int) {
			return p.moveToWaypoint(w, instructionValue), w, d
		},
	}
}
