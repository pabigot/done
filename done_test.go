// Copyright 2022 Peter Bigot Consulting, LLC
// SPDX-License-Identifier: Apache-2.0

package done

import (
	"errors"
	"testing"
)

func TestImplementation(t *testing.T) {
	di := Implementation{}
	if err := di.Err(); err != nil {
		t.Fatal(err.Error())
	}
	dc := di.Done()
	if dc == nil {
		t.Fatal("no Done channel")
	}
	select {
	case <-dc:
		t.Fatal("premature done")
	default:
	}

	errBad := errors.New("bad exit")

	di.Finalize(nil)
	select {
	case <-dc:
	default:
		t.Fatal("missing done")
	}
	if err := di.Err(); err != TerminatedOK {
		t.Fatalf("unexpected: %v", err)
	}

	di.Finalize(errBad)
	if err := di.Err(); err != TerminatedOK {
		t.Fatalf("unexpected: %v", err)
	}

	di = Implementation{}
	di.Finalize(errBad)
	if err := di.Err(); err != errBad {
		t.Fatalf("unexpected: %v", err)
	}
	select {
	case <-di.Done():
	default:
		t.Fatal("missing done")
	}
}
