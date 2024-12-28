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
