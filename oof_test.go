package oof

import (
	"errors"
	"testing"
)

type CustomError struct{}

func (ce *CustomError) Error() string {
	return "It's a custom error"
}

func LibCodeNested() error {
	x := 1
	x += 2
	return &CustomError{}
}

func ApplicationLevelCaller() error {
	err := LibCodeNested()
	if err != nil {
		return Trace(err)
	}

	return nil
}

func LibCodeNested1() error {
	x := 1
	x += 2
	return &CustomError{}
}

func ApplicationLevelCaller1() error {
	err := LibCodeNested()
	if err != nil {
		return Tracef("Hello, this is my error: %w", err)
	}

	return nil
}

func TestTrace(t *testing.T) {
	err := ApplicationLevelCaller()
	if !errors.Is(err, OofErrorInstance) {
		t.Errorf("Expecting an oof error instance")
		return
	}

	if !errors.Is(err, &CustomError{}) {
		t.Errorf("Expected custom error")
		return
	}
}

func TestTracef(t *testing.T) {
	err := ApplicationLevelCaller1()
	if !errors.Is(err, OofErrorInstance) {
		t.Errorf("Expecting an oof error instance")
		return
	}

	if !errors.Is(err, &CustomError{}) {
		t.Errorf("Expected custom error")
		return
	}
}
