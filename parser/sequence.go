package parser

import (
	"github.com/tealang/tea-go/lexer/tokens"
	"github.com/tealang/tea-go/runtime"
	"github.com/tealang/tea-go/runtime/nodes"
)

func newSequenceParser(substitute bool) *sequenceParser {
	sp := &sequenceParser{substitute: substitute}
	sp.handlers = map[*tokens.Type]func() error{
		tokens.LeftBlock:  sp.handleLeftBlock,
		tokens.Identifier: sp.handleIdentifier,
	}
	return sp
}

type sequenceParser struct {
	substitute  bool
	statement   bool
	index, size int
	sequence    *nodes.Sequence
	active      tokens.Token
	input       []tokens.Token
	handlers    map[*tokens.Type]func() error
}

func (sp *sequenceParser) handleLeftBlock() error {
	item, n, err := newSequenceParser(true).Parse(sp.inputSegment(1))
	if err != nil {
		return err
	}
	// ignore both closing and opening block
	sp.index += n + 2
	sp.sequence.AddBack(item)
	sp.statement = false
	return nil
}

func (sp *sequenceParser) inputSegment(offset int) []tokens.Token {
	return sp.input[sp.index+offset:]
}

func (sp *sequenceParser) checkForAssignment() bool {
	for i := sp.index; i < sp.size; i++ {
		sp.active = sp.input[i]
		switch sp.active.Type {
		case tokens.Identifier, tokens.Separator:
		case tokens.Operator:
			if sp.active.Value != "=" {
				return false
			}
			return true
		default:
			return false
		}
	}
	return false
}

func (sp *sequenceParser) handleIdentifier() error {
	switch sp.active.Value {
	case variableKeyword, constantKeyword:
		stmt, n, err := newDeclarationParser().Parse(sp.inputSegment(0))
		if err != nil {
			return err
		}
		sp.sequence.AddBack(stmt)
		sp.index += n
	case returnKeyword:
		stmt, n, err := newReturnParser().Parse(sp.inputSegment(0))
		if err != nil {
			return err
		}
		sp.sequence.AddBack(stmt)
		sp.index += n
	case breakKeyword:
		sp.sequence.AddBack(nodes.NewController(runtime.BehaviorBreak))
		sp.index++
	case continueKeyword:
		sp.sequence.AddBack(nodes.NewController(runtime.BehaviorContinue))
		sp.index++
	case ifKeyword:
		stmt, n, err := newBranchParser().Parse(sp.inputSegment(0))
		if err != nil {
			return err
		}
		sp.sequence.AddBack(stmt)
		sp.index += n
		sp.statement = false
	default:
		if sp.checkForAssignment() {
			stmt, n, err := newAssignmentParser().Parse(sp.inputSegment(0))
			if err != nil {
				return err
			}
			sp.sequence.AddBack(stmt)
			sp.index += n
		} else {
			return sp.handleTerm()
		}
	}
	return nil
}

func (sp *sequenceParser) handleTerm() error {
	term, n, err := newTermParser().Parse(sp.inputSegment(0))
	if err != nil {
		return err
	}
	if term != nil {
		sp.sequence.AddBack(term)
	}
	sp.index += n
	return nil
}

func (sp *sequenceParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	sp.index, sp.size = 0, len(input)
	sp.sequence = nodes.NewSequence(sp.substitute)
	sp.input = input
	for sp.index < sp.size {
		sp.statement = true
		sp.active = sp.input[sp.index]

		//fmt.Printf("[%d:%d] %s\n", sp.index, sp.size, sp.active)
		switch sp.active.Type {
		case nil, tokens.RightBlock:
			return sp.sequence, sp.index, nil
		default:
			handler, ok := sp.handlers[sp.active.Type]
			if !ok {
				handler = sp.handleTerm
			}
			if err := handler(); err != nil {
				return sp.sequence, sp.index, err
			}
		}
		if sp.index < sp.size && sp.statement {
			if sp.input[sp.index].Type != tokens.Statement {
				return sp.sequence, sp.index, ParseException{"Expected end statement"}
			}
			sp.index++
		}
	}
	return sp.sequence, sp.index, nil
}
