// CLI interface for s1t.

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jvoung/s1t"
)

func main() {
	startTime := time.Now()
	flag.Parse()
	remaining := flag.Args()
	var input *os.File
	if len(remaining) == 0 {
		input = os.Stdin
	} else if len(remaining) == 1 {
		var err error
		input, err = os.Open(remaining[0])
		if err != nil {
			panic(err)
		}
		defer input.Close()
	} else {
		fmt.Printf(
			"Can specify at most one input as argument (file or stdin), but "+
				"given %d arguments\n",
			len(remaining))
		os.Exit(1)
	}
	problem, err := s1t.ParseDimacs(input)
	if err != nil {
		fmt.Printf("Error parsing input %v: %e\n", input, err)
		os.Exit(1)
	}
	fmt.Printf("c Processing %d vars, %d clauses (parsed input in %f s)\n",
		problem.Spec.NumVariables, problem.Spec.NumClauses,
		time.Since(startTime).Seconds())
	solution := s1t.Solve(problem)
	fmt.Print(solution.Output(problem))
	fmt.Printf("t %s %d %d %f\n",
		problem.Spec.Format, problem.Spec.NumVariables, problem.Spec.NumVariables,
		time.Since(startTime).Seconds())
}
