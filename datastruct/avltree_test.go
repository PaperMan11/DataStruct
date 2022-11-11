package datastruct

import (
	"fmt"
	"testing"
)

func TestAvlTree(t *testing.T) {
	var tree AvlTree
	for i := 1; i <= 10; i++ {
		tree.treeRoot = tree.InsertTreeNode(tree.treeRoot, ElemType(i))
	}
	tree.PreTraversal(tree.treeRoot)
	fmt.Println()
	tree.InTraversal(tree.treeRoot)
	fmt.Println()
	tree.SqcTraversal(tree.treeRoot)
	tree.DeleteTreeNode(tree.treeRoot, ElemType(5))
	tree.DeleteTreeNode(tree.treeRoot, ElemType(6))
	tree.DeleteTreeNode(tree.treeRoot, ElemType(7))
	tree.DeleteTreeNode(tree.treeRoot, ElemType(4))
	fmt.Println()
	tree.SqcTraversal(tree.treeRoot)
}
