package stats

import (
	"fmt"
	"sync"

	pop "github.com/franciscobonand/symb-regr-gp/population"
)

func PrintStats(wg *sync.WaitGroup, gen, evals int, p pop.Population, e pop.Evaluator) {
    b, w, m, r := getStats(gen, evals, p, e)
    fmt.Printf("gen=%d evals=%d repeated=%d best=%.4f worst=%.4f mean=%.4f\n", gen, evals, r, b, w, m)
    wg.Done()
}

func getStats(gen, evals int, p pop.Population, e pop.Evaluator) (float64, float64, float64, int) {
    b, w, m := p.FitnessStats()
    r := p.GetRepeatedIndividuals()
    return b, w, m, r
}
