package lfu

import "container/heap"

type queue []*entry

type entry struct {
	key   string
	value interface{}

	//表示entry在queue中的優先順序(存取次數越多越高)
	weight int

	//表示在heap中的索引
	index int
}

func (q queue) Len() int {
	return len(q)
}

func (q queue) Less(i, j int) bool {
	return q[i].weight < q[j].weight
}

func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

//heap.Interface
func (q *queue) Push(x interface{}) {
	n := len(*q)
	en := x.(*entry)
	en.index = n
	*q = append(*q, en)
}

//heap.Interface
func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	en := old[n-1]
	old[n-1] = nil
	en.index = -1
	*q = old[0 : n-1]
	return en
}

//更新priority queue
func (q *queue) update(en *entry, value interface{}, weight int) {
	en.value = value
	en.weight = weight
	heap.Fix(q, en.index)
}
