package stats

import (
	"fmt"
	"sync"

	pop "github.com/franciscobonand/symb-regr-gp/population"
)

func PrintStats(wg *sync.WaitGroup, gen, evals, bCxChild, wCxChild int, p pop.Population, e pop.Evaluator) {
    s := p.GetStats()
    fmt.Printf("%d,%d,%d,%.4f,%.4f,%.4f,%d,%d,%d,%d,%d\n",
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

