package s1t

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ParseDimacs parses input in DIMACS format.
func ParseDimacs(in io.Reader) (Problem, error) {
	s := bufio.NewScanner(in)
	var spec ProblemSpec
	var clauses []Clause
	prevClause := newClause()
	for s.Scan() {
		line := s.Text()
		if spec.Format == "" {
			err := parseSpec(line, &spec)
			if err != nil {
				return Problem{}, err
			}
		} else {
			err := parseCnfClause(line, spec, &prevClause, &clauses)
			if err != nil {
				return Problem{}, err
			}
		}
	}
	// 0 terminator is not required for the last clause, so just add if there.
	if len(prevClause.Negatives) != 0 || len(prevClause.Positives) != 0 {
		clauses = append(clauses, prevClause)
	}
	if len(clauses) != spec.NumClauses {
		return Problem{}, fmt.Errorf("Expected %d clauses, but got %d",
			spec.NumClauses, len(clauses))
	}
	return Problem{Spec: spec, Clauses: clauses}, nil
}

func parseSpec(line string, spec *ProblemSpec) error {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil
	}
	if fields[0] == "c" {
		return nil
	}
	if fields[0] != "p" {
		return fmt.Errorf("Spec line starts with unknown char %q: %q",
			fields[0], line)
	}
	if len(fields) != 4 {
		return fmt.Errorf("Expected 4 fields of spec but got %d fields: %q",
			len(fields), line)
	}
	fileFormat := fields[1]
	if !strings.Contains(fileFormat, "cnf") {
		return fmt.Errorf("Expected \"cnf\" format but got %q: %q",
			fileFormat, line)
	}
	vars, err := strconv.Atoi(fields[2])
	if err != nil || vars < 0 {
		return fmt.Errorf("Expected non-negative integer for number of vars: %q, %e",
			fields[2], err)
	}
	clauses, err := strconv.Atoi(fields[3])
	if err != nil || clauses < 0 {
		return fmt.Errorf("Expected non-negative integer for number of clauses: %q, %e",
			fields[3], err)
	}
	*spec = ProblemSpec{
		Format:       fileFormat,
		NumVariables: vars,
		NumClauses:   clauses,
	}
	return nil
}

const clauseTerminatorNum = 0

func parseCnfClause(line string, spec ProblemSpec,
	prevClause *Clause, clauses *[]Clause) error {
	fields := strings.Fields(line)
	for _, f := range fields {
		num, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("Failed to parse var %q in clause %q %e",
				f, line, err)
		}
		if num == clauseTerminatorNum {
			*clauses = append(*clauses, *prevClause)
			*prevClause = newClause()
			continue
		}
		if intAbs(num) > spec.NumVariables {
			return fmt.Errorf("Variable number %d goes beyond pre-declared num vars %d",
				intAbs(num), spec.NumVariables)
		}
		if num < 0 {
			prevClause.Negatives[VarNum(-num)] = true
		} else {
			prevClause.Positives[VarNum(num)] = true

		}
	}
	return nil
}

func newClause() Clause {
	return Clause{
		Positives: map[VarNum]bool{},
		Negatives: map[VarNum]bool{},
	}
}

func intAbs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}
