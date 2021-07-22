package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("push back to empty", func(t *testing.T) {
		l := NewList() // empty list
		item := l.PushBack(10)
		require.Equal(t, item, l.Back())
		require.Equal(t, item, l.Front())
	})
}

func TestNextPrevLink(t *testing.T) {
	l := NewList()
	require.Nil(t, l.Front())
	require.Nil(t, l.Back())

	// nil <- (Front) item1 (Back) -> nil
	item1 := l.PushFront(10)
	compareList("push 10", t, l, item1)

	// nil <- (Front) item2 <-> item1 (Back) -> nil
	item2 := l.PushFront(20)
	compareList("push 20", t, l, item2, item1)

	// nil <- (Front) item3 <-> item2 <-> item1 (Back) -> nil
	item3 := l.PushFront(30)
	compareList("push 30", t, l, item3, item2, item1)

	// nil <- (Front) item1 <-> item3 <-> item2 (Back) -> nil
	l.MoveToFront(item1)
	compareList("move to front 10", t, l, item1, item3, item2)

	// nil <- (Front) item3 <-> item1 <-> item2 (Back) -> nil
	l.MoveToFront(item3)
	compareList("move to front 30", t, l, item3, item1, item2)

	// nil <- (Front) item3 <-> item2 (Back) -> nil
	l.Remove(item1)
	compareList("remove 10", t, l, item3, item2)

	// nil <- (Front) item3 (Back) -> nil
	l.Remove(item2)
	compareList("remove 20", t, l, item3)

	// nil <- (Front) (Back) -> nil
	l.Remove(item3)
	compareList("remove 30", t, l)
}

func compareList(name string, t *testing.T, l List, exp ...*ListItem) {
	t.Run(name, func(t *testing.T) {
		if len(exp) == 0 && l.Len() == 0 {
			return
		}
		if len(exp) != l.Len() {
			t.Fatalf("List length is not not equal. Expected: %d, actual: %d", len(exp), l.Len())
		}
		// пробегаем по списку от фронта к бэку
		j := 0
		for i := l.Front(); i != nil; i = i.Next {
			if i != exp[j] {
				t.Fatalf("List item is not equal. Expected: %v, actual: %v", exp[j], i)
			}
			require.Equal(t, i, exp[j])
			j++
		}
		// и в обратную сторону
		j = len(exp) - 1
		for i := l.Back(); i != nil; i = i.Prev {
			require.Equal(t, i, exp[j])
			j--
		}
	})
}
