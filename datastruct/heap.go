package datastruct

import "fmt"

// 最大堆的实现
// 一个最大堆（完全二叉树）最大堆要求根节点始终大于左右子节点
type Heap struct {
	Size  int
	Array []int
}

// 初始化一个堆
func NewHeap(array []int) *Heap {
	h := new(Heap)
	h.Array = array
	return h
}

// 最大堆插入元素
func (h *Heap) Push(x int) {
	if h.Size == 0 {
		h.Array[0] = x
		h.Size++
		return
	}

	// i 是要插入节点的下标
	i := h.Size
	// 如果下标存在
	// 将小的值 x 一直上浮
	for i > 0 {
		// parent 为该元素的父节点下标
		parent := (i - 1) / 2
		// 如果插入的值小于等于父节点，那么可以直接退出循环，因为父节点仍然是最大的
		if x <= h.Array[parent] {
			break
		}
		// 否则将父亲节点与该节点互换，然后向上翻转，将最大的元素一直往上推
		h.Array[i] = h.Array[parent]
		i = parent
	}
	// 将该值 x 放在不会再翻转的位置
	h.Array[i] = x
	h.Size++
}

func (h *Heap) Pop() int {
	if h.Size == 0 {
		return -1
	}
	// 取出根节点
	ret := h.Array[0]
	// 因为根节点要被删除了，将最后一个节点放到根节点的位置上
	h.Size--
	x := h.Array[h.Size]  // 将最后一个节点的值拿出来
	h.Array[h.Size] = ret // 将移除的元素放到最后

	// 对根节点进行向下翻转，小的值 x 一直下沉
	i := 0
	for {
		a := 2*i + 1
		b := 2*i + 2
		if a >= h.Size { // 左节点下标超出了，表示没有左子树，那么右子树也没有
			break
		}

		// 有右子树，拿到两个子节点中最大节点的下标
		if b < h.Size && h.Array[b] > h.Array[a] {
			a = b
		}

		// 父亲节点的值都 >= 两个儿子节点最大的那一个，不需要继续向下翻转，返回
		if x >= h.Array[a] {
			break
		}

		// 将较大儿子与父亲节点交换
		h.Array[i] = h.Array[a]
		i = a // 继续向下操作
	}
	// 将最后一个元素的值 x 放在不会翻转的位置
	h.Array[i] = x
	return ret
}

func (h *Heap) Println() {
	for i := 0; i < h.Size; i++ {
		fmt.Print(h.Array[i], " ")
	}
	fmt.Println()
}
