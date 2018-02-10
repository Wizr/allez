package errorf

import "fmt"

// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}
func Newf(format string, a ...interface{}) error {
	text := fmt.Sprintf(format, a)
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
