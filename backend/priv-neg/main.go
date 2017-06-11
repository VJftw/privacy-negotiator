package main

import (
	"fmt"
	"os"
)

func main() {
	if os.Getenv("TYPE") == "API" {
		NewPrivNegAPI()
	} else if os.Getenv("TYPE") == "WORKER" {
		NewPrivNegWorker(os.Getenv("QUEUE"))
	} else {
		panic(fmt.Sprintf("Invalid backend type: %s. Must be api|worker", os.Getenv("TYPE")))
	}
}
