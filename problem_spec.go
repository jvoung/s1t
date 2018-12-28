package s1t

import (
	"strconv"
)

// Problem is a specification of a problem instance.
type Problem struct {
	Spec    ProblemSpec
	Clauses []Clause
}

// ProblemSpec represents the shape of the input problem.
type ProblemSpec struct {
	Format       string // TODO(jvoung): handle "sat" format, not just "cnf"
	NumVariables int
	NumClauses   int
}

// ClauseNum is the index of a clause in the Problem's clause list (0 - NumClauses)
type ClauseNum uint

// VarNum is the identifier for a problem variable (0 - NumVariables from the ProblemSpec).
type VarNum uint

// Literal is a var or a negated var
type Literal int

// Positive returns a Literal for a VarNum atom
func Positive(v VarNum) Literal {
	return Literal((v << 1) + 1)
}

// Negative returns a negated variable Literal
func Negative(v VarNum) Literal {
	return Literal(v << 1)
}

// Var returns the variable number from a Literal (assuming we don't handle constant true/false).
func (l Literal) Var() VarNum {
	return VarNum(l >> 1)
}

// AsInt returns converts the literal to 0 or 1 if negated, or not.
func (l Literal) AsInt() int {
	return int(l) & 1
}

// Negate returns another literal that is the negation of this.
func (l Literal) Negate() Literal {
	return Literal(l ^ 1)
}

func (l Literal) String() string {
	if l.AsInt() == 0 {
		return "¬" + l.Var().String()
	}
	return l.Var().String()
}

func (v VarNum) String() string {
	return "v" + strconv.Itoa(int(v))
}

// Clause is a collection of Literals in a disjunction.
type Clause struct {
	Literals []Literal
}

// Empty returns true if this is an empty Clause.
func (c *Clause) Empty() bool {
	return len(c.Literals) == 0
}
