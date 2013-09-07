package futures

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrTimeout = errors.New("timed out")
)

type FutureGroup struct {
	to      time.Duration
	futures map[string]*Future
}

func NewFutureGroup(to time.Duration) (fg *FutureGroup) {
	fg = &FutureGroup{
		to:      to,
		futures: make(map[string]*Future),
	}

	return
}

func (fg *FutureGroup) Add(s string, fn func() error) {
	fg.futures[s] = NewFuture(fg.to, fn)
}

func (fg *FutureGroup) Call() error {
	var err error
	var errs []error
	var errstrs []string
	var identical bool

	for s, f := range fg.futures {
		if err = f.Call(); err != nil {
			errs = append(errs, err)
			errstrs = append(errstrs, fmt.Sprintf("%s(%s)", s, err))
		}
	}

	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	}

	identical = true
	for i := 1; i < len(errs); i++ {
		if errs[i] != errs[0] {
			identical = false
			break
		}
	}

	if identical {
		return errs[0]
	}

	return errors.New(fmt.Sprintf("%d errors occured: %s", len(errstrs), strings.Join(errstrs, ", ")))
}

type Future struct {
	to <-chan time.Time
	ch chan error
	fn func() error
}

func NewFuture(to time.Duration, fn func() error) (f *Future) {
	f = &Future{
		to: time.After(to),
		ch: make(chan error),
		fn: fn,
	}

	go func() {
		f.ch <- fn()
	}()

	return f
}

func (f *Future) Call() (err error) {
	for {
		time.Sleep(1 * time.Millisecond)
		select {
		case err = <-f.ch:
			return
		default:
			select {
			case <-f.to:
				err = ErrTimeout
				return
			default:
				continue
			}
		}
	}

	return
}
