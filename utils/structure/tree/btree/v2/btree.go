// Copyright 2020 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package btree

const maxItems = 255
const minItems = maxItems * 40 / 100

type node[T comparable] struct {
	leaf     bool
	numItems int16
	items    [maxItems]T
	children *[maxItems + 1]*node[T]
}

type justaLeaf[T comparable] struct {
	leaf     bool
	numItems int16
	items    [maxItems]T
}

// BTree is an ordered set items
type BTree[T comparable] struct {
	root   *node[T]
	length int
	less   func(a, b T) bool
	lnode  *node[T]
	zero   T
}

func newNode[T comparable](leaf bool) *node[T] {
	n := &node[T]{leaf: leaf}
	if !leaf {
		n.children = new([maxItems + 1]*node[T])
	}
	return n
}

// PathHint is a utility type used with the *Hint() functions. Hints provide
// faster operations for clustered keys.
type PathHint struct {
	path [8]uint8
}

// New returns a new BTree
func New[T comparable](less func(a, b T) bool) *BTree[T] {
	if less == nil {
		panic("nil less")
	}
	tr := new(BTree[T])
	tr.less = less
	return tr
}

// Less is a convenience function that performs a comparison of two items
// using the same "less" function provided to New.
func (tr *BTree[T]) Less(a, b T) bool {
	return tr.less(a, b)
}

