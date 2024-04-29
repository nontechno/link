// Copyright 2024 The NonTechno Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package link

import (
	"testing"
)

var count Counter

func add() int {
	return count.Add()
}

func remove() {
	count.Remove()
}

func recCall() int {
	current := add()
	defer remove()

	if current <= maxDepth {
		return recCall()
	}
	return current
}

func TestCounter(t *testing.T) {

	result := recCall()
	if result != maxDepth+1 {
		t.Errorf("got wrong ref count: %v instead of %v", result, maxDepth+1)
	}

	now := add()
	remove()
	if now != 1 {
		t.Errorf("got wrong ref count: %v instead of %v", now, 1)
	}
}
