// Copyright 2024 The NonTechno Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package link

import (
	"reflect"
	"runtime"
	"strings"
)

type Entry struct {
	Name      string // registration name of the function
	Signature string // signature of the function
	Func      string // source.code name of the function
	File      string // file containing the function
	Line      int    // line # containing the function
}

// Available returns available (registered) entries
func Available() []Entry {
	var entries []Entry

	if registryGuard.TryLock() {
		defer registryGuard.Unlock()

		for fullName, operator := range registryStore {
			entry := Entry{}

			parts := strings.Split(fullName, separator)
			entry.Name = getNameSubstitute(parts[0])
			if len(parts) > 1 {
				entry.Signature = parts[1]
			}

			value := reflect.ValueOf(operator)
			pc := value.Pointer()
			rf := runtime.FuncForPC(pc)
			if rf != nil {
				file, line := rf.FileLine(pc)
				entry.Func = rf.Name()
				entry.File = file
				entry.Line = line
			}

			entries = append(entries, entry)
		}
	}

	return entries
}

func getNameSubstitute(origin string) string {
	switch origin {
	case ReportFunc:
		return packageName + ".ReportFunc"
	case WarningFunc:
		return packageName + ".WarningFunc"
	case ErrorFunc:
		return packageName + ".ErrorFunc"
	case TerminateFunc:
		return packageName + ".TerminateFunc"
	}

	return "\"" + origin + "\""
}
