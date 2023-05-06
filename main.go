package main

import (
	crand "crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"sync"

	"github.com/franciscobonand/symb-regr-gp/datasets"
	"github.com/franciscobonand/symb-regr-gp/operator"
	pop "github.com/franciscobonand/symb-regr-gp/population"
	"github.com/franciscobonand/symb-regr-gp/stats"
)

var (
    popSize, tournamentSize, threads, generations, nElitism int
    file, sel string
    crossProb, mutProb float64
    seed int64
)

func main() {
    // ./symb-regr-gp -popsize 20 -selector tour -toursize 2 -gens 20 -threads 1 -file "abcd.csv" -cxprob 0.9 -mutprob 0.05 -elitism 0 -seed 4132
    initializeFlags()

    if !allPositiveInts(popSize, threads, generations) {
        panic("Invalid value for popsize, gens or threads, must be a positive integer")
    }
    if nElitism < 0 {
        panic("Elitism size must be at least 0")
    }
    if sel == "tour" && tournamentSize < 2 {
        panic("Tournament size must be at least 2")
    }
    if crossProb < 0.0 || mutProb < 0.0 || crossProb > 1.0 || mutProb > 1.0 {
        panic("Genetic operators probability must be between 0.0 and 1.0")
    }

    // fmt.Println(popSize, tournamentSize, threads, file, crossProb, mutProb, seed)
    ds, err := dataset.Read(file)
    if err != nil {
        panic(err.Error())
    }

    setSeed(seed)
    opset := operator.CreateOpSet(ds.Variables...)
    gen := pop.NewRampedGenerator(opset, 1, 6)
    rmse := pop.RMSE{ DS: ds }
    // Create initial population
    p := pop.CreatePopulation(popSize, gen)
    // Define selection method and genetic operators
    var selector pop.Selector
    if sel == "rol" {
        selector = pop.RouletteSelector(nElitism, rmse) 
    } else if sel == "tour" {
        selector = pop.TournamentSelector(nElitism, tournamentSize, rmse)
    } else if sel == "lex" {
        selector = pop.LexicaseSelector(nElitism, threads, rmse, ds.Copy())
    } else {
        selector = pop.RandomSelector(nElitism)
    }
    mut := pop.MutationOp(gen)
    cross := pop.CrossoverOp()
    // Run Fitness for initial population
    p, e := p.Evaluate(rmse, threads)

    var wg sync.WaitGroup
    var betterCxChild, worseCxChild int
    wg.Add(generations + 1)
    fmt.Println("gen,evals,repeated,bestfit,worstfit,meanfit,maxsize,minsize,meansize,betterCxChild,worseCxChild")
    go stats.PrintStats(&wg, 0, e, betterCxChild, worseCxChild, p, rmse)
    for i := 0; i < generations; i++ {
        // Selects new population
        children := selector.Select(p, len(p))
        // appleis genetic operators
        p, betterCxChild, worseCxChild = pop.ApplyGeneticOps(children, cross, mut, crossProb, mutProb)
        p,e = p.Evaluate(rmse, threads)
        // print new population stats
        go stats.PrintStats(&wg, i+1, e, betterCxChild, worseCxChild, p, rmse)
    }

    wg.Wait()
    best := p.Best(rmse)
    fmt.Println(best)
}

func initializeFlags() {
    flag.IntVar(&popSize, "popsize", 20, "population size")
    flag.IntVar(&nElitism, "elitism", 0, "number of best members of elitism")
    flag.IntVar(&tournamentSize, "toursize", 2, "tournament size")
    flag.StringVar(&sel, "selector", "tour", "defines the selection method ('rol', 'tour', 'lex', or 'rand')")
    flag.IntVar(&generations, "gens", 10, "number of generations to run")
    flag.IntVar(&threads, "threads", 1, "quantity of threads to be used when evaluating")
    flag.StringVar(&file, "file", "datasets/synth1/synth1-train.csv", "csv file containing data to be processed")
    flag.Float64Var(&crossProb, "cxprob", 0.9, "crossover probability")
    flag.Float64Var(&mutProb, "mutprob", 0.05, "mutation probability")
    flag.Int64Var(&seed, "seed", 1, "seed for generating the initial population")
    flag.Parse()
}

// setSeed sets the given number as seed, or a random value if seed is <= 0
func setSeed(seed int64) int64 {
    if seed <= 0 {
        max := big.NewInt(2<<31 - 1)
        rseed, _ := crand.Int(crand.Reader, max)
        seed = rseed.Int64()
    }
    fmt.Println("random seed:", seed)
    rand.Seed(seed)
    return seed
}

func allPositiveInts(nums... int) bool {
    for _, n := range nums {
        if n <= 0 {
            return false
        }
    }
    return true
}
