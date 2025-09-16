package datastruct

import (
	"fmt"
)

/***************AVL Tree****************/

type ElemType int
type TreeNode struct {
	height int
	data   ElemType
	left   *TreeNode
	right  *TreeNode
}
type AvlTree struct {
	treeRoot *TreeNode
}

func (tree *AvlTree) nodeHeight(node *TreeNode) int {
	if node == nil {
		return -1
	}
	return node.height
}

func (tree *AvlTree) nodeBf(node *TreeNode) int {
	return (tree.nodeHeight(node.left) - tree.nodeHeight(node.right))
}

func (tree *AvlTree) rightRotate(node *TreeNode) (l *TreeNode) {
	l = node.left
	node.left = l.right
	l.right = node

	/*if tree.nodeHeight(node.left) < tree.nodeHeight(node.right) {
		node.height = tree.nodeHeight(node.right) + 1
	} else {
		node.height = tree.nodeHeight(node.left) + 1
	}*/
	node.height = max(tree.nodeHeight(node.left), tree.nodeHeight(node.right)) + 1

	/*if tree.nodeHeight(l.left) < tree.nodeHeight(l.right) {
		l.height = tree.nodeHeight(l.right) + 1
	} else {
		l.height = tree.nodeHeight(l.left) + 1
	}*/
	l.height = max(tree.nodeHeight(node.left), tree.nodeHeight(node.right)) + 1

	return l
}

func (tree *AvlTree) leftRotate(node *TreeNode) (r *TreeNode) {
	r = node.right
	node.right = r.left
	r.left = node

	if tree.nodeHeight(node.left) < tree.nodeHeight(node.right) {
		node.height = tree.nodeHeight(node.right) + 1
	} else {
		node.height = tree.nodeHeight(node.left) + 1
	}

	if tree.nodeHeight(r.left) < tree.nodeHeight(r.right) {
		r.height = tree.nodeHeight(r.right) + 1
	} else {
		r.height = tree.nodeHeight(r.left) + 1
	}

	return r
}

// InsertTreeNode 插入节点
func (tree *AvlTree) InsertTreeNode(node *TreeNode, e ElemType) *TreeNode {
	if node == nil {
		newTreeNode := new(TreeNode)
		newTreeNode.right = nil
		newTreeNode.left = nil
		newTreeNode.height = 0
		newTreeNode.data = e
		return newTreeNode
	} else {
		if e < node.data {
			node.left = tree.InsertTreeNode(node.left, e)
			if tree.nodeBf(node) == 2 { // 插入完后，可能需要调整树结构
				if tree.nodeBf(node.left) < 0 {
					node.left = tree.leftRotate(node.left)
					node = tree.rightRotate(node.right)
				} else {
					node = tree.rightRotate(node)
				}
			}
		} else if e > node.data {
			node.right = tree.InsertTreeNode(node.right, e)
			if tree.nodeBf(node) == -2 {
				if tree.nodeBf(node.right) > 0 {
					node.right = tree.rightRotate(node.right)
					node = tree.leftRotate(node)
				} else {
					node = tree.leftRotate(node)
				}
			}
		} else {
			fmt.Println("InsertTreeNode failed")
			return node
		}
	}
	if node != nil { // 调整完树结构后，节点的高度可能发生改变
		if (node.right != nil && node.left == nil) ||
			(node.right.height > node.left.height) {
			node.height = node.right.height + 1
		} else if (node.right == nil && node.left != nil) ||
			(node.left.height > node.right.height) {
			node.height = node.left.height + 1
		}
	}
	return node
}

// DeleteTreeNode 删除节点
func (tree *AvlTree) DeleteTreeNode(node *TreeNode, e ElemType) *TreeNode {
	if node == nil {
		return nil
	}
	if e < node.data {
		node.left = tree.DeleteTreeNode(node.left, e)
		if tree.nodeBf(node) == -2 { // 删除完后，可能需要调整树结构
			if tree.nodeBf(node.right) > 0 {
				node.right = tree.rightRotate(node.right)
				node = tree.leftRotate(node)
			} else {
				node = tree.leftRotate(node)
			}
		}
	} else if e > node.data {
		node.right = tree.DeleteTreeNode(node.right, e)
		if tree.nodeBf(node) == 2 {
			if tree.nodeBf(node.left) < 0 {
				node.left = tree.leftRotate(node.left)
				node = tree.rightRotate(node.right)
			} else {
				node = tree.rightRotate(node)
			}
		}
	} else {
		if node.left == nil && node.right == nil {
			return nil
		} else if node.left == nil {
			return node.right
		} else {
			l := node.left
			lr := l.right
			for lr != nil && lr.right != nil {
				l = lr
				lr = lr.right
			}
			if lr == nil {
				node.data = l.data
				node.left = l.left
			} else {
				node.data = lr.data
				l.right = nil
			}
		}
	}
	if node != nil { // 调整完树结构后，节点的高度可能发生改变
		if (node.right != nil && node.left == nil) ||
			(node.right.height > node.left.height) {
			node.height = node.right.height + 1
		} else if (node.right == nil && node.left != nil) ||
			(node.left.height > node.right.height) {
			node.height = node.left.height + 1
		}
	}
	return node
}

// PreTraversal 先序遍历
func (tree *AvlTree) PreTraversal(root *TreeNode) {
	if root != nil {
		fmt.Print(root.data, " ")
		tree.PreTraversal(root.left)
		tree.PreTraversal(root.right)
	}
}

// InTraversal 中序遍历
func (tree *AvlTree) InTraversal(root *TreeNode) {
	if root != nil {
		tree.InTraversal(root.left)
		fmt.Print(root.data, " ")
		tree.InTraversal(root.right)
	}
}

// SqcTraversal 层序遍历
func (tree *AvlTree) SqcTraversal(root *TreeNode) {
	if root == nil {
		return
	}
	queue := NewQueue(10)
	queue.Push(root)
	for queue.Size() != 0 {
		size := queue.Size()
		for i := 0; i < size; i++ {
			val := queue.Pop().(*TreeNode)
			fmt.Print(val.data, " ")
			if val.left != nil {
				queue.Push(val.left)
			}
			if val.right != nil {
				queue.Push(val.right)
			}
		}
	}
}

func (tree *AvlTree) IsBalanced() bool {
	return tree.isBalancedTree(tree.treeRoot)
}

func (tree *AvlTree) isBalancedTree(root *TreeNode) bool {
	if root == nil {
		return true
	}
	if tree.nodeBf(root) > 1 || tree.nodeBf(root) < -1 {
		return false
	}
	return tree.isBalancedTree(root.left) && tree.isBalancedTree(root.right)
}
