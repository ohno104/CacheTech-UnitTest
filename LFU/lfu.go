package lfu

import (
	"container/heap"

	cache "felix.bs.com/felix/BeStrongerInGO/CacheTech"
)

type lfu struct {
	//Cache最大的容量(0表示不限制最大容量)
	//若是groupcache 使用的是最大儲存的entry個數
	maxBytes int

	//當一個entry從Cache中移除時呼叫此回呼函數
	onEvicted func(key string, value interface{})

	//已使用的位元組數
	usedBytes int

	queue *queue

	cache map[string]*entry
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	q := make(queue, 0, 1024)
	return &lfu{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		queue:     &q,
		cache:     make(map[string]*entry),
	}
}

func (l *lfu) Set(key string, value interface{}) {
	if e, ok := l.cache[key]; ok {
		l.usedBytes = l.usedBytes - cache.CalcLen(e.value) + cache.CalcLen(value)
		l.queue.update(e, value, e.weight+1)
		//fmt.Println("l.usedBytes: ", l.usedBytes)
		//fmt.Println("CalcLen: ", cache.CalcLen(value))
		return
	}

	en := &entry{
		key:   key,
		value: value,
	}
	heap.Push(l.queue, en)
	l.cache[key] = en

	l.usedBytes += en.Len()
	//fmt.Println("new l.usedBytes: ", l.usedBytes)
	//fmt.Println("l.maxBytes: ", l.maxBytes)
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		//超過限制大小,剔除最少使用的entry
		l.removeElement(heap.Pop(l.queue))

		//fmt.Println("removeElement! ")
	}
}

func (l *lfu) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {
		l.queue.update(e, e.value, e.weight+1)
		return e.value
	}

	return nil
}

func (l *lfu) Del(key string) {
	if e, ok := l.cache[key]; ok {
		heap.Remove(l.queue, e.index)
		l.removeElement(e)
	}
}

func (l *lfu) DelOldest() {
	if l.queue.Len() == 0 {
		return
	}
	l.removeElement(heap.Pop(l.queue))
}

func (l *lfu) removeElement(x interface{}) {
	if x == nil {
		return
	}

	en := x.(*entry)
	delete(l.cache, en.key)
	l.usedBytes -= en.Len()
	if l.onEvicted != nil {
		l.onEvicted(en.key, en.value)
	}
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}

func (l *lfu) Len() int {
	return l.queue.Len()
}
