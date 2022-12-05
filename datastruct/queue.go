package datastruct

import "fmt"

type Queue struct {
	arr      []interface{}
	size     int
	capacity int
}

func NewQueue(capacity int) *Queue {
	return &Queue{
		arr:      make([]interface{}, capacity),
		size:     0,
		capacity: capacity,
	}
}

func (q *Queue) Push(v interface{}) bool {
	if q == nil {
		return false
	}
	if q.size >= q.capacity {
		q.capacity = q.capacity * 2
		newArr := make([]interface{}, q.capacity)
		copy(newArr[:], q.arr[:])
		q.arr = newArr
	}
	q.arr[q.size] = v
	q.size++
	return true
}

func (q *Queue) Pop() interface{} {
	if q == nil || q.size == 0 {
		return nil
	}
	res := q.arr[0]
	copy(q.arr[0:], q.arr[1:])
	q.size--
	return res
}

func (q *Queue) Size() int {
	if q == nil || q.arr == nil {
		return 0
	}
	return q.size
}

func (q *Queue) PrintInLog() {
	if q == nil || q.arr == nil {
		return
	}
	fmt.Println(q.arr...)
}
