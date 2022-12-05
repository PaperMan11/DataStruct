package datastruct

import "testing"

func TestQueue(t *testing.T) {
	q := NewQueue(3)
	q.PrintInLog()
	for i := 0; i < 10; i++ {
		q.Push(i)
	}
	q.PrintInLog()
	q.Pop()
	q.PrintInLog()
}
