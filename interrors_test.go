// Copyright 2017 Dmitry Frank <mail@dmitryfrank.com>
// Licensed under the BSD, see LICENSE file for details.

package interrors

import (
	"fmt"
	"testing"

	"github.com/juju/errors"
)

func TestInternalError(t *testing.T) {
	errOrig := errors.Errorf("some internal error: %s", "foo")

	err := errOrig
	err = errors.Annotatef(err, "annotation 1")
	err = errors.Annotatef(err, "annotation 2")

	errPub := errors.Errorf("my public error: %s", "bar")
	errPubWrap := WrapInternalError(err, errPub)

	err2 := errors.Annotatef(errPubWrap, "pub annotation 1")
	err2 = errors.Annotatef(err2, "pub annotation 2")

	for _, v := range []struct {
		got, want interface{}
		descr     string
	}{
		{errPubWrap.Error(), "my public error: bar", fmt.Sprintf("errPubWrap.Error()")},
		{errors.Cause(err2), errPub, fmt.Sprintf("errors.Cause(%v)", err2)},
		{InternalCause(err2), errOrig, fmt.Sprintf("InternalCause(%v)", err2)},
		{InternalCause(err), errOrig, fmt.Sprintf("InternalCause(%v)", err)},
		{InternalErr(err2), err, fmt.Sprintf("InternalErr(%v)", err2)},
		{InternalErr(err), err, fmt.Sprintf("InternalErr(%v)", err)},
		{IsInternalError(err2), true, fmt.Sprintf("IsInternalError(%v)", err2)},
		{IsInternalError(errPubWrap), true, fmt.Sprintf("IsInternalError(%v)", errPubWrap)},
		{IsInternalError(err), false, fmt.Sprintf("IsInternalError(%v)", err)},
	} {
		if v.got != v.want {
			t.Errorf("%s: want: %v, got: %v", v.descr, v.want, v.got)
		}
	}
}

func TestDoubleInternal(t *testing.T) {
	errOrig := errors.Errorf("some internal error: %s", "foo")
	err := errors.Annotatef(errOrig, "annotation 1")
	err = errors.Annotatef(err, "annotation 2")

	errPub := errors.Errorf("my public error: %s", "bar")
	errPubWrap := WrapInternalError(err, errPub)
	err2 := errors.Annotatef(errPubWrap, "pub annotation 1")
	err2 = errors.Annotatef(err2, "pub annotation 2")

	errPub2 := WrapInternalErrorf(err2, "my public error2: %s", "baz")
	err3 := errors.Annotatef(errPub2, "pub2 annotation 1")
	err3 = errors.Annotatef(err3, "pub2 annotation 2")

	if want, got := errPub, errors.Cause(err3); want != got {
		t.Errorf("errors.Cause(%v): want: %q, got %q", err3, want, got)
	}

	if want, got := err, InternalErr(err3); want != got {
		t.Errorf("InternalErr(%v): want: %q, got %q", err3, want, got)
	}
}
