package parser

import (
	"github.com/tealang/tea-go/lexer/tokens"
	"github.com/tealang/tea-go/runtime/nodes"
)

type FunctionCallParser struct{}

func (FunctionCallParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	return nil, 0, nil
}

func NewTermParser() *TermParser {
	return &TermParser{
		Operands:  make([]nodes.Node, 0),
		Operators: make([]nodes.Node, 0),
	}
}

type TermParser struct {
	Operands, Operators []nodes.Node
}

func (TermParser) PriorityOf(symbol string, previous tokens.Token) (int, error) {
	switch symbol {
	case "&", "|":
		return 8, nil
	case "!":
		return 7, nil
	case "^":
		return 6, nil
	case "*", "/":
		return 5, nil
	case "+", "-", ":":
		return 4, nil
	case "%":
		return 3, nil
	case "<", ">", ">=", "<=", "=>", "=<", "!=", "==":
		return 2, nil
	case "&&", "||", "^|":
		return 1, nil
	}
	return 0, nil
}

func (tp *TermParser) PopOperator() (nodes.Node, error) {
	operatorStackSize := len(tp.Operators)
	if operatorStackSize < 1 {
		return nil, ParseException{"Operator stack is empty"}
	}
	operatorStackItem := tp.Operators[operatorStackSize-1]
	tp.Operators = tp.Operators[:operatorStackSize-1]
	operation, ok := operatorStackItem.(*nodes.Operation)
	if !ok {
		return nil, ParseException{"Operator stack item is no operation"}
	}
	operandStackSize := len(tp.Operands)
	for i := 0; i < operation.ArgCount; i++ {
		operand, err := tp.PopOperand()
		if err != nil {
			return nil, err
		}
		operation.AddFront(operand)
	}
	return operation, nil
}

func (tp *TermParser) PopOperand() (nodes.Node, error) {
	operandStackSize := len(tp.Operands)
	if operandStackSize < 1 {
		return nil, ParseException{"Operand stack is empty"}
	}
	operand := tp.Operands[operandStackSize-1]
	tp.Operands = tp.Operands[:operandStackSize-1]
	return operand, nil
}

func (tp *TermParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	var (
		index       int
		level       int
		activeToken tokens.Token
	)

	for index = 0; index < len(input); index++ {
		//lastToken = activeToken
		activeToken = input[index]
		/*

			switch activeToken.Type {
			case tokens.Separator:
				for len(operators) > 0 {
					err := popOperator()
					if err != nil {
						return nil, 0, err
					}
				}
				if level > 0 {
					return operands[0], index, nil
				}
			}*/
	}

	if len(tp.Operands) != 1 {
		return nil, 0, ParseException{"Operator stack should have size 1"}
	}

	return tp.Operands[0], index, nil
}
