package helper

import (
	"log"
)

// LogF logs fatal "msg: err" in case of error
func LogF(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// LogE logs error "msg: err" in case of error
func LogE(msg string, err error) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
