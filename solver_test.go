package s1t

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type solverTestCase struct {
	desc             string
	inputLines       []string
	expectedSolution Solution
}

// Simple tests that can be solved by unit propagation alone, etc.
func TestSimpleSat(t *testing.T) {
	cases := []solverTestCase{
		{
			desc: "One unit clause (sat pos)",
			inputLines: []string{
				"p cnf 1 1",
				"1 0",
			},
			expectedSolution: sat([]int{1}),
		},
		{
			desc: "One unit clause (sat neg)",
			inputLines: []string{
				"p cnf 1 1",
				"-1 0",
			},
			expectedSolution: sat([]int{0}),
		},
		{
			desc: "Has empty clause (unsat)",
			inputLines: []string{
				"p cnf 1 2",
				"1 0", "0",
			},
			expectedSolution: unsat(),
		},
		{
			desc: "Three unit clauses (sat)",
			inputLines: []string{
				"p cnf 3 3",
				"1 0", "-2 0", "3 0",
			},
			expectedSolution: sat([]int{1, 0, 1}),
		},
		{
			desc: "Three unit clauses (unsat)",
			inputLines: []string{
				"p cnf 2 3",
				"1 0", "-2 0", "-1 0",
			},
			expectedSolution: unsat(),
		},
		{
			desc: "Pure unit propagate (sat)",
			inputLines: []string{
				"p cnf 3 3",
				"1 -2 0", "-2 0", "-2 3 0",
			},
			expectedSolution: sat([]int{1, 0, 1}),
		},
	}
	for _, c := range cases {
		problem := inputToProblem(c.inputLines, t)
		t.Logf("Testing case %v: %v\n", c.desc, problem)
		solution := Solve(problem)
		if !equalSolution(solution, c.expectedSolution) {
			t.Errorf("Case %q, expected solution %v, but got %v",
				c.desc, c.expectedSolution, solution)
		}
	}
}

// Tests that need some backtracking.
func TestBacktrack(t *testing.T) {
	cases := []solverTestCase{
		{
			desc: "Two vars sat (1 1)",
			inputLines: []string{
				"p cnf 2 3",
				// X1 => X2 and X2 => X1 (so X1 <==> X2)
				"-1 2 0", "-2 1 0", "1 2 0",
			},
			expectedSolution: sat([]int{1, 1}),
		},
		{
			desc: "Two vars sat (0 0)",
			inputLines: []string{
				"p cnf 2 3",
				"-1 2 0", "-2 1 0", "-2 -1 0",
			},
			expectedSolution: sat([]int{0, 0}),
		},
		{
			desc: "Two vars sat (0 1)",
			inputLines: []string{
				"p cnf 2 3",
				"1 2 0", "-1 -2 0", "-1 2 0",
			},
			expectedSolution: sat([]int{0, 1}),
		},
		{
			desc: "Two vars sat (1 0)",
			inputLines: []string{
				"p cnf 2 3",
				"1 2 0", "-1 -2 0", "1 -2 0",
			},
			expectedSolution: sat([]int{1, 0}),
		},
		{
			desc: "2-towers sat",
			inputLines: []string{
				"c 1 == x11, 2 == x12, 3 == x21, 4 == x22",
				"c x11 => -x12, x11 => -x21, x12 => -x11, x12 => -x22, etc.",
				"p cnf 4 10",
				"-1 -2 0",
				"-1 -3 0",
				"-2 -1 0",
				"-2 -4 0",
				"-3 -1 0",
				"-3 -4 0",
				"-4 -2 0",
				"-4 -3 0",
				"1 2 0",
				"3 4 0",
			},
			expectedSolution: sat([]int{1, 0, 0, 1}), // {0, 1, 1, 0} also possible
		},
	}
	for _, c := range cases {
		problem := inputToProblem(c.inputLines, t)
		t.Logf("Testing problem %v: %v\n", c.desc, problem)
		solution := Solve(problem)
		if !equalSolution(solution, c.expectedSolution) {
			t.Errorf("Case %q, expected solution %v, but got %v",
				c.desc, c.expectedSolution, solution)
		}
	}
}

// Randomly generated subset sum problem from http://toughsat.appspot.com/
func TestSubsetSum2(t *testing.T) {
	expectedSolution := sat([]int{1, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0})
	testFromFile(t, "test_cnf/subsetsum2.cnf", expectedSolution)
}

func TestSubsetSum2b(t *testing.T) {
	expectedSolution := sat([]int{0, 1, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0})
	testFromFile(t, "test_cnf/subsetsum2b.cnf", expectedSolution)
}

