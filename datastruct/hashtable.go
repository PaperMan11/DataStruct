package datastruct

import (
	"fmt"
	"math"
	"sync"

	"github.com/OneOfOne/xxhash"
)

const (
	// 扩容因子
	HASHMAP_LOAD_FACTOR  = 0.7 // 扩容因子
	HASHMAP_DEFAULT_SIZE = 16  // 默认大小
	REHASH_STEP          = 10  // 每次操作迁移的节点数量
)

// 哈希表
type HashMap struct {
	array        []*keyPairs // 哈希表数组，每个元素是一个键值对
	capacity     int         // 数组容量
	len          int         // 已添加键值对数量
	capacityMask int         // 掩码，等于 capacity-1

	// 考虑并发安全
	lock sync.Mutex
}

// 键值对，连成一个链表
type keyPairs struct {
	key   string      // 键
	value interface{} // 值
	next  *keyPairs
}

// 初始化哈希表
func NewHashMap(capacity int) *HashMap {
	// 默认大小 16
	defaultCapacity := 1 << 4
	if capacity <= defaultCapacity {
		capacity = defaultCapacity
	} else {
		// 否则，实际大小为大于 capacity 的第一个 2^k
		capacity = 1 << (int(math.Ceil(math.Log2(float64(capacity)))))
	}

	// 新建一个哈希表
	m := new(HashMap)
	m.array = make([]*keyPairs, capacity)
	m.capacity = capacity
	m.capacityMask = capacity - 1
	return m
}

// 返回哈希表已添加元素的数量
func (m *HashMap) Len() int {
	return m.len
}
func (m *HashMap) Capacity() int {
	return m.capacity
}

// 将一个键进行 hash
var hashAlgorithm = func(key []byte) uint64 {
	h := xxhash.New64()
	h.Write(key)
	return h.Sum64()
}

/*
根据公式 hash(key) & (2^x-1)，使用 xxhash 哈希算法来计算键 key 的哈希值，
并且和容量掩码 mask 进行 & 求得数组的下标，用来定位键值对该放在数组的哪个下标下。
*/
// 对键进行哈希求值，并计算下标
func (m *HashMap) hashIndex(key string, mask int) int {
	// 求 hash
	hash := hashAlgorithm([]byte(key))
	// 求下标
	index := hash & uint64(mask)
	return int(index)
}

// 哈希表添加键值对
func (m *HashMap) Put(key string, value interface{}) {
	// 实现并发安全
	m.lock.Lock()
	defer m.lock.Unlock()

	// 键值对要放的哈希数组的下标
	index := m.hashIndex(key, m.capacityMask)
	// 哈希表数组下标的元素
	element := m.array[index]

	// 元素为空，表示空链表，没有哈希冲突，直接赋值
	if element == nil {
		m.array[index] = &keyPairs{
			key:   key,
			value: value,
		}
	} else {
		// 链表最后一个键值对
		var lastPairs *keyPairs

		// 遍历链表查看元素是否存在，存在则替换值，否则找到最后一个键值对
		for element != nil {
			if element.key == key {
				element.value = value
				return
			}
			lastPairs = element
			element = element.next
		}

		// 找不到键值对，将新键值对添加到链表的尾端
		lastPairs.next = &keyPairs{
			key:   key,
			value: value,
		}
	}
	// 新的哈希表数量
	newLen := m.len + 1

	// 如果超出扩容因子，需要扩容
	if float64(newLen)/float64(m.capacity) >= HASHMAP_LOAD_FACTOR {
		// 新建一个原来两倍大小的哈希表
		newM := new(HashMap)
		newM.array = make([]*keyPairs, 2*m.capacity)
		newM.capacity = 2 * m.capacity
		newM.capacityMask = 2*m.capacity - 1

		// 遍历老的哈希表，将键值对重新哈希到新哈希表
		for _, pairs := range m.array {
			for pairs != nil {
				// 直接递归 Put
				newM.Put(pairs.key, pairs.value)
				pairs = pairs.next
			}
		}

		// 替换老的哈希表
		m.array = newM.array
		m.capacity = newM.capacity
		m.capacityMask = newM.capacityMask
	}
	m.len = newLen
}

// 获取键值对
func (m *HashMap) Get(key string) (value interface{}, ok bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// 键值对要放的哈希表数组下标
	index := m.hashIndex(key, m.capacityMask)

	// 哈希表数组下标元素
	element := m.array[index]

	// 遍历链表查看元素是否存在
	for element != nil {
		if element.key == key {
			return element.value, true
		}
		element = element.next
	}
	return
}

// 删除键值对
func (m *HashMap) Delete(key string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	// 键值对要放的哈希表数组下标
	index := m.hashIndex(key, m.capacityMask)

	// 哈希表数组下标的元素
	element := m.array[index]

	// 空链表，不删除
	if element == nil {
		return
	}

	// 链表的第一个元素就是要删除的元素
	if element.key == key {
		// 将第一个元素后面的键值对链上
		m.array[index] = element.next
		m.len = m.len - 1
		return
	}

	// 下一个键值对
	nextElement := element.next
	for nextElement != nil {
		if nextElement.key == key {
			// 键值对匹配到，将该键值对从链中去掉
			element.next = nextElement.next
			m.len = m.len - 1
			return
		}
		element = nextElement
		nextElement = nextElement.next
	}
}

// 遍历打印哈希表
func (m *HashMap) Range() {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, pairs := range m.array {
		for pairs != nil {
			fmt.Printf("'%v'='%v', ", pairs.key, pairs.value)
			pairs = pairs.next
		}
	}
	fmt.Println()
}
