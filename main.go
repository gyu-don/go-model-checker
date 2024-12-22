package main

import (
	"log"
	"os"
)

func main() {
	sys := Mutex()
	model, err := KripkeModel(sys)
	if err != nil {
		log.Fatal(err)
	}
	result := model.VerifyInvariantLT("critical", 2)
	model.WriteAsDot(os.Stdout, &result)
}
