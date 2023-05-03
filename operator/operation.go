package operator

import (
	"strings"
)

// BaseFunc is an abstract implementation of the Opcode interface for embedding in concrete types
type BaseFunc struct {
    OpName  string
    OpArity int
}

func (f *BaseFunc) Arity() int { return f.OpArity }

func (f *BaseFunc) String() string { return f.OpName }

func (f *BaseFunc) Eval(args ...float64) float64 { panic("abstract method!") }

func (f *BaseFunc) Format(args ...string) string {
    if len(args) > 0 {
        return f.OpName + "(" + strings.Join(args, ", ") + ")"
    } else {
        return f.OpName
    }
}

// binary operator type
type binOp struct{ *BaseFunc }

// Operator returns an opcode that represents a function that takes two arguments
func Operator(name string) Opcode {
    return binOp{&BaseFunc{name, 2}}
}

func (b binOp) Format(args ...string) string {
    return "(" + args[0] + " " + b.OpName + " " + args[1] + ")"
}

// variable type
type variable struct {
    *BaseFunc
    Narg int
}

// Variable returns an opcode that represents a variable (a leaf in the tree)
func Variable(name string, narg int) Opcode {
    return variable{&BaseFunc{name, 0}, narg}
}

func (v variable) Eval(input ...float64) float64 { return input[v.Narg] }

const ZEROISH = 1e-10

// numOp defines a numeric binary operator type
type numOp struct {
	Opcode
	fun func(a, b float64) float64
}

func (o numOp) Eval(args ...float64) float64 {
	return o.fun(args[0], args[1])
}

var Add numOp = numOp{
    Operator("+"),
    func(a, b float64) float64 { return a + b },
}

var Sub numOp = numOp{
    Operator("-"),
    func(a, b float64) float64 { return a - b },
}
 
var Mul numOp = numOp{
    Operator("*"),
    func(a, b float64) float64 { return a * b },
}

var Div numOp = numOp{
    Operator("/"),
    func(a, b float64) float64 { 
        if b > -ZEROISH && b < ZEROISH {
            return 0
        }
        return a / b
    },
}
