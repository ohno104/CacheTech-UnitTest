package lru

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetGetWithAssert(t *testing.T) {
	cache := New(24, nil)
	cache.Set("k1", 1)
	assert.Equal(t, cache.Get("k1"), 1)

	cache.Del("k1")
	assert.Equal(t, cache.Len(), 0)

	assert.Panics(t, func() { cache.Set("k2", time.Now()) })
}

func TestOnEvictedWithAssert(t *testing.T) {
	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}

	cache := New(16, onEvicted)
	cache.Set("k1", 1)
	cache.Set("k2", 2)
	cache.Get("k1")
	cache.Set("k3", 3)
	cache.Get("k1")
	cache.Set("k4", 4)

	expected := []string{"k2", "k3"}
	assert.Equal(t, keys, expected)

	assert.Equal(t, 2, cache.Len())
}
