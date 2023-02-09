package unionfind

type UnionFind struct {
	rank   []int // rank[i]表示以i为根的集合所表示的树的层数
	parent []int // parent[i]表示第i个元素所指向的父节点
	count  int   // 数据个数
}

func NewUnionFind(count int) *UnionFind {
	uf := &UnionFind{
		rank:   make([]int, count),
		parent: make([]int, count),
		count:  count,
	}

	// init
	for i := 0; i < count; i++ {
		uf.parent[i] = i
		uf.rank[i] = 1
	}
	return uf
}

// Find
// 查找过程, 查找元素p所对应的集合编号
// O(h)复杂度, h为树的高度
func (u *UnionFind) Find(p int) int {
	if p < 0 || p >= u.count {
		panic("range [0, count)")
	}
	if p != u.parent[p] {
		// compress: 减小树的高度 rank 优化（变宽）
		u.parent[p] = u.Find(u.parent[p])
	}
	return u.parent[p]
}

// IsConnected
// 查看元素p和元素q是否所属一个集合
// O(h)复杂度, h为树的高度
func (u *UnionFind) IsConnected(p, q int) bool {
	return u.Find(p) == u.Find(q)
}

// UnionElements
// 合并元素p和元素q所属的集合
// O(h)复杂度, h为树的高度
func (u *UnionFind) UnionElements(p, q int) {
	pRoot := u.Find(p)
	qRoot := u.Find(q)

	if pRoot == qRoot {
		return
	}

	// 减小树的高度 size 优化（短树加入到长树）
	if u.rank[pRoot] < u.rank[qRoot] {
		u.parent[pRoot] = qRoot
	} else if u.rank[pRoot] > u.rank[qRoot] {
		u.parent[qRoot] = pRoot
	} else {
		u.parent[pRoot] = qRoot
		u.rank[qRoot]++
	}
}
