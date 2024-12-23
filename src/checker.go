package modelchecker

type verificationResult struct {
	targets worldIDs
}

type worldIDs map[worldID]struct{}

func (ids worldIDs) member(id worldID) bool {
	_, ok := ids[id]
	return ok
}

func (ids worldIDs) insert(id worldID) {
	ids[id] = struct{}{}
}

func (model kripkeModel) VerifyInvariantLT(name varName, val int) verificationResult {
	init := model.initial

	visited := worldIDs{}
	visited.insert(init)
	violated := worldIDs{}

	stack := []worldID{init}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		wld := model.worlds[current]
		if !wld.valuationLT(name, val) {
			violated.insert(current)
		}
		nexts, ok := model.accessible[current]
		if ok {
			for _, next := range nexts {
				if !visited.member(next) {
					visited.insert(next)
					stack = append(stack, next)
				}
			}
		}
	}
	return verificationResult{targets: violated}
}

func (model kripkeModel) VerifyDeadlockFreedom() verificationResult {
	init := model.initial

	visited := worldIDs{}
	visited.insert(init)
	violated := worldIDs{}

	stack := []worldID{init}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		wld := model.worlds[current]
		nexts, ok := model.accessible[current]
		if ok {
			if len(nexts) == 0 {
				for _, stmts := range wld.programCounters {
					if len(stmts) > 0 {
						violated.insert(current)
					}
				}
			} else {
				for _, next := range nexts {
					if !visited.member(next) {
						visited.insert(next)
						stack = append(stack, next)
					}
				}
			}
		}
	}
	return verificationResult{targets: violated}
}
