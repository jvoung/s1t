// Basic Solver using Watched Literals.

package s1t

const (
	none = -1
)

// Solve determines if a given problem is unsat or sat (with an assignment).
func Solve(problem Problem) Solution {
	clauses := problem.Clauses
	if hasEmptyClauses(clauses) {
		return unsat()
	}
	assignments := initialAssignments(problem.Spec.NumVariables)
	watchedLiterals := pickWatchedLiterals(clauses)
	if !initialUnitPropagate(clauses, assignments, watchedLiterals) {
		return unsat()
	}
	if searchSolution(clauses, assignments, watchedLiterals, 0, 0) {
		return sat(assignments)
	}
	return unsat()
}

func searchSolution(clauses []Clause, assignments []int, wls watchedLiterals, depth int, searchFrom VarNum) bool {
	varNum, hasUnassigned := nextUnassignedVariable(assignments, int(searchFrom))
	if !hasUnassigned {
		return true
	}
	assignedAtDepth := []VarNum{}
	doRollbacks := func() {
		for _, v := range assignedAtDepth {
			assignments[v] = none
		}
	}
	if tryAssign(clauses, varNum, Positive(varNum), assignments, wls, &assignedAtDepth) {
		if searchSolution(clauses, assignments, wls, depth+1, varNum) {
			return true
		}
	}
	doRollbacks()

	assignedAtDepth = []VarNum{}
	if tryAssign(clauses, varNum, Negative(varNum), assignments, wls, &assignedAtDepth) {
		if searchSolution(clauses, assignments, wls, depth+1, varNum) {
			return true
		}
	}
	doRollbacks()
	return false
}

func nextUnassignedVariable(assignments []int, searchFrom int) (VarNum, bool) {
	for i := searchFrom; i < len(assignments); i++ {
		if assignments[i] == none {
			return VarNum(i), true
		}
	}
	return 0, false
}

// Try assigning v to satisfy literal l and return true if successful (false if contradicted).
//
// Updates assignments, wls and assigned in-out arguments.
// May trigger a chain of unit clause propagation whose assignments
// are tracked by assigned.
func tryAssign(clauses []Clause, v VarNum, l Literal,
	assignments []int, wls watchedLiterals, assigned *[]VarNum) bool {
	if assignments[v] == none {
		assignments[v] = l.AsInt()
		*assigned = append(*assigned, v)
	} else {
		if assignments[v] != l.AsInt() {
			return false
		}
	}
	// l is fine, since it becomes true and satisfies watched literal invariants.
	// ¬l will become false, so need to watch a different literal.
	negatedL := l.Negate()
	affectedClauses := wls.literalToClause[negatedL]
	for i := 0; i < len(affectedClauses); {
		cnum := affectedClauses[i]
		watchedForC := wls.clauseToLiteral[cnum]
		newLit := findNewWatchedLiteral(clauses, assignments, watchedForC, cnum, negatedL)
		if newLit == none {
			// Only the other watched literal is available, in which case
			// we'll violate the invariant and just keep the watch at the same spot.
			// If we ever backtrack, the invariant will be restored.
			// If the other watched literal
			// - is true, the clause is already satisfied
			// - is unassigned, then we've found a unit clause so attempt to
			//   recursively do unit propagation
			otherWatchedLit := watchedForC.otherWatched(negatedL)
			otherWatchedV := otherWatchedLit.Var()
			otherA := assignments[otherWatchedV]
			if otherA == none {
				// it's unit clause -- try to propagate more
				if !tryAssign(clauses, otherWatchedV, otherWatchedLit, assignments, wls, assigned) {
					return false
				}
			} else {
				if otherA != otherWatchedLit.AsInt() {
					return false
				}
			}
			i++
		} else {
			// switch watches to keep invariant
			affectedClauses[i] = affectedClauses[len(affectedClauses)-1]
			affectedClauses = affectedClauses[:len(affectedClauses)-1]
			wls.literalToClause[negatedL] = affectedClauses

			wls.literalToClause[newLit] = append(wls.literalToClause[newLit], cnum)
			watchedForC.replaceOne(negatedL, newLit)
		}
	}
	return true
}

func findNewWatchedLiteral(
	clauses []Clause, assignments []int, watchedForC *twoWatchedLiterals,
	cnum ClauseNum, negatedL Literal) Literal {
	// Either:
	// - find a new true literal
	// - find a new unassigned literal
	// - no literal that matches invariants -- should check other watched lit
	otherWatchedLit := watchedForC.otherWatched(negatedL)
	c := clauses[cnum]
	for _, candidate := range c.Literals {
		if candidate == negatedL {
			continue
		}
		if candidate == otherWatchedLit {
			continue
		}
		v := candidate.Var()
		a := assignments[v]
		if a == none {
			return candidate
		}
		if a == candidate.AsInt() {
			return candidate
		}
	}
	return none
}

// Initialize assignments to "none"
func initialAssignments(numVars int) []int {
	assignments := make([]int, numVars)
	for i := 0; i < numVars; i++ {
		assignments[i] = none
	}
	return assignments
}

type twoWatchedLiterals struct {
	first  Literal
	second Literal
}

type watchedLiterals struct {
	literalToClause map[Literal][]ClauseNum
	clauseToLiteral []*twoWatchedLiterals
}

func (wl *twoWatchedLiterals) otherWatched(l Literal) Literal {
	if l == wl.first {
		return wl.second
	}
	return wl.first
}

func (wl *twoWatchedLiterals) replaceOne(l Literal, newL Literal) {
	if l == wl.first {
		wl.first = newL
	} else {
		wl.second = newL
	}
}

// Need to initialize Watched Literals, two per clause if not unit clauses
func pickWatchedLiterals(clauses []Clause) watchedLiterals {
	l2c := make(map[Literal][]ClauseNum)
	c2l := make([]*twoWatchedLiterals, len(clauses))
	for i, clause := range clauses {
		cnum := ClauseNum(i)
		if len(clause.Literals) < 2 {
			continue
		}
		l1 := clause.Literals[0]
		l2 := clause.Literals[1]
		l2c[l1] = append(l2c[l1], cnum)
		l2c[l2] = append(l2c[l2], cnum)
		c2l[i] = &twoWatchedLiterals{l1, l2}
	}
	return watchedLiterals{l2c, c2l}
}

// Does initial unit clause propagation and returns true if formula is not falsified.
// Since pickWatchLiterals skipped unit clauses, need to flush out the initial unit clauses.
func initialUnitPropagate(clauses []Clause, assignments []int,
	watchedLiterals watchedLiterals) bool {
	for _, clause := range clauses {
		if len(clause.Literals) == 1 {
			l := clause.Literals[0]
			v := l.Var()
			assigned := []VarNum{}
			if !tryAssign(clauses, v, l, assignments, watchedLiterals, &assigned) {
				return false
			}
		}
	}
	return true
}

func hasEmptyClauses(clauses []Clause) bool {
	for _, clause := range clauses {
		if clause.Empty() {
			return true
		}
	}
	return false
}
