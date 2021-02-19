package cache

type Getter interface {
	Get(key string) interface{}
}

//GetFunc是一個函數的Type, 該函數為func(key string)會回傳interface{}
type GetFunc func(key string) interface{}

//Type GetFunc有一個Get方法
//由於Type GetFunc實作了Get(key string)方法, 因此Type GetFunc為Getter interface
func (f GetFunc) Get(key string) interface{} {
	//f => GentFunc => func(key) interface{}
	return f(key)
}

type TourCache struct {
	mainCache *safeCache
	getter    Getter
}

func NewTourCache(getter Getter, cache Cache) *TourCache {
	return &TourCache{
		mainCache: newSafeCache(cache),
		getter:    getter,
	}
}

func (t *TourCache) Get(key string) interface{} {
	val := t.mainCache.get(key)
	if val != nil {
		return val
	}

	if t.getter != nil {
		val = t.getter.Get(key)
		if val == nil {
			return nil
		}

		t.mainCache.set(key, val)
		return val
	}

	return nil
}

func (t *TourCache) Stat() *Stat {
	return t.mainCache.stat()
}
