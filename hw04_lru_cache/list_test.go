package hw04lrucache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMy(t *testing.T) {
	l := NewList()
	l.PushFront(30)
	helper(l)
	l.PushFront(20)
	helper(l)
	i10 := l.PushFront(10)
	helper(l)

	i40 := l.PushBack(40)
	helper(l)
	i50 := l.PushBack(50)
	helper(l)
	i60 := l.PushBack(60)
	helper(l)

	l.Remove(i10)
	helper(l)

	l.Remove(i60)
	helper(l)

	l.MoveToFront(i40)
	helper(l)

	l.MoveToFront(i50)
	helper(l)
}

func helper(l List) {
	var e *ListItem
	e = l.Front()
	fmt.Printf("(%d): ", l.Len())
	for e != nil {
		fmt.Print(e.Value, " -> ")
		e = e.Next
	}
	fmt.Println()
}

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
}
