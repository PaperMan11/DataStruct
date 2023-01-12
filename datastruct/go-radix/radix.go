package goradix

import "strings"

// WalkFn 用于处理匹配到的节点
// true:  只匹配一次
// false: 匹配所有
type WalkFn func(k string, v interface{}) bool

type Tree struct {
	root *node // 根结点不存储数据
	size int
}

// New 新建空树
func New() *Tree {
	return NewFromMap(nil)
}

// NewFromMap 根据 Map 创建树
func NewFromMap(m map[string]interface{}) *Tree {
	t := &Tree{root: &node{}}
	for k, v := range m {
		t.Insert(k, v)
	}
	return t
}

func (t *Tree) Len() int {
	return t.size
}

// Insert
// 用于添加新条目或更新
// 一个已存在的条目。如果更新了现有记录，则返回（原值，true）
func (t *Tree) Insert(s string, v interface{}) (interface{}, bool) {
	var parent *node
	n := t.root // 插入点
	search := s
	for {
		// 1. match end
		if len(search) == 0 {
			if n.isLeaf() {
				old := n.leaf.val
				n.leaf.val = v
				return old, true
			}

			// create leaf node
			n.leaf = &leafNode{
				key: s,
				val: v,
			}
			t.size++
			return nil, false
		}
		// 2. walk
		parent = n
		n = n.getEdge(search[0])
		if n == nil { // no match, create one
			newEdge := edge{
				label: search[0],
				node: &node{
					leaf: &leafNode{
						key: s,
						val: v,
					},
					prefix: search,
				},
			}
			parent.addEdge(newEdge)
			t.size++
			return nil, false
		}
		// 3. match prefix
		commonPrefix := longestPrefix(search, n.prefix)
		if commonPrefix == len(n.prefix) { // match
			search = search[commonPrefix:]
			continue
		}
		// 4. split node
		// parent update edge
		t.size++
		child := &node{
			prefix: search[:commonPrefix],
		}
		parent.updateEdge(search[0], child)

		// child add first edge
		child.addEdge(edge{
			label: n.prefix[commonPrefix],
			node:  n,
		})
		n.prefix = n.prefix[commonPrefix:]

		// child add second edge
		leaf := &leafNode{
			key: s,
			val: v,
		}
		search := search[commonPrefix:]
		if len(search) == 0 {
			child.leaf = leaf
			return nil, false
		}
		child.addEdge(edge{
			label: search[0],
			node: &node{
				leaf:   leaf,
				prefix: search,
			},
		})
		return nil, false
	}
}

// Delete
// 用于删除一个键，返回前一个键 value，如果被删除
func (t *Tree) Delete(s string) (interface{}, bool) {
	var parent *node
	var label byte
	n := t.root
	search := s
	for {
		// match end
		if len(search) == 0 {
			if !n.isLeaf() {
				break
			}
			goto DELETE
		}
		// walk
		parent = n
		label = search[0]
		n = n.getEdge(label)
		if n == nil {
			break
		}
		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
		} else {
			break
		}
	}
	return nil, false
DELETE:
	// Delete the leaf
	leaf := n.leaf
	n.leaf = nil
	t.size--

	// Check if we should delete this node from the parent
	if parent != nil && len(n.edges) == 0 {
		parent.delEdge(label)
	}

	// compresse
	// Check if we should merge this node
	if n != t.root && len(n.edges) == 1 {
		n.mergeChild()
	}

	// Check if we should merge the parent's other child
	if parent != nil && parent != t.root && len(parent.edges) == 1 && !parent.isLeaf() {
		parent.mergeChild()
	}
	return leaf.val, true
}

// Get 返回 s 对应的值
func (t *Tree) Get(s string) (interface{}, bool) {
	search := s
	n := t.root
	for {
		if len(search) == 0 {
			if n.isLeaf() {
				return n.leaf.val, true
			}
			break
		}
		n = n.getEdge(search[0])
		if n == nil {
			break
		}

		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
		} else {
			break
		}
	}
	return nil, false
}

// LongestPrefix 返回最长前缀（模糊匹配）
func (t *Tree) LongestPrefix(s string) (string, interface{}, bool) {
	var last *leafNode
	n := t.root
	search := s
	for {
		if n.isLeaf() {
			last = n.leaf
		}
		if len(search) == 0 {
			break
		}

		n = n.getEdge(search[0])
		if n == nil {
			break
		}

		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
		} else {
			break
		}
	}
	if last != nil {
		return last.key, last.val, true
	}
	return "", nil, false
}

