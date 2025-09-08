package datastruct

type IHeap[T any] interface {
	Peek() T
	GetSize() int
	IsEmpty() bool
	IsFull() bool
	Insert(x T)
	Extract() T
}

type c[U any] func(i, j U) int

type Heap2[T any] struct {
	capacity   int
	size       int
	array      []T
	comparable c[T]
}

var _ IHeap[int] = (*Heap2[int])(nil)

func NewHeap2[T any](capacity int, comparable c[T]) *Heap2[T] {
	return &Heap2[T]{
		capacity:   capacity,
		array:      make([]T, 0, capacity),
		comparable: comparable,
		size:       0,
	}
}

func (*Heap2[T]) parent(i int) int {
	return (i - 1) / 2
}

func (*Heap2[T]) leftChild(i int) int {
	return 2*i + 1
}

func (*Heap2[T]) rightChild(i int) int {
	return 2*i + 2
}

func (h *Heap2[T]) swap(i, j int) {
	h.array[i], h.array[j] = h.array[j], h.array[i]
}

func (h *Heap2[T]) heapifyUp(i int) {
	curIndex := i
	for curIndex > 0 && h.comparable(h.array[h.parent(curIndex)], h.array[curIndex]) > 0 {
		h.swap(curIndex, h.parent(curIndex))
		curIndex = h.parent(curIndex)
	}
}

func (h *Heap2[T]) heapifyDown(i int) {
	curIndex := i
	for {
		leftChild := h.leftChild(curIndex)
		rightChild := h.rightChild(curIndex)
		smallest := curIndex
		if leftChild < h.size && h.comparable(h.array[leftChild], h.array[smallest]) < 0 {
			smallest = leftChild
		}
		if rightChild < h.size && h.comparable(h.array[rightChild], h.array[smallest]) < 0 {
			smallest = rightChild
		}
		if curIndex == smallest {
			break
		}
		h.swap(curIndex, smallest)
		curIndex = smallest
	}
}

func (h *Heap2[T]) Peek() T {
	if h.size == 0 {
		var zero T
		return zero
	}
	return h.array[0]
}

func (h *Heap2[T]) GetSize() int {
	return h.size
}

func (h *Heap2[T]) IsEmpty() bool {
	return h.size == 0
}

func (h *Heap2[T]) IsFull() bool {
	return h.size == h.capacity
}

func (h *Heap2[T]) Insert(x T) {
	if h.IsFull() {
		return
	}
	h.size++
	h.array = append(h.array, x)
	h.heapifyUp(h.size - 1)
}

func (h *Heap2[T]) Extract() T {
	if h.IsEmpty() {
		var zero T
		return zero
	}
	ret := h.array[0]
	h.swap(0, h.size-1)
	h.array = h.array[:h.size-1]
	h.size--
	h.heapifyDown(0)
	return ret
}
