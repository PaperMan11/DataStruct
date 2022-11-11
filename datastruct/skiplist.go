package datastruct

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// 仿 redis zset

type Node struct {
	Key   string
	Score float64 // 分数
	Val   interface{}
	Next  []*Node // 记录该节点后面链表的头节点(每一层)
}

func NewNode(key string, score float64, val interface{}, maxLevel uint16) *Node {
	return &Node{
		Key:   key,
		Score: score,
		Val:   val,
		Next:  make([]*Node, maxLevel),
	}
}

// 自定义 error

// 分数小于最小
type ScoreOutOfRangeError struct {
	key   string
	score float64
	min   float64
}

var _ error = (*ScoreOutOfRangeError)(nil)

func NewScoreOutOfRangeError(key string, score, min float64) *ScoreOutOfRangeError {
	return &ScoreOutOfRangeError{
		key:   key,
		score: score,
		min:   min,
	}
}

func (e *ScoreOutOfRangeError) Error() string {
	return fmt.Sprintf("score out of range, min(include)=%v, score=%v, key=%s", e.min, e.score, e.key)
}

// 键值错误
type InvalidKeyError struct {
	key    string
	reason string
}

var _ error = (*InvalidKeyError)(nil)

func NewInvalidKeyError(key, reason string) *InvalidKeyError {
	return &InvalidKeyError{
		key:    key,
		reason: reason,
	}
}

func (e *InvalidKeyError) Error() string {
	return fmt.Sprintf("invalid key: %s, reason: %s", e.key, e.reason)
}

// ------------------------------- skip list ---------------------------------

const (
	defaultSkipLinkedP = 0.5
	defaultMaxLevel    = 64
)

type SkipLinks struct {
	scoreMap    map[string]float64 // key score 映射表
	curLevel    uint16             // 当前层数
	maxLevel    uint16             // 最大层数
	minScore    float64            // 最小分数值
	skipLinkedP float64            // 上升索引概率
	head        *Node              // 头节点
	rand        *rand.Rand         // 生成随机数，用于与skipLinkedP比较，上升索引
}

// NewSkipLinks 创建跳跃表
// 		@maxLevel 设置最大层数
// 		@minScore 设置最小分数
// 		@p 自定义上升概率（0<p<1 默认为0.5）
func NewSkipLinked(maxLevel uint16, minScore float64, p ...float64) *SkipLinks {
	skipLinkedP := defaultSkipLinkedP
	if len(p) != 0 && p[0] > float64(0) && p[0] < float64(1) {
		skipLinkedP = p[0]
	}
	if maxLevel <= 0 || maxLevel > defaultMaxLevel {
		maxLevel = defaultMaxLevel
	}
	return &SkipLinks{
		scoreMap:    make(map[string]float64),
		curLevel:    1,
		maxLevel:    maxLevel,
		minScore:    minScore,
		skipLinkedP: skipLinkedP,
		head:        NewNode("", minScore, nil, maxLevel),
		rand:        rand.New(rand.NewSource(time.Now().Unix())),
	}
}

// 通过概率计算level，用于后续判断是否增加索引（并且具体加到第几层）
func (l *SkipLinks) getLevel() uint16 {
	level := uint16(1)
	for l.rand.Float64() < l.skipLinkedP && level < l.maxLevel {
		level++
	}
	return level
}

// 比较大小，用于排序
func (l *SkipLinks) compare(score1, score2 float64, key1, key2 string) int {
	if score1 > score2 {
		return 1
	} else if score1 < score2 {
		return -1
	} else {
		return strings.Compare(key1, key2)
	}
}

// Search 查找对应 key 的 value
func (l *SkipLinks) Search(key string) (bool, interface{}) {
	score, exist := l.scoreMap[key]
	if !exist {
		return false, nil
	}
	exist, node := l.search(key, score, l.head, l.curLevel)
	if !exist {
		return false, nil
	}
	return true, node
}

