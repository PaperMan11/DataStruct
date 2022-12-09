package btree

import (
	"DataStruct/datastruct"
	"fmt"
	"math"
)

type BTree struct {
	root   *BNode // 根结点
	order  int    // 阶数
	minNum int    // 结点最少存在 key 的个数(除root外)
}

// NewBTree 创建 B 树
func NewBTree(order int) *BTree {
	return &BTree{
		root:   nil,
		order:  order,
		minNum: int(math.Ceil(float64(order)/2) - 1),
	}
}

// Insert 插入
func (btree *BTree) Insert(key int64, value interface{}) {
	kvnode := newKvNode(key, value)
	if btree.root == nil {
		root := newBNode(btree.order)
		root.kvNodes[0] = kvnode
		root.num++
		btree.root = root
		return
	}
	if btree.root.num == 0 {
		btree.root.kvNodes[0] = kvnode
		btree.root.num++
	}
	btree.insert(btree.root, 0, kvnode)
}

// insert 找到插入点
func (btree *BTree) insert(pnode *BNode, index int, kvnode *KvNode) {
	childNum := len(pnode.childNodes)
	if kvnode.key < pnode.kvNodes[index].key {
		if childNum > 0 {
			btree.insert(pnode.childNodes[index], 0, kvnode) // 左孩子找
		} else {
			// 插入到结点中(头 或 中间)
			// 后移
			copy(pnode.kvNodes[index+1:], pnode.kvNodes[index:])
			// 插入
			pnode.kvNodes[index] = kvnode
			pnode.num++
			btree.insertCheck(pnode)
		}
	} else if kvnode.key > pnode.kvNodes[index].key {
		if index+1 < pnode.num {
			btree.insert(pnode, index+1, kvnode) // 后移
		} else {
			if childNum > 0 {
				btree.insert(pnode.childNodes[index+1], 0, kvnode) // 对最后一个成员的右子树遍历
			} else {
				// 插入到结点最后(一定为叶子结点)
				pnode.kvNodes[index+1] = kvnode
				pnode.num++
				btree.insertCheck(pnode)
			}
		}
	} else {
		pnode.kvNodes[index].value = kvnode.value
	}
}

// insertCheck 插入校验
func (btree *BTree) insertCheck(node *BNode) {
	if node.num < btree.order {
		return
	}

	// 分裂，提升结点
	mid := node.num / 2
	var (
		newKvNode    = newKvNode(node.kvNodes[mid].key, node.kvNodes[mid].value)
		newLeftNode  = newBNode(btree.order)
		newRightNode = newBNode(btree.order)
	)
	// 分裂
	if !node.isleaf {
		newLeftNode.isleaf = false
		newRightNode.isleaf = false
	}
	newLeftNode.num = mid
	newRightNode.num = btree.order - mid - 1
	newLeftNode.pindex = node.pindex
	newRightNode.pindex = node.pindex + 1
	copy(newLeftNode.kvNodes[:], node.kvNodes[:mid])
	copy(newRightNode.kvNodes[:], node.kvNodes[mid+1:])
	if len(node.childNodes) > 0 {
		newLeftNode.childNodes = append(newLeftNode.childNodes, node.childNodes[:mid+1]...)
		newRightNode.childNodes = append(newRightNode.childNodes, node.childNodes[mid+1:]...)
		// 重新赋父结点
		for i, c := range newLeftNode.childNodes {
			c.parent = newLeftNode
			c.pindex = i
		}
		for i, c := range newRightNode.childNodes {
			c.parent = newRightNode
			c.pindex = i
		}
	}

	// 提升
	if node == btree.root { // 根结点
		var newRoot = newBNode(btree.order)

		newLeftNode.pindex = 0
		newRightNode.pindex = 1
		newLeftNode.parent = newRoot
		newRightNode.parent = newRoot

		newRoot.kvNodes[0] = newKvNode
		newRoot.num = 1
		newRoot.childNodes = append(newRoot.childNodes, newLeftNode)
		newRoot.childNodes = append(newRoot.childNodes, newRightNode)
		newRoot.isleaf = false
		btree.root = newRoot
		node.freeBNode()
	} else { // 非根结点
		var (
			parent = node.parent
			i      = node.pindex // 插入点
		)
		node.freeBNode()
		newLeftNode.parent = parent
		newRightNode.parent = parent
		if i < parent.num { // 中间 or 头
			copy(parent.kvNodes[i+1:], parent.kvNodes[i:])
			parent.kvNodes[i] = newKvNode
			// 先删除再添加
			parent.deletePosNode(i)
			// 注意顺序
			parent.insertPosNode(newRightNode, i)
			parent.insertPosNode(newLeftNode, i)
			// 重新赋 parent.childNodes[].pindex
			for j := i + 2; j < len(parent.childNodes) && parent.childNodes[j] != nil; j++ {
				parent.childNodes[j].pindex += 1
			}
			parent.num++
		} else { // 最后
			parent.kvNodes[i] = newKvNode
			// 先删除再添加
			parent.deletePosNode(i)
			parent.childNodes = append(parent.childNodes, newLeftNode)
			parent.childNodes = append(parent.childNodes, newRightNode)
			parent.num++
		}
		btree.insertCheck(parent) // 继续递归处理
	}
}

