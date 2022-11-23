package rbtree

import "errors"

const (
	RED         = iota  // 红色
	BLACK               // 黑色
	LEFTROTATE  = true  // 左旋
	RIGHTROTATE = false // 右旋
)

// RBNode 红黑树节点
type RBNode struct {
	value               int64
	color               int
	left, right, parent *RBNode
}

// NewRBNode 新建节点
func NewRBNode(value int64) *RBNode {
	return &RBNode{
		value: value,
		color: RED,
	}
}

// getGrandParent 获取父节点的父节点
func (rbnode *RBNode) getGrandParent() *RBNode {
	if rbnode == nil || rbnode.parent == nil {
		return nil
	}
	return rbnode.parent.parent
}

// getSibling 获取兄弟节点
func (rbnode *RBNode) getSibling() *RBNode {
	if rbnode == nil || rbnode.parent == nil {
		return nil
	}
	if rbnode == rbnode.parent.left {
		return rbnode.parent.right
	} else {
		return rbnode.parent.left
	}
}

// getUncle 获取父节点的兄弟节点
func (rbnode *RBNode) getUncle() *RBNode {
	if rbnode.getGrandParent() == nil {
		return nil
	}
	grandParent := rbnode.getGrandParent()
	if rbnode.parent == grandParent.left {
		return grandParent.right
	} else {
		return grandParent.left
	}
}

// rotate 左旋/右旋(true/false)
// 若根节点发生了改变则返回根节点
func (rbnode *RBNode) rotate(isRotateLeft bool) (*RBNode, error) {
	var root *RBNode
	if rbnode == nil {
		return root, nil
	}
	if !isRotateLeft && rbnode.left == nil {
		return root, errors.New("右旋左节点不能为空")
	} else if isRotateLeft && rbnode.right == nil {
		return root, errors.New("左旋右节点不能为空")
	}

	parent := rbnode.parent
	var isleft bool // 用判断在哪个子树进行的操作
	if parent != nil {
		isleft = rbnode == parent.left
	}

	if isRotateLeft { // 左旋
		grandson := rbnode.right.left
		rbnode.right.left = rbnode
		rbnode.parent = rbnode.right
		rbnode.right = grandson
		if grandson != nil {
			grandson.parent = rbnode
		}
	} else { // 右旋
		grandson := rbnode.left.right
		rbnode.left.right = rbnode
		rbnode.parent = rbnode.left
		rbnode.left = grandson
		if grandson != nil {
			grandson.parent = rbnode
		}
	}

	if parent == nil { // 旋转完毕找到根节点
		rbnode.parent.parent = nil
		root = rbnode.parent
	} else {
		if isleft {
			parent.left = rbnode.parent
		} else {
			parent.right = rbnode.parent
		}
		rbnode.parent.parent = parent
	}
	return root, nil
}
