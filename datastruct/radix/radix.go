package radix

import "sync"

// Description:
//		前缀基数树-radix
//		以多叉树的形式实现,根据'/'进行string分割,将分割后的string数组进行段存储
//		不存储其他元素,仅对string进行分段存储和模糊匹配
//		使用互斥锁实现并发控制
type radix struct {
	root *node
	size int
	m    sync.Mutex
}

func New() *radix {
	return &radix{
		root: newNode(""),
		size: 0,
	}
}

func (r *radix) Size() int {
	if r == nil || r.root == nil {
		return 0
	}
	return r.size
}

func (r *radix) Clear() {
	if r == nil {
		return
	}
	r.m.Lock()
	defer r.m.Unlock()
	r.root = newNode("")
	r.size = 0
}

func (r *radix) Empty() bool {
	if r == nil {
		return true
	}
	return r.size == 0
}

func (r *radix) Insert(s string) (ok bool) {
	if r == nil {
		return false
	}
	parts, s := analysis(s)
	r.m.Lock()
	if r.root == nil {
		r.root = newNode("")
	}
	ok = r.root.insert(s, parts, 0)
	if ok {
		r.size++
	}
	r.m.Unlock()
	return ok
}

// Erase 删除某个 string
func (r *radix) Erase(s string) (ok bool) {
	if r.root == nil || r.Empty() || len(s) == 0 {
		return false
	}
	parts, _ := analysis(s)
	r.m.Lock()
	ok = r.root.erase(parts, 0)
	if ok {
		r.size--
		if r.size == 0 {
			r.root = nil
		}
	}
	r.m.Unlock()
	return ok
}

// Delete 删除前缀为 s 的所有 string
func (r *radix) Delete(s string) (num int) {
	if r.root == nil || r.Empty() || len(s) == 0 {
		return 0
	}
	parts, _ := analysis(s)
	r.m.Lock()
	num = r.root.delete(parts, 0)
	if num > 0 {
		r.size -= num
		if r.size == 0 {
			r.root = nil
		}
	}
	r.m.Unlock()
	return num
}

// Count 前缀为 s 的 string 的数量
func (r *radix) Count(s string) (num int) {
	if r.root == nil || r.Empty() || len(s) == 0 {
		return 0
	}
	parts, _ := analysis(s)
	r.m.Lock()
	num = r.root.count(parts, 0)
	r.m.Unlock()
	return num
}

// Mate 对具体的路由进行参数映射
// eg.
// 		/api/v1/user/:id/post
// 		/api/v1/user/1/post
//		map {"id": "1"}
func (r *radix) Mate(s string) (m map[string]string, ok bool) {
	if r.root == nil || r.Empty() || len(s) == 0 {
		return nil, false
	}
	return r.root.mate(s)
}
