package main

import (
	"gitlab.com/flimzy/ale/errors"

	"foo/bar"
)

func main() {
	sErr := bar.SomeError{}
	// sErr := &bar.SomeError{}
	err := bar.Fail()

	if errors.As(err, &sErr) && sErr.Code() == 42 {
		panic("omg that's a big number")
	}
	panic("phew all is fine")
}
