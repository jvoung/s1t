package s1t

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

// VarNum is the identifier for a problem variable.
type VarNum uint

// Clause is a set of positive or negative variables in a disjunction.
type Clause struct {
	// TODO(jvoung): could make these bitmaps.
	Positives map[VarNum]bool
	Negatives map[VarNum]bool
}
