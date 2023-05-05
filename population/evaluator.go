package pop

import (
	"math"
	"sync"

	dataset "github.com/franciscobonand/symb-regr-gp/datasets"
	"github.com/franciscobonand/symb-regr-gp/operator"
)

// Evaluator is an interface that provides methods to calculate the fitness of an individual
type Evaluator interface {
    GetFitness(code operator.Expr) (float64, bool)
    CompareFitness(a, b float64) bool
}

// Evaluate calls the eval Evaluator to calculate the fitness for each individual.
// The evaluation process can be done in parallel
func (pop Population) Evaluate(eval Evaluator, threads int) (Population, int) {
	todo := make([]int, 0, len(pop))
	for i, ind := range pop {
		if !ind.FitnessValid {
			todo = append(todo, i)
		}
	}
	evals := len(todo)
	chunkSize := evals / threads
	if chunkSize < 1 {
		chunkSize = 1
        threads = 1
	}
	start := 0
	end := chunkSize
	var wg sync.WaitGroup
	wg.Add(threads)
	for chunk := 0; chunk < threads; chunk++ {
		if chunk == threads-1 {
			end = evals
		}
		go func(indices []int) {
			for _, i := range indices {
				pop[i].Fitness, pop[i].FitnessValid = eval.GetFitness(pop[i].Code)
			}
			wg.Done()
		}(todo[start:end])
		start += chunkSize
		end += chunkSize
	}
	wg.Wait()
	return pop, evals
}

// RMSE defines the root mean squared error evaluator (fitness closer to 0.0 is better)
type RMSE struct {
    DS *dataset.Dataset
}

func (e RMSE) GetFitness(code operator.Expr) (float64, bool) {
    var acc, count float64
    for i, input := range e.DS.Input {
        acc += math.Pow(code.Eval(input...) - e.DS.Output[i], 2)
        count++
    }
    if count == 0.0 {
        return -1, false
    }
    return math.Sqrt(acc / count), true
}

func (e RMSE) CompareFitness(a, b float64) bool {
    return a < b
}
