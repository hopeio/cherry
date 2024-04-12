package trietree

import (
	stringsi "github.com/hopeio/cherry/utils/strings"
	"sort"
)

/**
 * @author     ：lbyi
 * @date       ：2020/12/17 18:21
 * @description：
 */
// TODO:
// untest don't use
type node[T any] struct {
	key      []byte
	indices  []byte
	children []*node[T]
	exists   bool
	value    T
}

func longestCommonPrefix(a, b []byte) int {
	i := 0
	n := min(len(a), len(b))
	for i < n && a[i] == b[i] {
		i++
	}
	return i
}

func (n *node[T]) putValue(key string, v T) {
	fullPath := key
walk:
	for {

		if key == stringsi.BytesToString(n.key) {
			n.exists = true
			n.value = v
			return
		}

		i := longestCommonPrefix(stringsi.StringToBytes(key), n.key)

		if i == len(n.key) {
			for idx, c := range n.indices {
				if c == key[i] {
					n = n.children[idx]
					continue walk
				}
			}
		}
		if i < len(n.key) {
			n.key = []byte(key[:i])

			child := &node[T]{
				key:      []byte(key[i:]),
				indices:  n.indices,
				children: n.children,
				exists:   true,
				value:    v,
			}
			n.children = []*node[T]{child}
			n.indices = []byte{n.key[i]}
			n.exists = false
		}

		if i < len(key) {
			key = key[i:]
			idxc := key[0]

			// Check if a child with the next path byte exists
			for idx, c := range n.indices {
				if c == idxc {
					n = n.children[idx]
					continue walk
				}
			}

			// []byte for proper unicode char conversion, see #65
			n.indices = append(n.indices, idxc)
			child := &node[T]{
				key: []byte(key),
			}
			n.children = append(n.children, child)
			n = child
			n.insertChild(key, fullPath, v)
		}

	}
}

func (n *node[T]) insertChild(key, fullKey string, v T) {
	i := longestCommonPrefix(stringsi.StringToBytes(key), n.key)
	n.key = []byte(key[:i])
	child := &node[T]{
		key:    []byte(key[i:]),
		exists: true,
		value:  v,
	}
	n.children = []*node[T]{child}
}

// 排序
func (n *node[T]) sortIndices() {
	sort.Slice(n.indices, func(i, j int) bool {
		return n.indices[i] < n.indices[j]
	})
	sort.Slice(n.children, func(i, j int) bool {
		return n.children[i].key[0] < n.children[j].key[0]
	})
}

// Shift bytes in array by n bytes left
func shiftNRuneBytes(rb [4]byte, n int) [4]byte {
	switch n {
	case 0:
		return rb
	case 1:
		return [4]byte{rb[1], rb[2], rb[3], 0}
	case 2:
		return [4]byte{rb[2], rb[3]}
	case 3:
		return [4]byte{rb[3]}
	default:
		return [4]byte{}
	}
}

type trietree struct {
}
