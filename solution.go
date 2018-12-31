// Basic representation of a solution.

package s1t

import (
	"fmt"
	"strings"
)

// Solution for a Problem.
type Solution struct {
	IsSat      bool
	Assignment []int // List from 0 to NumVars with the true/false/none assignment.
}

func unsat() Solution {
	return Solution{}
}

func sat(a []int) Solution {
	return Solution{
		IsSat:      true,
		Assignment: a,
	}
}

// Output returns the DIMACS format output for a solution of a problem.
func (s *Solution) Output(problem Problem) string {
	satNum := 1
	if !s.IsSat {
		satNum = 0
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("s %s %d %d %d\n",
		problem.Spec.Format, satNum,
		problem.Spec.NumVariables, problem.Spec.NumClauses))
	if !s.IsSat {
		return b.String()
	}
	for varNum, v := range s.Assignment {
		outNum := int(varNum + 1)
		if v == 0 {
			outNum = -outNum
		}
		b.WriteString(fmt.Sprintf("v %d\n", outNum))
	}
	return b.String()
}

// Satisfies sanity checks if the solution satifies the problem, returning true if so.
// Otherwise it returns false plus the first falsified clause.
func (s *Solution) Satisfies(p Problem) (bool, *Clause) {
	if !s.IsSat {
		panic("Checking Satisfied() on an unsat solution")
	}
	for _, c := range p.Clauses {
		c := c
		clauseSatisfied := false
		for _, l := range c.Literals {
			v := l.Var()
			if s.Assignment[v] == l.AsInt() {
				clauseSatisfied = true
				break
			}
		}
		if !clauseSatisfied {
			return false, &c
		}
	}
	return true, nil
}
