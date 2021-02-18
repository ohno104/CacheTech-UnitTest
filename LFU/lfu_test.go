package lfu

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetGetWithAssert(t *testing.T) {
	cache := New(24, nil)
	cache.DelOldsest()
	cache.Set("k1", "A")
	assert.Equal(t, cache.Get("k1"), "A")
	cache.Del("k1")
	assert.Equal(t, cache.Len(), 0)

	assert.Panics(t, func() { cache.Set("k2", time.Now()) })
}

func TestOnEvictedWithAssert(t *testing.T) {
	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}

	cache := New(24, onEvicted)
	cache.Set("A", 1)
	cache.Set("B", 2)
	cache.Set("C", 3)
	cache.Set("D", 4)
	expected := []string{"A"}
	assert.Equal(t, expected, keys)

	cache.Get("D")
	cache.Get("D")
	cache.Get("C")
	cache.DelOldsest()
	expected = []string{"A", "B"}
	assert.Equal(t, expected, keys)

	cache.Del("D")
	expected = []string{"A", "B", "D"}
	assert.Equal(t, expected, keys)

	assert.Equal(t, 1, cache.Len())
}
