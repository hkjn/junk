package bar

type SomeError struct{
	code int
}

func (se SomeError) Error() string {
	return "zomg all is fire and pain"
}

func (se SomeError) GetCode() int {
	return se.code
}

func Fail() error {
	return SomeError{code: 44}
}


