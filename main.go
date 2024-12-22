package main

import (
	"log"
	"os"
)

func main() {
	sys := DiningGoodPhilosophers()
	model, err := KripkeModel(sys)
	if err != nil {
		log.Fatal(err)
	}
	//result := model.VerifyInvariantLT("critical", 2)
	result := model.VerifyDeadlockFreedom()
	model.WriteAsDot(os.Stdout, &result)
}
