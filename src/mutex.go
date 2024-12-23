package modelchecker

func badThread(name procName) process {
	return Process(name,
		For(
			Case(When(True()),
				Assign("critical", Add(Var("critical"), Int(1))),
				// critical session
				Assign("critical", Sub(Var("critical"), Int(1))),
			),
		),
	)
}

func goodThread(name procName) process {
	return Process(name,
		For(
			Case(Lock("mutex"),
				Assign("critical", Add(Var("critical"), Int(1))),
				// critical session
				Assign("critical", Sub(Var("critical"), Int(1))),
				Unlock("mutex"),
			),
		),
	)
}

func BadMutex() system {
	return System(
		Variables{"critical": 0},
		Locks{},
		badThread("A"),
		badThread("B"),
		badThread("C"),
	)
}

func GoodMutex() system {
	return System(
		Variables{"critical": 0},
		Locks{"mutex"},
		goodThread("A"),
		goodThread("B"),
		goodThread("C"),
	)
}
