package btree

import (
	"fmt"
	"testing"
)

func TestSlicInsert(t *testing.T) {
	b1 := &BNode{kvNodes: make([]*KvNode, 0)}
	b1.kvNodes = append(b1.kvNodes, &KvNode{key: 1, value: 1})
	b2 := &BNode{kvNodes: make([]*KvNode, 0)}
	b2.kvNodes = append(b2.kvNodes, &KvNode{key: 2, value: 2})
	b3 := &BNode{kvNodes: make([]*KvNode, 0)}
	b3.kvNodes = append(b3.kvNodes, &KvNode{key: 3, value: 1})
	B := newBNode(5)
	B.childNodes = append(B.childNodes, b1)
	fmt.Println(len(B.childNodes), B.childNodes[0].kvNodes[0].key)
	B.insertPosNode(b2, 0)
	fmt.Println(len(B.childNodes), B.childNodes[0].kvNodes[0].key)
	B.insertPosNode(b3, 0)
	fmt.Println(len(B.childNodes), B.childNodes[0].kvNodes[0].key, B.childNodes[1].kvNodes[0].key, B.childNodes[2].kvNodes[0].key)
}
