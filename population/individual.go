package pop

import (
	"fmt"
	"github.com/franciscobonand/symb-regr-gp/operator"
)


// Individual is a member of the population. Code represents its genome
type Individual struct {
	Code         operator.Expr
	Fitness      float64
	FitnessValid bool
	depth        int
}

// Create constructor produces a new individual with copy of given Code (genome)
func Create(code operator.Expr) *Individual {
	return &Individual{Code: code.Clone()}
}

// Clone returns a deep copy of the given individual
func (ind *Individual) Clone() *Individual {
	return &Individual{
		Code:         ind.Code.Clone(),
		Fitness:      ind.Fitness,
		FitnessValid: ind.FitnessValid,
	}
}

// String returns a textual representation of the individual
func (ind Individual) String() string {
	if ind.FitnessValid {
		return fmt.Sprintf("%6.3f  %s", ind.Fitness, ind.Code.Format())
	} else {
		return fmt.Sprintf("%6s  %s", "????", ind.Code.Format())
	}
}

// Size returns the length of the individual's genome
func (ind *Individual) Size() int {
	return len(ind.Code)
}

// Depth returns the depth of the code tree for the individual
func (ind *Individual) Depth() int {
	if len(ind.Code) == 0 {
		return 0
	}
	if ind.depth == 0 {
		ind.depth = ind.Code.Depth()
	}
	return ind.depth
}
