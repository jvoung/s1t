package main

import (
	"bufio"
	"io"
	"math"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSolvableBoard(t *testing.T) {
	input, err := os.Open("test_data/board1.in")
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	board := ParseBoard(input)
	boardStrBuilder := strings.Builder{}
	PrintBoard(board, "", &boardStrBuilder)
	boardStr := boardStrBuilder.String()
	expectedBoard := `4 0 0 |0 0 0 |8 0 5
0 3 0 |0 0 0 |0 0 0
0 0 0 |7 0 0 |0 0 0
------+------+------
0 2 0 |0 0 0 |0 6 0
0 0 0 |0 8 0 |4 0 0
0 0 0 |0 1 0 |0 0 0
------+------+------
0 0 0 |6 0 3 |0 7 0
5 0 0 |2 0 0 |0 0 0
1 0 4 |0 0 0 |0 0 0
`
	if boardStr != expectedBoard {
		t.Errorf("Parse failed, got %v instead of %v", boardStr, expectedBoard)
	}
	solvedBoard := solveBoard(board)
	boardStrBuilder.Reset()
	PrintBoard(solvedBoard, "", &boardStrBuilder)
	solvedBoardStr := boardStrBuilder.String()
	expectedSolvedBoard := `4 1 7 |3 6 9 |8 2 5
6 3 2 |1 5 8 |9 4 7
9 5 8 |7 2 4 |3 1 6
------+------+------
8 2 5 |4 3 7 |1 6 9
7 9 1 |5 8 6 |4 3 2
3 4 6 |9 1 2 |7 5 8
------+------+------
2 8 9 |6 4 3 |5 7 1
5 7 3 |2 9 1 |6 8 4
1 6 4 |8 7 5 |2 9 3
`
	if solvedBoardStr != expectedSolvedBoard {
		t.Errorf("Solution failed, got %v instead of %v", solvedBoardStr, expectedSolvedBoard)
	}
}

func BenchmarkTop95(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testFromFileLines(b, "test_data/top95.txt")
	}
}

func BenchmarkHardest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testFromFileLines(b, "test_data/hardest.txt")
	}
}

func BenchmarkMultiSolution(b *testing.B) {
	// The random sudoku board test case where:
	//"Unfortunately, this is not a true Sudoku puzzle because it has multiple solutions".
	for i := 0; i < b.N; i++ {
		testFromFileLines(b, "test_data/multi_solution.txt")
	}
}

func testFromFileLines(b *testing.B, fname string) {
	input, err := os.Open(fname)
	if err != nil {
		b.Fatal(err)
	}
	defer input.Close()
	scanner := bufio.NewScanner(input)
	min := math.MaxFloat64
	max := 0.0
	total := 0.0
	numLines := 0
	for scanner.Scan() {
		s := time.Now()
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		parseAndSolve(strings.NewReader(line))
		duration := time.Since(s).Seconds()
		min = math.Min(min, duration)
		max = math.Max(max, duration)
		total += duration
		numLines++
	}
	b.Logf("Solved: %d, Min: %f, Max: %f, Avg: %f", numLines, min, max, total/float64(numLines))
}

func parseAndSolve(input io.Reader) Board {
	board := ParseBoard(input)
	return solveBoard(board)
}
