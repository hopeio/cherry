package node

type Node[T any] struct {
	data     T
	previous *Node[T]
	next     *Node[T]
}
