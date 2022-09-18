package oof

import (
	"errors"
	"fmt"
	"runtime/debug"
)

type OofError struct {
	OrigError error
	stack     []byte
}

var OofErrorInstance = &OofError{}

// Error returns a string representation of the error
func (e *OofError) Error() string {
	return fmt.Sprintf("%v\n%s", e.OrigError, string(e.stack))
}

func (e *OofError) Is(target error) bool {
	_, ok := target.(*OofError)
	return ok
}

// Trace wraps the error in an OofError and captures the stack, along with the original error
func Trace(err error) error {
	switch {
	case errors.Is(err, OofErrorInstance):
		return fmt.Errorf("%w", err)
	}

	// Create a stack trace and attach it
	return &OofError{
		OrigError: err,
		stack:     debug.Stack(),
	}
}

// Tracef wraps the error in an OofError and captures the stack, along with the original error and provides annotation
func Tracef(fmtString string, args ...any) error {
	wrap := true
	var err error
	errIdx := 0
	for i, arg := range args {
		switch t := arg.(type) {
		case error:
			if err != nil {
				err = fmt.Errorf("Tracef: Can only wrap a single error: %w", err)
			}
			errIdx = i
			err = t
		}
	}
	if err == nil {
		err = fmt.Errorf(fmtString, args...)
		wrap = false
	}

	switch {
	case errors.Is(err, OofErrorInstance):
		return fmt.Errorf(fmtString, args...)
	}

	// Create a stack trace and attach it
	oofError := &OofError{
		OrigError: err,
		stack:     debug.Stack(),
	}

	if !wrap {
		return oofError
	}

	newArgs := make([]any, len(args))
	for i, arg := range args {
		if i == errIdx {
			newArgs[i] = oofError
		} else {
			newArgs[i] = arg
		}
	}
	return fmt.Errorf(fmtString, newArgs...)
}

// GetOrigError returns the original error under an OofError (if err is not an OofError, err will be returned)
func GetOrigError(err error) error {
	var target *OofError
	result := errors.As(err, &target)
	if !result {
		return err
	}

	return target.OrigError
}
