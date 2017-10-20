package nodes

import (
	"fmt"

	"github.com/tealang/tea-go/runtime"
)

// Assignment assigns one or multiple existing variables a new value (for each variable).
type Assignment struct {
	BasicNode
	Alias []string
}

func (Assignment) Name() string {
	return "Assignment"
}

func (a *Assignment) Graphviz(uid string) []string {
	a.Metadata["label"] = fmt.Sprintf("Assignment (%s)", a.Alias)
	return a.BasicNode.Graphviz(uid)
}

func (a *Assignment) Eval(c *runtime.Context) (runtime.Value, error) {
	if len(a.Childs) != len(a.Alias) {
		return runtime.Value{}, runtime.AssignmentMismatchException{}
	}
	var (
		value runtime.Value
		err   error
	)
	// Step 1: generate values
	results := make([]runtime.Value, len(a.Childs))
	for i, node := range a.Childs {
		value, err = node.Eval(c)
		if err != nil {
			return runtime.Value{}, err
		}
		results[i] = value
	}
	// Step 2: store them
	for i, value := range results {
		if err = c.Namespace.Update(value.Rename(a.Alias[i])); err != nil {
			return runtime.Value{}, err
		}
	}
	return value, nil
}

func NewMultiAssignment(alias []string, values ...Node) *Assignment {
	return &Assignment{
		BasicNode: NewBasic(values...),
		Alias:     alias,
	}
}

func NewAssignment(alias string, value Node) *Assignment {
	return NewMultiAssignment([]string{alias}, value)
}
