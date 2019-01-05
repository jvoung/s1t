// Reads a sudoku board and uses CNF constraints + SAT solver to solve the board.
// Input format is like Peter Norvig's from: http://norvig.com/sudoku.html
// TODO(jvoung): generate random (solvable) boards.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jvoung/s1t"
)

var toCnf = flag.Bool("cnf", false, "convert unsolved board to cnf")
var toBoard = flag.Bool("board", false, "parse assignment and print solved board")
var endToEnd = flag.Bool("all", false, "unsolved board to cnf => solve => print")

func main() {
	flag.Parse()
	args := flag.Args()
	var input *os.File
	if len(args) == 0 {
		input = os.Stdin
	} else if len(args) == 1 {
		var err error
		input, err = os.Open(args[0])
		if err != nil {
			panic(err)
		}
		defer input.Close()
	} else {
		panic(fmt.Sprintf(
			"Can specify at most one input as argument (file or stdin), but "+
				"given %d arguments\n",
			len(args)))
	}
	if *toCnf {
		board := ParseBoard(input)
		WriteCNF(board, os.Stdout)
	} else if *toBoard {
		board := ParseAssignments(input)
		PrintBoard(board, "", os.Stdout)
	} else if *endToEnd {
		startTime := time.Now()
		solveEndToEnd(input)
		fmt.Printf("Solved in %f\n", time.Since(startTime).Seconds())
	} else {
		flag.Usage()
	}
}

// Board holds the board values. 0 means blank, otherwise 1-9 are set.
type Board [][]int

func solveEndToEnd(input io.Reader) {
	board := ParseBoard(input)
	fmt.Println("Solving board:")
	PrintBoard(board, "", os.Stdout)
	fmt.Println("and got:")
	solvedBoard := solveBoard(board)
	PrintBoard(solvedBoard, "", os.Stdout)
}

func solveBoard(board Board) Board {
	buf := strings.Builder{}
	WriteCNF(board, &buf)
	problem, err := s1t.ParseDimacs(strings.NewReader(buf.String()))
	if err != nil {
		panic(err)
	}
	solution := s1t.Solve(problem)
	solutionStr := solution.Output(problem)
	solvedBoard := ParseAssignments(strings.NewReader(solutionStr))
	return solvedBoard
}

// ParseBoard parses a board givee format like in: http://norvig.com/sudoku.html
func ParseBoard(input io.Reader) Board {
	b := make(Board, 9)
	for r := 0; r < 9; r++ {
		b[r] = make([]int, 9)
	}
	i := 0
	scanner := bufio.NewScanner(input)
	acceptedChars := ".0123456789"
	for scanner.Scan() {
		line := scanner.Text()
		for j, c := range line {
			cindex := strings.IndexRune(acceptedChars, c)
			if cindex < 0 {
				continue
			}
			// Convert "." to "0"
			if c == '.' {
				c = '0'
			}
			v, err := strconv.Atoi(string(c))
			if err != nil {
				panic(fmt.Sprintf("Failed to parse digit %v: %v", c, err))
			}
			r, c := i/9, i%9
			b[r][c] = v
			i++
			if i > len(b)*len(b[0]) {
				panic(fmt.Sprintf(
					"Got more than %d digits. Last digit was %v on index %v of line %v",
					len(b)*len(b[0]), c, j, line))
			}
		}
	}
	return b
}

// PrintBoard pretty prints a Board.
func PrintBoard(b Board, prefix string, w io.Writer) {
	dashes := strings.Repeat("-", 6)
	dashArray := []string{dashes, dashes, dashes}
	dashLine := strings.Join(dashArray, "+")
	for r := 0; r < len(b); r++ {
		if r == 3 || r == 6 {
			fmt.Fprintf(w, "%v%v\n", prefix, dashLine)
		}
		fmt.Fprintf(w, "%v", prefix)
		for c := 0; c < len(b[0]); c++ {
			if c == 3 || c == 6 {
				fmt.Fprint(w, " |")
			}
			if c%3 == 0 {
				fmt.Fprintf(w, "%d", b[r][c])
			} else {
				fmt.Fprintf(w, " %d", b[r][c])
			}
		}
		fmt.Fprintln(w)
	}
}

// WriteCNF writes CNF constraints of a given board to stdout.
func WriteCNF(b Board, w io.Writer) {
	PrintBoard(b, "c ", w)
	out := strings.Builder{}
	clauses := 0
	clauses += writePreassigned(b, &out)
	clauses += writeRowConstraints(b, &out)
	clauses += writeColConstraints(b, &out)
	clauses += writeCellConstraints(b, &out)
	clauses += writeBlockConstraints(b, &out)
	clauses += writeSlopVariables(b, &out)
	fmt.Fprintf(w, "p cnf %d %d\n", literalForCell(len(b)-1, len(b[0])-1, 9), clauses)
	fmt.Fprint(w, out.String())
}

func writePreassigned(b Board, out *strings.Builder) int {
	clauses := 0
	for r := 0; r < len(b); r++ {
		for c := 0; c < len(b[0]); c++ {
			if b[r][c] != 0 {
				out.WriteString(fmt.Sprintf("%d 0\n", literalForCell(r, c, b[r][c])))
				clauses++
			}
		}
	}
	return clauses
}

