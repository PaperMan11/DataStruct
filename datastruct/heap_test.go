package datastruct

import "testing"

func TestHeap(t *testing.T) {
	heap := NewHeap(make([]int, 10))
	for i := 0; i < 10; i++ {
		heap.Push(i)
	}
	heap.Println()
	heap.Pop()
	heap.Println()
}