func TestSubsetSum3(t *testing.T) {
	expectedSolution := sat([]int{
		1, 1, 0, 1, 0, 1, 1, 0, 1, 1,
		0, 1, 0, 1, 1, 0, 0, 0, 0, 1,
		0, 0, 0, 1, 1, 1, 0, 0, 1, 0,
		0, 0, 0})
	testFromFile(t, "test_cnf/subsetsum3.cnf", expectedSolution)
}

func TestQueen3(t *testing.T) {
	expectedSolution := unsat()
	testFromFile(t, "test_cnf/queen3.cnf", expectedSolution)
}

func Test4Queens(t *testing.T) {
	// Other solutions possible too (e.g., reflections):
	expectedSolution := sat([]int{
		0, 0, 1, 0,
		1, 0, 0, 0,
		0, 0, 0, 1,
		0, 1, 0, 0})
	expectedSolution2 := sat([]int{
		0, 1, 0, 0,
		0, 0, 0, 1,
		1, 0, 0, 0,
		0, 0, 1, 0})
	testFromFile(t, "test_cnf/queen4.cnf", expectedSolution, expectedSolution2)
}

// Some tests from https://www.cs.ubc.ca/~hoos/SATLIB/benchm.html
// which are known to run quickly for sanity testing.
func TestMediumSat(t *testing.T) {
	testFromFile(t, "test_cnf/hole6.cnf", unsat())
	testFromFileSelfCheck(t, "test_cnf/blocksworld_anomaly.cnf")
	testFromFileSelfCheck(t, "test_cnf/blocksworld_medium.cnf")
	testFromFileSelfCheck(t, "test_cnf/RTI_k3_n100_m429_0.cnf")
	testFromFileSelfCheck(t, "test_cnf/RTI_k3_n100_m429_499.cnf")
}

// Some benchmarks from https://www.cs.ubc.ca/~hoos/SATLIB/benchm.html
// which currently run slowly.
func BenchmarkPigeonHole(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testFromFile(b, "test_cnf_slow/hole8.cnf", unsat())
		testFromFile(b, "test_cnf_slow/hole9.cnf", unsat())
	}
}

func BenchmarkHanoi4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testFromFileSelfCheck(b, "test_cnf_slow/hanoi4.cnf")
	}
}
func BenchmarkBMS_k3_n100_m429(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testFromFileSelfCheck(b, "test_cnf_slow/BMS_k3_n100_m429_0.cnf")
		testFromFileSelfCheck(b, "test_cnf_slow/BMS_k3_n100_m429_499.cnf")
	}
}

func BenchmarkGraphColoring(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testFromFileSelfCheck(b, "test_cnf_slow/flat175-1.cnf")
		testFromFileSelfCheck(b, "test_cnf_slow/flat200-1.cnf")
	}
}
func BenchmarkBlockworld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testFromFileSelfCheck(b, "test_cnf_slow/bw_large.a.cnf")
		testFromFileSelfCheck(b, "test_cnf_slow/bw_large.b.cnf")
	}
}

func testFromFile(tb testing.TB, relativePath string, expectedSolutions ...Solution) {
	checker := func(tb testing.TB, problem Problem, solution Solution) {
		matched := false
		for _, es := range expectedSolutions {
			if equalSolution(solution, es) {
				matched = true
				break
			}
		}
		if !matched {
			tb.Errorf("Case %q, none of the expected solutions match (%v |%d|). Got %v",
				relativePath, expectedSolutions, len(expectedSolutions), solution)
		}
	}
	testFromFileWithChecker(tb, relativePath, checker)
}

func testFromFileSelfCheck(tb testing.TB, relativePath string) {
	checker := func(tb testing.TB, problem Problem, solution Solution) {
		sat, failedClause := solution.Satisfies(problem)
		if !sat {
			tb.Errorf("Case %q, solution does not satisfy clause %v (sol=%v)",
				relativePath, failedClause, solution)
		}
	}
	testFromFileWithChecker(tb, relativePath, checker)
}

func testFromFileWithChecker(
	tb testing.TB, relativePath string, checker func(testing.TB, Problem, Solution)) {
	input, err := os.Open(relativePath)
	if err != nil {
		tb.Fatalf("Failed to open test file: %v", err)
	}
	defer input.Close()
	problem := parseOrDie(input, tb)
	solution := Solve(problem)
	checker(tb, problem, solution)
}

func equalSolution(s1, s2 Solution) bool {
	return cmp.Equal(s1, s2)
}

func inputToProblem(lines []string, tb *testing.T) Problem {
	return parseOrDie(strings.NewReader(strings.Join(lines, "\n")), tb)
}

func parseOrDie(in io.Reader, tb testing.TB) Problem {
	problem, err := ParseDimacs(in)
	if err != nil {
		tb.Fatalf("Failed to parse input: %v", err)
	}
	return problem
}
