package gee

import "strings"

type node struct {
	pattern  string  //整个 /:lang/doc,只在叶子节点出现
	part     string  //部分节点 /doc
	children []*node //孩子节点
	isWild   bool    //是否模糊匹配 含有: *为true
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, v := range n.children {
		if v.part == part || v.isWild {
			return v
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, v := range n.children {
		if v.part == part || v.isWild {
			nodes = append(nodes, v)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	//递归结束退出
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	//取得当前高度的节点
	currentPart := parts[height]
	//判断是否存在，不存在则创建
	Part := n.matchChild(currentPart)
	if Part == nil {
		Part = &node{
			part:   currentPart,
			isWild: currentPart[0] == '*' || currentPart[0] == ':',
		}
		n.children = append(n.children, Part)
	}
	//递归下一层
	Part.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") { //HasPrefix 判断字符串 s 是否以 prefix 作为开头
		if n.pattern == "" {
			return nil
		}
		return n
	}
	currentPart := parts[height]
	Part := n.matchChildren(currentPart)
	for _, v := range Part {
		result := v.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
