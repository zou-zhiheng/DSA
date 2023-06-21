package Radix

import (
	"bytes"
	"fmt"
)

type RaxNode struct {
	bitLen int
	bitVal []byte
	left   *RaxNode
	right  *RaxNode
	val    interface{}
}

type Radix struct {
	root [256]*RaxNode
}

func NewRadixNode() *RaxNode {
	node := &RaxNode{}
	node.bitLen = 0
	node.bitVal = nil
	node.left = nil
	node.right = nil
	node.val = nil
	return node
}

func (r *RaxNode) pathMerge(bPos int) bool {

	//若当前节点存在值，则不能合并
	if r.val != nil {
		return false
	}

	//若当前节点有2个子结点，则不能合并
	if r.left != nil && r.right != nil {
		return false
	}

	//若当前结点没有子结点，则不能合并
	if r.left == nil && r.right == nil {
		return false
	}

	//获取当前结点的子结点
	child := r.left
	if r.right != nil {
		child = r.right
	}

	//判断当前结点最后一个字节是否是完整的字节
	//若不是完整的字节，需要与子结点的第一个字节进行合并
	if bPos%8 != 0 {
		charLen := len(r.bitVal)
		charLast := r.bitVal[charLen-1]
		char0000 := child.bitVal[0]
		child.bitVal = child.bitVal[1:]
		r.bitVal[charLen-1] = charLast | char0000
	}

	//合并当前结点以及子结点
	r.val = child.val
	r.bitVal = append(r.bitVal, child.bitVal...)
	r.bitLen += child.bitLen
	r.left = child.left
	r.right = child.right

	return true
}



//打印节点信息,用于调试
func (r *RaxNode)GetNodeInfo(bbeg int) string {
	buff := new(bytes.Buffer)

	bend := bbeg + int(r.bitLen)
	//起始和终止字节的位置
	cBeg := bbeg / 8; cend := bend / 8
	//起始和终止字节的偏移量
	oBeg := bbeg % 8; oend := bend % 8
	for bb := bbeg; bb < bend; {
		//获取两个数组的当前字节位置
		dci := bb / 8
		nci := dci - cBeg
		byteNode := r.bitVal[nci]

		//获取数据的当前字节以及循环步长
		step := 8
		if nci == 0 && oBeg > 0 {
			step = 8-oBeg
		}
		if dci == cend && oend > 0 {
			step = oend
		}
		if cBeg == cend {
			step = int(r.bitLen)
		}

		if step != 8 {
			buff.WriteString(fmt.Sprintf("(%08b:%d)", byteNode, byteNode))
		} else {
			buff.WriteByte(byteNode)
		}
		bb += step
	}

	if r.val != nil {
		buff.WriteString(fmt.Sprintf("=%v", r.val))
	}

	return buff.String()
}


//递归打印节点信息，用于调试
func (r *Radix) getNodeInfo(cur *RaxNode,pos int,data map[string]interface{}){
	data["info"] = cur.GetNodeInfo(pos)
	pos += int(cur.bitLen)

	if cur.left != nil {
		tmp := make(map[string]interface{})
		data["left"] = tmp
		r.getNodeInfo(cur.left, pos, tmp)
	}

	if cur.right != nil {
		tmp := make(map[string]interface{})
		data["right"] = tmp
		r.getNodeInfo(cur.right, pos, tmp)
	}
}