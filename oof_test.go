package oof

import (
	"errors"
	"fmt"
	"testing"
)

var SpecialError = errors.New("SpecialError")

func LibCodeNested() error {
	x := 1
	x += 2
	return fmt.Errorf("Special error occurred: %w", SpecialError)
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
	return fmt.Errorf("Special error occurred: %w", SpecialError)
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

	origErr := GetOrigError(err)
	if !errors.Is(origErr, SpecialError) {
		t.Errorf("Expected special error")
		return
	}
}

func TestTracef(t *testing.T) {
	err := ApplicationLevelCaller1()
	if !errors.Is(err, OofErrorInstance) {
		t.Errorf("Expecting an oof error instance")
		return
	}

	origErr := GetOrigError(err)
	if !errors.Is(origErr, SpecialError) {
		t.Errorf("Expected special error")
		return
	}
}
