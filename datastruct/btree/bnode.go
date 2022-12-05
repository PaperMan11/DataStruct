package btree

// BNode B树的结点
type BNode struct {
	isleaf           bool      // 是否为叶子结点
	num              int       // 当前结点 key 的个数
	kvNodes          []*KvNode // KvNode 数组
	lparent, rparent *KvNode   // 左右父结点 KvNode
	parent           *BNode    // 父结点 BNode
}

// NewBNode 新建 B- 树结点
// 		@order: 树的阶数
func NewBNode(order int) *BNode {
	return &BNode{
		isleaf:  true,
		num:     0,
		kvNodes: make([]*KvNode, order),
		lparent: nil,
		rparent: nil,
	}
}

// GetLeftSib 获取左兄弟
func (bnode *BNode) GetLeftSib() *BNode {
	if bnode == nil || bnode.lparent == nil {
		return nil
	}
	return bnode.lparent.left
}

// GetRightSib 获取右兄弟
func (bnode *BNode) GetRightSib() *BNode {
	if bnode == nil || bnode.rparent == nil {
		return nil
	}
	return bnode.rparent.right
}
