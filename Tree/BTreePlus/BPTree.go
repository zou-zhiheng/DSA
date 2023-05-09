package BTreePlus

import "sync"

type BPTree struct {
	Mutex sync.RWMutex //读写锁
	Root  *BPNode      //表示B+树的根节点
	Width int          //树的阶
	HalfW int          //用于[M/2]=ceil(M/2)
}


// Get B+树的查找
func (t *BPTree) Get(key int64) interface{} {

	t.Mutex.Lock()         //上锁，保证读写安全
	defer t.Mutex.Unlock() //释放锁

	node := t.Root
	//找到目标子结点
	for i := 0; i < len(node.Nodes); i++ {
		if key <= node.Nodes[i].MaxKey {
			node = node.Nodes[i]
			i = 0
		}
	}

	//没有找到叶子结点
	if len(node.Nodes) > 0 {
		return nil
	}

	for i:=0;i<len(node.Items);i++{
		if node.Items[i].Key==key{
			return node.Items[i].Val
		}
	}

	return nil
}

func (t *BPTree) deleteItem(parent *BPNode, node *BPNode, key int64) {
	for i:=0; i < len(node.Nodes); i++ {
		if key <= node.Nodes[i].MaxKey {
			t.deleteItem(node, node.Nodes[i], key)
			break
		}
	}

	if  len(node.Nodes) < 1 {
		//删除记录后若结点的子项<m/2，则从兄弟结点移动记录，或者合并结点
		node.deleteItem(key)
		if len(node.Items) < t.HalfW {
			t.itemMoveOrMerge(parent, node)
		}
	} else {
		//若结点的子项<m/2，则从兄弟结点移动记录，或者合并结点
		node.MaxKey = node.Nodes[len(node.Nodes)-1].MaxKey
		if len(node.Nodes) < t.HalfW {
			t.childMoveOrMerge(parent, node)
		}
	}
}


func (t *BPTree) Set(key int64, value interface{}) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	t.setValue(nil, t.Root, key, value)
}

func (t *BPTree) setValue(parent *BPNode, node *BPNode, key int64, value interface{}) {
	for i:=0; i < len(node.Nodes); i++ {
		if key <= node.Nodes[i].MaxKey || i== len(node.Nodes)-1 {
			t.setValue(node, node.Nodes[i], key, value)
			break
		}
	}

	//叶子结点，添加数据
	if len(node.Nodes) < 1 {
		node.setValue(key, value)
	}

	//结点分裂
	node_new := t.splitNode(node)
	if node_new != nil {
		//若父结点不存在，则创建一个父节点
		if parent == nil {
			parent = NewIndexNode(t.Width)
			parent.addChild(node)
			t.Root = parent
		}
		//添加结点到父亲结点
		parent.addChild(node_new)
	}
}

func (t *BPTree) childMoveOrMerge(parent *BPNode, node *BPNode) {
	if parent == nil {
		return
	}

	//获取兄弟结点
	var node1 *BPNode = nil
	var node2 *BPNode = nil
	for i:=0; i < len(parent.Nodes); i++ {
		if parent.Nodes[i] == node {
			if i < len(parent.Nodes)-1 {
				node2 = parent.Nodes[i+1]
			} else if i > 0 {
				node1 = parent.Nodes[i-1]
			}
			break
		}
	}

	//将左侧结点的子结点移动到删除结点
	if node1 != nil && len(node1.Nodes) > t.HalfW {
		item := node1.Nodes[len(node1.Nodes)-1]
		node1.Nodes = node1.Nodes[0:len(node1.Nodes)-1]
		node.Nodes = append([]*BPNode{item}, node.Nodes...)
		return
	}

	//将右侧结点的子结点移动到删除结点
	if node2 != nil && len(node2.Nodes) > t.HalfW {
		item := node2.Nodes[0]
		node2.Nodes = node1.Nodes[1:]
		node.Nodes = append(node.Nodes, item)
		return
	}

	if node1 != nil && len(node1.Nodes) + len(node.Nodes) <= t.Width {
		node1.Nodes = append(node1.Nodes, node.Nodes...)
		parent.deleteChild(node)
		return
	}

	if node2 != nil && len(node2.Nodes) + len(node.Nodes) <= t.Width {
		node.Nodes = append(node.Nodes, node2.Nodes...)
		parent.deleteChild(node2)
		return
	}
}