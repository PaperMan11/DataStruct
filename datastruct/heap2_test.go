package datastruct

import (
	"testing"
)

func TestNewHeap2(t *testing.T) {
	comp := func(i, j int) int {
		if i < j {
			return -1
		} else if i == j {
			return 0
		} else {
			return 1
		}
	}

	heap := NewHeap2[int](5, comp)
	if heap == nil {
		t.Error("Expected new heap to be created, got nil")
	}

	if heap.IsEmpty() != true {
		t.Error("Expected new heap to be empty")
	}

	if heap.IsFull() != false {
		t.Error("Expected new heap to not be full")
	}

	if heap.GetSize() != 0 {
		t.Errorf("Expected heap size to be 0, got %d", heap.GetSize())
	}
}

func TestHeap2InsertAndExtract(t *testing.T) {
	// 创建一个最小堆的比较函数
	comparable := func(i, j int) int {
		a := i
		b := j
		if a < b {
			return -1
		} else if a == b {
			return 0
		} else {
			return 1
		}
	}

	heap := NewHeap2(5, comparable)

	// 插入元素
	heap.Insert(3)
	heap.Insert(1)
	heap.Insert(4)
	heap.Insert(1)
	heap.Insert(5)

	if heap.IsEmpty() {
		t.Error("Expected heap not to be empty after insertions")
	}

	if heap.GetSize() != 5 {
		t.Errorf("Expected heap size to be 5, got %d", heap.GetSize())
	}

	// 提取元素并验证顺序
	expected := []int{1, 1, 3, 4, 5}
	for i, exp := range expected {
		val := heap.Extract()
		if val != exp {
			t.Errorf("Extract %d: expected %d, got %d", i, exp, val)
		}
	}

	if heap.IsEmpty() != true {
		t.Error("Expected heap to be empty after all extractions")
	}

	if heap.GetSize() != 0 {
		t.Errorf("Expected heap size to be 0, got %d", heap.GetSize())
	}
}

func TestHeap2Peek(t *testing.T) {
	comparable := func(i, j int) int {
		if i < j {
			return -1
		} else if i == j {
			return 0
		} else {
			return 1
		}
	}

	heap := NewHeap2(5, comparable)

	// 测试空堆peek
	val := heap.Peek()
	if val != 0 {
		t.Errorf("Expected peek on empty heap to return 0, got %v", val)
	}

	// 插入元素后测试peek
	heap.Insert(3)
	heap.Insert(1)
	heap.Insert(2)

	val = heap.Peek()
	if val != 1 {
		t.Errorf("Expected peek to return 1 (minimum), got %v", val)
	}

	// 确保peek不删除元素
	if heap.GetSize() != 3 {
		t.Errorf("Expected heap size to remain 3 after peek, got %d", heap.GetSize())
	}
}

func TestHeap2IsFull(t *testing.T) {
	comparable := func(i, j int) int {
		if i < j {
			return -1
		} else if i == j {
			return 0
		} else {
			return 1
		}
	}

	heap := NewHeap2(3, comparable)

	if heap.IsFull() {
		t.Error("Expected heap not to be full initially")
	}

	heap.Insert(1)
	heap.Insert(2)
	heap.Insert(3)

	if !heap.IsFull() {
		t.Error("Expected heap to be full after 3 insertions")
	}

	// 尝试插入到已满的堆中
	heap.Insert(4)
	if heap.GetSize() != 3 {
		t.Error("Expected heap size to remain 3 after trying to insert into full heap")
	}
}

func TestHeap2StringValues(t *testing.T) {
	comparable := func(i, j int) int {
		// 按字符串长度比较
		a := i
		b := j
		if a < b {
			return -1
		} else if a == b {
			return 0
		} else {
			return 1
		}
	}

	heap := NewHeap2(5, comparable)

	heap.Insert(5) // 长度为5
	heap.Insert(2) // 长度为2
	heap.Insert(8) // 长度为8
	heap.Insert(1) // 长度为1
	heap.Insert(3) // 长度为3

	// 提取元素并验证顺序
	expected := []int{1, 2, 3, 5, 8}
	for i, exp := range expected {
		val := heap.Extract()
		if val != exp {
			t.Errorf("Extract %d: expected %d, got %d", i, exp, val)
		}
	}
}
