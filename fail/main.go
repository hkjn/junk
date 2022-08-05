package main

import (
	"fmt"

	"gitlab.com/flimzy/ale/errors"

	"foo/bar"
)

func main() {
	sErr := bar.SomeError{}
	// sErr := &bar.SomeError{}
	err := bar.Fail()

	if errors.As(err, &sErr) && sErr.Code() == 42 {
		panic(fmt.Sprintf("omg that's a big number: %v", sErr))
	}
	panic("phew all is fine")
}
