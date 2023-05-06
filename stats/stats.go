package stats

import (
	"fmt"
	"sync"

	pop "github.com/franciscobonand/symb-regr-gp/population"
)

func PrintRunStats(wg *sync.WaitGroup, gen, evals , bCxChild, wCxChild float64, p pop.Population, e pop.Evaluator) {
    s := p.GetStats()
    fmt.Printf("%.1f,%.1f,%.1f,%.3f,%.3f,%.3f,%.1f,%.1f,%.1f,%.1f,%.1f\n",
        gen,
        evals,
        s.Repeated,
        s.BestFit,
        s.WorstFit,
        s.MeanFit,
        s.MaxSize,
        s.MinSize,
        s.MeanSize,
        bCxChild,
        wCxChild,
    )
    wg.Done()
}

func GetRunStats(gen, evals , bCxChild, wCxChild float64, p pop.Population, e pop.Evaluator) []float64 {
    s := p.GetStats()
    data := []float64{
        gen,
        evals,
        s.Repeated,
        s.BestFit,
        s.WorstFit,
        s.MeanFit,
        s.MaxSize,
        s.MinSize,
        s.MeanSize,
        bCxChild,
        wCxChild,
    }
    return data
}

