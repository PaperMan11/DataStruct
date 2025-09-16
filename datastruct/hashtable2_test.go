package datastruct

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

// 测试在并发环境下Get方法的边界条件，特别是索引越界问题
func TestHashMap2_GetBoundaryConditions(t *testing.T) {
	// 创建一个小容量的HashMap以快速触发扩容
	hm := NewHashMap2[string, int](WithInitialCapacity[string, int](4))

	// 插入足够多的数据以触发扩容
	for i := 0; i < 100; i++ {
		hm.Put("key"+strconv.Itoa(i), i)
	}

	// 在扩容过程中并发执行Get操作
	var wg sync.WaitGroup
	errCount := 0
	var mu sync.Mutex

	// 启动多个goroutine并发读取
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := "key" + strconv.Itoa((threadID*100+j)%100)
				_, ok := hm.Get(key)
				if !ok {
					mu.Lock()
					errCount++
					mu.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()

	// 检查是否有错误发生
	if errCount > 0 {
		t.Errorf("并发Get操作中出现 %d 次错误", errCount)
	}
}

// 测试在扩容过程中访问不存在的键
func TestHashMap2_GetNonExistentDuringResize(t *testing.T) {
	hm := NewHashMap2[string, int](WithInitialCapacity[string, int](4))

	// 插入数据触发扩容
	for i := 0; i < 100; i++ {
		hm.Put("key"+strconv.Itoa(i), i)
	}

	// 在扩容过程中尝试访问不存在的键
	for i := 0; i < 1000; i++ {
		_, ok := hm.Get("nonexistent" + strconv.Itoa(i))
		if ok {
			t.Errorf("不应该找到键: %s", "nonexistent"+strconv.Itoa(i))
		}
	}
}

// 测试在并发Put和Get操作中访问buckets数组边界
func TestHashMap2_ConcurrentPutGetBucketsBoundary(t *testing.T) {
	hm := NewHashMap2[int, string](WithInitialCapacity[int, string](16))

	var wg sync.WaitGroup

	// 并发执行Put操作
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := id*1000 + j
				hm.Put(key, "value"+strconv.Itoa(key))
			}
		}(i)
	}

	// 并发执行Get操作
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := id*1000 + j
				hm.Get(key)
			}
		}(i)
	}

	wg.Wait()

	// 验证插入的数据都能正确获取
	for i := 0; i < 5; i++ {
		for j := 0; j < 100; j++ {
			key := i*1000 + j
			_, ok := hm.Get(key)
			if !ok {
				t.Errorf("找不到应该存在的键: %d", key)
			}
		}
	}
}

// 测试在各种操作中oldBuckets访问的边界条件
func TestHashMap2_OldBucketsBoundaryConditions(t *testing.T) {
	hm := NewHashMap2[string, int](WithInitialCapacity[string, int](4))

	// 插入数据触发扩容
	for i := 0; i < 50; i++ {
		hm.Put("key"+strconv.Itoa(i), i)
	}

	// 在扩容过程中混合执行各种操作
	var wg sync.WaitGroup

	// Get操作
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				hm.Get("key" + strconv.Itoa(j%50))
			}
		}()
	}

	// Put操作
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				hm.Put("newkey"+strconv.Itoa(j), j)
			}
		}()
	}

	// Remove操作
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 25; j++ {
				hm.Remove("key" + strconv.Itoa(j))
			}
		}()
	}

	wg.Wait()
}

// 测试自定义哈希函数导致的边界情况
func TestHashMap2_CustomHashBoundary(t *testing.T) {
	// 使用返回固定值的哈希函数，导致大量冲突
	customHash := func(key string) uint64 {
		return 1000000 // 固定返回一个大值
	}

	hm := NewHashMap2[string, int](WithHashAlgorithm[string, int](customHash))

	// 插入大量数据，全部冲突到同一个bucket
	for i := 0; i < 100; i++ {
		hm.Put("key"+strconv.Itoa(i), i)
	}

	// 验证所有数据都能正确获取
	for i := 0; i < 100; i++ {
		val, ok := hm.Get("key" + strconv.Itoa(i))
		if !ok || val != i {
			t.Errorf("获取键值对失败: key=key%d, expected=%d, actual=%d", i, i, val)
		}
	}

	// 在存在大量冲突的情况下测试Remove
	for i := 0; i < 50; i++ {
		ok := hm.Remove("key" + strconv.Itoa(i))
		if !ok {
			t.Errorf("删除键失败: key%d", i)
		}
	}

	// 验证剩余数据
	for i := 50; i < 100; i++ {
		val, ok := hm.Get("key" + strconv.Itoa(i))
		if !ok || val != i {
			t.Errorf("删除操作后获取键值对失败: key=key%d, expected=%d, actual=%d", i, i, val)
		}
	}
}

