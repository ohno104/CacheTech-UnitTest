package cache_test

import (
	"sync"
	"testing"

	cache "felix.bs.com/felix/BeStrongerInGO/CacheTech"
	lru "felix.bs.com/felix/BeStrongerInGO/CacheTech/LRU"
	_ "github.com/allegro/bigcache/v2"
	_ "github.com/coocood/freecache"
	"github.com/stretchr/testify/assert"
)

func TestTourCacheGet(t *testing.T) {
	mock := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
	}

	getter := cache.GetFunc(func(key string) interface{} {
		//log.Println("[From mock DB] find key", key)

		if val, ok := mock[key]; ok {
			return val
		}
		return nil
	})

	tourCache := cache.NewTourCache(getter, lru.New(0, nil))

	var wg sync.WaitGroup

	for k, v := range mock {
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()
			assert.Equal(t, tourCache.Get(k), v)
			assert.Equal(t, tourCache.Get(k), v)
		}(k, v)
	}
	wg.Wait()

	assert.Equal(t, tourCache.Get("X"), nil)
	assert.Equal(t, tourCache.Get("X"), nil)

	assert.Equal(t, tourCache.Stat().NGet, 10)

	//第一次tourCache.Get的快取為空,呼叫getter方法將資料放入快取中(未命中)
	//第二次直接從快取讀取資料(命中)
	assert.Equal(t, tourCache.Stat().NHit, 4)
}
