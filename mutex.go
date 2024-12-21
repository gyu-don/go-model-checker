package main

func badThread(name procName) process {
	/*
		process {name}
		for {
		True:
			critical = critical + 1
			critical = critical - 1
		}
	*/
	// write above with our syntax
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

func Mutex() system {
	return System(
		Variables{"critical": 0},
		badThread("A"),
		badThread("B"),
		badThread("C"),
	)
}
