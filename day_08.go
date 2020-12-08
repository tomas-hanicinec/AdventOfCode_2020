package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	bootCode := NewBootCode()

	finished, result, _ := bootCode.run(0, -1)
	if finished {
		panic(fmt.Errorf("given boot code was not supposed to finish without correction"))
	}
	fmt.Printf("Accumulator value before infinite loop: %d\n", result) // Part I.

	for instructionIndexToRepair := range bootCode.instructions {
		if !bootCode.canRepairInstruction(instructionIndexToRepair) {
			continue // This instruction cannot be repaired, no need to run the whole boot code
		}
		finished, result, _ := bootCode.run(0, instructionIndexToRepair)
		if finished {
			fmt.Printf("Repaired instruction #%d, boot code finished with accumulator value: %d\n", instructionIndexToRepair, result) // Part II.
			return
		}
	}

	panic(fmt.Errorf("no fix to the given boot code found (always ends in infinite loop)"))
}

type BootCode struct {
	instructions []*Instruction
}

func NewBootCode() BootCode {
	lines := ReadLines("inputs/day_08.txt")

	bc := BootCode{
		instructions: make([]*Instruction, len(lines)),
	}
	for i, line := range lines {
		parsed := strings.Split(line, " ")
		argument, err := strconv.Atoi(parsed[1])
		if err != nil {
			panic(fmt.Errorf("invalid argument [%s] for operation [%s]", parsed[1], parsed[0]))
		}
		bc.instructions[i] = &Instruction{
			operationCode: parsed[0],
			argument:      argument,
		}
	}

	return bc
}

func (bc BootCode) run(instructionIndex int, instructionIndexToRepair int) (bool, int, []string) {
	visitedInstructions := make([]bool, len(bc.instructions))
	accumulator := 0
	currentInstructionIndex := instructionIndex
	runLog := make([]string, 0)
	for currentInstructionIndex < len(bc.instructions) {
		if visitedInstructions[currentInstructionIndex] == true {
			return false, accumulator, runLog // We already processed this instruction -> infinite loop
		}
		repairCurrentInstruction := instructionIndexToRepair == currentInstructionIndex
		instructionOffset, accOffset := bc.instructions[currentInstructionIndex].process(repairCurrentInstruction)
		runLog = append(runLog, fmt.Sprintf(
			"instruction #%d [%s], acc [%d->%d], index[%d->%d]",
			currentInstructionIndex,
			bc.instructions[currentInstructionIndex].operationCode,
			accumulator, accumulator+accOffset,
			currentInstructionIndex, currentInstructionIndex+instructionOffset,
		))
		visitedInstructions[currentInstructionIndex] = true
		accumulator += accOffset
		currentInstructionIndex = bc.getNextInstructionIndex(currentInstructionIndex, instructionOffset)
	}

	return true, accumulator, runLog // currentInstructionIndex at the end of instruction file -> successful finish
}

func (bc BootCode) getNextInstructionIndex(currentInstructionIndex int, instructionOffset int) int {
	newIndex := currentInstructionIndex + instructionOffset
	if newIndex < 0 || newIndex > len(bc.instructions) {
		panic(fmt.Errorf("instruction index overflow (from [%d] jump by [%d] to [%d])", currentInstructionIndex, instructionOffset, newIndex))
	}

	return newIndex
}

func (bc BootCode) canRepairInstruction(instructionIndexToRepair int) bool {
	return bc.instructions[instructionIndexToRepair].getOperationDefinition().repairedCode != ""
}

type Instruction struct {
	operationCode string
	argument      int
}

func (i Instruction) process(repairInstruction bool) (int, int) {
	instructionToProcess := i
	if repairInstruction {
		instructionToProcess = i.getRepaired()
	}
	return instructionToProcess.getOperationDefinition().processor(instructionToProcess.argument)
}

func (i Instruction) getRepaired() Instruction {
	definition := i.getOperationDefinition()
	if !definition.canRepair() {
		panic(fmt.Errorf("instruction [%s %d] cannot be repaired", i.operationCode, i.argument))
	}
	return Instruction{
		operationCode: definition.repairedCode,
		argument:      i.argument,
	}
}

func (i Instruction) getOperationDefinition() InstructionOperationDefinition {
	definitions := getInstructionOperationDefinitions()
	definition, exists := definitions[i.operationCode]
	if !exists {
		panic(fmt.Errorf("unknown instruction [%s]", i.operationCode))
	}

	return definition
}

type InstructionOperationDefinition struct {
	processor    InstructionOperationProcessor
	repairedCode string
}

type InstructionOperationProcessor func(argument int) (int, int)

func getInstructionOperationDefinitions() map[string]InstructionOperationDefinition {
	return map[string]InstructionOperationDefinition{
		"nop": {
			processor: func(argument int) (int, int) {
				return 1, 0
			},
			repairedCode: "jmp",
		},
		"acc": {
			processor: func(argument int) (int, int) {
				return 1, argument
			},
			repairedCode: "", // Not possible to repair
		},
		"jmp": {
			processor: func(argument int) (int, int) {
				return argument, 0
			},
			repairedCode: "nop",
		},
	}
}

func (iod InstructionOperationDefinition) canRepair() bool {
	return iod.repairedCode != ""
}
