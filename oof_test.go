package oof

import (
	"errors"
	"testing"
	"time"
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

func TestTotalOofs(t *testing.T) {
	totalOofs.Store(0)
	loops := uint64(100)
	err := errors.New("Something")
	for i := uint64(0); i < loops; i++ {
		go func() {
			Trace(err)
		}()
	}
	wt := time.Now().Add(time.Second * 5)
	for GetTotalOofs() < loops || time.Until(wt) <= time.Duration(0) {
		time.Sleep(time.Millisecond)
	}
	if GetTotalOofs() != loops {
		t.Error("Wrong number of oofs")
		return
	}
	localOof := Trace(err)
	for i := uint64(0); i < loops; i++ {
		go func() {
			Trace(localOof)
		}()
	}
	wt = time.Now().Add(time.Second * 5)
	for GetTotalOofs() < loops+1 || time.Until(wt) <= time.Duration(0) {
		time.Sleep(time.Millisecond)
	}
	if GetTotalOofs() != loops+1 {
		t.Error("Wrong number of oofs")
		return
	}
	for i := uint64(0); i < loops; i++ {
		q := i
		go func() {
			Tracef("Something %d", q)
		}()
	}
	wt = time.Now().Add(time.Second * 5)
	for GetTotalOofs() < (loops*2)+1 || time.Until(wt) <= time.Duration(0) {
		time.Sleep(time.Millisecond)
	}
	if GetTotalOofs() != (loops*2)+1 {
		t.Error("Wrong number of oofs")
		return
	}
}
