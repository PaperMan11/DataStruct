package btree

import (
	"fmt"
	"testing"
)

func TestBTree(t *testing.T) {
	btree := NewBTree(5)
	// for i := 0; i < 11; i++ {
	// 	btree.Insert(int64(i), i)
	// }
	btree.Insert(39, 39)
	btree.Insert(22, 22)
	btree.Insert(97, 97)
	btree.Insert(41, 41)
	btree.Insert(53, 53)
	btree.Insert(13, 13)
	btree.Insert(21, 21)
	btree.Insert(40, 40)
	btree.Insert(30, 30)
	btree.Insert(27, 27)
	btree.Insert(33, 33)
	btree.Insert(36, 36)
	btree.Insert(35, 35)
	btree.Insert(34, 34)
	btree.Insert(24, 24)
	btree.Insert(29, 29)
	btree.Insert(26, 26)
	btree.Insert(17, 17)
	btree.Insert(28, 28)
	btree.Insert(23, 23)
	btree.Insert(31, 31)
	btree.Insert(32, 32)
	btree.PrintTreeInLog()

	fmt.Println("--------------删除后-----------------")
	btree.Delete(39)
	btree.Delete(22)
	btree.Delete(97)
	btree.Delete(41)
	btree.Delete(53)
	// btree.Delete(13)
	// btree.Delete(21)
	// btree.Delete(40)
	// btree.Delete(30)
	// btree.Delete(27)
	// btree.Delete(33)
	// btree.Delete(36)
	// btree.Delete(35)
	// btree.Delete(34)
	// btree.Delete(24)
	// btree.Delete(29)
	// btree.Delete(26)
	// btree.Delete(17)
	// btree.Delete(28)
	// btree.Delete(23)
	// btree.Delete(31)
	// btree.Delete(32)
	btree.PrintTreeInLog()
	fmt.Println("--------------添加后-----------------")
	// 添加重复的 key 直接修改 value
	btree.Insert(39, 39)
	btree.Insert(22, 22)
	btree.Insert(97, 97)
	btree.Insert(41, 41)
	btree.Insert(53, 53)
	btree.Insert(13, -1)
	btree.Insert(21, -1)
	btree.Insert(40, -1)
	btree.PrintTreeInLog()
}

func TestSlice(t *testing.T) {
	s1 := make([]int, 5)
	s2 := make([]int, 5)
	s1[0] = 7
	s2[0] = 1
	s2[1] = 2
	s2[2] = 3
	s2[3] = 4
	s2[4] = 5
	copy(s1[:], s2[:0])
	copy(s1[1:], s2[0:])
	fmt.Println(s1)
	copy(s1[4:], s1[5:])
	fmt.Println(s1)
}

