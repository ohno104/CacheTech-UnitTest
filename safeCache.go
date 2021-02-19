package cache

import (
	"log"
	"sync"
)

//預設允許佔用的最大記憶體
const DefaultMaxByte = 1 << 29

//平行處理的安全快取
type safeCache struct {
	m     sync.RWMutex
	cache Cache

	//快取命中次數
	nhit int
	//快取取得次數
	nget int
}

type Stat struct {
	NHit int
	NGet int
}

func newSafeCache(cache Cache) *safeCache {
	return &safeCache{
		cache: cache,
	}
}

func (sc *safeCache) set(key string, value interface{}) {
	sc.m.Lock()
	defer sc.m.Unlock()
	sc.cache.Set(key, value)
}

func (sc *safeCache) get(key string) interface{} {
	sc.m.Lock()
	defer sc.m.Unlock()
	sc.nget++

	if sc.cache == nil {
		return nil
	}

	v := sc.cache.Get(key)
	if v != nil {
		log.Println("[TourCache] hit")
		sc.nhit++
	}

	return v
}

func (sc *safeCache) stat() *Stat {
	sc.m.Lock()
	defer sc.m.Unlock()
	return &Stat{
		NHit: sc.nhit,
		NGet: sc.nget,
	}
}