// DeletePrefix 删除前缀为 s 的 string
// 返回删除的个数
func (t *Tree) DeletePrefix(s string) int {
	return t.deletePrefix(nil, t.root, s)
}

func (t *Tree) deletePrefix(parent, n *node, prefix string) int {
	if len(prefix) == 0 {
		subTreeSize := 0
		recursiveWalk(n, func(k string, v interface{}) bool {
			subTreeSize++
			return false
		})
		if n.isLeaf() {
			n.leaf = nil
		}
		n.edges = nil // delete all

		// merge
		if parent != nil && parent != t.root && len(n.edges) == 1 && !parent.isLeaf() {
			parent.mergeChild()
		}
		t.size -= subTreeSize
		return subTreeSize
	}
	label := prefix[0]
	child := n.getEdge(label)
	if child == nil || (!strings.HasPrefix(child.prefix, prefix) && !strings.HasPrefix(prefix, child.prefix)) {
		return 0
	}

	if len(child.prefix) > len(prefix) {
		prefix = prefix[len(prefix):] // 匹配完成
	} else {
		prefix = prefix[len(child.prefix):]
	}
	return t.deletePrefix(n, child, prefix)
}

// Minimum is used to return the minimum value in the tree
func (t *Tree) Minimum() (string, interface{}, bool) {
	n := t.root
	for {
		if n.isLeaf() {
			return n.leaf.key, n.leaf.val, true
		}
		if len(n.edges) > 0 {
			n = n.edges[0].node
		} else {
			break
		}
	}
	return "", nil, false
}

// Maximum is used to return the maximum value in the tree
func (t *Tree) Maximum() (string, interface{}, bool) {
	n := t.root
	for {
		if num := len(n.edges); num > 0 {
			n = n.edges[num-1].node
			continue
		}
		if n.isLeaf() {
			return n.leaf.key, n.leaf.val, true
		}
		break
	}
	return "", nil, false
}

// Walk 遍历整棵树
func (t *Tree) Walk(fn WalkFn) {
	_ = recursiveWalk(t.root, fn)
}

// WalkPrefix
// prefix --> end (fn)
func (t *Tree) WalkPrefix(prefix string, fn WalkFn) {
	search := prefix
	n := t.root
	for {
		if len(search) == 0 {
			recursiveWalk(n, fn)
			return
		}

		n = n.getEdge(search[0])
		if n == nil {
			return
		}

		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
			continue
		}

		if strings.HasPrefix(n.prefix, search) {
			recursiveWalk(n, fn) // Child may be under our search prefix
		}
		return
	}
}

// WalkPrefix
// root --> prefix (fn)
func (t *Tree) WalkPath(prefix string, fn WalkFn) {
	search := prefix
	n := t.root
	for {
		if n.isLeaf() && fn(n.leaf.key, n.leaf.val) {
			return
		}

		if len(search) == 0 {
			return
		}

		n = n.getEdge(search[0])
		if n == nil {
			return
		}

		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
		} else {
			return
		}
	}
}

// ToMap 遍历树转化为 map
func (t *Tree) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	t.Walk(func(k string, v interface{}) bool {
		out[k] = v
		return false
	})
	return out
}

// longestPrefix 返回 s1 s2 最长公共前缀的索引
func longestPrefix(s1, s2 string) int {
	max := len(s1)
	if l := len(s2); l < max {
		max = l
	}

	var i int
	for i = 0; i < max; i++ {
		if s1[i] != s2[i] {
			break
		}
	}
	return i
}

// recursiveWalk 递归遍历
// true 单个匹配
// false 全匹配
func recursiveWalk(n *node, fn WalkFn) bool {
	if n.isLeaf() && fn(n.leaf.key, n.leaf.val) {
		return true
	}

	i := 0
	k := len(n.edges)
	for i < k {
		e := n.edges[i]
		if recursiveWalk(e.node, fn) {
			return true
		}
		// fn() 可能会 add 或 delete
		// It is a possibility that the WalkFn modified the node we are
		// iterating on. If there are no more edges, mergeChild happened,
		// so the last edge became the current node n, on which we'll
		// iterate one last time.
		if len(n.edges) == 0 {
			return recursiveWalk(n, fn)
		}

		if len(n.edges) >= k {
			i++
		}
		k = len(n.edges)
	}
	return false
}
