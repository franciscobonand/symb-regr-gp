package operator

import "math/rand"

// Opcode is an individual's genome smallest part
type Opcode interface {
	Arity() int
	Eval(...float64) float64
	String() string
	Format(...string) string
}

// Expr is slice of Opcodes, which represents the genome
type Expr []Opcode


// Clone makes a deep copy of an Expr
func (e Expr) Clone() Expr {
	return append([]Opcode{}, e...)
}

// Traverse walks the expression tree using a depth first traversal starting at element pos.
// If not nil then tfunc is called for each leaf node and nfunc is called for each node
// having one or more child nodes.
func (e Expr) Traverse(pos int, nfunc, tfunc func(Opcode)) int {
	op := e[pos]
	arity := op.Arity()
	if arity == 0 {
		if tfunc != nil {
			tfunc(op)
		}
	} else {
		for i := 0; i < arity; i++ {
			pos = e.Traverse(pos+1, nfunc, tfunc)
		}
		if nfunc != nil {
			nfunc(op)
		}
	}
	return pos
}

// Eval evaluates an expression for given input values by calling the Eval method on each Opcode.
func (e Expr) Eval(input ...float64) float64 {
	var doEval func() float64
	pos := -1
	doEval = func() float64 {
		pos++
		op := e[pos]
		arity := op.Arity()
		switch arity {
		case 0:
			return op.Eval(input...)
		case 1:
			return op.Eval(doEval())
		case 2:
			return op.Eval(doEval(), doEval())
		default:
			args := make([]float64, arity)
			for i := range args {
				args[i] = doEval()
			}
			return op.Eval(args...)
		}
	}
	return doEval()
}

// Format returns a string representation of an expression.
// It calls the Format method on each Opcode to return a result in infix notation.
func (e Expr) Format() string {
	list := []string{}
	node := func(op Opcode) {
		end := len(list) - op.Arity()
		list = append(list[:end], op.Format(list[end:]...))
	}
	term := func(op Opcode) {
		list = append(list, op.Format())
	}
	e.Traverse(0, node, term)
	return list[0]
}

// Depth returns the maximum height of the code tree from the root.
func (e Expr) Depth() int {
	stack := make([]int, 1, len(e))
	maxDepth, depth := 0, 0
	stack[0] = 0
	for _, op := range e {
		end := len(stack) - 1
		depth, stack = stack[end], stack[:end]
		if depth > maxDepth {
			maxDepth = depth
		}
		for i := 0; i < op.Arity(); i++ {
			stack = append(stack, depth+1)
		}
	}
	return maxDepth
}

// ReplaceSubtree replaces the code at pos with subtree without updating the subtree argument.
func (e Expr) ReplaceSubtree(pos int, subtree Expr) Expr {
	end := e.Traverse(pos, nil, nil)
	tail := subtree
	if end < len(e)-1 {
		tail = append(tail.Clone(), e[end+1:]...)
	}
	return append(e[:pos], tail...)
}

// RandomSubtree returns postion and a copy of nodes in randomly selected subtree of code
func (e Expr) RandomSubtree() (pos int, subtree Expr) {
	pos = rand.Intn(len(e))
	end := e.Traverse(pos, nil, nil)
	subtree = e[pos : end+1].Clone()
	return
}
