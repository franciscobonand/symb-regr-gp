package stats

import (
	"fmt"
	"sync"

	pop "github.com/franciscobonand/symb-regr-gp/population"
)

func PrintStats(wg *sync.WaitGroup, gen, evals int, p pop.Population, e pop.Evaluator) {
    best, b, w, m := getStats(gen, evals, p, e)
    fmt.Printf("gen=%d evals=%d fit=%.4f\n", gen, evals, best.Fitness)
    fmt.Printf("gen=%d best=%.4f worst=%.4f mean=%.4f\n", gen, b, w, m)
    wg.Done()
}

func getStats(gen, evals int, p pop.Population, e pop.Evaluator) (*pop.Individual, float64, float64, float64) {
    best := p.Best(e)
    b, w, m := p.FitnessStats()
    return best, b, w, m
}
