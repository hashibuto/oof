package oof

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"sync/atomic"
)

type OofError struct {
	OrigError error
	stack     []byte
}

var totalOofs = atomic.Uint64{}

var OofErrorInstance = &OofError{}

// GetTotalOofs returns the total number of times oof.Trace has been called
func GetTotalOofs() uint64 {
	return totalOofs.Load()
}

// Error returns a string representation of the error
func (e *OofError) Error() string {
	if len(e.stack) == 0 {
		return fmt.Sprintf("%v", e.OrigError)
	}
	return fmt.Sprintf("%v\n%s", e.OrigError, string(e.stack))
}

func (e *OofError) Is(target error) bool {
	_, ok := target.(*OofError)
	return ok
}

// Unwrap will recursively unwrap the error, returning the original error
func (e *OofError) Unwrap() error {
	return e.OrigError
}

// StripTrace strips the stack trace from an error
func (e *OofError) StripTrace() {
	e.stack = []byte{}
}

// Trace wraps the error in an OofError and captures the stack, along with the original error
// If the supplied error is nil, Trace will return nil
func Trace(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, OofErrorInstance):
		return fmt.Errorf("%w", err)
	}

	// Add to total oofs
	totalOofs.Add(1)
	// Create a stack trace and attach it
	return &OofError{
		OrigError: err,
		stack:     debug.Stack(),
	}
}

// Tracef wraps the error in an OofError and captures the stack, along with the original error and provides annotation
// If the supplied error is nil, Tracef will return nil
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
	// Add to total oofs
	totalOofs.Add(1)
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

// Fatal will log a fatal error message and cause the process to exit if and only if the provided err is non-nil.  A full strack trace will be included.
func Fatal(err error) {
	if err == nil {
		return
	}

	log.Fatalf("%s", Trace(err).Error())
}

// Fatalf will log a fatal error message and cause the process to exit if and only if the provided err is non-nil.  Additionally, a user specified message will be logged.  A full strack trace will be included.
func Fatalf(fmtString string, args ...any) {
	var err error
	for _, arg := range args {
		switch t := arg.(type) {
		case error:
			if err != nil {
				err = fmt.Errorf("Tracef: Can only wrap a single error: %w", err)
			}
			err = t
		}
	}
	if err == nil {
		return
	}

	log.Fatalf("%s", Tracef(fmtString, args...).Error())
}

// StripTrace strips the stack trace from an error
func StripTrace(err error) error {
	var e *OofError
	if errors.As(err, &e) {
		e.StripTrace()
		return e
	}

	return err
}
