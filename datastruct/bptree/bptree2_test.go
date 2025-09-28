package bptree

import (
	"fmt"
	"testing"
	"time"
)

// 测试基本插入和查询功能
func TestInsertAndSearch(t *testing.T) {
	tree := NewBpTree[int, string](3)

	// 测试插入
	tree.Insert(1, "one")
	tree.Insert(2, "two")
	tree.Insert(3, "three")

	// 测试查询存在的键
	val, ok := tree.Search(2)
	if !ok || val != "two" {
		t.Errorf("Search(2) 期望得到 (two, true)，实际得到 (%v, %v)", val, ok)
	}

	// 测试查询不存在的键
	val, ok = tree.Search(4)
	if ok {
		t.Errorf("Search(4) 期望得到 (零值, false)，实际得到 (%v, %v)", val, ok)
	}

	// 测试更新值
	tree.Insert(2, "second")
	val, ok = tree.Search(2)
	if !ok || val != "second" {
		t.Errorf("Search(2) 更新后期望得到 (second, true)，实际得到 (%v, %v)", val, ok)
	}

	// 测试大小
	if tree.Size() != 3 {
		t.Errorf("Size() 期望得到 3，实际得到 %d", tree.Size())
	}
	tree.PrintTree()
}

// 测试节点分裂功能
func TestNodeSplit(t *testing.T) {
	// 使用阶为3的B+树，超过3个键会触发分裂
	tree := NewBpTree[int, int](3)

	// 插入4个键，应该触发分裂
	for i := 1; i <= 4; i++ {
		tree.Insert(i, i*10)
	}
	tree.PrintTree()

	// 检查大小
	if tree.Size() != 4 {
		t.Errorf("Size() 期望得到 4，实际得到 %d", tree.Size())
	}

	// 检查所有键是否都能查到
	for i := 1; i <= 4; i++ {
		val, ok := tree.Search(i)
		if !ok || val != i*10 {
			t.Errorf("Search(%d) 期望得到 (%d, true)，实际得到 (%v, %v)", i, i*10, val, ok)
		}
	}

	// 插入更多键，测试多层分裂
	for i := 5; i <= 10; i++ {
		tree.Insert(i, i*10)
	}

	t.Log("-----------------------------\n")
	tree.PrintTree()

	if tree.Size() != 10 {
		t.Errorf("Size() 期望得到 10，实际得到 %d", tree.Size())
	}
}

// 测试删除功能
func TestRemove(t *testing.T) {
	tree := NewBpTree[int, string](3)

	// 插入测试数据
	for i := 1; i <= 5; i++ {
		tree.Insert(i, fmt.Sprintf("value%d", i))
	}

	// 测试删除存在的键
	ok := tree.Remove(3)
	if !ok || tree.Size() != 4 {
		t.Errorf("Remove(3) 期望成功，实际结果 %v，大小 %d", ok, tree.Size())
	}

	// 测试删除后无法查询
	_, ok = tree.Search(3)
	if ok {
		t.Errorf("删除后 Search(3) 应该返回 false，实际返回 %v", ok)
	}

	// 测试删除不存在的键
	ok = tree.Remove(10)
	if ok {
		t.Errorf("Remove(10) 期望失败，实际结果 %v", ok)
	}

	// 测试删除导致的合并
	for i := 1; i <= 2; i++ {
		tree.Remove(i)
	}

	if tree.Size() != 2 {
		t.Errorf("多次删除后 Size() 期望得到 2，实际得到 %d", tree.Size())
	}

	// 测试删除所有键
	for i := 4; i <= 5; i++ {
		tree.Remove(i)
	}

	if !tree.IsEmpty() {
		t.Errorf("删除所有键后 IsEmpty() 应该返回 true，实际返回 false")
	}
}

