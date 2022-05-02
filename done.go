// Copyright 2022 Peter Bigot Consulting, LLC
// SPDX-License-Identifier: Apache-2.0

// Package done provides an concurrency-safe implementation of the Done() and
// Err() methods as used in context.Context to communication completion status
// to an application.
package done

import (
	"errors"
	"sync"
	"sync/atomic"
)

// Interface provides functions used by clients to determine whether some
// activity has completed and the cause of the termination.  It is inspired by
// that functionality of standard library context.Context.
type Interface interface {
	// Done() returns a channel that is closed when some activity is
	// completed.  The function is safe for concurrent use.
	Done() <-chan struct{}

	// Err() returns nil if the companion Done is not yet closed.
	// Otherwise it returns a non-nil value associated with the close
	// event.  Normal (error-free) termination is indicated by
	// TerminatedOK.
	//
	// The function is safe for concurrent use.
	Err() error
}

// TerminatedOK is returned by ErrDone.Err() when the activity has completed
// without error.
var TerminatedOK = errors.New("ok")

var closedChan = func() chan struct{} {
	c := make(chan struct{})
	close(c)
	return c
}()

// Implementation provides an implementation of Interface along with a
// Finalize function that is used to record information about completion.
type Implementation struct {
	mu   sync.Mutex
	done atomic.Value
	err  error
}

// Finalize records the disposition of some activity.  If err is nil, it is
// replaced by ErrDoneOK.  err is recorded so it will subsequently be provided
// by Err(), and the Done channel is closed.
//
// This function is safe for concurrent use.  Only the first invocation to
// proceed will have any effect.
func (di *Implementation) Finalize(err error) {
	if err == nil {
		err = TerminatedOK
	}
	di.mu.Lock()
	defer di.mu.Unlock()
	if di.err != nil { // already finalized
		return
	}
	di.err = err
	d, _ := di.done.Load().(chan struct{})
	if d == nil {
		di.done.Store(closedChan)
	} else {
		close(d)
	}
}

// Err returns nil if the companion Done channel has not been closed.
// Otherwise it returns an error instance describing the cause of completion,
// with ErrDoneOK used if the termination was not associated with an error.
//
// A nil return may be obsolete by the time the caller inspects it.  A non-nil
// return value will never change.
//
// This function is safe for concurrent use.
func (di *Implementation) Err() error {
	di.mu.Lock()
	err := di.err
	di.mu.Unlock()
	return err
}

// Done returns a channel that is closed when some activity completes.
//
// This function is safe for concurrent use.
func (di *Implementation) Done() <-chan struct{} {
	dv := di.done.Load()
	if dv != nil {
		return dv.(chan struct{})
	}
	di.mu.Lock()
	defer di.mu.Unlock()
	dv = di.done.Load()
	if dv == nil {
		dv = make(chan struct{})
		di.done.Store(dv)
	}
	return dv.(chan struct{})
}
