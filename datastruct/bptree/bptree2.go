package bptree

import (
	"DataStruct/datastruct"
	"cmp"
	"fmt"
)

type keyPairs[K cmp.Ordered, V any] struct {
	Key   K
	Value V
}

// b+树节点
type BpNode[K cmp.Ordered, V any] struct {
	isLeaf   bool              // 是否是叶子节点
	keyPairs []*keyPairs[K, V] // 键值对
	children []*BpNode[K, V]   // 子节点
	next     *BpNode[K, V]     // 叶子节点的下一个节点指针
	order    int               // 阶
	num      int               // 节点的键值对数量
	maxKey   K                 // 节点最大键值
}

func NewBpNode[K cmp.Ordered, V any](order int, isLeaf bool) *BpNode[K, V] {
	return &BpNode[K, V]{
		isLeaf:   isLeaf,
		keyPairs: make([]*keyPairs[K, V], order+1),
		children: make([]*BpNode[K, V], order+1),
		order:    order,
		num:      0,
	}
}

// 叶子节点插入
func (node *BpNode[K, V]) InsertKeyPair(key K, value V) bool {
	i := 0
	// 找到插入节点
	for i < node.num {
		if node.keyPairs[i].Key == key {
			node.keyPairs[i].Value = value
			return false
		} else if node.keyPairs[i].Key > key {
			break
		}
		i++
	}
	copy(node.keyPairs[i+1:], node.keyPairs[i:])
	node.keyPairs[i] = &keyPairs[K, V]{
		Key:   key,
		Value: value,
	}
	node.num++
	node.maxKey = node.keyPairs[node.num-1].Key
	return true
}

// 叶子节点删除
func (node *BpNode[K, V]) RemoveKeyPair(key K) bool {
	i := 0
	for i < node.num {
		if node.keyPairs[i].Key == key {
			break
		}
		i++
	}
	if i == node.num {
		return false
	}
	copy(node.keyPairs[i:], node.keyPairs[i+1:])
	node.num--
	if node.num <= 0 {
		var zero K
		node.maxKey = zero
	} else {
		node.maxKey = node.keyPairs[node.num-1].Key
	}
	return true
}

func (node *BpNode[K, V]) GetFirstKeyPair() *keyPairs[K, V] {
	return node.keyPairs[0]
}

func (node *BpNode[K, V]) GetLastKeyPair() *keyPairs[K, V] {
	return node.keyPairs[node.num-1]
}

// 非叶子节点插入子节点
func (node *BpNode[K, V]) InsertChild(child *BpNode[K, V]) bool {
	if child == nil {
		return false
	}

	i := 0
	// 找到插入节点
	for i < node.num && node.children[i].maxKey < child.maxKey {
		i++
	}
	copy(node.children[i+1:], node.children[i:])
	node.children[i] = child
	node.num++
	node.maxKey = node.children[node.num-1].maxKey
	return true
}

// 非叶子节点删除子节点
func (node *BpNode[K, V]) RemoveChild(child *BpNode[K, V]) bool {
	if child == nil {
		return false
	}

	i := 0
	for i < node.num {
		if node.children[i] == child {
			break
		}
		i++
	}
	if i == node.num {
		return false
	}
	copy(node.children[i:], node.children[i+1:])
	node.num--
	node.maxKey = node.children[node.num-1].maxKey
	return true
}

func (node *BpNode[K, V]) GetFirstChild() *BpNode[K, V] {
	return node.children[0]
}

func (node *BpNode[K, V]) GetLastChild() *BpNode[K, V] {
	return node.children[node.num-1]
}

type BpTree[K cmp.Ordered, V any] struct {
	root    *BpNode[K, V] // 根节点
	order   int           // B+树的阶
	size    int           // 树中键值对的数量
	minKeys int           // 非根节点的最小键值对数量
}

func NewBpTree[K cmp.Ordered, V any](order int) *BpTree[K, V] {
	if order < 3 {
		panic("order must be greater than or equal to 3")
	}
	return &BpTree[K, V]{
		root:    NewBpNode[K, V](order, true),
		order:   order,
		size:    0,
		minKeys: (order + 1) / 2, // math.ceil(order/2)
	}
}

