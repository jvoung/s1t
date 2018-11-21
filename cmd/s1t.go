// CLI interface for s1t.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jvoung/s1t"
)

func main() {
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
				"given %d arguments",
			len(remaining))
		os.Exit(1)
	}
	problem, err := s1t.ParseDimacs(input)
	if err != nil {
		fmt.Printf("Error parsing input %v: %e", input, err)
		os.Exit(1)
	}
	fmt.Printf("Processing problem with %d vars, %d clauses",
		problem.Spec.NumVariables, problem.Spec.NumClauses)
}
