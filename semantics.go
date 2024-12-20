package main

import "fmt"

func (expr intValueExpression) eval(_ map[varName]int) (int, error) {
	return expr.value, nil
}

func (expr variableExpression) eval(vars map[varName]int) (int, error) {
	if val, ok := vars[expr.name]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("undefined variable: %s", expr.name)
}

func (expr addExpression) eval(vars map[varName]int) (int, error) {
	l, err := expr.left.eval(vars)
	if err != nil {
		return 0, err
	}
	r, err := expr.right.eval(vars)
	if err != nil {
		return 0, err
	}
	return l + r, nil
}

func (expr subExpression) eval(vars map[varName]int) (int, error) {
	l, err := expr.left.eval(vars)
	if err != nil {
		return 0, err
	}
	r, err := expr.right.eval(vars)
	if err != nil {
		return 0, err
	}
	return l - r, nil
}

func (expr boolValueExpression) eval(_ map[varName]int) (bool, error) {
	return expr.value, nil
}

func (expr eqExpression) eval(vars map[varName]int) (bool, error) {
	l, err := expr.left.eval(vars)
	if err != nil {
		return false, err
	}
	r, err := expr.right.eval(vars)
	if err != nil {
		return false, err
	}
	return l == r, nil
}

func (expr ltExpression) eval(vars map[varName]int) (bool, error) {
	l, err := expr.left.eval(vars)
	if err != nil {
		return false, err
	}
	r, err := expr.right.eval(vars)
	if err != nil {
		return false, err
	}
	return l < r, nil
}

func (expr notExpression) eval(vars map[varName]int) (bool, error) {
	val, err := expr.expression.eval(vars)
	if err != nil {
		return false, err
	}
	return !val, nil
}

func (expr orExpression) eval(vars map[varName]int) (bool, error) {
	l, err := expr.left.eval(vars)
	if err != nil {
		return false, err
	}
	r, err := expr.right.eval(vars)
	if err != nil {
		return false, err
	}
	return l || r, nil
}

type environment struct {
	variables map[varName]int
}

type localState struct {
	environment environment
	statements  []statement
}

func (stmt assignStatement) execute(env environment, pname procName, cont []statement) ([]localState, error) {
	if _, ok := env.variables[stmt.varName]; !ok {
		return []localState{}, fmt.Errorf("undefined variable: %s", stmt.varName)
	}

	vars := map[varName]int{}
	for name, val := range env.variables {
		if name == stmt.varName {
			// Evaluate as rvalue
			newVal, err := stmt.expression.eval(env.variables)
			if err != nil {
				return []localState{}, err
			}
			vars[name] = newVal
		} else {
			vars[name] = val
		}
	}
	state := localState{
		environment: environment{variables: vars},
		statements:  cont,
	}
	return []localState{state}, nil
}

func (grd whenGuard) execute(
	env environment, pname procName) (environment, bool, error) {
	condition, err := grd.expression.eval(env.variables)
	if err != nil {
		return environment{}, false, err
	}
	if !condition {
		return env, false, nil
	}
	return env, true, nil
}

func (stmt switchStatement) execute(env environment, pname procName, cont []statement) ([]localState, error) {
	states := []localState{}
	for _, c := range stmt.cases {
		newEnv, condition, err := c.guard.execute(env, pname)
		if err != nil {
			return []localState{}, err
		}
		if condition {
			stmts := append(c.statements, cont...)
			state := localState{
				environment: newEnv,
				statements:  stmts,
			}
			states = append(states, state) // TODO: write continue
		}
	}
	return states, nil
}

func (stmt forStatement) execute(env environment, pname procName, cont []statement) ([]localState, error) {
	states := []localState{}
	for _, c := range stmt.cases {
		newEnv, condition, err := c.guard.execute(env, pname)
		if err != nil {
			return []localState{}, err
		}
		if condition {
			stmts := append(c.statements, stmt)
			stmts = append(stmts, cont...)
			state := localState{
				environment: newEnv,
				statements:  stmts,
			}
			states = append(states, state)
		}
	}
	return states, nil
}
