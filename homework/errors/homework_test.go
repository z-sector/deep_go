package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// go test -v homework_test.go

type MultiError struct {
	errs []error
}

func (e *MultiError) Error() string {
	if e == nil || len(e.errs) == 0 {
		return "<nil>"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d errors occured:\n\t", len(e.errs)))
	for i, err := range e.errs {
		if i > 0 {
			sb.WriteString("\t")
		}
		sb.WriteString("* " + err.Error())
	}
	sb.WriteString("\n")

	return sb.String()
}

func (e *MultiError) Is(target error) bool {
	for _, err := range e.errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (e *MultiError) As(target any) bool {
	for _, err := range e.errs {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

func (e *MultiError) Unwrap() error {
	if e == nil || len(e.errs) == 0 {
		return nil
	}
	return e.errs[0]
}

func Append(err error, errs ...error) *MultiError {
	if err == nil && len(errs) == 0 {
		return nil
	}

	me := &MultiError{}
	if err != nil {
		var multiErr *MultiError
		if errors.As(err, &multiErr) {
			me = multiErr
		} else {
			me = &MultiError{errs: []error{err}}
		}
	}

	for _, e := range errs {
		if e != nil {
			me.errs = append(me.errs, e)
		}
	}

	return me
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestMultiError_Is(t *testing.T) {
	baseErr := errors.New("base error")
	err := Append(baseErr, errors.New("other error"))
	assert.True(t, errors.Is(err, baseErr), "errors.Is should find baseErr inside MultiError")
}

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func TestMultiError_As(t *testing.T) {
	custErr := &customError{msg: "custom"}
	err := Append(custErr, errors.New("something else"))

	var target *customError
	ok := errors.As(err, &target)

	require.True(t, ok, "errors.As should find customError inside MultiError")
	assert.Equal(t, "custom", target.msg)
}

func TestMultiError_Unwrap(t *testing.T) {
	err1 := errors.New("first")
	err2 := errors.New("second")
	me := Append(err1, err2)

	unwrapped := errors.Unwrap(me)
	require.NotNil(t, unwrapped)
	assert.Equal(t, err1, unwrapped)
}
