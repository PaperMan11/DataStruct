package consistent

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"
	"sync/atomic"

	blake2b "github.com/minio/blake2b-simd"
)

const replicationFactor = 10 // 虚拟节点数量

var ErrNoHosts = errors.New("no hosts added")

type Host struct {
	Name string
	Load int64
}

type Consistent struct {
	hosts     map[uint64]string // 哈希与主机地址映射
	sortedSet []uint64          // 哈希环
	loadMap   map[string]*Host  // 主机地址与主机映射（不包括虚拟节点）
	totalLoad int64             // 总访问量

	sync.RWMutex
}

func New() *Consistent {
	return &Consistent{
		hosts:     make(map[uint64]string),
		sortedSet: make([]uint64, 0),
		loadMap:   make(map[string]*Host),
	}
}

// Add host to the ring
func (c *Consistent) Add(host string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.loadMap[host]; ok {
		return
	}

	c.loadMap[host] = &Host{Name: host, Load: 0}
	for i := 0; i < replicationFactor; i++ { // 添加虚拟节点
		hashV := c.hash(fmt.Sprintf("%s%d", host, i))
		c.hosts[hashV] = host
		c.sortedSet = append(c.sortedSet, hashV)
	}

	// 升序
	sort.Slice(c.sortedSet, func(i, j int) bool {
		return c.sortedSet[i] < c.sortedSet[j]
	})
}

// Deletes host from the ring
func (c *Consistent) Remove(host string) bool {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.loadMap[host]; !ok {
		return false
	}

	for i := 0; i < replicationFactor; i++ {
		h := c.hash(fmt.Sprintf("%s%d", host, i))
		delete(c.hosts, h)
		c.delSlice(h)
	}
	delete(c.loadMap, host)
	return true
}

// delSlice 删除哈希环上的值
func (c *Consistent) delSlice(val uint64) {
	idx := -1
	l := 0
	r := len(c.sortedSet) - 1
	for l <= r {
		m := (l + r) / 2
		if c.sortedSet[m] == val {
			idx = m
			break
		} else if c.sortedSet[m] < val {
			l = m + 1
		} else if c.sortedSet[m] > val {
			r = m - 1
		}
	}
	if idx != -1 {
		c.sortedSet = append(c.sortedSet[:idx], c.sortedSet[idx+1:]...)
	}
}

// Get 根据 hash(host) 找到对应主机
func (c *Consistent) Get(host string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.hosts) == 0 {
		return "", ErrNoHosts
	}

	h := c.hash(host)
	idx := c.search(h)
	return c.hosts[c.sortedSet[idx]], nil
}

// GetLeast
// It uses Consistent Hashing With Bounded loads
// to pick the least loaded host that can serve the key
// It returns `ErrNoHosts` if the ring has no hosts in it.
func (c *Consistent) GetLeast(key string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.hosts) == 0 {
		return "", ErrNoHosts
	}

	h := c.hash(key)
	i := c.search(h)
	for {
		host := c.hosts[c.sortedSet[i]]
		if c.loadOK(host) {
			return host, nil
		}
		i = (i + 1) % len(c.hosts)
	}
}

// search 在哈希环上找到对应的主机
func (c *Consistent) search(key uint64) int {
	idx := sort.Search(len(c.sortedSet), func(i int) bool {
		return c.sortedSet[i] >= key
	})

	if idx >= len(c.sortedSet) {
		idx = 0
	}

	return idx
}

// Sets the load of `host` to the given `load`
func (c *Consistent) UpdateLoad(host string, load int64) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.loadMap[host]; !ok {
		return
	}
	c.totalLoad -= c.loadMap[host].Load
	c.loadMap[host].Load = load
	c.totalLoad += load
}

// Increments the load of host by 1
func (c *Consistent) Inc(host string) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.loadMap[host]; !ok {
		return
	}
	atomic.AddInt64(&c.loadMap[host].Load, 1)
	atomic.AddInt64(&c.totalLoad, 1)
}

// Decrements the load of host by 1
func (c *Consistent) Done(host string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.loadMap[host]; !ok {
		return
	}
	atomic.AddInt64(&c.loadMap[host].Load, -1)
	atomic.AddInt64(&c.totalLoad, -1)
}

// Return the list of hosts in the ring (real host)
func (c *Consistent) Hosts() (hosts []string) {
	c.RLock()
	defer c.RUnlock()
	for k := range c.loadMap {
		hosts = append(hosts, k)
	}
	return hosts
}

// GetLoads
// Returns the loads of all the hosts
func (c *Consistent) GetLoads() map[string]int64 {
	loads := map[string]int64{}

	for k, v := range c.loadMap {
		loads[k] = v.Load
	}
	return loads
}

// MaxLoad
// Returns the maximum load of the single host
// which is:
// (total_load/number_of_hosts)*1.25
// total_load = is the total number of active requests served by hosts
func (c *Consistent) MaxLoad() int64 {
	if c.totalLoad == 0 {
		c.totalLoad = 1
	}
	var avgLoadPerNode float64
	avgLoadPerNode = float64(c.totalLoad / int64(len(c.loadMap)))
	if avgLoadPerNode == 0 {
		avgLoadPerNode = 1
	}
	avgLoadPerNode = math.Ceil(avgLoadPerNode * 1.25)
	return int64(avgLoadPerNode)
}

func (c *Consistent) loadOK(host string) bool {
	// a safety check if someone performed c.Done more than needed
	if c.totalLoad < 0 {
		c.totalLoad = 0
	}

	var avgLoadPerNode float64
	avgLoadPerNode = float64((c.totalLoad + 1) / int64(len(c.loadMap)))
	if avgLoadPerNode == 0 {
		avgLoadPerNode = 1
	}
	avgLoadPerNode = math.Ceil(avgLoadPerNode * 1.25)

	bhost, ok := c.loadMap[host]
	if !ok {
		panic(fmt.Sprintf("given host(%s) not in loadsMap", bhost.Name))
	}

	if float64(bhost.Load)+1 <= avgLoadPerNode {
		return true
	}

	return false
}

func (c *Consistent) hash(key string) uint64 {
	out := blake2b.Sum512([]byte(key))
	return binary.LittleEndian.Uint64(out[:])
}