func (tree *BpTree[K, V]) Size() int {
	return tree.size
}

func (tree *BpTree[K, V]) IsEmpty() bool {
	return tree.size == 0
}

func (tree *BpTree[K, V]) Insert(key K, value V) {
	splitNode := tree.insert(tree.root, key, value)
	if splitNode != nil {
		newRoot := NewBpNode[K, V](tree.order, false)
		newRoot.InsertChild(tree.root)
		newRoot.InsertChild(splitNode)
		tree.root = newRoot
	}
}

func (tree *BpTree[K, V]) insert(node *BpNode[K, V], key K, value V) (newNode *BpNode[K, V]) {
	if node.isLeaf {
		if node.InsertKeyPair(key, value) {
			tree.size++
		}
		// 节点分裂
		if node.num > tree.order {
			return tree.splitNode(node)
		}
	} else {
		// 找到插入节点
		i := 0
		// 为什么是num-1: 因为插入key的值比树中所有的key都大(即:i=node.num -> node.children[i]==nil)
		for i < node.num-1 && node.children[i].maxKey < key {
			i++
		}
		splitNode := tree.insert(node.children[i], key, value)
		if splitNode != nil {
			if node.InsertChild(splitNode) {
				// 节点分裂
				if node.num > tree.order {
					return tree.splitNode(node)
				}
			}
		}
		node.maxKey = node.children[node.num-1].maxKey
	}
	return nil
}

func (tree *BpTree[K, V]) Remove(key K) bool {
	ok := tree.remove(tree.root, key)
	if ok && tree.root.num == 1 && !tree.root.isLeaf {
		if tree.root.children[0] == nil {
			tree.root = NewBpNode[K, V](tree.order, true)
		} else {
			tree.root = tree.root.children[0]
		}
	}
	return ok
}

func (tree *BpTree[K, V]) remove(node *BpNode[K, V], key K) bool {
	if node.isLeaf {
		if node.RemoveKeyPair(key) {
			tree.size--
			return true
		}
		return false
	} else {
		// 找到删除节点
		i := 0
		for i < node.num && node.children[i].maxKey < key {
			i++
		}
		if i >= node.num || node.children[i] == nil {
			return false
		}

		child := node.children[i]
		if !tree.remove(child, key) {
			return false
		}
		// 节点平衡
		if child.num < tree.minKeys {
			var (
				leftSibling  *BpNode[K, V]
				rightSibling *BpNode[K, V]
			)
			if i > 0 {
				leftSibling = node.children[i-1]
			}
			if i < node.num-1 {
				rightSibling = node.children[i+1]
			}

			if leftSibling != nil && leftSibling.num > tree.minKeys {
				// 尝试从左兄弟借
				tree.borrowFromLeft(child, leftSibling)
			} else if rightSibling != nil && rightSibling.num > tree.minKeys {
				// 尝试从右兄弟借
				tree.borrowFromRight(child, rightSibling)
			} else if i > 0 {
				// 与左兄弟合并
				tree.mergeNodes(leftSibling, child)
				node.RemoveChild(child)
			} else if i < node.num-1 {
				// 与右兄弟合并
				tree.mergeNodes(child, rightSibling)
				node.RemoveChild(rightSibling)
			}
		}
		node.maxKey = node.children[node.num-1].maxKey
		return true
	}
}

func (tree *BpTree[K, V]) splitNode(node *BpNode[K, V]) (newNode *BpNode[K, V]) {
	if node.isLeaf {
		newNode = NewBpNode[K, V](tree.order, true)
		mid := tree.order / 2
		copy(newNode.keyPairs[:], node.keyPairs[mid+1:])
		newNode.num = node.num - mid - 1
		newNode.maxKey = newNode.keyPairs[newNode.num-1].Key
		node.num = mid + 1
		node.maxKey = node.keyPairs[node.num-1].Key
		newNode.next = node.next
		node.next = newNode
	} else {
		newNode = NewBpNode[K, V](tree.order, false)
		mid := tree.order / 2
		copy(newNode.children[:], node.children[mid+1:])
		newNode.num = node.num - mid - 1
		newNode.maxKey = newNode.children[newNode.num-1].maxKey
		node.num = mid + 1
		node.maxKey = node.children[node.num-1].maxKey
	}
	return
}

