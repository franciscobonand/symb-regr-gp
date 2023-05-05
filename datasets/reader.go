package dataset

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Dataset defines the variables (x0, x1, x2...), the inputs which are the values
// a variable can assume and the expected output given a set of inputs
type Dataset struct {
    Input [][]float64
    Output []float64
    Variables []string
}

func (ds Dataset) Copy() Dataset {
    input := make([][]float64, len(ds.Input))
    input = append(input, ds.Input...)
    output := []float64{}
    output = append(output, ds.Output...)
    return Dataset{
        Input: input,
        Output: output,
    }
}

// Read reads a file resided in the given path.
// The path is relative to the directory the program is executed
func Read(fpath string) (*Dataset, error) {
    f, err := os.Open(fpath)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    ds := Dataset{}
    ds.Input = [][]float64{}
    scanner := bufio.NewScanner(f)
    count := 0
    for scanner.Scan() {
        inputs := []float64{}
        line := scanner.Text()
        items := strings.Split(line, ",")
        for i, str := range items {
            num, err := strconv.ParseFloat(str, 64)
            if err != nil {
                return nil, err
            }
            if i == len(items) - 1 {
                ds.Output = append(ds.Output, num) 
                break
            }
            inputs = append(inputs, num)
        }
        ds.Input = append(ds.Input, inputs)
        count++;
    }

    for i := range ds.Input[0] {
        ds.Variables = append(ds.Variables, fmt.Sprintf("x%d", i))
    }

    return &ds, nil
}
