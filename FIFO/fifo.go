package fifo

import (
	"container/list"

	cache "felix.bs.com/felix/BeStrongerInGO/CacheTech"
)

type fifo struct {
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
func (f *fifo) Set(key string, value interface{}) {
	if e, ok := f.cache[key]; ok {
		f.ll.MoveToBack(e)
		en := e.Value.(*entry)
		//減去原本entry佔用的容量(en),加上新的entry佔用的容量(value)
		f.usedBytes = f.usedBytes - cache.CalcLen(en.value) + cache.CalcLen(value)
		en.value = value
		return
	}

	en := &entry{
		key:   key,
		value: value,
	}
	e := f.ll.PushBack(en)
	f.cache[key] = e

	f.usedBytes += en.Len()

	//超過限制大小,剔除無用的entry
	if f.maxBytes > 0 && f.usedBytes > f.maxBytes {
		f.DelOldsest()
	}
}

//Get: 取得Cache中所對應的entry
func (f *fifo) Get(key string) interface{} {
	if e, ok := f.cache[key]; ok {
		return e.Value.(*entry).value
	}
	return nil
}

//Del: 從Cache中刪除所對應的entry
func (f *fifo) Del(key string) {
	if e, ok := f.cache[key]; ok {
		f.removeElement(e)
	}
}

//DelOldsest: 剔除無用的entry (FIFO:為首節點)
func (f *fifo) DelOldsest() {
	f.removeElement(f.ll.Front())
}

func (f *fifo) Len() int {
	return f.ll.Len()
}

func (f *fifo) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	f.ll.Remove(e)
	en := e.Value.(*entry)
	f.usedBytes -= en.Len()
	delete(f.cache, en.key)

	if f.onEvicted != nil {
		f.onEvicted(en.key, en.value)
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
	return &fifo{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}
