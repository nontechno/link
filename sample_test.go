// Copyright 2024 The NonTechno Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package link

import (
	"testing"
)

func TestBasic(t *testing.T) {
	const linkage = "fjsjdhfjsdfjksdfhksjdhf"

	some := func(s string) string {
		return "[" + s + "]"
	}

	var other func(s string) string
	Link(&other, linkage, nil)
	Register(some, linkage)

	input := ""
	output := other(input)

	t.Logf("(%s) ===> (%s)", input, output)
}

func TestMultiple(t *testing.T) {
	const name = "some.unique.name"

	one := func() {
		t.Logf("---one")
	}
	Register(one, name)

	two := func(string) {
		t.Logf("---two")
	}
	Register(two, name)

	three := func() string {
		t.Logf("---three")
		return "#3"
	}
	Register(three, name)

	var other func(s string)
	Link(&other, name, nil)
	other("")

	var another func() string
	Link(&another, name, nil)
	another()
}

func TestWarning(t *testing.T) {
	defer func() {
		if recover() != nil {
			t.Logf("catching termination...\n")

			report(t)
		}
	}()

	warn := func(f string, a ...interface{}) {
		t.Logf("======> "+f+"\n", a...)
	}

	term := func() {
		t.Logf("terminating...\n")
		panic("longjump")
	}

	Register(warn)
	Register(warn, ReportFunc, WarningFunc, ErrorFunc)
	Register(term, TerminateFunc)

	fallback := func() {
		t.Logf("fallback")
	}

	var other, another func()
	Link(&other, "bs", fallback)
	Link(&another, "bs", nil)

	other()
	another()
}

func report(t *testing.T) {
	t.Logf("------------------------------ here is what's available:\n")
	for _, entry := range Available() {
		t.Logf("\t%v\n", entry.Name)
		t.Logf("\t%v\n", entry.Signature)
		t.Logf("\t%v\n", entry.Func)
		t.Logf("\t%v\n", entry.File)
		t.Logf("\t%v\n", entry.Line)
	}
}
