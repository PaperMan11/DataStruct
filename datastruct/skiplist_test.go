package datastruct

import (
	"fmt"
	"testing"
)

func TestSkipList(t *testing.T) {
	sl := NewSkipLinked[string](20, 1)
	for i := 1; i < 10; i++ {
		kv := fmt.Sprintf("%d", i)
		sl.Add(kv, float64(i+10), kv)
	}
	sl.Println()
	if b, i := sl.Search("1"); b {
		fmt.Println(i.(*Node[string]).Val)
	}

	sl.Erase("1")

	fmt.Println("-----------------------------")

	sl.Println()
	if b, i := sl.Search("1"); b {
		fmt.Println(i)
	} else {
		fmt.Println("not found")
	}
}
