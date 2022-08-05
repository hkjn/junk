package bar

type SomeError struct{}

func (se SomeError) Error() string {
	return "zomg all is fire and pain"
}

func (se SomeError) Code() int {
	return 42
}

func Fail() error {
	return SomeError{}
}


