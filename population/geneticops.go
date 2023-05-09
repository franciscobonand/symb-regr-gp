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
func MutationOp(gen Generator, eval Evaluator) Variation {
	mutate := func(ind Population) Population {
		tree := ind[0].Code.Clone()
		pos := rand.Intn(len(tree))
		newtree := gen.Generate().Code
        newcode := tree.ReplaceSubtree(pos, newtree)
        newfit, _ := eval.GetFitness(newcode)
        if newfit < ind[0].Fitness {
            ind[0] = Create(newcode)
        }
		return ind
	}
	return &variation{mutate, fmt.Sprintf("Mutation(%s)", gen)}
}

// CrossoverOp returns a crossover variation
func CrossoverOp(eval Evaluator) Variation {
	cross := func(ind Population) Population {
		if ind[0].Size() < 2 || ind[1].Size() < 2 {
			return ind
		}
		pos1, subtree1 := ind[0].Code.RandomSubtree()
		pos2, subtree2 := ind[1].Code.RandomSubtree()
        newcode := ind[0].Code.Clone().ReplaceSubtree(pos1, subtree2)
        newfit, _ := eval.GetFitness(newcode)
        if newfit < ind[0].Fitness {
            ind[0] = Create(newcode)
        }
        newcode = ind[1].Code.Clone().ReplaceSubtree(pos2, subtree1)
        newfit, _ = eval.GetFitness(newcode)
        if newfit < ind[1].Fitness {
            ind[1] = Create(newcode)
        }
		return ind
	}
	return &variation{cross, "Crossover"}
}

// ApplyGeneticOps applies crossover and/or mutation operators based on their probability.
// Both operators can be applied in the same individual
func ApplyGeneticOps(pop Population, cross, mutate Variation, cxProb, mutProb float64) (Population, float64, float64) {
    var betterchild, worsechild float64
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