func writeSlopVariables(b Board, out *strings.Builder) int {
	// Some variables we don't actually use (we multiply by 10 and 100 in literalToCell):
	// * v == 0
	// * c == 9
	// May should hardwire them to false to save some time backtracking on variables
	// that are essentially "don't care".
	clauses := 0
	for r := 0; r < len(b); r++ {
		for c := 0; c < len(b[0]); c++ {
			if c == 0 && r == 0 {
				continue
			}
			out.WriteString(fmt.Sprintf("-%d 0\n", literalForCell(r, c, 0)))
			clauses++
		}
	}
	for r := 0; r < len(b)-1; r++ {
		for v := 0; v <= 9; v++ {
			out.WriteString(fmt.Sprintf("-%d 0\n", literalForCell(r, 9, v)))
			clauses++
		}
	}
	return clauses
}

func writeRowConstraints(b Board, out *strings.Builder) int {
	clauses := 0
	for r := 0; r < len(b); r++ {
		for v := 1; v <= 9; v++ {
			// At least one of the cols in the row have 'v' from 1-9
			for c := 0; c < len(b[0]); c++ {
				out.WriteString(fmt.Sprintf("%d ", literalForCell(r, c, v)))
			}
			out.WriteString("0\n")
			clauses++
			// At most one of the cols in the row have 'v' from 1-9
			// Together, exactly one of the cols in the row have 'v'.
			for c := 0; c < len(b[0])-1; c++ {
				for c2 := c + 1; c2 < len(b[0]); c2++ {
					out.WriteString(fmt.Sprintf("-%d -%d 0\n",
						literalForCell(r, c, v),
						literalForCell(r, c2, v)))
					clauses++
				}
			}
		}
	}
	return clauses
}

func writeColConstraints(b Board, out *strings.Builder) int {
	clauses := 0
	for c := 0; c < len(b[0]); c++ {
		for v := 1; v <= 9; v++ {
			for r := 0; r < len(b); r++ {
				out.WriteString(fmt.Sprintf("%d ", literalForCell(r, c, v)))
			}
			out.WriteString("0\n")
			clauses++
			for r := 0; r < len(b)-1; r++ {
				for r2 := r + 1; r2 < len(b); r2++ {
					out.WriteString(fmt.Sprintf("-%d -%d 0\n",
						literalForCell(r, c, v),
						literalForCell(r2, c, v)))
					clauses++
				}
			}
		}
	}
	return clauses
}

func writeCellConstraints(b Board, out *strings.Builder) int {
	clauses := 0
	for r := 0; r < len(b); r++ {
		for c := 0; c < len(b[0]); c++ {
			for v := 1; v <= 9; v++ {
				out.WriteString(fmt.Sprintf("%d ", literalForCell(r, c, v)))
			}
			out.WriteString("0\n")
			clauses++
			for v := 1; v <= 8; v++ {
				for v2 := v + 1; v2 <= 9; v2++ {
					out.WriteString(fmt.Sprintf("-%d -%d 0\n",
						literalForCell(r, c, v),
						literalForCell(r, c, v2)))
					clauses++
				}
			}
		}
	}
	return clauses
}

func writeBlockConstraints(b Board, out *strings.Builder) int {
	clauses := 0
	toBlockRC := func(rb, cb, subx int) (int, int) {
		r := rb*3 + (subx / 3)
		c := cb*3 + (subx % 3)
		return r, c
	}
	for rb := 0; rb < 3; rb++ {
		for cb := 0; cb < 3; cb++ {
			for v := 1; v <= 9; v++ {
				for subx := 0; subx < 9; subx++ {
					r, c := toBlockRC(rb, cb, subx)
					out.WriteString(fmt.Sprintf("%d ", literalForCell(r, c, v)))
				}
				out.WriteString("0\n")
				clauses++
				for subx := 0; subx < 8; subx++ {
					r, c := toBlockRC(rb, cb, subx)
					for sub2x := subx + 1; sub2x < 9; sub2x++ {
						r2, c2 := toBlockRC(rb, cb, sub2x)
						if r2 == r || c2 == c {
							// Same row/col already covered by row/col constraints
							continue
						}
						out.WriteString(fmt.Sprintf("-%d -%d 0\n",
							literalForCell(r, c, v),
							literalForCell(r2, c2, v)))
						clauses++
					}
				}
			}
		}
	}
	return clauses
}

func literalForCell(r int, c int, v int) int {
	return r*100 + c*10 + v
}

func cellValForLit(lit int) (int, int, int) {
	v := lit % 10
	c := (lit / 10) % 10
	r := lit / 100
	return r, c, v
}

// ParseAssignments reads SAT solver variable assignments and returns the solved board.
func ParseAssignments(input io.Reader) Board {
	b := make(Board, 9)
	for r := 0; r < 9; r++ {
		b[r] = make([]int, 9)
	}
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if fields[0] == "v" {
			lit, err := strconv.Atoi(fields[1])
			if err != nil {
				panic(fmt.Sprintf("Failed to parse variable line %v: %v", line, err))
			}
			if lit > 0 {
				r, c, v := cellValForLit(lit)
				if v > 0 {
					b[r][c] = v
				}
			}
		}
	}
	return b
}
