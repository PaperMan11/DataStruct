package datastruct

import (
	"encoding/json"
	"fmt"
	"github.com/OneOfOne/xxhash"
	"reflect"
	"sync"
	"sync/atomic"
)

// 每个桶包含独立的锁和节点链表
type bucket[U comparable, T any] struct {
	head  *hashMapNode[U, T]
	mutex sync.RWMutex
}

type hashMapNode[U comparable, T any] struct {
	key   U
	value T
	next  *hashMapNode[U, T]
}

type HashMap2[U comparable, T any] struct {
	buckets       []*bucket[U, T]    // 新桶数组
	oldBuckets    []*bucket[U, T]    // 旧桶数组（扩容时使用）
	capacity      int32              // 当前容量（2的幂次）
	size          int32              // 元素数量
	capacityMask  int32              // 容量掩码 (capacity-1)
	hashAlgorithm func(key U) uint64 // 哈希算法
	rehashIndex   int32              // 扩容迁移索引
	resizingNum   int32              // 正在迁移桶数
	isResizing    atomic.Bool        // 扩容状态标记
	globalLock    sync.RWMutex       // 仅用于保护扩容元数据
}

type HashMapOption[U comparable, T any] func(*HashMap2[U, T])

func WithHashAlgorithm[U comparable, T any](hashAlgorithm func(key U) uint64) HashMapOption[U, T] {
	return func(m *HashMap2[U, T]) {
		if hashAlgorithm != nil {
			m.hashAlgorithm = hashAlgorithm
		}
	}
}

func WithInitialCapacity[U comparable, T any](capacity int) HashMapOption[U, T] {
	return func(m *HashMap2[U, T]) {
		if capacity > 0 {
			// 确保容量是2的幂次且不小于最小值
			cap := pow2(capacity)
			if cap < HASHMAP_DEFAULT_SIZE {
				cap = HASHMAP_DEFAULT_SIZE
			}
			m.capacity = int32(cap)
			m.capacityMask = int32(cap - 1)
			m.buckets = make([]*bucket[U, T], cap)
			for i := 0; i < cap; i++ {
				m.buckets[i] = &bucket[U, T]{}
			}
		}
	}
}

func NewHashMap2[U comparable, T any](options ...HashMapOption[U, T]) *HashMap2[U, T] {
	hashMap := &HashMap2[U, T]{
		capacity:      int32(HASHMAP_DEFAULT_SIZE),
		capacityMask:  int32(HASHMAP_DEFAULT_SIZE - 1),
		hashAlgorithm: defaultHashAlgorithm[U],
	}

	// 初始化默认桶
	hashMap.buckets = make([]*bucket[U, T], hashMap.capacity)
	for i := int32(0); i < hashMap.capacity; i++ {
		hashMap.buckets[i] = &bucket[U, T]{}
	}

	hashMap.isResizing.Store(false)

	for _, option := range options {
		option(hashMap)
	}

	return hashMap
}

