package pop

import (
	"fmt"
	"math/rand"
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
    for _, indiv := range pop {
        fitSum += 1/indiv.Fitness
    }
    if s.elitismSize > 0 {
        chosen = pop.NBest(s.elitismSize)
    }
    for i := 0; i < num - s.elitismSize; i++ {
        val := rand.Float64() * fitSum
        for idx := range pop {
            val -= 1/pop[idx].Fitness
            if val <= 0 {
                chosen = append(chosen, pop[idx])
                break
            }
        }
    }
    return chosen
}

// randomSel defines a structure to select individuals at random
type randomSel struct{
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