// 测试范围查询
func TestRangeSearch(t *testing.T) {
	tree := NewBpTree[int, int](3)

	// 插入测试数据
	for i := 1; i <= 10; i++ {
		tree.Insert(i, i)
	}

	// 测试正常范围
	result := tree.BinaryRangeSearch(3, 7)
	if len(result) != 5 {
		t.Errorf("RangeSearch(3,7) 期望得到 5 个结果，实际得到 %d", len(result))
	}

	// 验证结果
	expected := []int{3, 4, 5, 6, 7}
	for i, val := range result {
		if val != expected[i] {
			t.Errorf("RangeSearch 结果在位置 %d 期望得到 %d，实际得到 %d", i, expected[i], val)
		}
	}

	// 测试边界范围
	result = tree.BinaryRangeSearch(1, 1)
	if len(result) != 1 || result[0] != 1 {
		t.Errorf("RangeSearch(1,1) 期望得到 [1]，实际得到 %v", result)
	}

	// 测试超出范围
	result = tree.BinaryRangeSearch(11, 20)
	if len(result) != 0 {
		t.Errorf("RangeSearch(11,20) 期望得到空结果，实际得到 %d 个结果", len(result))
	}
}

// 测试根节点分裂和合并
func TestRootSplitAndMerge(t *testing.T) {
	tree := NewBpTree[int, string](3)

	// 插入足够多的键使根节点分裂
	for i := 1; i <= 10; i++ {
		tree.Insert(i, fmt.Sprintf("val%d", i))
	}

	// 验证根节点已分裂（不再是叶子节点）
	if tree.root.isLeaf {
		t.Error("插入足够多键后，根节点应该是非叶子节点")
	}

	tree.PrintTree()
	t.Log("-----------------------------\n")
	// 删除键直到根节点合并
	for i := 10; i >= 1; i-- {
		t.Logf("删除键 %d, %v", i, tree.Remove(i))

		tree.PrintTree()
		t.Log("-----------------------------\n")
	}

	tree.PrintTree()

	// 验证根节点已合并（变回叶子节点）
	if !tree.root.isLeaf {
		t.Error("删除足够多键后，根节点应该是叶子节点")
	}

	tree.Insert(1, "val1")
	tree.PrintTree()
}

// 测试字符串键
func TestStringKeys(t *testing.T) {
	tree := NewBpTree[string, int](3)

	// 插入字符串键
	tree.Insert("apple", 1)
	tree.Insert("banana", 2)
	tree.Insert("cherry", 3)
	tree.Insert("date", 4)

	// 测试查询
	val, ok := tree.BinarySearch("banana")
	if !ok || val != 2 {
		t.Errorf("Search(banana) 期望得到 (2, true)，实际得到 (%v, %v)", val, ok)
	}

	// 测试范围查询
	result := tree.BinaryRangeSearch("banana", "date")
	if len(result) != 3 {
		t.Errorf("RangeSearch 期望得到 3 个结果，实际得到 %d", len(result))
	}
}

// 测试大量数据插入和查询
func TestLargeDataset(t *testing.T) {
	tree := NewBpTree[int, int](5)
	count := 100000

	// 插入大量数据
	for i := 1; i <= count; i++ {
		tree.Insert(i, i)
	}

	//tree.PrintTree()

	if tree.Size() != count {
		t.Errorf("插入 %d 条数据后，大小应该是 %d，实际是 %d", count, count, tree.Size())
	}

	// 随机查询一些数据
	for i := 1; i <= count; i += 1 {
		val, ok := tree.BinarySearch(i)
		if !ok || val != i {
			t.Errorf("Search(%d) 期望得到 (%d, true)，实际得到 (%v, %v)", i, i, val, ok)
		}
	}

	// 删除一部分数据
	for i := 1; i <= count/2; i++ {
		tree.Remove(i)
	}

	if tree.Size() != count/2 {
		t.Errorf("删除一半数据后，大小应该是 %d，实际是 %d", count/2, tree.Size())
	}
}

type User struct {
	Id        int
	Name      string
	Age       int
	Addr      string
	CreatedAt int64
}

func TestMap(t *testing.T) {
	user1 := User{
		Id:        1,
		Name:      "张三",
		Age:       18,
		Addr:      "北京",
		CreatedAt: time.Now().Unix(),
	}
	user2 := User{
		Id:        2,
		Name:      "李四",
		Age:       19,
		Addr:      "上海",
		CreatedAt: time.Now().Unix(),
	}
	m := make(map[User]int)
	arr := []User{user1, user2, user1, user1, user2, user1}
	for _, v := range arr {
		m[v]++
	}
	t.Log(m)
}
