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
		root := NewBNode(btree.order)
		root.kvNodes[0] = kvnode
		root.num++
		btree.root = root
		return
	}
	btree.insert(btree.root, 0, kvnode)
}

// insert 找到插入点
func (btree *BTree) insert(pnode *BNode, index int, kvnode *KvNode) {
	if kvnode.key < pnode.kvNodes[index].key {
		if pnode.kvNodes[index].left != nil {
			btree.insert(pnode.kvNodes[index].left, 0, kvnode) // 左孩子找
		} else {
			// 插入到结点中(头 或 中间)
			kvnode.right = pnode.kvNodes[index].left
			if pnode.kvNodes[index].left != nil {
				pnode.kvNodes[index].left.lparent = kvnode
			}
			if index != 0 { // 中间
				kvnode.left = pnode.kvNodes[index-1].right
				if pnode.kvNodes[index-1].right != nil {
					pnode.kvNodes[index-1].right.rparent = kvnode
				}
			}

			newKvSlice := make([]*KvNode, btree.order)
			copy(newKvSlice[:], pnode.kvNodes[:index])
			newKvSlice[index] = kvnode
			copy(newKvSlice[index+1:], pnode.kvNodes[index:])

			pnode.kvNodes = newKvSlice
			pnode.num++
			btree.insertCheck(pnode)
		}
	} else if kvnode.key > pnode.kvNodes[index].key {
		if pnode.kvNodes[index+1] != nil {
			btree.insert(pnode, index+1, kvnode) // 后移 index+1 < node.num
		} else {
			if pnode.kvNodes[index].right != nil { // 对最后一个成员的右子树遍历
				btree.insert(pnode.kvNodes[index].right, 0, kvnode)
			} else {
				// 插入到结点最后(插入点一定没有右孩子)
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
	if int(node.num) < int(btree.order) {
		return
	}

	// 分裂
	mid := int(node.num) / 2
	midKvNode := node.kvNodes[mid]
	// 结点提升了，对应的子结点的父结点也发生变化
	if midKvNode.left != nil {
		midKvNode.left.rparent = nil
	}
	if midKvNode.right != nil {
		midKvNode.right.lparent = nil
	}
	// fmt.Println(node.kvNodes, midKvNode, node.num)

	newKvNode := newKvNode(midKvNode.key, midKvNode.value)
	newLeftNode := NewBNode(btree.order)
	copy(newLeftNode.kvNodes[:], node.kvNodes[:mid])
	newLeftNode.num = mid // TODO: 分裂新结点需要获取到结点中 key 的数量(num)
	newRightNode := NewBNode(btree.order)
	copy(newRightNode.kvNodes[:], node.kvNodes[mid+1:])
	newRightNode.num = btree.order - mid - 1
	if !node.isleaf {
		newLeftNode.isleaf = false
		newRightNode.isleaf = false
	}

	newKvNode.left = newLeftNode
	newLeftNode.rparent = newKvNode
	newKvNode.right = newRightNode
	newRightNode.lparent = newKvNode

	// 重新赋 父结点
	for _, kvn := range newLeftNode.kvNodes {
		if kvn == nil {
			break
		}
		if kvn.left != nil {
			kvn.left.parent = newLeftNode
		}
		if kvn.right != nil {
			kvn.right.parent = newLeftNode
		}
	}

	for _, kvn := range newRightNode.kvNodes {
		if kvn == nil {
			break
		}
		if kvn.left != nil {
			kvn.left.parent = newRightNode
		}
		if kvn.right != nil {
			kvn.right.parent = newRightNode
		}
	}

	if node == btree.root { // 根结点
		newNode := NewBNode(btree.order)
		newNode.kvNodes[0] = newKvNode

		newLeftNode.parent = newNode
		newRightNode.parent = newNode

		btree.root = newNode
		btree.root.isleaf = false // 新 root 肯定不是叶子结点
		newNode.num = 1
		node = nil
	} else { // 非根结点
		parent := node.parent
		node = nil
		newLeftNode.parent = parent
		newRightNode.parent = parent
		for i, v := range parent.kvNodes { // 找到插入点
			if v == nil { // 插入最后
				parent.kvNodes[i-1].right = newKvNode.left
				newKvNode.left.lparent = parent.kvNodes[i-1]
				parent.kvNodes[i] = newKvNode
				parent.num++
				break
			} else if v.key > newKvNode.key { // 插入到第一个或中间
				parent.kvNodes[i].left = newKvNode.right
				newKvNode.right.rparent = parent.kvNodes[i]
				if i != 0 { // 中间
					parent.kvNodes[i-1].right = newKvNode.left
					newKvNode.left.lparent = parent.kvNodes[i-1]
				}

				newParentSilce := make([]*KvNode, btree.order)
				copy(newParentSilce[:], parent.kvNodes[:i])
				newParentSilce[i] = newKvNode
				copy(newParentSilce[i+1:], parent.kvNodes[i:])

				parent.kvNodes = newParentSilce
				parent.num++
				break
			}
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
	if key < node.kvNodes[index].key {
		if node.kvNodes[index].left != nil {
			return btree.delete(node.kvNodes[index].left, 0, key)
		}
		return false
	} else if key > node.kvNodes[index].key {
		if node.kvNodes[index+1] != nil { // 后移 index+1 < node.num
			return btree.delete(node, index+1, key)
		} else {
			if node.kvNodes[index].right != nil { // 对最后一个成员的右子树遍历
				return btree.delete(node.kvNodes[index].right, 0, key)
			}
		}
		return false
	} else {
		// 用后继结点替换
		var delBNode *BNode
		if node.isleaf {
			delBNode = node
			// fmt.Println(delBNode.kvNodes[index], delBNode.num)
			btree.deleteCheck(delBNode, index, true)
		} else {
			// B 树性质每个kvn肯定有左右孩子（非叶子）。
			delBNode = node.kvNodes[index].right
			for !delBNode.isleaf {
				delBNode = delBNode.kvNodes[0].left
			}
			node.kvNodes[index].key = delBNode.kvNodes[0].key
			node.kvNodes[index].value = delBNode.kvNodes[0].value
			// fmt.Println(delBNode.kvNodes[0], delBNode.num)
			btree.deleteCheck(delBNode, 0, true)
		}
		return true
	}
}

// deleteCheck 删除校验
//		@delBNode: 当前指向的结点
// 		@index: 索引
// 		@flag: 是否执行删除操作(false表示调整阶段)
func (btree *BTree) deleteCheck(delBNode *BNode, index int, flag bool) {
	if flag {
		if delBNode.num > btree.minNum || delBNode == btree.root {
			// 直接删除
			copy(delBNode.kvNodes[index:], delBNode.kvNodes[index+1:])
			delBNode.kvNodes[delBNode.num-1] = nil
			delBNode.num--
			return
		}
	}
	// deleteCheck 递归结束
	if delBNode.num > btree.minNum || (delBNode == btree.root && !flag) {
		return
	}

	var (
		leftSib     *BNode = delBNode.GetLeftSib()
		rightSib    *BNode = delBNode.GetRightSib()
		leftSibNum         = 0
		rightSibNum        = 0
	)
	if leftSib != nil {
		leftSibNum = leftSib.num
	}
	if rightSib != nil {
		rightSibNum = rightSib.num
	}

	var parent = delBNode.parent
	if leftSibNum <= btree.minNum && rightSibNum <= btree.minNum {
		// 需要合并结点
		var (
			nc_kvn   *KvNode                         // 下移 kvnode 结点
			nc_bnode *BNode  = NewBNode(btree.order) // 新的 bnode 子结点
		)
		// 覆盖删除的结点
		if flag {
			copy(delBNode.kvNodes[index:], delBNode.kvNodes[index+1:])
			delBNode.kvNodes[delBNode.num-1] = nil
			delBNode.num--
		}
		if leftSib != nil {
			nc_kvn = delBNode.lparent
			// 找到下移结点的索引并删除
			for i, kvn := range parent.kvNodes {
				if kvn != nc_kvn {
					continue
				}
				// 删除
				if i != 0 {
					parent.kvNodes[i-1].right = nc_bnode
					nc_bnode.lparent = parent.kvNodes[i-1]
				}
				if i != parent.num-1 {
					parent.kvNodes[i+1].left = nc_bnode
					nc_bnode.rparent = parent.kvNodes[i+1]
				}
				nc_bnode.parent = parent
				// 覆盖
				copy(parent.kvNodes[i:], parent.kvNodes[i+1:])
				parent.kvNodes[parent.num-1] = nil
				parent.num--
				break
			}

			// 合并
			copy(nc_bnode.kvNodes[:leftSibNum], leftSib.kvNodes[:leftSibNum])
			nc_bnode.kvNodes[leftSibNum] = nc_kvn
			copy(nc_bnode.kvNodes[leftSibNum+1:], delBNode.kvNodes[:delBNode.num])
			nc_bnode.num = leftSibNum + delBNode.num + 1
			nc_bnode.isleaf = delBNode.isleaf
			if leftSibNum != 0 {
				nc_kvn.left = nc_bnode.kvNodes[leftSibNum-1].right
				if nc_bnode.kvNodes[leftSibNum-1].right != nil {
					nc_bnode.kvNodes[leftSibNum-1].right.rparent = nc_kvn
				}
			}
			if leftSibNum != nc_bnode.num-1 {
				nc_kvn.right = nc_bnode.kvNodes[leftSibNum+1].left
				if nc_bnode.kvNodes[leftSibNum+1].left != nil {
					nc_bnode.kvNodes[leftSibNum+1].left.lparent = nc_kvn
				}
			}

			leftSib.parent = nil
			leftSib.lparent = nil
			leftSib.rparent = nil
			leftSib.kvNodes = nil
			leftSib = nil
		} else {
			nc_kvn = delBNode.rparent
			// 找到下移结点的索引并删除
			for i, kvn := range parent.kvNodes {
				if kvn != nc_kvn {
					continue
				}
				// 删除
				if i != 0 {
					parent.kvNodes[i-1].right = nc_bnode
					nc_bnode.lparent = parent.kvNodes[i-1]
				}
				if i != parent.num-1 {
					parent.kvNodes[i+1].left = nc_bnode
					nc_bnode.rparent = parent.kvNodes[i+1]
				}
				nc_bnode.parent = parent
				// 覆盖
				copy(parent.kvNodes[i:], parent.kvNodes[i+1:])
				parent.kvNodes[parent.num-1] = nil
				parent.num--
				break
			}

			// 合并
			copy(nc_bnode.kvNodes[:delBNode.num], delBNode.kvNodes[:delBNode.num])
			nc_bnode.kvNodes[delBNode.num] = nc_kvn
			copy(nc_bnode.kvNodes[delBNode.num+1:], rightSib.kvNodes[:rightSibNum])
			nc_bnode.num = rightSibNum + delBNode.num + 1
			nc_bnode.isleaf = delBNode.isleaf
			if delBNode.num != 0 {
				nc_kvn.left = nc_bnode.kvNodes[delBNode.num-1].right
				if nc_bnode.kvNodes[delBNode.num-1].right != nil {
					nc_bnode.kvNodes[delBNode.num-1].right.rparent = nc_kvn
				}
			}
			if delBNode.num != nc_bnode.num-1 {
				nc_kvn.right = nc_bnode.kvNodes[delBNode.num+1].left
				if nc_bnode.kvNodes[delBNode.num+1].left != nil {
					nc_bnode.kvNodes[delBNode.num+1].left.lparent = nc_kvn
				}
			}

			rightSib.parent = nil
			rightSib.lparent = nil
			rightSib.rparent = nil
			rightSib.kvNodes = nil
			rightSib = nil
		}
		delBNode.parent = nil
		delBNode.lparent = nil
		delBNode.rparent = nil
		delBNode.kvNodes = nil
		delBNode = nil

		// 重新赋 父结点
		for _, kvn := range nc_bnode.kvNodes {
			if kvn == nil {
				break
			}
			if kvn.left != nil {
				kvn.left.parent = nc_bnode
			}
			if kvn.right != nil {
				kvn.right.parent = nc_bnode
			}
		}

		if parent.num == 0 {
			btree.root = nc_bnode
			btree.root.parent = nil
			parent = nil
			return
		}
		btree.deleteCheck(parent, 0, false) // 继续向上递归调整
	} else if leftSibNum > btree.minNum {
		nc_kvn := delBNode.lparent              // 下移 kvnode 结点
		np_kvn := leftSib.kvNodes[leftSibNum-1] // 上移 kvnode 结点
		nc_left := np_kvn.right                 // nc_kvn 的左子结点

		if np_kvn.left != nil {
			np_kvn.left.rparent = nil
		}
		np_kvn.left = nc_kvn.left
		nc_kvn.left.rparent = np_kvn
		np_kvn.right = nc_kvn.right
		nc_kvn.right.lparent = np_kvn

		nc_kvn.left = nc_left
		if nc_left != nil {
			nc_left.rparent = nc_kvn
			nc_left.lparent = nil
		}
		nc_kvn.right = delBNode.kvNodes[0].left
		if delBNode.kvNodes[0].left != nil {
			delBNode.kvNodes[0].left.lparent = nc_kvn
		}
		// 上移
		for i, kvn := range parent.kvNodes {
			if kvn == nc_kvn { // 找到 np_kvn 插入点替换
				parent.kvNodes[i] = np_kvn
				break
			}
		}
		leftSib.kvNodes[leftSibNum-1] = nil
		leftSib.num--
		// 下移
		if !flag { // 调整阶段不需要删除任何结点
			delBNode.num++
			copy(delBNode.kvNodes[1:], delBNode.kvNodes[:]) // 全部后移
		} else {
			copy(delBNode.kvNodes[1:index+1], delBNode.kvNodes[:index]) // 直接覆盖了删除点
		}
		delBNode.kvNodes[0] = nc_kvn
	} else {
		nc_kvn := delBNode.rparent    // 下移 kvnode 结点
		np_kvn := rightSib.kvNodes[0] // 上移 kvnode 结点
		nc_right := np_kvn.left       // nc_kvn 的右子结点

		np_kvn.left = nc_kvn.left
		nc_kvn.left.rparent = np_kvn
		if np_kvn.right != nil {
			np_kvn.right.lparent = nil
		}
		np_kvn.right = nc_kvn.right
		nc_kvn.right.lparent = np_kvn

		nc_kvn.right = nc_right
		if nc_right != nil {
			nc_right.lparent = nc_kvn
			nc_right.rparent = nil
		}
		nc_kvn.left = delBNode.kvNodes[delBNode.num-1].right
		if delBNode.kvNodes[delBNode.num-1].right != nil {
			delBNode.kvNodes[delBNode.num-1].right.rparent = nc_kvn
		}
		// 上移
		for i, kvn := range parent.kvNodes {
			if kvn == nc_kvn { // 找到 np_kvn 插入点替换
				parent.kvNodes[i] = np_kvn
				break
			}
		}
		copy(rightSib.kvNodes[0:], rightSib.kvNodes[1:])
		rightSib.kvNodes[rightSibNum-1] = nil
		rightSib.num--
		// 下移
		if !flag { // 调整阶段不需要删除任何结点
			delBNode.num++
		} else {
			copy(delBNode.kvNodes[index:], delBNode.kvNodes[index+1:]) // 直接覆盖了删除点
		}
		delBNode.kvNodes[delBNode.num-1] = nc_kvn
	}
}

// PrintTreeInLog 层序遍历
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
			// // 打印左右父结点
			// if node.lparent != nil {
			// 	fmt.Printf("[lp: %v]", node.lparent.key)
			// }
			// if node.rparent != nil {
			// 	fmt.Printf("[rp: %v]", node.rparent.key)
			// }
			// // 打印左右兄弟
			// leftSib := node.GetLeftSib()
			// rightSib := node.GetRightSib()
			// if leftSib != nil {
			// 	fmt.Printf("[左兄弟: %v]", leftSib.kvNodes[leftSib.num-1].key)
			// }
			// if rightSib != nil {
			// 	fmt.Printf("[右兄弟: %v]", rightSib.kvNodes[0].key)
			// }
			len := int(node.num)
			for index, kv := range node.kvNodes {
				if kv == nil {
					break
				}
				fmt.Printf("(%d, %d) ", kv.key, kv.value)
				if kv.left != nil {
					queue.Push(kv.left)
				}
				if index == len-1 {
					if kv.right != nil {
						queue.Push(kv.right)
					}
				}
			}
			fmt.Print("| ")
		}
		fmt.Println()
	}
}
