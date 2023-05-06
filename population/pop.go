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

type Stats struct {
    Repeated, MaxSize, MinSize, MeanSize, BestFit, WorstFit, MeanFit float64
}

// FitnessStats returns best, worst and mean fitness of a population
func (pop Population) GetStats() Stats {
    stats := Stats{}
    var worstfit, meanfit float64
    bestfit := math.MaxFloat64
    var maxsize, meansize float64
    minsize := math.MaxFloat64
    set := map[string]bool{}
    for _, ind := range pop {
        sz := float64(ind.Size())
        meansize += sz
        if sz < minsize {
            minsize = sz
        }
        if sz > maxsize {
            maxsize = sz
        }
        set[ind.String()] = true
        if ind.FitnessValid {
            meanfit += ind.Fitness
            if ind.Fitness < bestfit {
                bestfit = ind.Fitness
            }
            if ind.Fitness > worstfit {
                worstfit = ind.Fitness
            }
        }
    }
    stats.BestFit = bestfit
    stats.WorstFit = worstfit
    stats.MeanFit = meanfit/float64(len(pop))
    stats.Repeated = float64(len(pop) - len(set))
    stats.MaxSize = maxsize
    stats.MinSize = minsize
    stats.MeanSize = meansize/float64(len(pop))
    return stats
}
