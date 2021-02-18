package fifo

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetGet(t *testing.T) {
	cache := New(24, nil)
	cache.DelOldsest()
	cache.Set("k1", 1)
	v := cache.Get("k1")
	if v != 1 {
		t.Fatal("cache.Get test error!")
	}

	cache.Del("k1")
	if cache.Len() != 0 {
		t.Fatal("cache.Del test error!")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Implement \"Value interface{}\" test error!")
		}
	}()
	cache.Set("k2", time.Now())
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}

	cache := New(16, onEvicted)
	cache.Set("k1", 1)
	cache.Set("k2", 2)
	cache.Get("k1")
	cache.Set("k3", 3)
	cache.Get("k2")
	cache.Set("k4", 4)

	expected := []string{"k1", "k2"}
	if !reflect.DeepEqual(keys, expected) {
		t.Fatal("cache.onEvicted test error!")
	}

	if cache.Len() != 2 {
		t.Fatal("cache.Len test error!")
	}

}

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
	cache.Get("k2")
	cache.Set("k4", 4)

	expected := []string{"k1", "k2"}
	assert.Equal(t, keys, expected)

	assert.Equal(t, cache.Len(), 2)
}
