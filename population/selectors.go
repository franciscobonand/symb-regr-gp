package pop

import (
	"fmt"
	"math/rand"
	"sync"

	dataset "github.com/franciscobonand/symb-regr-gp/datasets"
)

// Selector is an interface for selecting individuals from population
type Selector interface {
	Select(pop Population, num int) Population
	String() string
}

// tournament selection defines a structure to select individuals using the tournament method
type tournament struct {
    tournamentSize int
    elitismSize int
    indivSelector Selector
    evaluator Evaluator
}

func TournamentSelector(elsize, tsize int, e Evaluator) Selector {
    return tournament{
        elitismSize: elsize,
        tournamentSize: tsize,
        indivSelector: randomSel{},
        evaluator: e,
    }
}

func (s tournament) String() string {
    return fmt.Sprintf("Tournament(%d)", s.tournamentSize)
}

func (s tournament) Select(pop Population, num int) Population {
    chosen := Population{}
    if s.elitismSize > 0 {
        chosen = pop.NBest(s.elitismSize)
    }
    for i := 0; i < num - s.elitismSize; i++ {
        group := s.indivSelector.Select(pop, s.tournamentSize)
        best := group.Best(s.evaluator)
        if !best.FitnessValid {
            panic("no best individual found!")
        }
        chosen = append(chosen, best)
    }
    return chosen
}

// roulette defines a structure to select individuals using the roulette method
type roulette struct {
    elitismSize int
    evaluator Evaluator
}

func RouletteSelector(elsize int, e Evaluator) Selector {
    return roulette{
        elitismSize: elsize,
        evaluator: e,
    }
}

func (s roulette) String() string {
    return "Roulette"
}

func (s roulette) Select(pop Population, num int) Population {
    chosen := Population{}
    fitSum := 0.0
    percSum := 0.0
    for _, indiv := range pop {
        fitSum += indiv.Fitness
    }
    for _, indiv := range pop {
        percSum += (1 - indiv.Fitness/fitSum)
    }
    if s.elitismSize > 0 {
        chosen = pop.NBest(s.elitismSize)
    }
    for i := 0; i < num - s.elitismSize; i++ {
        val := rand.Float64() * percSum
        for idx := range pop {
            val -= (1 - pop[idx].Fitness/fitSum)
            if val <= 0 {
                chosen = append(chosen, pop[idx])
                break
            }
        }
    }
    return chosen
}

// randomSel defines a structure to select individuals at random
type randomSel struct {
    elitismSize int
}

func RandomSelector(elsize int) Selector {
	return randomSel{
        elitismSize: elsize,
    }
}

func (s randomSel) String() string {
	return "RandomSelection"
}

func (s randomSel) Select(pop Population, num int) Population {
	chosen := Population{}
    if s.elitismSize > 0 {
        chosen = pop.NBest(s.elitismSize)
    }
	for i := 0; i < num - s.elitismSize; i++ {
		chosen = append(chosen, pop[rand.Intn(len(pop))])
	}
	return chosen
}

// lexicase defines a structure to select individuals using the lexicase method
type lexicase struct {
    elitismSize, threads int
    ds dataset.Dataset
    evaluator Evaluator
}

func LexicaseSelector(elsize, t int, e Evaluator, ds dataset.Dataset) Selector {
    return lexicase{
        elitismSize: elsize,
        threads: t,
        evaluator: e,
        ds: ds,
    }
}

func (s lexicase) String() string {
    return "LexicaseSelection"
}

func (s lexicase) Select(pop Population, num int) Population {
    chosen := Population{}
    if s.elitismSize > 0 {
        chosen = pop.NBest(s.elitismSize)
    }

    threads := s.threads
	chunkSize := (num - s.elitismSize) / threads
	if chunkSize < 1 {
		chunkSize = 1
        threads = 1
	}
	start := 0
	end := chunkSize
    cn := make(chan *Individual, num - s.elitismSize)
	var wg sync.WaitGroup
	wg.Add(threads)
    for chunk := 0; chunk < threads; chunk++ {
        if chunk == threads - 1 {
            end = (num - s.elitismSize)
        }
        go s.lexSelection(&wg, start, end, &pop, cn)
        start += chunkSize
        end += chunkSize
    }
    wg.Wait()
    close(cn)
    for ind := range cn {
        chosen = append(chosen, ind)
    }

    return chosen
}

func (s lexicase) lexSelection(wg *sync.WaitGroup, start, end int, pop *Population, cn chan *Individual) {
    for i := start; i < end; i++ {
        cases := s.ds.Copy()
        tempCandidates := pop.Clone()
        for {
            // When there are no cases left, pick one indiv at random
            if len(cases.Output) == 0 {
                cn <- tempCandidates[rand.Intn(len(tempCandidates))]
                break
            }
            cIdx := rand.Intn(len(cases.Output))
            cin := [][]float64{ cases.Input[cIdx] }
            cout := []float64{ cases.Output[cIdx] }
            rmse := RMSE{
                &dataset.Dataset{
                    Input: cin,
                    Output: cout,
                },
            }
            auxcand, _ := tempCandidates.Evaluate(rmse, s.threads)
            best := tempCandidates.Best(s.evaluator)
            tempCandidates = Population{}
            // Remove all indiv with fitness worse than the best fitness for this case
            for _, ind := range auxcand {
                if ind.Fitness == best.Fitness {
                    tempCandidates = append(tempCandidates, ind)
                }
            }
            // If there's only one candidate, it's chosen as parent
            if len(tempCandidates) == 1 {
                cn <- tempCandidates[0]
                break
            }
            // Remove used case
            cases.Input = append(cases.Input[:cIdx], cases.Input[cIdx+1:]...)
            cases.Output = append(cases.Output[:cIdx], cases.Output[cIdx+1:]...)
        }
    }
    wg.Done()
}
