package main

import (
	"flag"
	"fmt"

	"github.com/franciscobonand/symb-regr-gp/operator"
	"github.com/franciscobonand/symb-regr-gp/datasets"
	pop "github.com/franciscobonand/symb-regr-gp/population"
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

    pop.SetSeed(seed)
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
    }else {
        selector = pop.RandomSelector(nElitism)
    }
    mut := pop.MutationOp(gen)
    cross := pop.CrossoverOp()
    // Run Fitness for initial population
    p, e := p.Evaluate(rmse, threads)
    best := p.Best(rmse)
    b, w, m := p.FitnessStats()
    fmt.Printf("gen=%d evals=%d fit=%.4f\n", 0, e, best.Fitness)
    fmt.Printf("gen=%d best=%.4f worst=%.4f mean=%.4f\n", 0, b, w, m)

    ds.Copy()

    for i := 0; i < generations; i++ {
        // Selects new population
        children := selector.Select(p, len(p))
        // appleis genetic operators
        p, e = pop.ApplyGeneticOps(children, cross, mut, crossProb, mutProb).Evaluate(rmse, threads)
        // get best individual
        best = p.Best(rmse)
        b, w, m = p.FitnessStats()
        fmt.Printf("gen=%d evals=%d fit=%.4f\n", i, e, best.Fitness)
        fmt.Printf("gen=%d best=%.4f worst=%.4f mean=%.4f\n", i, b, w, m)
    }

    fmt.Println(best)
    // p.Print()
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

func allPositiveInts(nums... int) bool {
    for _, n := range nums {
        if n <= 0 {
            return false
        }
    }
    return true
}
