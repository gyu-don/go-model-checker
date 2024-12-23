package main

import (
	"log"
	"os"

	"modelchecker/src"
)

func main() {
	sys := modelchecker.DiningGoodPhilosophers()
	model, err := modelchecker.KripkeModel(sys)
	if err != nil {
		log.Fatal(err)
	}
	//result := model.VerifyInvariantLT("critical", 2)
	result := model.VerifyDeadlockFreedom()
	model.WriteAsDot(os.Stdout, &result)
}
