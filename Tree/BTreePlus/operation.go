package BTreePlus

func NewBPTree(width int) *BPTree {

	if width < 3 {
		width = 3
	}

	var bt = &BPTree{}
	bt.Root = NewLeafNode(width)
	bt.HalfW = (bt.Width + 1) / 2

	return bt
}

func NewLeafNode(width int) *BPNode {

	var node = &BPNode{}
	node.Items = make([]BPItem, width+1)
	node.Items = node.Items[0:0]

	return node
}

//申请width+1是因为插入时可能暂时出现节点key大于申请width的情况,待后期再分裂处理
func NewIndexNode(width int) *BPNode {
	var node = &BPNode{}
	node.Nodes = make([]*BPNode, width+1)
	node.Nodes = node.Nodes[0:0]
	return node
}
