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
	fmt.Println()
	fmt.Println("is balanced", tree.IsBalanced())
	tree.DeleteTreeNode(tree.treeRoot, ElemType(5))
	fmt.Println("is balanced", tree.IsBalanced())
	tree.DeleteTreeNode(tree.treeRoot, ElemType(6))
	tree.DeleteTreeNode(tree.treeRoot, ElemType(7))
	tree.DeleteTreeNode(tree.treeRoot, ElemType(4))
	fmt.Println("is balanced", tree.IsBalanced())
	fmt.Println()
	tree.SqcTraversal(tree.treeRoot)
}
