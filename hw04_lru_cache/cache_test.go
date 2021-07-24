package hw04lrucache

import (
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("overflow cache", func(t *testing.T) {
		c := NewCache(4)
		c.Set("k1", 1)
		c.Set("k2", 2)
		c.Set("k3", 3)
		c.Set("k4", 4)
		// в кеше должно быть: (front) 4,3,2,1 (back)

		c.Get("k1")
		c.Set("k2", 2)
		// чтение из кеша или перезапись должны продвигать элемент в начало
		// в кеше должно быть: (front) 2,1,4,3 (back)

		c.Set("k5", 5) // вытесняет k3
		c.Set("k6", 6) // вытесняет k4
		// добавление новых элементов сверх лимита должно выталкивать старые элементы
		// в кеше должно быть: (front) 6,5,2,1 (back)

		expected := []struct {
			key    Key
			exists bool
		}{
			{key: "k6", exists: true},
			{key: "k5", exists: true},
			{key: "k2", exists: true},
			{key: "k1", exists: true},
			{key: "k3", exists: false},
			{key: "k4", exists: false},
		}

		for _, exp := range expected {
			exp := exp
			t.Run(string(exp.key), func(t *testing.T) {
				_, exists := c.Get(exp.key)
				require.Equal(t, exp.exists, exists)
			})
		}
	})

	t.Run("clear", func(t *testing.T) {
		c := NewCache(2)
		c.Set("k1", strings.Repeat("A", 1024))
		c.Set("k2", strings.Repeat("B", 1024))

		c.Clear()

		k1, ok := c.Get("k1")
		require.False(t, ok)
		require.Nil(t, k1)

		k2, ok := c.Get("k2")
		require.False(t, ok)
		require.Nil(t, k2)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
