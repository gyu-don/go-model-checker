package main

import (
	"log"
	"os"

	//lint:ignore ST1001 for DSL
	. "modelchecker/src"
)

func main() {
	sys := DiningGoodPhilosophers()
	model, err := KripkeModel(sys)
	if err != nil {
		log.Fatal(err)
	}
	// result := model.VerifyInvariantLT("critical", 2)
	// result := model.VerifyDeadlockFreedom()
	f := EG(LT("hold2", 2))
	result := model.VerifyCTL(f)
	model.WriteAsDot(os.Stdout, &result)
}
