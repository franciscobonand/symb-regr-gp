package pop

import (
	"fmt"
	"math"
	"sort"
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
