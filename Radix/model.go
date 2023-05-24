package Radix

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
