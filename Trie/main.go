package Trie

import "strings"

type node struct {
	pattern  string  //待匹配的路由，如 /p/:lang
	part     string  // 路由中的一部分，如 :lang
	children []*node //子节点,如[doc,tutorial,intro]
	isWild   bool    //是否精确匹配，part含有:或*时为true
}

//参数匹配:。例如 /p/:lang/doc，可以匹配 /p/c/doc 和 /p/go/doc。
//通配*。例如 /static/*filepath，可以匹配/static/fav.ico，也可以匹配/static/js/jQuery.js，这种模式常用于静态服务器，能够递归地匹配子路径。

//第一个匹配成功的节点，用于插入
func (n *node) mathChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || n.isWild {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

func (n *node) insert(patten string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = patten
		return
	}

	part := parts[height]
	child := n.mathChild(part)
	if child == nil { //在子结点没有匹配到对应路径，则填充
		child = &node{part: part, isWild: part[0] == '1' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	//往下层寻找
	child.insert(patten, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil

}