// 测试空键值的处理
func TestHashMap2_EmptyKeyHandling(t *testing.T) {
	hm := NewHashMap2[string, string]()

	// 插入空键
	hm.Put("", "empty key value")

	// 获取空键的值
	val, ok := hm.Get("")
	if !ok || val != "empty key value" {
		t.Errorf("处理空键失败: expected='empty key value', actual='%s'", val)
	}

	// 删除空键
	ok = hm.Remove("")
	if !ok {
		t.Error("删除空键失败")
	}

	// 确认空键已被删除
	_, ok = hm.Get("")
	if ok {
		t.Error("空键应该已被删除")
	}
}

// 测试并发安全：多线程读写+扩容，验证数据一致性（无竞态、无丢失）
func TestHashMap2_ConcurrentSafety(t *testing.T) {
	const (
		numWriters = 5    // 写协程数量
		numReaders = 10   // 读协程数量
		numOps     = 1000 // 每个写协程插入/删除次数
	)

	hm := NewHashMap2[string, int]()
	var wg sync.WaitGroup
	var (
		writeErrCount int32
		readErrCount  int32
	)

	// 1. 启动写协程：插入/删除数据（模拟并发修改）
	for w := 0; w < numWriters; w++ {
		wg.Add(1)
		go func(writerID int) {
			defer wg.Done()
			for i := 0; i < numOps; i++ {
				key := strconv.Itoa(writerID*numOps + i) // 唯一 key（避免写冲突）
				val := writerID*numOps + i

				// 先插入
				hm.Put(key, val)
				// 验证插入
				storedVal, ok := hm.Get(key)
				if !ok || storedVal != val {
					t.Log("插入失败", key, val, storedVal, ok)
					atomic.AddInt32(&writeErrCount, 1)
				}

				// 每 10 次操作删除一次（模拟随机删除）
				if i%10 == 0 {
					ok := hm.Remove(key)
					if !ok {
						t.Log("删除失败", key)
						atomic.AddInt32(&writeErrCount, 1)
					}
				}
			}
		}(w)
	}

	// 2. 启动读协程：持续读取数据（模拟并发查询）
	for r := 0; r < numReaders; r++ {
		wg.Add(1)
		go func(readerID int) {
			defer wg.Done()
			for i := 0; i < numOps*2; i++ { // 读操作次数多于写操作
				key := strconv.Itoa(i % (numWriters * numOps)) // 循环读取可能存在的 key
				_, ok := hm.Get(key)
				// 读操作无错误（即使 key 不存在，ok=false 也是正常结果）
				_ = ok // 仅验证无 panic，无需计数
			}
		}(r)
	}

	// 3. 等待所有协程完成
	wg.Wait()

	// 4. 验证无并发错误
	if writeErrCount > 0 {
		t.Errorf("并发写操作中出现 %d 次错误", writeErrCount)
	}
	if readErrCount > 0 {
		t.Errorf("并发读操作中出现 %d 次错误", readErrCount)
	}

	// 5. 验证最终数据一致性（剩余数据量应在合理范围）
	finalLen := hm.Len()
	expectedMinLen := numWriters*numOps - (numWriters*numOps)/10 // 减去 1/10 的删除量
	if finalLen < expectedMinLen || finalLen > numWriters*numOps {
		t.Errorf("最终数据量异常：预期 [%d, %d]，实际 %d",
			expectedMinLen, numWriters*numOps, finalLen)
	}
}
