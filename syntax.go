package main

type system struct {
	variables Variables
	locks     Locks
	processes []process
}

func System(vars Variables, locks Locks, procs ...process) system {
	return system{
		variables: vars,
		locks:     locks,
		processes: procs,
	}
}

type varName string
type Variables map[varName]int

type lockName string
type Locks []lockName

type procName string
type process struct {
	name       procName
	statements []statement
}

func Process(name procName, stmts ...statement) process {
	return process{name: name, statements: stmts}
}

type intExpression interface {
	eval(map[varName]int) (int, error)
}

var _ intExpression = intValueExpression{}
var _ intExpression = variableExpression{}
var _ intExpression = addExpression{}
var _ intExpression = subExpression{}

type intValueExpression struct {
	value int
}

func Int(val int) intValueExpression {
	return intValueExpression{value: val}
}

type variableExpression struct {
	name varName
}

func Var(name varName) variableExpression {
	return variableExpression{name: name}
}

type addExpression struct {
	left  intExpression
	right intExpression
}

func Add(l intExpression, r intExpression) addExpression {
	return addExpression{left: l, right: r}
}

type subExpression struct {
	left  intExpression
	right intExpression
}

func Sub(l intExpression, r intExpression) subExpression {
	return subExpression{left: l, right: r}
}

type boolExpression interface {
	eval(map[varName]int) (bool, error)
}

var _ boolExpression = boolValueExpression{}
var _ boolExpression = eqExpression{}
var _ boolExpression = ltExpression{}
var _ boolExpression = notExpression{}
var _ boolExpression = orExpression{}

type boolValueExpression struct {
	value bool
}

func True() boolValueExpression {
	return boolValueExpression{
		value: true,
	}
}

type eqExpression struct {
	left  intExpression
	right intExpression
}

func Eq(l intExpression, r intExpression) eqExpression {
	return eqExpression{
		left:  l,
		right: r,
	}
}

type ltExpression struct {
	left  intExpression
	right intExpression
}

func Lt(l intExpression, r intExpression) ltExpression {
	return ltExpression{
		left:  l,
		right: r,
	}
}

type notExpression struct {
	expression boolExpression
}

func Not(bexpr boolExpression) notExpression {
	return notExpression{
		expression: bexpr,
	}
}

type orExpression struct {
	left  boolExpression
	right boolExpression
}

func Or(l boolExpression, r boolExpression) orExpression {
	return orExpression{
		left:  l,
		right: r,
	}
}

func Neq(l intExpression, r intExpression) notExpression {
	return Not(Eq(l, r))
}

func Le(l intExpression, r intExpression) orExpression {
	return Or(Lt(l, r), Eq(l, r))
}

func Gt(l intExpression, r intExpression) notExpression {
	return Not(Le(l, r))
}

func Ge(l intExpression, r intExpression) notExpression {
	return Not(Lt(l, r))
}

func And(l boolExpression, r boolExpression) notExpression {
	return Not(Or(Not(l), Not(r)))
}

type statement interface {
	execute(env environment, pname procName, cont []statement) ([]localState, error)
}

var _ statement = assignStatement{}
var _ statement = switchStatement{}
var _ statement = forStatement{}

type assignStatement struct {
	varName    varName
	expression intExpression
}

func Assign(x varName, iexpr intExpression) assignStatement {
	return assignStatement{
		varName:    x,
		expression: iexpr,
	}
}

type guard interface {
	execute(env environment, pname procName) (environment, bool, error)
}

var _ guard = whenGuard{}
var _ guard = lockGuard{}

type whenGuard struct {
	expression boolExpression
}

func When(bexpr boolExpression) whenGuard {
	return whenGuard{
		expression: bexpr,
	}
}

type lockGuard struct {
	name lockName
}

func Lock(name lockName) lockGuard {
	return lockGuard{name: name}
}

type guardedCase struct {
	guard      guard
	statements []statement
}

func Case(grd guard, stmts ...statement) guardedCase {
	return guardedCase{
		guard:      grd,
		statements: stmts,
	}
}

type switchStatement struct {
	cases []guardedCase
}

func Switch(cases ...guardedCase) switchStatement {
	return switchStatement{cases: cases}
}

type forStatement struct {
	cases []guardedCase
}

func For(cases ...guardedCase) forStatement {
	return forStatement{
		cases: cases,
	}
}

var _ statement = unlockStatement{}

type unlockStatement struct {
	name lockName
}

func Unlock(name lockName) unlockStatement {
	return unlockStatement{name: name}
}
