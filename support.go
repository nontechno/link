// Copyright 2024 The NonTechno Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package link

import (
	"fmt"
	"os"
	"reflect"
	"sync/atomic"
)

func callRemote(remote string, f interface{}, format string, args ...interface{}) bool {

	// a simple registryGuard against endless recursion ...
	current := remoteCounter.Add()
	defer remoteCounter.Remove()
	if current > maxDepth {
		// this is too deep - most likely an endless recursion
		return false
	}

	if target, found := getRegistered(remote, getSignature(f)); found && target != nil {

		switch operation := target.(type) {
		case func(string, ...interface{}):
			operation(format, args...)
			return true
		default:
		}
	}
	return false
}

// default (non-overwritten) version of "onReport" function
func localOnReport(format string, args ...interface{}) {
	callRemote(ReportFunc, localOnReport, format, args...)
}

// default (non-overwritten) version of "onWarning" function
func localOnWarning(format string, args ...interface{}) {
	callRemote(WarningFunc, localOnWarning, format, args...)
}

// default (non-overwritten) version of "onError" function
func localOnError(format string, args ...interface{}) {
	if !callRemote(ErrorFunc, localOnError, format, args...) {
		// there was no overload - let's 'report' something before quitting
		fmt.Fprintf(os.Stderr, "Error: "+format+"\nQuitting...", args...)
	}
	// force exit, since this is an error
	onTerminate()
}

func localOnTerminate() {

	// a simple registryGuard against endless recursion ...
	current := remoteCounter.Add()
	defer remoteCounter.Remove()
	if current <= maxDepth {

		if target, found := getRegistered(TerminateFunc, getSignature(localOnTerminate)); found && target != nil {

			switch operation := target.(type) {
			case func():
				operation()
			default:
			}
		}
	}

	// no over-write was found, just quit
	// force exit, since this is an error
	os.Exit(onErrorExitCode)
}

// returns 'true' (if `what` is a pointer) and what it points to
func isPointer(what interface{}) (bool, reflect.Type) {
	if what != nil {
		var t reflect.Type
		if tt, okay := what.(reflect.Type); okay {
			t = tt
		} else {
			t = reflect.TypeOf(what)
		}

		if t.Kind() == reflect.Pointer {
			return true, t.Elem()
		}
	}
	return false, nil
}

func isFunction(what interface{}) (bool, int) {
	if what != nil {
		var t reflect.Type
		if tt, okay := what.(reflect.Type); okay {
			t = tt
		} else {
			t = reflect.TypeOf(what)
		}

		if t.Kind() == reflect.Func {
			return true, t.NumOut()
		}
	}
	return false, 0
}

func getSignature(what interface{}) string {
	if what != nil {
		var t reflect.Type
		if tt, okay := what.(reflect.Type); okay {
			t = tt
		} else {
			t = reflect.TypeOf(what)
		}

		return t.String()
	}
	return "nil"
}

func getFullname(name, signature string) string {
	return name + separator + signature
}

type Counter int32

func (c *Counter) Add() int {
	current := atomic.AddInt32((*int32)(c), 1)
	return int(current)
}

func (c *Counter) Remove() {
	atomic.AddInt32((*int32)(c), -1)
}