func (tree *BpTree[K, V]) mergeNodes(node *BpNode[K, V], sibling *BpNode[K, V]) {
	if node.isLeaf {
		copy(node.keyPairs[node.num:], sibling.keyPairs[:sibling.num])
		node.num += sibling.num
		node.maxKey = node.keyPairs[node.num-1].Key
		node.next = sibling.next
	} else {
		copy(node.children[node.num:], sibling.children[:sibling.num])
		node.num += sibling.num
	}
}

func (tree *BpTree[K, V]) borrowFromLeft(node *BpNode[K, V], left *BpNode[K, V]) {
	if node.isLeaf {
		lastKeyPair := left.GetLastKeyPair()
		node.InsertKeyPair(lastKeyPair.Key, lastKeyPair.Value)
		left.RemoveKeyPair(lastKeyPair.Key)
	} else {
		lastChild := left.GetLastChild()
		node.InsertChild(lastChild)
		left.RemoveChild(lastChild)
	}
}

func (tree *BpTree[K, V]) borrowFromRight(node *BpNode[K, V], right *BpNode[K, V]) {
	if node.isLeaf {
		firstKeyPair := right.GetFirstKeyPair()
		node.InsertKeyPair(firstKeyPair.Key, firstKeyPair.Value)
		right.RemoveKeyPair(firstKeyPair.Key)
	} else {
		firstChild := right.GetFirstChild()
		node.InsertChild(firstChild)
		right.RemoveChild(firstChild)
	}
}

func (tree *BpTree[K, V]) Search(key K) (V, bool) {
	var zero V
	current := tree.root
	for !current.isLeaf {
		i := 0
		for i < current.num && current.children[i].maxKey < key {
			i++
		}
		if i == current.num || current.children[i] == nil {
			return zero, false
		}
		current = current.children[i]
	}

	for i := 0; i < current.num; i++ {
		if current.keyPairs[i].Key == key {
			return current.keyPairs[i].Value, true
		}
	}
	return zero, false
}

func (tree *BpTree[K, V]) RangeSearch(start K, end K) []V {
	res := make([]V, 0)
	current := tree.root
	// 找到第一个节点
	for !current.isLeaf {
		i := 0
		for i < current.num && current.children[i].maxKey < start {
			i++
		}
		if i == current.num || current.children[i] == nil {
			return res
		}
		current = current.children[i]
	}
	// 遍历叶子节点
	for current != nil {
		for i := 0; i < current.num; i++ {
			if current.keyPairs[i].Key >= start && current.keyPairs[i].Key <= end {
				res = append(res, current.keyPairs[i].Value)
			}
		}
		current = current.next
	}
	return res
}

// 打印树结构（用于调试）
func (tree *BpTree[K, V]) PrintTree() {
	if tree.size == 0 {
		return
	}

	queue := datastruct.NewQueue(100)
	queue.Push(tree.root)
	level := 0
	for queue.Size() != 0 {
		levelSize := queue.Size()
		levelStr := fmt.Sprintf("Level %d: -> ", level)
		level++
		for i := 0; i < levelSize; i++ {
			node := queue.Pop().(*BpNode[K, V])
			levelStr += fmt.Sprintf("(num:%v, maxKey:%v, isLeaf:%v)  [ ", node.num, node.maxKey, node.isLeaf)
			for j := 0; j < node.num; j++ {
				if node.isLeaf {
					levelStr += fmt.Sprintf("<%v,%v> ", node.keyPairs[j].Key, node.keyPairs[j].Value)
				} else {
					levelStr += fmt.Sprintf("%v ", node.children[j].maxKey)
					queue.Push(node.children[j])
				}
			}
			levelStr += "]\t"
		}
		fmt.Println(levelStr)
	}
}
