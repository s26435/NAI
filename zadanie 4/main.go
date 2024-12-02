package main

import (
	"zad4/models"
)

//
//wybór numeru setu
//
const num_set int = 3

func main() {
	models.ShowSVM(num_set)
	models.ShowTree(num_set)
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
