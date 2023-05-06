package pop

import (
    "math/rand"
	"github.com/franciscobonand/symb-regr-gp/operator"
)

// A Generator is used to generate new individuals from the provided operations set
type Generator interface {
    Generate() *Individual
    String() string
}

// genBase is the base generator to be embedded by the other generators
type genBase struct {
    pset      *operator.OpSet
    min, max  int
    condition func(height, depth int) bool
    name      string
}

func (g genBase) String() string {
    return g.name
}

// Generate defines the core logic of the generators
func (g genBase) Generate() *Individual {
    code := operator.Expr{}
    height := rand.Intn(1+g.max-g.min) + g.min
    stack := []int{0}
    depth := 0
    for len(stack) > 0 {
        depth, stack = stack[len(stack)-1], stack[:len(stack)-1]
        if g.condition(height, depth) {
            op := randomOp(g.pset.Terminals)
            code = append(code, op)
        } else {
            op := randomOp(g.pset.Primitives)
            code = append(code, op)
            for i := 0; i < op.Arity(); i++ {
                stack = append(stack, depth+1)
            }
        }
    }
    return &Individual{Code: code}
}

// NewGrowGenerator returns a generator to produce individuals with irregular expression trees
func NewGrowGenerator(ops *operator.OpSet, min, max int) Generator {
    terms, prims := len(ops.Terminals), len(ops.Primitives)
    terminalRatio := float64(terms) / float64(terms+prims)
    return genBase{
        ops, min, max,
        func(height, depth int) bool {
            return depth == height || (depth >= min && rand.Float64() < terminalRatio)
        },
        "GrowGenerator",
    }
}

// NewFullGenerator returns a generator to produce individuals with balanced expression trees
func NewFullGenerator(ops *operator.OpSet, min, max int) Generator {
    return genBase{
        ops, min, max,
        func(height, depth int) bool {
            return depth == height 
        },
        "FullGenerator",
    }
}

type rampedGenerator struct {
    grow, full Generator
}

// NewRampedGenerator returns a Ramped population generator (combination of Grow and Full)
func NewRampedGenerator(ops *operator.OpSet, min, max int) Generator {
    return rampedGenerator{
        NewGrowGenerator(ops, min, max),
        NewFullGenerator(ops, min, max),
    }
}

func (rg rampedGenerator) String() string {
    return "RampedGenerator"
}

func (rg rampedGenerator) Generate() *Individual {
    if rand.Float64() >= 0.5 {
        return rg.grow.Generate()
    }
    return rg.full.Generate()
}

func randomOp(list []operator.Opcode) operator.Opcode {
    return list[rand.Intn(len(list))]
}
