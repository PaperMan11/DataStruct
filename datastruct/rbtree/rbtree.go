package rbtree

import (
	"log"
)

type RBTree struct {
	root *RBNode
}

// NewRBTree 新建红黑树
func NewRBTree() *RBTree {
	return &RBTree{root: nil}
}

// 左旋
func (rbtree *RBTree) rotateLeft(node *RBNode) {
	if tempNode, err := node.rotate(LEFTROTATE); err == nil {
		// 根节点可能发生改动
		if tempNode != nil {
			rbtree.root = tempNode
		}
	}
}

// 右旋
func (rbtree *RBTree) rotateRight(node *RBNode) {
	if tempNode, err := node.rotate(RIGHTROTATE); err == nil {
		// 根节点可能发生改动
		if tempNode != nil {
			rbtree.root = tempNode
		}
	} else {
		log.Println(err)
	}
}

// Insert 新增节点
func (rbtree *RBTree) Insert(data int64) {
	if rbtree.root == nil {
		rootNode := NewRBNode(data)
		rootNode.color = BLACK
		rbtree.root = rootNode
	} else {
		rbtree.insertNode(rbtree.root, data)
	}
}

// insertNode 找到插入位置
func (rbtree *RBTree) insertNode(pnode *RBNode, data int64) {
	if data <= pnode.value {
		if pnode.left != nil {
			rbtree.insertNode(pnode.left, data)
		} else {
			node := NewRBNode(data)
			node.parent = pnode
			pnode.left = node
			rbtree.insertCheck(node) // 插入校验
		}
	} else {
		if pnode.right != nil {
			rbtree.insertNode(pnode.right, data)
		} else {
			node := NewRBNode(data)
			node.parent = pnode
			pnode.right = node
			rbtree.insertCheck(node) // 插入校验
		}
	}
}

// insertCheck 红黑树插入规则判断
func (rbtree *RBTree) insertCheck(node *RBNode) {
	if node.parent == nil {
		// 根节点，改变颜色为黑色
		rbtree.root = node
		rbtree.root.color = BLACK
		return
	}

	if node.parent.color == BLACK {
		// TODO: 父节点为黑色，直接添加不做处理
	} else {
		if node.getUncle() != nil && node.getUncle().color == RED { // 叔叔节点为红色
			node.parent.color = BLACK
			node.getUncle().color = BLACK
			node.getGrandParent().color = RED
			rbtree.insertCheck(node.getGrandParent()) // 继续向上递归处理
		} else { // 叔叔节点为黑色或没有
			isleft := node == node.parent.left
			isParentLeft := node.parent == node.getGrandParent().left
			if isParentLeft && isleft { // 左左
				node.parent.color = BLACK
				grandParent := node.getGrandParent()
				grandParent.color = RED
				rbtree.rotateRight(grandParent)
			} else if isParentLeft && !isleft { // 左右
				rbtree.rotateLeft(node.parent)
				rbtree.rotateRight(node.parent)
				node.color = BLACK
				node.right.color = RED
			} else if !isParentLeft && !isleft { // 右右
				node.parent.color = BLACK
				grandParent := node.getGrandParent()
				grandParent.color = RED
				rbtree.rotateLeft(grandParent)
			} else if !isParentLeft && isleft { // 右左
				rbtree.rotateRight(node.parent)
				rbtree.rotateLeft(node.parent)
				node.color = BLACK
				node.left.color = RED
			}
		}
	}
}

// Delete 删除节点
func (rbtree *RBTree) Delete(data int64) bool {
	return rbtree.deleteNode(rbtree.root, data)
}

