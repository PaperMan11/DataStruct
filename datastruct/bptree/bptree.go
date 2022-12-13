package bptree

import (
	"DataStruct/datastruct"
	"fmt"
)

type BPTree struct {
	root   *BPNode
	order  int // 阶数
	minNum int // 结点最少存在 key 的个数(除root外)
}

// NewBPTree 创建 b+ 树
func NewBPTree(order int) *BPTree {
	return &BPTree{
		root:   newLeafNode(order),
		order:  order,
		minNum: (order + 1) / 2,
	}
}

// Insert 添加指定的 key
func (t *BPTree) Insert(key int64, value interface{}) {
	if t == nil {
		panic("BPTree is null")
	}
	kvnode := newKvNode(key, value)
	t.insert(nil, t.root, kvnode)
}

// insert 递归插入调整
func (t *BPTree) insert(parent, node *BPNode, kvnode *KvNode) {
	// 找到插入结点
	for i := 0; !node.isLeaf && i < node.num; i++ {
		if kvnode.key <= node.childNodes[i].maxKey || i == node.num-1 {
			// 递归查找
			t.insert(node, node.childNodes[i], kvnode)
			break
		}
	}

	// 叶子结点插入数据
	if node.isLeaf {
		node.insertKvn(kvnode)
	}
	// 判断是否分裂了
	newNode := t.spliteNode(node)
	if newNode != nil {
		if parent == nil {
			parent = newIndexNode(t.order)
			parent.addChild(node)
			t.root = parent
		}
		parent.addChild(newNode)
	}
}

// spliteNode 判断是否分裂，分裂了返回分裂结点
func (t *BPTree) spliteNode(node *BPNode) *BPNode {
	if node.isLeaf && node.num > t.order { // 叶子结点
		// 创建新结点
		newNode := newLeafNode(t.order)
		mid := node.num / 2
		copy(newNode.kvNodes[:], node.kvNodes[mid:node.num])
		newNode.num = node.num - mid
		newNode.maxKey = newNode.kvNodes[newNode.num-1].key
		newNode.next = node.next

		// 修改原结点
		node.num = mid
		node.maxKey = node.kvNodes[node.num-1].key
		node.next = newNode
		return newNode
	} else if !node.isLeaf && node.num > t.order { // 索引结点
		// 创建新结点
		newNode := newIndexNode(t.order)
		mid := node.num / 2
		copy(newNode.childNodes[:], node.childNodes[mid:node.num])
		newNode.num = node.num - mid
		newNode.maxKey = newNode.childNodes[newNode.num-1].maxKey

		// 修改原结点
		node.num = mid
		node.maxKey = node.childNodes[node.num-1].maxKey
		return newNode
	}
	return nil
}

// Delete 删除指定的key
func (t *BPTree) Delete(key int64) bool {
	return t.delete(nil, t.root, key)
}

// delete 递归删除调整
func (t *BPTree) delete(parent, node *BPNode, key int64) (ok bool) {
	for i := 0; !node.isLeaf && i < node.num; i++ {
		if key <= node.childNodes[i].maxKey {
			ok = t.delete(node, node.childNodes[i], key)
			break
		}
	}

	// 删除并判断是否需要移动或合并
	if node.isLeaf {
		ok = node.deleteKvn(key)
		if node.num < t.minNum {
			t.kvnMoveOrMerge(parent, node)
		}
	} else {
		node.maxKey = node.childNodes[node.num-1].maxKey
		if node.num == 1 { // 保证 b+ 树的形态
			if parent != nil {
				parent.childNodes[0] = node.childNodes[0]
				parent.maxKey = node.childNodes[0].maxKey
			} else {
				t.root = node.childNodes[0]
			}
			node.freeIndexNode()
		} else if node.num < t.minNum {
			t.childMoveOrMerge(parent, node)
		}
	}

	return
}

