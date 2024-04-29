// Copyright 2024 The NonTechno Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package link

import (
	"sync"
)

var (
	registryGuard sync.RWMutex
	registryStore = map[string]interface{}{}
	onWarning     = localOnWarning
	onError       = localOnError
	onReport      = localOnReport
	onTerminate   = localOnTerminate
	entryCounter  Counter
	remoteCounter Counter
)
