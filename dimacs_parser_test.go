package s1t

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPreambleOnly(t *testing.T) {
	problem, err := ParseDimacs(strings.NewReader(
		strings.Join([]string{
			"c comment 123",
			"c another comment",
			"p cnf 0 0",
		}, "\n")))
	if err != nil {
		t.Errorf("Expected no errors but got %e", err)
	}
	if problem.Spec.NumClauses != 0 {
		t.Errorf("Expected no clauses but got %d", problem.Spec.NumClauses)
	}
	if problem.Spec.NumVariables != 0 {
		t.Errorf("Expected no variables but got %d", problem.Spec.NumVariables)
	}
	if problem.Spec.Format != "cnf" {
		t.Errorf("Expected cnf format but got %q", problem.Spec.Format)
	}
	if len(problem.Clauses) != 0 {
		t.Errorf("Expected no clauses but got %d", len(problem.Clauses))
	}
}

type errorTestCase struct {
	desc                 string
	lines                []string
	expectedErrSubstring string
}

func TestPreambleErrorCases(t *testing.T) {
	cases := []errorTestCase{
		{
			desc: "Too many problem statement lines",
			lines: []string{
				"c comment 123",
				"p cnf 0 0",
				"p cnf 0 0",
			},
		},
		{
			desc: "Unknown first char in preamble",
			lines: []string{
				"1 2 3",
				"p cnf 3 1",
			},
		},
		{
			desc:  "Truncated cnf spec",
			lines: []string{"p cnf 0"},
		},
		{
			desc:  "Negative number of variables",
			lines: []string{"p cnf -1 0"},
		},
		{
			desc:  "Negative number of clauses",
			lines: []string{"p cnf 0 -1"},
		},
	}

	for _, c := range cases {
		_, err := ParseDimacs(strings.NewReader(
			strings.Join(c.lines, "\n")))
		// Once it parses a problem spec, it's no longer using the preamble parser
		// so it's not looking for additional problem specs in particular.
		// Some downstream parser will fail, and it's hard to predict how.
		// Just check that there *is* and error instead of silently continuing.
		if err == nil {
			t.Errorf("case %q, expected an error but got none", c.desc)
		}
		if c.expectedErrSubstring != "" && !strings.Contains(err.Error(), c.expectedErrSubstring) {
			t.Errorf("case %q, expected err string %q but got %e",
				c.desc, c.expectedErrSubstring, err)
		}
	}
}

func TestUnknownFormatInPreamble(t *testing.T) {
	_, err := ParseDimacs(strings.NewReader(
		strings.Join([]string{
			"p magic 0 0",
		}, "\n")))
	if err == nil {
		t.Error("Expected an error but got none")
	}
}

type cnfTestCase struct {
	desc            string
	inputLines      []string
	expectedClauses []Clause
}

func TestCnfClauses(t *testing.T) {
	cases := []cnfTestCase{
		{
			desc: "Normal one line two vars, with 0 terminator",
			inputLines: []string{
				"p cnf 2 1",
				"1 -2 0",
			},
			expectedClauses: []Clause{
				Clause{Literals: []Literal{Positive(0), Negative(1)}},
			},
		},
		{
			desc: "Normal one line two vars, no 0 terminator",
			inputLines: []string{
				"p cnf 2 1",
				"1 -2",
			},
			expectedClauses: []Clause{
				Clause{Literals: []Literal{Positive(0), Negative(1)}},
			},
		},
		{
			desc: "Multiline with clause split across lines",
			inputLines: []string{
				"p cnf 4 3",
				"1 3 -4 0",
				"4 0 2",
				"-3",
			},
			expectedClauses: []Clause{
				Clause{Literals: []Literal{Positive(0), Positive(2), Negative(3)}},
				Clause{Literals: []Literal{Positive(3)}},
				Clause{Literals: []Literal{Positive(1), Negative(2)}},
			},
		},
		{
			desc: "Multiline with splitter on own line",
			inputLines: []string{
				"p cnf 4 3",
				"1 3 -4",
				"0",
				"4",
				"0",
				"2 -3",
				"0",
			},
			expectedClauses: []Clause{
				Clause{Literals: []Literal{Positive(0), Positive(2), Negative(3)}},
				Clause{Literals: []Literal{Positive(3)}},
				Clause{Literals: []Literal{Positive(1), Negative(2)}},
			},
		},
		{
			desc: "Duplicated literals in same clause",
			inputLines: []string{
				"p cnf 2 1",
				"1 -2 2 -1 2 1 -1 -2 0",
			},
			expectedClauses: []Clause{
				Clause{
					Literals: []Literal{Positive(0), Negative(1), Positive(1), Negative(0)},
				},
			},
		},
	}
	for _, c := range cases {
		problem, err := ParseDimacs(strings.NewReader(
			strings.Join(c.inputLines, "\n")))
		if err != nil {
			t.Errorf("Case %q, expected no errors but got %e",
				c.desc, err)
		}
		if !equalClauses(problem.Clauses, c.expectedClauses) {
			t.Errorf("Case %q, expected clauses %v, but got %v",
				c.desc, c.expectedClauses, problem.Clauses)
		}
	}
}

func equalClauses(c1, c2 []Clause) bool {
	if len(c1) != len(c2) {
		return false
	}
	for i, a := range c1 {
		b := c2[i]
		if !cmp.Equal(a, b) {
			return false
		}
	}
	return true
}

func TestCnfErrorCases(t *testing.T) {
	cases := []errorTestCase{
		{
			desc: "Variable number greater than predeclared",
			lines: []string{
				"p cnf 1 1",
				"1 -2 0",
			},
			expectedErrSubstring: "Variable number 2 goes beyond pre-declared",
		},
		{
			desc: "Clauses greater than predeclared",
			lines: []string{
				"p cnf 2 1",
				"1 -2 0",
				"1 0",
			},
			expectedErrSubstring: "Expected 1 clauses, but got 2",
		},
		{
			desc: "Clauses fewer than predeclared",
			lines: []string{
				"p cnf 2 3",
				"1 -2 0",
				"1 0",
			},
			expectedErrSubstring: "Expected 3 clauses, but got 2",
		},
	}

	for _, c := range cases {
		_, err := ParseDimacs(strings.NewReader(
			strings.Join(c.lines, "\n")))
		if err == nil {
			t.Errorf("case %q, expected an error but got none", c.desc)
		}
		if c.expectedErrSubstring != "" && !strings.Contains(err.Error(), c.expectedErrSubstring) {
			t.Errorf("case %q, expected err string %q but got %e",
				c.desc, c.expectedErrSubstring, err)
		}
	}
}