// kvnMoveOrMerge 移动或合并叶子结点
func (t *BPTree) kvnMoveOrMerge(parent, node *BPNode) {
	if parent == nil {
		return
	}
	var (
		leftSib  *BPNode
		rightSib *BPNode
	)
	for i := 0; i < parent.num; i++ {
		if parent.childNodes[i] == node {
			if i > 0 {
				leftSib = parent.childNodes[i-1]
			}
			if i < parent.num-1 {
				rightSib = parent.childNodes[i+1]
			}
			break
		}
	}

	// move
	if leftSib != nil && leftSib.num > t.minNum {
		copy(node.kvNodes[1:], node.kvNodes[0:])
		node.kvNodes[0] = leftSib.kvNodes[leftSib.num-1]
		leftSib.num--
		leftSib.maxKey = leftSib.kvNodes[leftSib.num-1].key
		node.num++
		// parent.maxKey = parent.childNodes[parent.num-1].maxKey
		return
	}
	if rightSib != nil && rightSib.num > t.minNum {
		node.kvNodes[node.num] = rightSib.kvNodes[0]
		node.maxKey = rightSib.kvNodes[0].key
		copy(rightSib.kvNodes[0:], rightSib.kvNodes[1:])
		rightSib.num--
		node.num++
		// parent.maxKey = parent.childNodes[parent.num-1].maxKey
		return
	}
	// merge
	if leftSib != nil {
		copy(leftSib.kvNodes[leftSib.num:], node.kvNodes[:node.num])
		leftSib.maxKey = node.maxKey
		leftSib.num += node.num
		leftSib.next = node.next
		parent.deleteChild(node)
		// parent.maxKey = parent.childNodes[parent.num-1].maxKey
		node.freeLeafNode()
		return
	}
	if rightSib != nil {
		copy(node.kvNodes[node.num:], rightSib.kvNodes[0:rightSib.num])
		node.maxKey = rightSib.maxKey
		node.num += rightSib.num
		node.next = rightSib.next
		parent.deleteChild(rightSib)
		// parent.maxKey = parent.childNodes[parent.num-1].maxKey
		rightSib.freeLeafNode()
		return
	}
}

// childMoveOrMerge 移动或合并索引结点
func (t *BPTree) childMoveOrMerge(parent, node *BPNode) {
	if parent == nil {
		return
	}
	var (
		leftSib  *BPNode
		rightSib *BPNode
	)
	for i := 0; i < parent.num; i++ {
		if parent.childNodes[i] == node {
			if i < parent.num-1 {
				rightSib = parent.childNodes[i+1]
			}
			if i > 0 {
				leftSib = parent.childNodes[i-1]
			}
		}
	}

	// move
	if leftSib != nil && leftSib.num > t.minNum {
		copy(node.childNodes[1:], node.childNodes[0:])
		node.childNodes[0] = leftSib.childNodes[leftSib.num-1]
		node.num++
		leftSib.num--
		leftSib.maxKey = leftSib.childNodes[leftSib.num-1].maxKey
		// parent.maxKey = parent.childNodes[parent.num-1].maxKey
		return
	}
	if rightSib != nil && rightSib.num > t.minNum {
		node.childNodes[node.num] = rightSib.childNodes[0]
		node.maxKey = rightSib.childNodes[0].maxKey
		node.num++
		rightSib.num--
		// parent.maxKey = parent.childNodes[parent.num-1].maxKey
		return
	}
	// merge
	if leftSib != nil {
		copy(leftSib.childNodes[leftSib.num:], node.childNodes[:node.num])
		leftSib.maxKey = node.maxKey
		leftSib.num += node.num
		leftSib.next = node.next
		parent.deleteChild(node)
		// parent.maxKey = parent.childNodes[parent.num-1].maxKey
		node.freeIndexNode()
		return
	}
	if rightSib != nil {
		copy(node.childNodes[node.num:], rightSib.childNodes[:rightSib.num])
		node.maxKey = rightSib.maxKey
		node.num += rightSib.num
		node.next = rightSib.next
		parent.deleteChild(rightSib)
		// parent.maxKey = parent.childNodes[parent.num-1].maxKey
		rightSib.freeIndexNode()
		return
	}
}

func (t *BPTree) printInLog() {
	if t.root == nil {
		return
	}
	var (
		flag = true
		list *BPNode
	)
	queue := datastruct.NewQueue(10)
	queue.Push(t.root)
	for queue.Size() != 0 {
		size := queue.Size()
		for i := 0; i < size; i++ {
			node, _ := queue.Pop().(*BPNode)
			// if !node.isLeaf {
			// 	fmt.Printf("maxchild:[%d] ", node.childNodes[node.num-1].maxKey)
			// } else {
			// 	fmt.Printf("maxKey:[%d] ", node.maxKey)
			// }
			for i := 0; i < node.num; i++ {
				if !node.isLeaf {
					fmt.Printf("%d ", node.childNodes[i].maxKey)
					queue.Push(node.childNodes[i])
				} else {
					if flag {
						list = node
						flag = false
					}
					fmt.Printf("[%d,%d] ", node.kvNodes[i].key, node.kvNodes[i].value)
				}
			}
			fmt.Print("| ")
		}
		fmt.Println()
	}

	fmt.Println("--------------leaf----------------")
	for list != nil {
		for i := 0; i < list.num; i++ {
			fmt.Printf("%d ", list.kvNodes[i].key)
		}
		fmt.Print("--> ")
		list = list.next
	}
	fmt.Println("nil")
}
