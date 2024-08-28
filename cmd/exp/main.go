package main

import (
	"errors"
	"fmt"
)

func main() {
	err := B()

	var customErr *CustomError
	if errors.As(err, &customErr) {
		fmt.Printf("Error is of type *CustomError: %v\n", customErr)
	} else {
		fmt.Println("Error is not of type *CustomError")
	}

}

// It is common for packages like database/sql to return
// an error that is predefined like this one.
var ErrNotFound = errors.New("not found")
var ErrMadiLoh = errors.New("madi loh")

func A() error {
	return &CustomError{msg: "wrapped error: not found"}
}

func B() error {
	err := A()
	return fmt.Errorf("b: %w", err)
}

type CustomError struct {
	msg string
}

func (e *CustomError) Error() string {
	return e.msg
}
