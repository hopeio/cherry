// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sync

import (
	"sync/atomic"
	"unsafe"
)

// Stack implements lock-free freelist based stack.
type Stack struct {
	top unsafe.Pointer
	len uint64
}

// NewStack creates a new lock-free queue.
func NewStack() *Stack {
	return &Stack{}
}

// Pop pops value from the top of the stack.
func (s *Stack) Pop() interface{} {
	var top, next unsafe.Pointer
	var item *directItem
	for {
		top = atomic.LoadPointer(&s.top)
		if top == nil {
			return nil
		}
		item = (*directItem)(top)
		next = atomic.LoadPointer(&item.next)
		if atomic.CompareAndSwapPointer(&s.top, top, next) {
			atomic.AddUint64(&s.len, ^uint64(0))
			return item.v
		}
	}
}

// Push pushes a value on top of the stack.
func (s *Stack) Push(v interface{}) {
	item := directItem{v: v}
	var top unsafe.Pointer
	for {
		top = atomic.LoadPointer(&s.top)
		item.next = top
		if atomic.CompareAndSwapPointer(&s.top, top, unsafe.Pointer(&item)) {
			atomic.AddUint64(&s.len, 1)
			return
		}
	}
}
