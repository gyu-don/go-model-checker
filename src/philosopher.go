package modelchecker

func badPhilosopher(name procName, right, left lockName, hold varName) process {
	return Process(name,
		For(
			Case(Lock(left), Assign(hold, Add(Var(hold), Int(1)))),
			Case(Lock(right), Assign(hold, Add(Var(hold), Int(1)))),
			Case(When(Eq(Var(hold), Int(2))),
				// eating
				Assign(hold, Int(0)),
				Unlock(left),
				Unlock(right),
			),
		),
	)
}

func goodPhilosopher(name procName, first, second lockName, hold varName) process {
	return Process(name,
		For(
			Case(Lock(first),
				Assign(hold, Add(Var(hold), Int(1))),
				Switch(
					Case(Lock(second),
						Assign(hold, Add(Var(hold), Int(1))),
						// eating
						Assign(hold, Int(0)),
						Unlock(second),
						Unlock(first),
					),
				),
			),
		),
	)
}

func DiningBadPhilosophers() system {
	return System(
		Variables{
			"hold1": 0, "hold2": 0,
		},
		Locks{"fork1", "fork2"},
		badPhilosopher("P1", "fork1", "fork2", "hold1"),
		badPhilosopher("P2", "fork2", "fork1", "hold2"),
	)
}

func DiningGoodPhilosophers() system {
	return System(
		Variables{
			"hold1": 0, "hold2": 0,
		},
		Locks{"fork1", "fork2"},
		goodPhilosopher("P1", "fork1", "fork2", "hold1"),
		goodPhilosopher("P2", "fork1", "fork2", "hold2"),
	)
}
