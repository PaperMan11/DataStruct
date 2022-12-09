package btree

// KvNode BNode中的成员
type KvNode struct {
	key   int64
	value interface{}
}

// newKvNode 新建KvNode成员
func newKvNode(key int64, value interface{}) *KvNode {
	return &KvNode{
		key:   key,
		value: value,
	}
}

