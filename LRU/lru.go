package lru

import (
	"container/list"

	cache "felix.bs.com/felix/BeStrongerInGO/CacheTech"
)

type lru struct {
	//Cache最大的容量(0表示不限制最大容量)
	//若是groupcache 使用的是最大儲存的entry個數
	maxBytes int

	//當一個entry從Cache中移除時呼叫此回呼函數
	onEvicted func(key string, value interface{})

	//已使用的位元組數
	usedBytes int

	//雙向link-list
	ll *list.List

	//Element中有一個value欄位為interface{},可以儲存任意類型
	cache map[string]*list.Element
}

//Set: 往Cache尾部增加一個元素(若已存在則移到尾部,並修改值)
func (l *lru) Set(key string, value interface{}) {
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToBack(e)
		en := e.Value.(*entry)
		//減去原本entry佔用的容量(en),加上新的entry佔用的容量(value)
		l.usedBytes = l.usedBytes - cache.CalcLen(en.value) + cache.CalcLen(value)
		en.value = value
		return
	}

	en := &entry{
		key:   key,
		value: value,
	}
	e := l.ll.PushBack(en)
	l.cache[key] = e

	l.usedBytes += en.Len()

	//超過限制大小,剔除無用的entry
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.DelOldest()
	}
}

//Get: 取得Cache中所對應的entry
func (l *lru) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {

		//資料最近被存取過,重新放入Tail
		l.ll.MoveToBack(e)
		return e.Value.(*entry).value
	}
	return nil
}

//Del: 從Cache中刪除所對應的entry
func (l *lru) Del(key string) {
	if e, ok := l.cache[key]; ok {
		l.removeElement(e)
	}
}

//DelOldest: 剔除無用的entry (FIFO:為首節點)
func (l *lru) DelOldest() {
	l.removeElement(l.ll.Front())
}

func (l *lru) Len() int {
	return l.ll.Len()
}

func (l *lru) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	l.ll.Remove(e)
	en := e.Value.(*entry)
	l.usedBytes -= en.Len()
	delete(l.cache, en.key)

	if l.onEvicted != nil {
		l.onEvicted(en.key, en.value)
	}
}

type entry struct {
	key   string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	return &lru{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}
