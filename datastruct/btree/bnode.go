package btree

// BNode B树的结点
type BNode struct {
	isleaf     bool      // 是否为叶子结点
	num        int       // 当前结点 key 的个数
	pindex     int       // 在父结点中的索引 bNodes[i]，-1代表根结点没有父结点
	kvNodes    []*KvNode // KvNode 数组
	childNodes []*BNode  // 孩子结点
	parent     *BNode    // 父结点 BNode
}

// newBNode 新建 BNode 结点
func newBNode(order int) *BNode {
	return &BNode{
		isleaf:     true,
		num:        0,
		pindex:     -1,
		kvNodes:    make([]*KvNode, order),
		childNodes: make([]*BNode, 0),
		parent:     nil,
	}
}

// freeBNode 释放 BNode 结点
func (bnode *BNode) freeBNode() {
	if bnode.parent != nil {
		bnode.parent = nil
	}
	if bnode.kvNodes != nil {
		bnode.kvNodes = nil
	}
	if bnode.childNodes != nil {
		bnode.childNodes = nil
	}
	bnode = nil
}

// insertSpPos 将 node 插入到 childNodes[] 指定位置
func (bnode *BNode) insertPosNode(node *BNode, i int) {
	bnode.childNodes = append(bnode.childNodes, nil)
	copy(bnode.childNodes[i+1:], bnode.childNodes[i:])
	bnode.childNodes[i] = node
}

// deleteSpPos 删除 childNodes[] 指定位置
func (bnode *BNode) deletePosNode(i int) {
	bnode.childNodes = append(bnode.childNodes[:i], bnode.childNodes[i+1:]...)
}

// getLeftSib 获得左兄弟
func (bnode *BNode) getLeftSib() *BNode {
	if bnode.parent == nil || bnode.pindex <= 0 {
		return nil
	}
	return bnode.parent.childNodes[bnode.pindex-1]
}

// getLeftSib 获得右兄弟
func (bnode *BNode) getRightSib() *BNode {
	if bnode.parent == nil || bnode.pindex >= bnode.parent.num {
		return nil
	}
	return bnode.parent.childNodes[bnode.pindex+1]
}

// mergeNode 合并，node被释放
// bnode.kNodes[] + kvn + node.kNodes[]
func (bnode *BNode) mergeNode(kvn *KvNode, node *BNode) {
	// 添加kvn
	bnode.kvNodes[bnode.num] = kvn
	copy(bnode.kvNodes[bnode.num+1:], node.kvNodes[:node.num])
	// 追加childNodes
	if len(node.childNodes) > 0 {
		bnode.childNodes = append(bnode.childNodes, node.childNodes...)
		for _, n := range node.childNodes {
			n.pindex += bnode.num + 1
			n.parent = bnode
		}
	}
	bnode.num += node.num + 1
	node.freeBNode()
}

// moveToRight 右移调整
func (bnode *BNode) moveToRight(parent, rightSib *BNode) {
	var (
		upKvNode   = bnode.kvNodes[bnode.num-1]   // 上移kvn
		donwKvNode = parent.kvNodes[bnode.pindex] // 下移kvn
	)
	parent.kvNodes[bnode.pindex] = upKvNode
	copy(rightSib.kvNodes[1:], rightSib.kvNodes[0:])
	rightSib.kvNodes[0] = donwKvNode
	if len(bnode.childNodes) > 0 { // 添加新的子结点
		for _, n := range rightSib.childNodes {
			n.pindex += 1
		}
		rightSibNewChild := bnode.childNodes[bnode.num]
		rightSibNewChild.pindex = 0
		rightSibNewChild.parent = rightSib
		bnode.deletePosNode(bnode.num)
		rightSib.insertPosNode(rightSibNewChild, 0)
	}
	rightSib.num++
	bnode.num--
}

// moveToRight 左移调整
func (bnode *BNode) moveToLeft(parent, leftSib *BNode) {
	var (
		upKvNode   = bnode.kvNodes[0]               // 上移kvn
		donwKvNode = parent.kvNodes[bnode.pindex-1] // 下移kvn
	)
	parent.kvNodes[bnode.pindex-1] = upKvNode
	leftSib.kvNodes[leftSib.num] = donwKvNode
	copy(bnode.kvNodes[0:], bnode.kvNodes[1:])
	if len(bnode.childNodes) > 0 { // 添加新的子结点
		leftSibNewChild := bnode.childNodes[0]
		leftSibNewChild.pindex = leftSib.num + 1
		leftSibNewChild.parent = leftSib
		bnode.deletePosNode(0)
		leftSib.insertPosNode(leftSibNewChild, leftSib.num+1)
	}
	leftSib.num++
	bnode.num--
}

