package BTreePlus

type BPNode struct {
	MaxKey int64     //子树的最大关键字
	Nodes  []*BPNode //节点的子树
	Items  []BPItem  //叶子结点的数据记录
	Next   *BPNode   //叶子结点中指向下一个节点，用于实现叶子结点链表
}

func (node *BPNode) deleteItem(key int64) bool {

	num := len(node.Items)
	for i := 0; i < num; i++ {
		if node.Items[i].Key > key {
			return false
		} else if node.Items[i].Key == key {

			copy(node.Items[i:], node.Items[i+1:])
			node.Items = node.Items[0 : len(node.Items)-1]
			node.MaxKey = node.Items[len(node.Items)-1].Key
			return true
		}
	}

	return false
}

func (node *BPNode) setValue(key int64, value interface{}) {

	item := BPItem{key, value}
	num := len(node.Items)

	if num < 1 {
		node.Items = append(node.Items, item)
		node.MaxKey = item.Key
		return
	} else if key < node.Items[0].Key {
		node.Items = append([]BPItem{item}, node.Items...)
		return
	} else if key > node.Items[num-1].Key {
		node.Items = append(node.Items, item)
	}

	for i := 0; i < num; i++ {
		if node.Items[i].Key > key {
			node.Items = append(node.Items, BPItem{})
			copy(node.Items[i+1:], node.Items[i:])
			node.Items[i] = item

			return
		} else if node.Items[i].Key == key {
			node.Items[i] = item
			return
		}
	}

}

func (t *BPTree) itemMoveOrMerge(parent, node *BPNode) {
	//获取兄弟结点
	var node1 *BPNode = nil
	var node2 *BPNode = nil
	for i := 0; i < len(parent.Nodes); i++ {
		if parent.Nodes[i] == node {
			if i < len(parent.Nodes)-1 {
				node2 = parent.Nodes[i+1]
			} else if i > 0 {
				node1 = parent.Nodes[i-1]
			}

			break
		}
	}

	//将左侧结点的记录移动到删除结点
	if node1 != nil && len(node1.Items) > t.HalfW {
		item := node1.Items[len(node1.Items)-1]
		node1.Items = node1.Items[0 : len(node1.Items)-1]
		node1.MaxKey = node1.Items[len(node1.Items)-1].Key
		node.Items = append([]BPItem{item}, node.Items...)
	}

	//将右侧结点的记录移动到删除结点
	if node2 != nil && len(node2.Items) > t.HalfW {
		item := node2.Items[0]
		node.Items = append(node.Items, item)
		node.MaxKey = node.Items[len(node.Items)-1].Key

		return
	}

	//与左侧结点进行合并
	if node1 != nil && len(node1.Items)+len(node.Items) <= t.Width {
		node1.Items = append(node1.Items, node.Items...)
		node1.Next = node.Next
		node1.MaxKey = node1.Items[len(node1.Items)-1].Key
		//parent.deleteChild(node)
		return
	}

	//与右侧结点进行合并
	if node2 != nil && len(node2.Items)+len(node.Items) <= t.Width {
		node.Items = append(node.Items, node2.Items...)
		node.Next = node2.Next
		node.MaxKey = node.Items[len(node.Items)-1].Key
		//parent.deleteChild(node2)
		return
	}

}

func (node *BPNode) addChild(child *BPNode) {

	num := len(node.Nodes)
	if num < 1 {
		node.Nodes = append(node.Nodes, child)
		node.MaxKey = child.MaxKey
		return
	} else if child.MaxKey < node.Nodes[0].MaxKey {
		node.Nodes = append([]*BPNode{child}, node.Nodes...)
		return
	} else if child.MaxKey > node.Nodes[num-1].MaxKey {
		node.Nodes = append(node.Nodes, child)
		node.MaxKey = child.MaxKey

		return
	}

	for i := 0; i < num; i++ {
		if node.Nodes[i].MaxKey > child.MaxKey {
			node.Nodes = append(node.Nodes, nil)
			copy(node.Nodes[i+1:], node.Nodes[i:])
			node.Nodes[i] = child

			return
		}
	}

}

func (node *BPNode) deleteChild(child *BPNode) bool {
	num := len(node.Nodes)
	for i := 0; i < num; i++ {
		if node.Nodes[i] == child {
			copy(node.Nodes[i:], node.Nodes[i+1:])
			node.Nodes = node.Nodes[0 : len(node.Nodes)-1]
			node.MaxKey = node.Nodes[len(node.Nodes)-1].MaxKey

			return true
		}
	}

	return false
}



func (t *BPTree) splitNode(node *BPNode) *BPNode {
	if len(node.Nodes) > t.Width {
		//创建新结点
		halfw := t.Width / 2 + 1
		node2 := NewIndexNode(t.Width)
		node2.Nodes = append(node2.Nodes, node.Nodes[halfw : len(node.Nodes)]...)
		node2.MaxKey = node2.Nodes[len(node2.Nodes)-1].MaxKey

		//修改原结点数据
		node.Nodes = node.Nodes[0:halfw]
		node.MaxKey = node.Nodes[len(node.Nodes)-1].MaxKey

		return node2
	} else if len(node.Items) > t.Width {
		//创建新结点
		halfw := t.Width / 2 + 1
		node2 := NewLeafNode(t.Width)
		node2.Items = append(node2.Items, node.Items[halfw: len(node.Items)]...)
		node2.MaxKey = node2.Items[len(node2.Items)-1].Key

		//修改原结点数据
		node.Next = node2
		node.Items = node.Items[0:halfw]
		node.MaxKey = node.Items[len(node.Items)-1].Key

		return node2
	}

	return nil
}