func (l *SkipLinks) search(key string, score float64, node *Node, level uint16) (bool, *Node) {
	next := node.Next[level-1] // 当前层 当前节点的 下一个节点
	for next != nil && l.compare(next.Score, score, next.Key, key) <= 0 {
		// 往后移
		node = next
		next = node.Next[level-1]
	}
	if l.compare(node.Score, score, node.Key, key) == 0 {
		return true, node
	}
	if level == 1 {
		return false, nil
	}
	return l.search(key, score, node, level-1) // 往下一层找
}

func (l *SkipLinks) Add(key string, score float64, val interface{}) error {
	if score < l.minScore {
		return NewScoreOutOfRangeError(key, score, l.minScore)
	}
	if key == "" {
		return NewInvalidKeyError(key, "key empty")
	}
	// 添加
	hight, newNode := l.add(key, score, val, l.head, l.curLevel)
	if hight > l.curLevel { // 增加索引，这里只改变 head.Next[]，成员的自己去改变
		l.head.Next[l.curLevel] = newNode
		l.curLevel++
	}
	l.scoreMap[key] = score
	return nil
}

func (l *SkipLinks) add(key string, score float64, val interface{}, node *Node, level uint16) (uint16, *Node) {
	next := node.Next[level-1]
	for next != nil && l.compare(next.Score, score, next.Key, key) <= 0 { // 一层一层比
		// 后移
		node = next
		next = node.Next[level-1]
	}
	if level == 1 { // 找到了
		newNode := NewNode(key, score, val, l.maxLevel)
		newNode.Next[level-1] = next // 肯定在 next 前面
		node.Next[level-1] = newNode
		return l.getLevel(), newNode
	}
	hight, newNode := l.add(key, score, val, node, level-1) // 往下找
	if hight >= level {
		newNode.Next[level-1] = next
		node.Next[level-1] = newNode
	}
	return hight, newNode
}

func (l *SkipLinks) Erase(key string) (bool, interface{}) {
	score, exist := l.scoreMap[key]
	if !exist {
		return false, nil
	}
	delete(l.scoreMap, key)
	// 可能链表中已经没有了，缓存却有
	exist, rmNode := l.erase(key, score, l.head, l.curLevel)
	if !exist {
		return false, nil
	}
	// 注意这种情况，1被删除了
	/*
		1
		1
		1 2 3 4 5
	*/
	for i := int(l.curLevel - 1); i > 0; i-- {
		if l.head.Next[i] == nil { // 最高一层的被删除了（改变head的next）
			// Erase之后，层数可能减少
			l.curLevel--
		} else {
			break
		}
	}
	return true, rmNode.Val
}

func (l *SkipLinks) erase(key string, score float64, node *Node, level uint16) (bool, *Node) {
	next := node.Next[level-1]
	for next != nil && l.compare(next.Score, score, next.Key, key) < 0 {
		node = next
		next = node.Next[level-1]
	}
	if level == 1 {
		if next != nil && l.compare(next.Score, score, next.Key, key) == 0 {
			node.Next[level-1] = next.Next[level-1]
			next.Next[level-1] = nil
			return true, next
		}
		return false, nil
	}
	exist, rmNode := l.erase(key, score, node, level-1)
	if exist && next != nil && l.compare(next.Score, score, next.Key, key) == 0 {
		// 每一层都删掉
		node.Next[level-1] = next.Next[level-1]
		next.Next[level-1] = nil
	}
	return exist, rmNode
}

// Println 打印跳跃表
func (l *SkipLinks) Println() {
	for i := int(l.curLevel - 1); i >= 0; i-- {
		var (
			content string = fmt.Sprintf("level%d", i)
			node    *Node  = l.head.Next[i]
		)
		for node != nil {
			content += fmt.Sprintf(" %s(%v)", node.Key, node.Score)
			node = node.Next[i]
		}
		fmt.Println(content)
	}
}
