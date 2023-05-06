package stats

import (
	"fmt"
	"sync"

	pop "github.com/franciscobonand/symb-regr-gp/population"
)

func PrintStats(wg *sync.WaitGroup, gen, evals int, p pop.Population, e pop.Evaluator) {
    s := p.GetStats()
    fmt.Printf("gen=%d evals=%d repeated=%d bestfit=%.4f worstfit=%.4f meanfit=%.4f maxsize=%d minsize=%d meansize=%d\n",
        gen,
        evals,
        s.Repeated,
        s.BestFit,
        s.WorstFit,
        s.MeanFit,
        s.MaxSize,
        s.MinSize,
        s.MeanSize,
    )
    wg.Done()
}

