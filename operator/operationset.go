package operator

import (
	"fmt"
	"strconv"
)

// OpSet represents the set of all available functions/variables.
// NumVars is the number of input variables, Terminals a list of all the variables
// and Primitives are the operators 
type OpSet struct {
	NumVars    int
	Terminals  []Opcode
	Primitives []Opcode
}

// CreateOpSet returns a set of available operators (Add, Sub, Mul and Div)
// and nvars input variables.
func CreateOpSet(varNames ...string) *OpSet {
	pset := &OpSet{}
    nvars := len(varNames)
	pset.NumVars = nvars
	pset.Terminals = make([]Opcode, nvars)
	for i := 0; i < nvars; i++ {
		var name string
		if len(varNames) > i {
			name = varNames[i]
		} else {
			name = "in" + strconv.Itoa(i)
		}
		pset.Terminals[i] = Variable(name, i)
	}
    pset.Primitives = []Opcode{
        Add,
        Sub,
        Mul,
        Div,
    }
	return pset
}

// String returns a string representation of the operators and variables
func (pset *OpSet) String() string {
    ops := append(pset.Terminals, pset.Primitives...)
	return fmt.Sprint(ops)
}

// Var returns the nth variable
func (pset *OpSet) Var(n int) Opcode {
	return pset.Terminals[n]
}
