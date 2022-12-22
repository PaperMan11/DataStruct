package radix

import (
	"strings"
)

type node struct {
	pattern string           // 到终点时为整个路径字符串，eg. /path/:id/som
	part    string           // 路径的一部分，eg. :id path
	num     int              // 以当前结点为前缀的数量
	sons    map[string]*node // 下属节点
	fuzzy   bool             // 模糊匹配?该结点首字符为':'或'*'为模糊匹配
}

func newNode(part string) *node {
	fuzzy := false
	if len(part) > 0 {
		fuzzy = part[0] == ':' || part[0] == '*'
	}
	return &node{
		pattern: "",
		part:    part,
		num:     0,
		sons:    make(map[string]*node),
		fuzzy:   fuzzy,
	}
}

// analysis 将 s 按 '/' 进行分段解析，同时按规则重组
func analysis(s string) (ss []string, newS string) {
	vs := strings.Split(s, "/")
	ss = make([]string, 0)
	for _, item := range vs {
		if item != "" {
			ss = append(ss, item)
			newS = newS + "/" + item
			if item[0] == '*' {
				break
			}
		}
	}
	return ss, newS
}

// inOrder 用于遍历
func (n *node) inOrder(s string) (es []interface{}) {
	if n == nil {
		return es
	}
	if n.pattern != "" {
		es = append(es, s+n.part)
	}
	for _, son := range n.sons {
		es = append(es, son.inOrder(s+n.part+"/")...)
	}
	return es
}

// insert 添加某条 string 以及所有的 part
func (n *node) insert(pattern string, ss []string, i int) (b bool) {
	if i == len(ss) {
		if n.pattern != "" {
			// 该节点承载了string
			return false
		}
		n.pattern = pattern
		n.num++
		return true
	}
	s := ss[i]
	son, ok := n.sons[s]
	if !ok {
		// 不存在
		son = newNode(s)
		n.sons[s] = son
	}
	b = son.insert(pattern, ss, i+1)
	if b {
		n.num++
	} else {
		//插入失败且该子节点为新建结点则需要删除该子结点
		if !ok {
			delete(n.sons, s)
		}
	}
	return b
}

// erase 删除某个 string
func (n *node) erase(ss []string, i int) (b bool) {
	if i == len(ss) {
		if n.pattern != "" {
			n.pattern = ""
			n.num--
			return true
		}
		return false
	}
	s := ss[i]
	son, ok := n.sons[s]
	if !ok || son == nil {
		return false
	}
	b = son.erase(ss, i+1)
	if b {
		n.num--
		//删除后子结点的num<=0即该节点无后续存储元素,可以销毁
		if son.num <= 0 {
			delete(n.sons, s)
		}
	}
	return b
}

// delete 删除前缀为 ss 的所有 string
func (n *node) delete(ss []string, i int) (num int) {
	if i == len(ss) {
		return n.num
	}

	s := ss[i]
	son, ok := n.sons[s]
	if !ok || son == nil {
		return 0
	}
	num = son.delete(ss, i+1)
	if num > 0 {
		son.num -= num
		if son.num <= 0 {
			delete(n.sons, s)
		}
	}
	return num
}

// count 前缀为 ss 的 string 的数量
func (n *node) count(ss []string, i int) (num int) {
	if i == len(ss) {
		return n.num
	}
	s := ss[i]
	son, ok := n.sons[s]
	if !ok || son == nil {
		return 0
	}
	return son.count(ss, i+1)
}

// find 精确匹配找到终点结点
func (n *node) find(ss []string, i int) *node {
	if i == len(ss) || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	s := ss[i]
	children := make([]*node, 0)
	for _, child := range n.sons {
		//局部string相同或动态匹配
		if child.part == s || child.fuzzy {
			children = append(children, child)
		}
	}

	for _, child := range children {
		res := child.find(ss, i+1)
		if res != nil {
			return res
		}
	}
	return nil
}

// mate 对具体的路由进行参数映射
// eg.
// 		/api/v1/user/:id/post
// 		/api/v1/user/1/post
//		map {"id": "1"}
func (n *node) mate(s string) (map[string]string, bool) {
	//解析url
	searchParts, _ := analysis(s)
	//从该请求类型中寻找对应的路由结点
	targetNode := n.find(searchParts, 0)
	if targetNode != nil {
		//解析该结点的pattern
		parts, _ := analysis(targetNode.pattern)
		//动态参数映射表
		params := make(map[string]string)
		for index, part := range parts {
			if part[0] == ':' {
				//动态匹配,将参数名和参数内容的映射放入映射表内
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				//通配符,将后续所有内容全部添加到映射表内同时结束遍历
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return params, true
	}
	return nil, false
}