func (rbtree *RBTree) deleteNode(node *RBNode, data int64) bool {
	if data < node.value {
		if node.left != nil {
			return rbtree.deleteNode(node.left, data)
		}
		return false
	} else if data > node.value {
		if node.right != nil {
			return rbtree.deleteNode(node.right, data)
		}
		return false
	} else {
		// 找到 delNode 后继节点(右子树最小的节点)
		var delNode *RBNode
		if node.right == nil {
			delNode = node
		} else {
			r := node.right
			rl := r.left
			for rl != nil && rl.left != nil {
				r = rl
				rl = rl.left
			}
			if rl == nil {
				delNode = r
			} else {
				delNode = rl
			}
		}
		node.value = delNode.value
		rbtree.deleteCheck(delNode)
		return true
	}
}

func (rbtree *RBTree) deleteCheck(delNode *RBNode) {
	if delNode.parent == nil {
		rbtree.root.color = BLACK
		return
	}

	sibNode := delNode.getSibling() // 删除节点的兄弟节点
	isleft := delNode == delNode.parent.left
	if sibNode != nil && sibNode.color == RED { // 删除节点的兄弟节点为红色
		// 3.2
		sibNode.color = BLACK
		if isleft {
			if sibNode.left != nil {
				sibNode.left.color = RED
			}
			rbtree.rotateLeft(delNode.parent)
			delNode.parent.left = delNode.left
		} else {
			if sibNode.right != nil {
				sibNode.right.color = RED
			}
			rbtree.rotateRight(delNode.parent)
			delNode.parent.right = delNode.right
		}
		delNode.parent = nil
		delNode = nil
	} else { // 删除节点的兄弟节点为黑色或没有
		if (sibNode == nil) || (sibNode.left == nil && sibNode.right == nil) {
			// 3.1.4
			delNode.parent.color = BLACK
			if sibNode != nil {
				sibNode.color = RED
			}
			if isleft {
				delNode.parent.left = delNode.left
			} else {
				delNode.parent.right = delNode.right
			}
		} else if sibNode.right != nil && sibNode.left == nil {
			// 3.1.1
			if isleft {
				sibNode.color = delNode.parent.color
				delNode.parent.color = BLACK
				sibNode.right.color = BLACK
				rbtree.rotateLeft(delNode.parent)
				delNode.parent.left = delNode.left
			} else {
				sibNode.color = RED
				sibNode.right.color = BLACK
				rbtree.rotateLeft(sibNode)
				rbtree.deleteCheck(delNode)
			}
		} else if sibNode.right == nil && sibNode.left != nil {
			// 3.1.2
			if isleft {
				sibNode.color = RED
				sibNode.left.color = BLACK
				rbtree.rotateRight(sibNode)
				rbtree.deleteCheck(delNode)
			} else {
				sibNode.color = delNode.parent.color
				delNode.parent.color = BLACK
				sibNode.left.color = BLACK
				rbtree.rotateRight(delNode.parent)
				delNode.parent.right = delNode.right
			}
		} else if sibNode.right != nil && sibNode.left != nil {
			// 3.1.3
			sibNode.color = delNode.parent.color
			delNode.parent.color = BLACK
			if isleft {
				sibNode.right.color = BLACK
				rbtree.rotateLeft(delNode.parent)
				delNode.parent.left = delNode.left
			} else {
				sibNode.left.color = BLACK
				rbtree.rotateRight(delNode.parent)
				delNode.parent.right = delNode.right
			}
		}
		delNode.parent = nil
		delNode = nil
	}
}

// log输出树
func printTreeInLog(n *RBNode, front string) {
	if n != nil {
		var colorstr string
		if n.color == RED {
			colorstr = "红"
		} else {
			colorstr = "黑"
		}
		log.Printf(front+"%d,%s\n", n.value, colorstr)
		// if n.parent != nil {
		// 	fmt.Printf("parent:%d--", n.parent.value)
		// }
		// if n.left != nil {
		// 	fmt.Printf("left:%d--", n.left.value)
		// }
		// if n.right != nil {
		// 	fmt.Printf("right:%d", n.right.value)
		// }
		// fmt.Println()
		printTreeInLog(n.left, front+"-(l)|")

		printTreeInLog(n.right, front+"-(r)|")
	}
}
