package pop

import (
	"fmt"
	"math/rand"
)

const MAX_DEPTH = 7

// Variation is an interface for applying genetic operations
type Variation interface {
	Variate(ind Population) Population
	String() string
}

// variation defines the base structure to be embedded by other genetic operators
type variation struct {
	vfunc      func(in Population) (out Population)
	name       string
}

func (v *variation) String() string {
	return v.name
}

func (v *variation) Variate(in Population) Population {
    out := v.vfunc(in.Clone())
    // returns parent if child exceed max depth
    for i := range in {
        if out[i].Depth() > MAX_DEPTH {
            out[i] = in[i] 
        }
    }
    return out
}

// MutationOp returns a mutation variation
func MutationOp(gen Generator) Variation {
	mutate := func(ind Population) Population {
		tree := ind[0].Code
		pos := rand.Intn(len(tree))
		newtree := gen.Generate().Code
		ind[0] = Create(tree.ReplaceSubtree(pos, newtree))
		return ind
	}
	return &variation{mutate, fmt.Sprintf("Mutation(%s)", gen)}
}

// CrossoverOp returns a crossover variation
func CrossoverOp() Variation {
	cross := func(ind Population) Population {
		if ind[0].Size() < 2 || ind[1].Size() < 2 {
			return ind
		}
		pos1, subtree1 := ind[0].Code.RandomSubtree()
		pos2, subtree2 := ind[1].Code.RandomSubtree()
		ind[0] = Create(ind[0].Code.ReplaceSubtree(pos1, subtree2))
		ind[1] = Create(ind[1].Code.ReplaceSubtree(pos2, subtree1))
		return ind
	}
	return &variation{cross, "Crossover"}
}

// ApplyGeneticOps applies crossover and/or mutation operators based on their probability.
// Both operators can be applied in the same individual
func ApplyGeneticOps(pop Population, cross, mutate Variation, cxProb, mutProb float64) (Population, int, int) {
    var betterchild, worsechild int
    cxindivs := Population{}
    totalfit := 0.0
	offspring := pop.Clone()
	for i := 1; i < len(pop); i += 2 {
		if rand.Float64() < cxProb {
			children := cross.Variate(offspring[i-1 : i+1])
			offspring[i-1], offspring[i] = children[0], children[1]
            cxindivs = append(cxindivs, children...)
		}
	}
	for i := 0; i < len(pop); i++ {
        totalfit += pop[i].Fitness
		if rand.Float64() < mutProb {
			children := mutate.Variate(offspring[i : i+1])
			offspring[i] = children[0]
		}
	}
    meanParentFit := totalfit / float64(len(pop))
    for _, ind := range cxindivs {
        if ind.Fitness > meanParentFit {
            betterchild++
        } else if ind.Fitness < meanParentFit {
            worsechild++
        }
    }
	return offspring, betterchild, worsechild
}