func defaultHashAlgorithm[T comparable](key T) uint64 {
	h := xxhash.New64()

	val := reflect.ValueOf(key)
	switch val.Kind() {
	case reflect.String:
		h.WriteString(val.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num := uint64(val.Int())
		h.Write([]byte{
			byte(num >> 56), byte(num >> 48), byte(num >> 40), byte(num >> 32),
			byte(num >> 24), byte(num >> 16), byte(num >> 8), byte(num),
		})
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		num := val.Uint()
		h.Write([]byte{
			byte(num >> 56), byte(num >> 48), byte(num >> 40), byte(num >> 32),
			byte(num >> 24), byte(num >> 16), byte(num >> 8), byte(num),
		})
	case reflect.Bool:
		if val.Bool() {
			h.Write([]byte("1"))
		} else {
			h.Write([]byte("0"))
		}
	case reflect.Ptr, reflect.Struct, reflect.Array, reflect.Slice:
		if bytes, err := json.Marshal(key); err == nil {
			h.Write(bytes)
		} else {
			addr := val.UnsafeAddr()
			h.WriteString(reflect.TypeOf(key).String())
			h.Write([]byte{
				byte(addr >> 56), byte(addr >> 48), byte(addr >> 40), byte(addr >> 32),
				byte(addr >> 24), byte(addr >> 16), byte(addr >> 8), byte(addr),
			})
		}
	default:
		h.WriteString(reflect.TypeOf(key).String())
		h.WriteString(fmt.Sprintf("%v", key))
	}

	return h.Sum64()
}

func pow2(n int) int {
	if n <= HASHMAP_DEFAULT_SIZE {
		return HASHMAP_DEFAULT_SIZE
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return n + 1
}

// 计算新桶索引（带边界检查）
func (m *HashMap2[U, T]) hashIndex(key U) int {
	mask := atomic.LoadInt32(&m.capacityMask)
	hash := m.hashAlgorithm(key)
	index := int(hash & uint64(mask))

	// 双重保险：确保索引在有效范围内
	if index < 0 || index >= len(m.buckets) {
		return int(hash % uint64(len(m.buckets)))
	}
	return index
}

// 计算旧桶索引（带边界检查）
func (m *HashMap2[U, T]) oldHashIndex(key U, oldCap int) int {
	if oldCap <= 0 {
		return -1
	}
	oldMask := oldCap - 1
	hash := m.hashAlgorithm(key)
	index := int(hash & uint64(oldMask))

	if index < 0 || index >= oldCap {
		return int(hash % uint64(oldCap))
	}
	return index
}

// 渐进式扩容
func (m *HashMap2[U, T]) rehash() {
	if !m.isResizing.Load() {
		return
	}
	migrated := 0
	oldBuckets := m.oldBuckets
	oldCap := len(oldBuckets)
	if oldCap == 0 {
		m.oldBuckets = nil
		m.isResizing.Store(false)
		atomic.StoreInt32(&m.rehashIndex, 0)
		return
	}

	for m.isResizing.Load() && migrated < REHASH_STEP {
		currentIdx := atomic.LoadInt32(&m.rehashIndex)
		if currentIdx >= int32(oldCap) {
			if atomic.LoadInt32(&m.resizingNum) == 0 {
				m.oldBuckets = nil
				m.isResizing.Store(false)
				atomic.StoreInt32(&m.rehashIndex, 0)
			}
			return
		}

		if !atomic.CompareAndSwapInt32(&m.rehashIndex, currentIdx, currentIdx+1) {
			continue
		}

		oldBucket := oldBuckets[currentIdx]
		oldBucket.mutex.Lock()
		if oldBucket == nil || oldBucket.head == nil {
			oldBucket.mutex.Unlock()
			migrated++
			continue
		}

		// 取出旧桶所有节点
		nodes := oldBucket.head
		oldBucket.head = nil
		oldBucket.mutex.Unlock()

		// 迁移到新桶
		newBuckets := m.buckets
		for nodes != nil {
			next := nodes.next
			newIdx := m.hashIndex(nodes.key)

			// 确保新桶索引有效
			if newIdx < 0 || newIdx >= len(newBuckets) {
				nodes = next
				continue
			}

			newBucket := newBuckets[newIdx]
			if newBucket == nil {
				nodes = next
				continue
			}

			newBucket.mutex.Lock()
			nodes.next = newBucket.head
			newBucket.head = nodes
			newBucket.mutex.Unlock()
			nodes = next
		}
		atomic.AddInt32(&m.resizingNum, -1)
		migrated++
	}
}

// 尝试扩容检查
func (m *HashMap2[U, T]) tryResize() {
	currentCap := atomic.LoadInt32(&m.capacity)
	currentSize := atomic.LoadInt32(&m.size)

	if float64(currentSize) <= float64(currentCap)*HASHMAP_LOAD_FACTOR ||
		m.isResizing.Load() {
		return
	}
	// 计算新容量（确保不溢出）
	newCap := int32(pow2(int(currentCap * 2)))
	// 初始化新桶
	newBuckets := make([]*bucket[U, T], newCap)
	for i := int32(0); i < newCap; i++ {
		newBuckets[i] = &bucket[U, T]{}
	}
	// 原子更新桶数组
	m.oldBuckets = m.buckets
	m.buckets = newBuckets
	atomic.StoreInt32(&m.capacity, newCap)
	atomic.StoreInt32(&m.capacityMask, newCap-1)
	atomic.StoreInt32(&m.rehashIndex, 0)
	atomic.StoreInt32(&m.resizingNum, int32(len(m.oldBuckets)))
	m.isResizing.Store(true)
}

// Get 操作（只锁定相关桶）
func (m *HashMap2[U, T]) Get(key U) (T, bool) {
	var zero T

	m.globalLock.RLock()
	defer m.globalLock.RUnlock()
	// 触发扩容迁移
	if m.isResizing.Load() {
		m.rehash()
	}

	// 1. 访问新桶
	newBuckets := m.buckets
	index := m.hashIndex(key)

	// 检查新桶索引有效性
	if index < 0 || index >= len(newBuckets) {
		return zero, false
	}

	newBucket := newBuckets[index]
	newBucket.mutex.RLock()
	current := newBucket.head
	for current != nil {
		if current.key == key {
			val := current.value
			newBucket.mutex.RUnlock()
			return val, true
		}
		current = current.next
	}
	newBucket.mutex.RUnlock()

	// 2. 若在扩容，访问旧桶
	if !m.isResizing.Load() {
		return zero, false
	}

	// 安全获取旧桶引用
	oldCap := len(m.oldBuckets)
	oldIndex := m.oldHashIndex(key, oldCap)
	if oldIndex < 0 || oldIndex >= oldCap {
		return zero, false
	}

	oldBucket := m.oldBuckets[oldIndex]
	oldBucket.mutex.RLock()
	current = oldBucket.head
	for current != nil {
		if current.key == key {
			val := current.value
			oldBucket.mutex.RUnlock()
			return val, true
		}
		current = current.next
	}
	oldBucket.mutex.RUnlock()

	return zero, false
}

// Put 操作（只锁定相关桶）
func (m *HashMap2[U, T]) Put(key U, value T) {
	m.globalLock.Lock()
	defer m.globalLock.Unlock()

	if m.isResizing.Load() {
		m.rehash()
	}

	m.tryResize()

	// 1. 检查新桶并更新
	newBuckets := m.buckets
	index := m.hashIndex(key)

	if index < 0 || index >= len(newBuckets) {
		return // 无效索引，放弃操作
	}

	newBucket := newBuckets[index]
	newBucket.mutex.Lock()
	current := newBucket.head
	for current != nil {
		if current.key == key {
			current.value = value
			newBucket.mutex.Unlock()
			return
		}
		current = current.next
	}
	newBucket.mutex.Unlock()

	// 2. 若在扩容，检查旧桶并更新
	var oldBucket *bucket[U, T]
	var oldIndex int

	if m.isResizing.Load() && m.oldBuckets != nil {
		oldBuckets := m.oldBuckets
		oldCap := len(oldBuckets)
		oldIndex = m.oldHashIndex(key, oldCap)
		if oldIndex >= 0 && oldIndex < oldCap {
			oldBucket = oldBuckets[oldIndex]
		}
	}

	if oldBucket != nil {
		oldBucket.mutex.Lock()
		current = oldBucket.head
		for current != nil {
			if current.key == key {
				current.value = value
				oldBucket.mutex.Unlock()
				return
			}
			current = current.next
		}
		oldBucket.mutex.Unlock()
	}

	// 3. 插入新节点
	newBucket.mutex.Lock()
	newNode := &hashMapNode[U, T]{
		key:   key,
		value: value,
		next:  newBucket.head,
	}
	newBucket.head = newNode
	newBucket.mutex.Unlock()

	atomic.AddInt32(&m.size, 1)
}

// Remove 操作（只锁定相关桶）
func (m *HashMap2[U, T]) Remove(key U) bool {
	m.globalLock.Lock()
	defer m.globalLock.Unlock()

	if m.isResizing.Load() {
		m.rehash()
	}

	// 1. 尝试从新桶删除
	newBuckets := m.buckets
	index := m.hashIndex(key)

	if index >= 0 && index < len(newBuckets) {
		newBucket := newBuckets[index]
		newBucket.mutex.Lock()

		var prev *hashMapNode[U, T]
		current := newBucket.head
		for current != nil {
			if current.key == key {
				if prev == nil {
					newBucket.head = current.next
				} else {
					prev.next = current.next
				}
				newBucket.mutex.Unlock()
				atomic.AddInt32(&m.size, -1)
				return true
			}
			prev = current
			current = current.next
		}
		newBucket.mutex.Unlock()
	}

	// 2. 若在扩容，尝试从旧桶删除
	var oldBucket *bucket[U, T]
	var oldIndex int

	if m.isResizing.Load() && m.oldBuckets != nil {
		oldBuckets := m.oldBuckets
		oldCap := len(oldBuckets)
		oldIndex = m.oldHashIndex(key, oldCap)
		if oldIndex >= 0 && oldIndex < oldCap {
			oldBucket = oldBuckets[oldIndex]
		}
	}

	if oldBucket != nil {
		oldBucket.mutex.Lock()
		var prev *hashMapNode[U, T]
		current := oldBucket.head
		for current != nil {
			if current.key == key {
				if prev == nil {
					oldBucket.head = current.next
				} else {
					prev.next = current.next
				}
				oldBucket.mutex.Unlock()
				atomic.AddInt32(&m.size, -1)
				return true
			}
			prev = current
			current = current.next
		}
		oldBucket.mutex.Unlock()
	}

	return false
}

// 其他辅助方法（Len, Capacity等）
func (m *HashMap2[U, T]) Len() int {
	return int(atomic.LoadInt32(&m.size))
}

func (m *HashMap2[U, T]) Capacity() int {
	return int(atomic.LoadInt32(&m.capacity))
}

func (m *HashMap2[U, T]) IsEmpty() bool {
	return atomic.LoadInt32(&m.size) == 0
}

func (m *HashMap2[U, T]) IsResizing() bool {
	return m.isResizing.Load()
}
