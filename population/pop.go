package pop

import (
	crand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sort"

	"github.com/franciscobonand/symb-regr-gp/operator"
)

// Population is a slice of individuals
type Population []*Individual

// CreatePopulation creates a new population of popsize using the provided generator
func CreatePopulation(popsize int, gen Generator) Population {
    pop := make(Population, popsize)
    for i := range pop {
        pop[i] = gen.Generate()
    }
    return pop
}

// Len defines the Len interface method for sorting a Population
func (pop Population) Len() int {
    return len(pop)
}

// Less defines the Less interface method for sorting a Population
func (pop Population) Less(i, j int) bool {
    if !pop[i].FitnessValid {
        return false
    }
    if !pop[j].FitnessValid {
        return true
    }
    // should use 'CompareFitness' method so this could be generic/customizable
    return pop[i].Fitness < pop[j].Fitness
}

// Swap defines the Swap interface method for sorting a Population
func (pop Population) Swap(i, j int) {
    pop[i], pop[j] = pop[j], pop[i]
}

// Print prints out every individual from a population
func (pop Population) Print() {
    for i, ind := range pop {
        fmt.Printf("%4d: %s\n", i, *ind)
    }
}

// Clone makes a deep copy of all of the individuals in the population
func (pop Population) Clone() Population {
    newpop := make(Population, len(pop))
    for i, ind := range pop {
        newpop[i] = ind.Clone()
    }
    return newpop
}

// Best returns the individual with the best fitness
func (pop Population) Best(e Evaluator) *Individual {
    best := &Individual{}
    for _, ind := range pop {
        if ind.FitnessValid && (!best.FitnessValid || e.CompareFitness(ind.Fitness, best.Fitness)) {
            best = ind
        }
    }
    return best
}

// NBest returns the first nind best individuals
func (pop Population) NBest(nind int) Population {
    clone := pop.Clone()
    sort.Sort(clone)
    if len(clone) < nind {
        nind = len(clone)
    }
    return clone[:nind]
}

// FitnessStats returns best, worst and mean fitness of a population
func (pop Population) FitnessStats() (float64, float64, float64) {
    var worst, mean float64
    best := math.MaxFloat64
    nIndiv := float64(len(pop))
    for _, ind := range pop {
        if ind.FitnessValid {
            mean += ind.Fitness
            if ind.Fitness < best {
                best = ind.Fitness
            }
            if ind.Fitness > worst {
                worst = ind.Fitness
            }
        }
    }
    return best, worst, mean/nIndiv
}

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

// SetSeed sets the given number as seed, or a random value if seed is <= 0
func SetSeed(seed int64) int64 {
    if seed <= 0 {
        max := big.NewInt(2<<31 - 1)
        rseed, _ := crand.Int(crand.Reader, max)
        seed = rseed.Int64()
    }
    fmt.Println("random seed:", seed)
    rand.Seed(seed)
    return seed
}
