package trietree

import "testing"

func TestTrieTree(t *testing.T) {
	trie := node[int]{}
	trie.putValue("abc", 1)
	trie.putValue("abcdef", 2)
	trie.putValue("aabcdef", 3)
}
