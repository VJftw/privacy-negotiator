package utils

import (
	"fmt"
	"log"
)

// FailOnError - Exits the application if the given error is non-nil.
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
