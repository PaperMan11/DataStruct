package bptree

// BNode B+树的结点
type BPNode struct {
	num        int       // 当前结点 key 的个数（leaf为kv的数量，index为child数量）
	maxKey     int64     // 当前结点中最大的 key
	kvNodes    []*KvNode // KvNode 数组
	childNodes []*BPNode // 孩子结点
	isLeaf     bool      // 是否为叶子结点
	next       *BPNode
}

// newIndexNode 新建索引结点
func newIndexNode(order int) *BPNode {
	return &BPNode{
		num:        0,
		childNodes: make([]*BPNode, order+1),
		isLeaf:     false,
	}
}

// newLeafNode 新建叶子结点
func newLeafNode(order int) *BPNode {
	return &BPNode{
		num:     0,
		kvNodes: make([]*KvNode, order+1),
		isLeaf:  true,
	}
}

func (node *BPNode) freeIndexNode() {
	node.childNodes = nil
}

func (node *BPNode) freeLeafNode() {
	node.kvNodes = nil
}

// insertKvn 叶子结点插入 kvn
func (node *BPNode) insertKvn(kvn *KvNode) {
	if node.num == 0 {
		node.kvNodes[0] = kvn
		node.maxKey = kvn.key
		node.num++
		return
	}
	if kvn.key > node.kvNodes[node.num-1].key {
		node.kvNodes[node.num] = kvn
		node.maxKey = kvn.key
		node.num++
		return
	}
	for i, v := range node.kvNodes {
		if v.key == kvn.key {
			v.value = kvn.value
			return
		} else if v.key > kvn.key {
			copy(node.kvNodes[i+1:], node.kvNodes[i:]) // 后移
			node.kvNodes[i] = kvn                      // 插入
			node.num++
			return
		}
	}
}

// deleteKvn 叶子结点删除 kvn
func (node *BPNode) deleteKvn(key int64) bool {
	for i := 0; i < node.num; i++ {
		if key == node.kvNodes[i].key {
			copy(node.kvNodes[i:], node.kvNodes[i+1:])
			node.num--
			if node.num == 0 {
				node.maxKey = 0 // 暂定为 0
			} else if i == node.num { // 删除了最后一个点
				node.maxKey = node.kvNodes[node.num-1].key
			}
			return true
		}
	}
	return false
}

// addChild 索引结点添加子结点
func (node *BPNode) addChild(child *BPNode) {
	if node.num == 0 {
		node.childNodes[0] = child
		node.maxKey = child.maxKey
		node.num++
		return
	} else if child.maxKey > node.childNodes[node.num-1].maxKey {
		node.childNodes[node.num] = child
		node.maxKey = child.maxKey
		node.num++
		return
	}

	for i, n := range node.childNodes {
		if child.maxKey < n.maxKey { // 后移并插入
			copy(node.childNodes[i+1:], node.childNodes[i:])
			node.childNodes[i] = child
			node.num++
			return
		}
	}
}

// deleteChild 索引结点删除子结点
func (node *BPNode) deleteChild(child *BPNode) {
	if child == node.childNodes[node.num-1] {
		node.num--
		node.maxKey = node.childNodes[node.num-1].maxKey
		return
	}

	for i, n := range node.childNodes {
		if child == n {
			copy(node.childNodes[i:], node.childNodes[i+1:])
			node.num--
			return
		}
	}
}
