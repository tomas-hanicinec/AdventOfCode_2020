package main

import (
	"fmt"
	"strconv"
)

const operatorAdd = "+"
const operatorMultiply = "*"

func main() {
	lines := ReadLines("inputs/day_18.txt")
	sumV1, sumV2 := 0, 0
	for _, line := range lines {
		expression := NewExpression(line)
		sumV1 += expression.evaluate()
		expressionV2 := expression.modifyForPrecedence([]Operator{operatorAdd, operatorMultiply})
		sumV2 += expressionV2.evaluate()
	}
	fmt.Printf("Sum of the given expressions: %d\n", sumV1)                                // Part I.
	fmt.Printf("Sum of the given expressions with given operator precedence: %d\n", sumV2) // Part II.
}

type Operator string

type Expression struct {
	number      int
	expressions []Expression
	operators   []Operator
}

func NewExpression(line string) Expression {
	expression, _ := readExpression(line, 0)
	return expression
}

func newNumberExpression(number int) Expression {
	return Expression{
		number:      number,
		expressions: nil,
		operators:   nil,
	}
}

func newEmptyExpression() Expression {
	return Expression{
		number:      0,
		expressions: make([]Expression, 0),
		operators:   make([]Operator, 0),
	}
}

func readExpression(line string, index int) (Expression, int) {
	result := Expression{
		number:      0,
		expressions: make([]Expression, 0),
		operators:   make([]Operator, 0),
	}
	startIndex := index
	for index < len(line) {
		switch line[index] {
		case ' ':
			index++
		case '(':
			index++
			subExpression, length := readExpression(line, index)
			result.expressions = append(result.expressions, subExpression)
			index += length
		case ')':
			index++
			return result, index - startIndex
		case operatorAdd[0]:
			result.operators = append(result.operators, operatorAdd)
			index++
		case operatorMultiply[0]:
			result.operators = append(result.operators, operatorMultiply)
			index++
		default:
			number, err := strconv.Atoi(string(line[index]))
			if err != nil {
				panic(fmt.Errorf("invalid input [%s] - [%s] is not a number", line, string(line[index])))
			}
			result.expressions = append(result.expressions, newNumberExpression(number))
			index++
		}
	}

	return result, index - startIndex // end of line reached
}

func (e Expression) modifyForPrecedence(precedence []Operator) Expression {
	if e.isNumberExpression() {
		return e // no modification required
	}

	result := e.copy() // start with the main expression
	for opi := 0; opi < len(precedence)-1; opi++ {
		precedenceOperator := precedence[opi] // take operators one by one in the order of their precedence

		newResult := newEmptyExpression()
		currentExpression := result.expressions[0].modifyForPrecedence(precedence) // start with the first expression of the result
		for i := 0; i < len(result.operators); {
			if e.operators[i] == precedenceOperator {
				// replace the current operation (2+ expressions, 1+ precedenceOperator operators) with a single sub-expression (will be evaluated first)
				subExpression := newEmptyExpression()
				for ; i < len(e.operators) && e.operators[i] == precedenceOperator; i++ {
					// keep adding operands until the operator is precedenceOperator
					subExpression.expressions = append(subExpression.expressions, result.expressions[i].modifyForPrecedence(precedence))
					subExpression.operators = append(subExpression.operators, result.operators[i])
				}
				subExpression.expressions = append(subExpression.expressions, result.expressions[i].modifyForPrecedence(precedence)) // add the last operand too
				currentExpression = subExpression                                                                                    // save for next iteration
			} else {
				// keep the current operation in the result (but use the currentExpression which could have been modified in the previous iteration)
				newResult.expressions = append(newResult.expressions, currentExpression)
				newResult.operators = append(newResult.operators, result.operators[i])
				i++
				currentExpression = result.expressions[i].modifyForPrecedence(precedence)
			}
		}
		newResult.expressions = append(newResult.expressions, currentExpression) // add the last expression (loop ends with the last operator)
		result = newResult
	}

	return result
}

func (e Expression) isNumberExpression() bool {
	return e.expressions == nil && e.operators == nil
}

func (e Expression) evaluate() int {
	if e.isNumberExpression() {
		return e.number
	}

	// compound expression
	if len(e.expressions) != len(e.operators)+1 {
		panic(fmt.Errorf("invalid expression, inconsistent number of operators and operands"))
	}

	result := e.expressions[0].evaluate()
	for i := 1; i < len(e.expressions); i++ {
		switch e.operators[i-1] {
		case operatorAdd:
			result += e.expressions[i].evaluate()
		case operatorMultiply:
			result *= e.expressions[i].evaluate()
		default:
			panic(fmt.Errorf("unknown operator [%s]", e.operators[i-1]))
		}
	}

	return result
}

func (e Expression) copy() Expression {
	newExpressions := make([]Expression, len(e.expressions))
	copy(newExpressions, e.expressions)
	newOperators := make([]Operator, len(e.operators))
	copy(newOperators, e.operators)
	return Expression{
		number:      e.number,
		expressions: newExpressions,
		operators:   newOperators,
	}
}
