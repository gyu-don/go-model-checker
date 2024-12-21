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
	model.WriteAsDot(os.Stdout)
}
