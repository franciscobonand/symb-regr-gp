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
    indivSelector Selector
    evaluator Evaluator
}

func TournamentSelector(tsize int, e Evaluator) Selector {
    return tournament{
        tsize,
        randomSel{},
        e,
    }
}

func (s tournament) String() string {
    return fmt.Sprintf("Tournament(%d)", s.tournamentSize)
}

func (s tournament) Select(pop Population, num int) Population {
    chosen := Population{}
    for i := 0; i < num; i++ {
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
type roulette struct{}

func RouletteSelector() Selector {
    return roulette{}
}

func (s roulette) String() string {
    return "Roulette"
}

func (s roulette) Select(pop Population, num int) Population {
    chosen := Population{}
    fitSum := 0.0
    for _, indiv := range pop {
        fitSum += indiv.Fitness
    }

    for i := 0; i < num; i++ {
        val := rand.Float64() * fitSum
        for idx, indiv := range pop {
            val -= indiv.Fitness
            if val <= 0 {
                chosen = append(chosen, pop[idx])
            }
        }
    }

    return chosen
}

// randomSel defines a structure to select individuals at random
type randomSel struct{}

func RandomSelector() Selector {
	return randomSel{}
}

func (s randomSel) String() string {
	return "RandomSelection"
}

func (s randomSel) Select(pop Population, num int) Population {
	chosen := Population{}
	for i := 0; i < num; i++ {
		chosen = append(chosen, pop[rand.Intn(len(pop))])
	}
	return chosen
}
