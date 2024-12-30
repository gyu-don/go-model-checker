package modelchecker

type ctlFormula interface {
	satisfyingSet(model kripkeModel) worldIDs
}

var _ ctlFormula = ltCTLFormula{}
var _ ctlFormula = notCTLFormula{}
var _ ctlFormula = orCTLFormula{}
var _ ctlFormula = existsNextCTLFormula{}
var _ ctlFormula = existsGloballyCTLFormula{}
var _ ctlFormula = existsUntilCTLFormula{}

type ltCTLFormula struct {
	varName varName
	value   int
}

func LT(name varName, val int) ltCTLFormula {
	return ltCTLFormula{varName: name, value: val}
}

func (f ltCTLFormula) satisfyingSet(model kripkeModel) worldIDs {
	sat := worldIDs{}
	for _, wld := range model.worlds {
		if wld.valuationLT(f.varName, f.value) {
			sat.insert(wld.id)
		}
	}
	return sat
}

type notCTLFormula struct {
	formula ctlFormula
}

func NOT(f ctlFormula) notCTLFormula {
	return notCTLFormula{formula: f}
}

func (f notCTLFormula) satisfyingSet(model kripkeModel) worldIDs {
	sat := worldIDs{}
	unsat := f.formula.satisfyingSet(model)
	for _, wld := range model.worlds {
		if !unsat.member(wld.id) {
			sat.insert(wld.id)
		}
	}
	return sat
}

type orCTLFormula struct {
	left  ctlFormula
	right ctlFormula
}

func OR(l ctlFormula, r ctlFormula) orCTLFormula {
	return orCTLFormula{left: l, right: r}
}

func (f orCTLFormula) satisfyingSet(model kripkeModel) worldIDs {
	sat := worldIDs{}
	for id := range f.left.satisfyingSet(model) {
		sat.insert(id)
	}
	for id := range f.right.satisfyingSet(model) {
		sat.insert(id)
	}
	return sat
}

type existsNextCTLFormula struct {
	formula ctlFormula
}

func EX(f ctlFormula) existsNextCTLFormula {
	return existsNextCTLFormula{formula: f}
}

func (f existsNextCTLFormula) satisfyingSet(model kripkeModel) worldIDs {
	// It returns {w \in model.worlds |
	// 					\exists wld \in model.worlds,
	//							(wld is accessible from w) && wld \in S(f)}.

	target := f.formula.satisfyingSet(model)
	sat := worldIDs{}
	for _, wld := range model.worlds {
		for _, nextID := range model.accessible[wld.id] {
			if target.member(nextID) {
				sat.insert(wld.id)
				break
			}
		}
	}
	return sat
}

type existsGloballyCTLFormula struct {
	formula ctlFormula
}

func EG(f ctlFormula) existsGloballyCTLFormula {
	return existsGloballyCTLFormula{formula: f}
}

type existsUntilCTLFormula struct {
	left  ctlFormula
	right ctlFormula
}

func EU(l ctlFormula, r ctlFormula) existsUntilCTLFormula {
	return existsUntilCTLFormula{left: l, right: r}
}

func (f existsUntilCTLFormula) satisfyingSet(model kripkeModel) worldIDs {
	targetL := f.left.satisfyingSet(model)

	sat := worldIDs{}
	visited := worldIDs{} // visited必要か? 今の実装ではsatと同じなはず
	queue := []worldID{}

	for id := range f.right.satisfyingSet(model) {
		sat.insert(id)
		visited.insert(id)
		queue = append(queue, id)
	}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, prev := range model.reverse[current] {
			if !visited.member(prev) && targetL.member(prev) {
				sat.insert(prev)
				visited.insert(prev)
				queue = append(queue, prev)
			}
		}
	}
	return sat
}

func backtrackOrder(accessible map[worldID][]worldID) []worldID {
	order := []worldID{}
	visited := worldIDs{}
	ordered := worldIDs{}

	for start := range accessible {
		if visited.member(start) {
			continue
		}
		stack := []worldID{start}
		for len(stack) > 0 {
			current := stack[len(stack)-1]
			if !visited.member(current) {
				visited.insert(current)
				for _, next := range accessible[current] {
					if !visited.member(next) {
						stack = append(stack, next)
					}
				}
			} else {
				if !ordered.member(current) {
					order = append(order, current)
					ordered.insert(current)
				}
				stack = stack[:len(stack)-1]
			}
		}
	}
	return order
}

// strongly connected component
func scc(accessible, reverse map[worldID][]worldID) [][]worldID {
	order := backtrackOrder(accessible)
	revOrder := make([]worldID, len(order))
	for i, id := range order {
		revOrder[len(order)-i-1] = id
	}

	compos := [][]worldID{}
	visited := worldIDs{}

	for _, start := range revOrder {
		if visited.member(start) {
			continue
		}
		compo := []worldID{start}
		visited.insert(start)
		stack := []worldID{start}

		for len(stack) > 0 {
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, prev := range reverse[current] {
				if !visited.member(prev) {
					compo = append(compo, prev)
					visited.insert(prev)
					stack = append(stack, prev)
				}
			}
		}
		compos = append(compos, compo)
	}
	return compos
}

func restrict(model kripkeModel, target worldIDs) (map[worldID][]worldID, map[worldID][]worldID) {
	accs := map[worldID][]worldID{}
	for from, tos := range model.accessible {
		if target.member(from) {
			acc := []worldID{}
			for _, to := range tos {
				if target.member(to) {
					acc = append(acc, to)
				}
			}
			accs[from] = acc
		}
	}
	revs := map[worldID][]worldID{}
	for to, froms := range model.reverse {
		rev := []worldID{}
		for _, from := range froms {
			if target.member(from) {
				rev = append(rev, from)
			}
		}
		revs[to] = rev
	}
	return accs, revs
}

func (f existsGloballyCTLFormula) satisfyingSet(model kripkeModel) worldIDs {
	target := f.formula.satisfyingSet(model)

	sat := worldIDs{}
	visited := worldIDs{}
	queue := []worldID{}

	accs, revs := restrict(model, target)
	sccs := scc(accs, revs)
	for _, scc := range sccs {
		if len(scc) >= 2 {
			for _, id := range scc {
				sat.insert(id)
				visited.insert(id)
				queue = append(queue, id)
			}
		}
	}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, prev := range model.reverse[current] {
			if !visited.member(prev) && target.member(prev) {
				sat.insert(prev)
				visited.insert(prev)
				queue = append(queue, prev)
			}
		}
	}
	return sat
}

func (model kripkeModel) VerifyCTL(f ctlFormula) verificationResult {
	return verificationResult{targets: f.satisfyingSet(model)}
}