// Delete 删除结点
func (btree *BTree) Delete(key int64) bool {
	if btree == nil || btree.root == nil {
		return false
	}
	return btree.delete(btree.root, 0, key)
}

// delete 找到删除点
func (btree *BTree) delete(node *BNode, index int, key int64) bool {
	childNum := len(node.childNodes)
	if key < node.kvNodes[index].key {
		if childNum > 0 {
			return btree.delete(node.childNodes[index], 0, key)
		}
		return false
	} else if key > node.kvNodes[index].key {
		if index+1 < node.num {
			return btree.delete(node, index+1, key)
		} else {
			if childNum > 0 {
				return btree.delete(node.childNodes[index+1], 0, key)
			}
		}
		return false
	} else {
		if node.isleaf {
			// 直接覆盖
			copy(node.kvNodes[index:], node.kvNodes[index+1:])
			node.num--
			btree.deleteCheck(node)
		} else {
			// 用后继结点替换
			delNode := node.childNodes[index+1]
			for !delNode.isleaf {
				delNode = delNode.childNodes[0]
			}
			node.kvNodes[index].key = delNode.kvNodes[0].key
			node.kvNodes[index].value = delNode.kvNodes[0].value
			copy(delNode.kvNodes[0:], delNode.kvNodes[1:])
			delNode.num--
			btree.deleteCheck(delNode)
		}
		return true
	}
}

// deleteCheck 删除校验
func (btree *BTree) deleteCheck(node *BNode) {
	if node == btree.root || node.num >= btree.minNum {
		if btree.root.num == 0 && len(btree.root.childNodes) > 0 {
			btree.root = btree.root.childNodes[0]
			btree.root.pindex = -1
			btree.root.parent = nil
		}
		return
	}
	var (
		leftSib     = node.getLeftSib()
		rightSib    = node.getRightSib()
		leftSibNum  = 0
		rightSibNum = 0
	)
	if leftSib != nil {
		leftSibNum = leftSib.num
	}
	if rightSib != nil {
		rightSibNum = rightSib.num
	}

	var parent = node.parent
	if leftSibNum <= btree.minNum && rightSibNum <= btree.minNum {
		// 需要合并结点
		var downKvn *KvNode = nil
		if leftSib != nil {
			downKvn = parent.kvNodes[leftSib.pindex]
			copy(parent.kvNodes[leftSib.pindex:], parent.kvNodes[leftSib.pindex+1:])
			leftSib.mergeNode(downKvn, node)
			for i := node.pindex + 1; i <= parent.num; i++ {
				// TODO: 合并后记得更改后面结点的pindex
				parent.childNodes[i].pindex--
			}
			parent.deletePosNode(node.pindex)
			node = nil
		} else {
			downKvn = parent.kvNodes[node.pindex]
			copy(parent.kvNodes[node.pindex:], parent.kvNodes[node.pindex+1:])
			node.mergeNode(downKvn, rightSib)
			for i := rightSib.pindex + 1; i <= parent.num; i++ {
				// TODO: 合并后记得更改后面结点的pindex
				parent.childNodes[i].pindex--
			}
			parent.deletePosNode(rightSib.pindex)
			rightSib = nil
		}
		parent.num--
		btree.deleteCheck(parent)
	} else if leftSibNum > btree.minNum {
		leftSib.moveToRight(parent, node)
	} else {
		rightSib.moveToLeft(parent, node)
	}
}

func (btree *BTree) PrintTreeInLog() {
	if btree.root == nil {
		return
	}
	queue := datastruct.NewQueue(10)
	queue.Push(btree.root)
	for queue.Size() != 0 {
		size := queue.Size()
		for i := 0; i < size; i++ {
			node, _ := queue.Pop().(*BNode)
			// // 打印父结点
			// if node.parent != nil {
			// 	fmt.Print("parent:", node.parent.kvNodes[0].key)
			// }
			// // 打印结点key的个数
			// fmt.Print("[", node.isleaf, node.num, "]")
			// // 打印孩子结点个数 和 pindex
			// fmt.Print("[", len(node.childNodes), node.pindex, "]")
			// // 打印左右兄弟
			// leftSib := node.getLeftSib()
			// rightSib := node.getRightSib()
			// if leftSib != nil {
			// 	fmt.Printf("[左兄弟: %v]", leftSib.kvNodes[leftSib.num-1].key)
			// }
			// if rightSib != nil {
			// 	fmt.Printf("[右兄弟: %v]", rightSib.kvNodes[0].key)
			// }
			len1 := int(node.num)
			for index, kv := range node.kvNodes {
				if kv == nil || index >= node.num {
					break
				}
				fmt.Printf("(%d, %d) ", kv.key, kv.value)
				if len(node.childNodes) != 0 {
					queue.Push(node.childNodes[index])
				}
				if index == len1-1 {
					if len(node.childNodes) != 0 {
						queue.Push(node.childNodes[index+1])
					}
				}
			}
			fmt.Print("| ")
		}
		fmt.Println()
	}
}
