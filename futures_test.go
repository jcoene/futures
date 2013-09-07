package futures

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestFutureScope(t *testing.T) {
	name := "bob"
	num := 5
	flag := false

	f := NewFuture(2*time.Second, func() error {
		name = "bill"
		num = 8
		flag = true

		return nil
	})

	if err := f.Call(); err != nil {
		t.Error(err)
	}

	if name != "bill" {
		t.Error(fmt.Sprintf("name should be bill, got %s", name))
	}

	if num != 8 {
		t.Error(fmt.Sprintf("num should be 8, got %d", num))
	}

	if flag != true {
		t.Error(fmt.Sprintf("flag should be true, got %v", flag))
	}

	return
}

func TestFutureError(t *testing.T) {
	dumb := errors.New("dumb error")

	f := NewFuture(2*time.Second, func() error {
		return dumb
	})

	if err := f.Call(); err != dumb {
		t.Error(fmt.Sprintf("expected %v, got %v", dumb, err))
	}

	return
}

func TestFutureTimeout(t *testing.T) {
	f := NewFuture(10*time.Millisecond, func() error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	if err := f.Call(); err != ErrTimeout {
		t.Error(fmt.Sprintf("expected %v, got %v", ErrTimeout, err))
	}

	return
}

func TestFutureGroupOneError(t *testing.T) {
	fg := NewFutureGroup(2 * time.Second)
	dumb := errors.New("dumb error")

	fg.Add("error_a", func() error {
		return dumb
	})

	fg.Add("no_error", func() error {
		return nil
	})

	if err := fg.Call(); err != dumb {
		t.Error(fmt.Sprintf("expected %v, got %v", dumb, err))
	}

	return
}

func TestFutureGroupIdenticalErrors(t *testing.T) {
	fg := NewFutureGroup(2 * time.Second)
	dumb := errors.New("dumb error")

	fg.Add("error_a", func() error {
		return dumb
	})

	fg.Add("error_b", func() error {
		return dumb
	})

	fg.Add("error_c", func() error {
		return dumb
	})

	fg.Add("no_error", func() error {
		return nil
	})

	if err := fg.Call(); err != dumb {
		t.Error(fmt.Sprintf("expected %v, got %v", dumb, err))
	}

	return
}

func TestFutureGroupDifferentErrors(t *testing.T) {
	fg := NewFutureGroup(2 * time.Second)
	dumb := errors.New("dumb error")
	dumber := errors.New("dumber error")

	fg.Add("error_a", func() error {
		return dumb
	})

	fg.Add("error_b", func() error {
		return dumber
	})

	fg.Add("no_error", func() error {
		return nil
	})

	err := fg.Call()
	errstr := fmt.Sprintf("%s", err)
	if strings.Index(errstr, "2 errors") != 0 {
		t.Error(fmt.Sprintf("unexpected error: %v", err))
	}

	return
}

func TestFutureGroupTimeout(t *testing.T) {
	fg := NewFutureGroup(10 * time.Millisecond)

	fg.Add("slowpoke", func() error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	fg.Add("its_fine", func() error {
		return nil
	})

	if err := fg.Call(); err != ErrTimeout {
		t.Error(fmt.Sprintf("expected %v, got %v", ErrTimeout, err))
	}

	return
}