func (n *node[T]) find(key T, less func(a, b T) bool,
	hint *PathHint, depth int,
) (index int16, found bool) {
	low := int16(0)
	high := n.numItems - 1
	if hint != nil && depth < 8 {
		index = int16(hint.path[depth])
		if index > n.numItems-1 {
			index = n.numItems - 1
		}
		if less(key, n.items[index]) {
			high = index - 1
		} else if less(n.items[index], key) {
			low = index + 1
		} else {
			found = true
			goto done
		}
	}
	for low <= high {
		mid := low + ((high+1)-low)/2
		if !less(key, n.items[mid]) {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	if low > 0 && !less(n.items[low-1], key) &&
		!less(key, n.items[low-1]) {
		index = low - 1
		found = true
	} else {
		index = low
		found = false
	}
done:
	if hint != nil && depth < 8 {
		if n.leaf && found {
			hint.path[depth] = byte(index + 1)
		} else {
			hint.path[depth] = byte(index)
		}
	}
	return index, found
}

// SetHint sets or replace a value for a key using a path hint
func (tr *BTree[T]) SetHint(item T, hint *PathHint) (prev T) {
	if item == nil {
		panic("nil item")
	}
	if tr.root == nil {
		tr.root = newNode[T](true)
		tr.root.items[0] = item
		tr.root.numItems = 1
		tr.length = 1
		return
	}
	prev = tr.root.set(item, tr.less, hint, 0)
	if prev != nil {
		return prev
	}
	tr.lnode = nil
	if tr.root.numItems == maxItems {
		n := tr.root
		right, median := n.split()
		tr.root = newNode[T](false)
		tr.root.children[0] = n
		tr.root.items[0] = median
		tr.root.children[1] = right
		tr.root.numItems = 1
	}
	tr.length++
	return prev
}

// Set or replace a value for a key
func (tr *BTree[T]) Set(item T) (prev T) {
	return tr.SetHint(item, nil)
}

func (n *node[T]) split() (right *node[T], median T) {
	right = newNode[T](n.leaf)
	median = n.items[maxItems/2]
	copy(right.items[:maxItems/2], n.items[maxItems/2+1:])
	if !n.leaf {
		copy(right.children[:maxItems/2+1], n.children[maxItems/2+1:])
	}
	right.numItems = maxItems / 2
	if !n.leaf {
		for i := maxItems/2 + 1; i < maxItems+1; i++ {
			n.children[i] = nil
		}
	}
	for i := maxItems / 2; i < maxItems; i++ {
		n.items[i] = *new(T)
	}
	n.numItems = maxItems / 2
	return right, median
}

func (n *node[T]) set(item T, less func(a, b T) bool,
	hint *PathHint, depth int,
) (prev interface{}) {
	i, found := n.find(item, less, hint, depth)
	if found {
		prev = n.items[i]
		n.items[i] = item
		return prev
	}
	if n.leaf {
		copy(n.items[i+1:n.numItems+1], n.items[i:n.numItems])
		n.items[i] = item
		n.numItems++
		return nil
	}
	prev = n.children[i].set(item, less, hint, depth+1)
	if prev != nil {
		return prev
	}
	if n.children[i].numItems == maxItems {
		right, median := n.children[i].split()
		copy(n.children[i+1:], n.children[i:])
		copy(n.items[i+1:], n.items[i:])
		n.items[i] = median
		n.children[i+1] = right
		n.numItems++
	}
	return prev
}

func (n *node[T]) scan(iter func(item T) bool) bool {
	if n.leaf {
		for i := int16(0); i < n.numItems; i++ {
			if !iter(n.items[i]) {
				return false
			}
		}
		return true
	}
	for i := int16(0); i < n.numItems; i++ {
		if !n.children[i].scan(iter) {
			return false
		}
		if !iter(n.items[i]) {
			return false
		}
	}
	return n.children[n.numItems].scan(iter)
}

// Get a value for key
func (tr *BTree[T]) Get(key interface{}) interface{} {
	return tr.GetHint(key, nil)
}

// GetHint gets a value for key using a path hint
func (tr *BTree[T]) GetHint(key interface{}, hint *PathHint) interface{} {
	if tr.root == nil || key == nil {
		return nil
	}
	depth := 0
	n := tr.root
	for {
		i, found := n.find(key, tr.less, hint, depth)
		if found {
			return n.items[i]
		}
		if n.leaf {
			return nil
		}
		n = n.children[i]
		depth++
	}
}

// Len returns the number of items in the tree
func (tr *BTree[T]) Len() int {
	return tr.length
}

// Delete a value for a key
func (tr *BTree[T]) Delete(key interface{}) interface{} {
	return tr.DeleteHint(key, nil)
}

// DeleteHint deletes a value for a key using a path hint
func (tr *BTree[T]) DeleteHint(key interface{}, hint *PathHint) interface{} {
	if tr.root == nil || key == nil {
		return nil
	}
	prev := tr.root.delete(false, key, tr.less, hint, 0)
	if prev == nil {
		return nil
	}
	tr.lnode = nil
	if tr.root.numItems == 0 && !tr.root.leaf {
		tr.root = tr.root.children[0]
	}
	tr.length--
	if tr.length == 0 {
		tr.root = nil
	}
	return prev
}

func (n *node[T]) delete(max bool, key T,
	less func(a, b T) bool, hint *PathHint, depth int,
) interface{} {
	var i int16
	var found bool
	if max {
		i, found = n.numItems-1, true
	} else {
		i, found = n.find(key, less, hint, depth)
	}
	if n.leaf {
		if found {
			prev := n.items[i]
			// found the items at the leaf, remove it and return.
			copy(n.items[i:], n.items[i+1:n.numItems])
			n.items[n.numItems-1] = *new(T)
			n.numItems--
			return prev
		}
		return nil
	}

	var prev interface{}
	if found {
		if max {
			i++
			prev = n.children[i].delete(true, *new(T), less, nil, 0)
		} else {
			prev = n.items[i]
			maxItem := n.children[i].delete(true, *new(T), less, nil, 0)
			n.items[i] = maxItem
		}
	} else {
		prev = n.children[i].delete(max, key, less, hint, depth+1)
	}
	if prev == nil {
		return nil
	}
	if n.children[i].numItems < minItems {
		if i == n.numItems {
			i--
		}
		if n.children[i].numItems+n.children[i+1].numItems+1 < maxItems {
			// merge left + item + right
			n.children[i].items[n.children[i].numItems] = n.items[i]
			copy(n.children[i].items[n.children[i].numItems+1:],
				n.children[i+1].items[:n.children[i+1].numItems])
			if !n.children[0].leaf {
				copy(n.children[i].children[n.children[i].numItems+1:],
					n.children[i+1].children[:n.children[i+1].numItems+1])
			}
			n.children[i].numItems += n.children[i+1].numItems + 1
			copy(n.items[i:], n.items[i+1:n.numItems])
			copy(n.children[i+1:], n.children[i+2:n.numItems+1])
			n.items[n.numItems] = *new(T)
			n.children[n.numItems+1] = nil
			n.numItems--
		} else if n.children[i].numItems > n.children[i+1].numItems {
			// move left -> right
			copy(n.children[i+1].items[1:],
				n.children[i+1].items[:n.children[i+1].numItems])
			if !n.children[0].leaf {
				copy(n.children[i+1].children[1:],
					n.children[i+1].children[:n.children[i+1].numItems+1])
			}
			n.children[i+1].items[0] = n.items[i]
			if !n.children[0].leaf {
				n.children[i+1].children[0] =
					n.children[i].children[n.children[i].numItems]
			}
			n.children[i+1].numItems++
			n.items[i] = n.children[i].items[n.children[i].numItems-1]
			n.children[i].items[n.children[i].numItems-1] = *new(T)
			if !n.children[0].leaf {
				n.children[i].children[n.children[i].numItems] = nil
			}
			n.children[i].numItems--
		} else {
			// move right -> left
			n.children[i].items[n.children[i].numItems] = n.items[i]
			if !n.children[0].leaf {
				n.children[i].children[n.children[i].numItems+1] =
					n.children[i+1].children[0]
			}
			n.children[i].numItems++
			n.items[i] = n.children[i+1].items[0]
			copy(n.children[i+1].items[:],
				n.children[i+1].items[1:n.children[i+1].numItems])
			if !n.children[0].leaf {
				copy(n.children[i+1].children[:],
					n.children[i+1].children[1:n.children[i+1].numItems+1])
			}
			n.children[i+1].numItems--
		}
	}
	return prev
}

// Ascend the tree within the range [pivot, last]
// Pass nil for pivot to scan all item in ascending order
// Return false to stop iterating
func (tr *BTree[T]) Ascend(pivot T, iter func(item T) bool) {
	if tr.root == nil {
		return
	}
	if pivot == nil {
		tr.root.scan(iter)
	} else if tr.root != nil {
		tr.root.ascend(pivot, tr.less, nil, 0, iter)
	}
}

func (n *node[T]) ascend(pivot T, less func(a, b T) bool,
	hint *PathHint, depth int, iter func(item T) bool,
) bool {
	i, found := n.find(pivot, less, hint, depth)
	if !found {
		if !n.leaf {
			if !n.children[i].ascend(pivot, less, hint, depth+1, iter) {
				return false
			}
		}
	}
	for ; i < n.numItems; i++ {
		if !iter(n.items[i]) {
			return false
		}
		if !n.leaf {
			if !n.children[i+1].scan(iter) {
				return false
			}
		}
	}
	return true
}

func (n *node[T]) reverse(iter func(item T) bool) bool {
	if n.leaf {
		for i := n.numItems - 1; i >= 0; i-- {
			if !iter(n.items[i]) {
				return false
			}
		}
		return true
	}
	if !n.children[n.numItems].reverse(iter) {
		return false
	}
	for i := n.numItems - 1; i >= 0; i-- {
		if !iter(n.items[i]) {
			return false
		}
		if !n.children[i].reverse(iter) {
			return false
		}
	}
	return true
}

// Descend the tree within the range [pivot, first]
// Pass nil for pivot to scan all item in descending order
// Return false to stop iterating
func (tr *BTree[T]) Descend(pivot T, iter func(item T) bool) {
	if tr.root == nil {
		return
	}
	if pivot == nil {
		tr.root.reverse(iter)
	} else if tr.root != nil {
		tr.root.descend(pivot, tr.less, nil, 0, iter)
	}
}

func (n *node[T]) descend(pivot T, less func(a, b T) bool,
	hint *PathHint, depth int, iter func(item T) bool,
) bool {
	i, found := n.find(pivot, less, hint, depth)
	if !found {
		if !n.leaf {
			if !n.children[i].descend(pivot, less, hint, depth+1, iter) {
				return false
			}
		}
		i--
	}
	for ; i >= 0; i-- {
		if !iter(n.items[i]) {
			return false
		}
		if !n.leaf {
			if !n.children[i].reverse(iter) {
				return false
			}
		}
	}
	return true
}

// Load is for bulk loading pre-sorted items
func (tr *BTree[T]) Load(item T) (T, bool) {
	if item == tr.zero {
		panic("nil item")
	}
	if tr.lnode != nil && tr.lnode.numItems < maxItems-2 {
		if tr.less(tr.lnode.items[tr.lnode.numItems-1], item) {
			tr.lnode.items[tr.lnode.numItems] = item
			tr.lnode.numItems++
			tr.length++
			return tr.zero, false
		}
	}
	prev := tr.Set(item)
	if prev != nil {
		return prev, true
	}
	n := tr.root
	for {
		if n.leaf {
			tr.lnode = n
			break
		}
		n = n.children[n.numItems]
	}
	return tr.zero, false
}

// Min returns the minimum item in tree.
// Returns nil if the tree has no items.
func (tr *BTree[T]) Min() (T, bool) {
	if tr.root == nil {
		return tr.zero, false
	}
	n := tr.root
	for {
		if n.leaf {
			return n.items[0], true
		}
		n = n.children[0]
	}
}

// Max returns the maximum item in tree.
// Returns nil if the tree has no items.
func (tr *BTree[T]) Max() (T, bool) {
	if tr.root == nil {
		return tr.zero, false
	}
	n := tr.root
	for {
		if n.leaf {
			return n.items[n.numItems-1], true
		}
		n = n.children[n.numItems]
	}
}

// PopMin removes the minimum item in tree and returns it.
// Returns nil if the tree has no items.
func (tr *BTree[T]) PopMin() (T, bool) {
	if tr.root == nil {
		return tr.zero, false
	}
	tr.lnode = nil
	n := tr.root
	for {
		if n.leaf {
			item := n.items[0]
			if n.numItems == minItems {
				return tr.Delete(item), true
			}
			copy(n.items[:], n.items[1:])
			n.items[n.numItems-1] = tr.zero
			n.numItems--
			tr.length--
			return item, true
		}
		n = n.children[0]
	}
}

// PopMax removes the minimum item in tree and returns it.
// Returns nil if the tree has no items.
func (tr *BTree[T]) PopMax() (T, bool) {
	if tr.root == nil {
		return tr.zero, false
	}
	tr.lnode = nil
	n := tr.root
	for {
		if n.leaf {
			item := n.items[n.numItems-1]
			if n.numItems == minItems {
				return tr.Delete(item), true
			}
			n.items[n.numItems-1] = tr.zero
			n.numItems--
			tr.length--
			return item, true
		}
		n = n.children[n.numItems]
	}
}

// Height returns the height of the tree.
// Returns zero if tree has no items.
func (tr *BTree[T]) Height() int {
	var height int
	if tr.root != nil {
		n := tr.root
		for {
			height++
			if n.leaf {
				break
			}
			n = n.children[n.numItems]
		}
	}
	return height
}

// Walk iterates over all items in tree, in order.
// The items param will contain one or more items.
func (tr *BTree[T]) Walk(iter func(item []T)) {
	if tr.root != nil {
		tr.root.walk(iter)
	}
}

func (n *node[T]) walk(iter func(item []T)) {
	if n.leaf {
		iter(n.items[:n.numItems])
	} else {
		for i := int16(0); i < n.numItems; i++ {
			n.children[i].walk(iter)
			iter(n.items[i : i+1])
		}
		n.children[n.numItems].walk(iter)
	}
}
