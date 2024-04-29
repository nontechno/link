// Copyright 2024 The NonTechno Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package link

import (
	"reflect"
)

func Link(fptr interface{}, linkage string, fallback interface{}) {
	signature := ""
	if pointer, elem := isPointer(fptr); pointer {
		if yes, _ := isFunction(elem); !yes {
			onError("supplied parameter is not a pointer to a function (%v)", fptr)
		}
		signature = getSignature(elem)
	} else {
		onError("supplied parameter is not a pointer (%v)", fptr)
	}

	if fallback != nil {
		if fn, _ := isFunction(fallback); !fn {
			onError("supplied fallback is not a function (%v)", fptr)
		}
	}

	fn := reflect.ValueOf(fptr).Elem()
	variadic := fn.Type().IsVariadic()

	universal := func(args []reflect.Value) []reflect.Value {
		variadic := variadic
		// a simple guard against an endless recursion
		current := entryCounter.Add()
		defer entryCounter.Remove()

		if current > maxDepth {
			onError("too many hops (%v) to resolve linkage (%s)", current, linkage)
			return []reflect.Value{}
		}

		if target := resolve(fn, linkage, signature); target != nil {
			onReport("resolved linkage (%s)", linkage)
			return call(target, args, variadic)
		} else {
			if fallback != nil {
				arg := reflect.ValueOf(fallback)
				fn.Set(arg)

				onReport("unresolved linkage (%s), using fallback", linkage)
				return call(fallback, args, variadic)
			}

			onError("failed to resolve linkage (%s)", linkage)
			return []reflect.Value{}
		}
	}

	v := reflect.MakeFunc(fn.Type(), universal)
	fn.Set(v)
}

func call(target any, args []reflect.Value, variadic bool) []reflect.Value {
	if variadic {
		return reflect.ValueOf(target).CallSlice(args) // this can be a chain call
	}
	return reflect.ValueOf(target).Call(args) // this can be a chain call
}
