package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	l.pushFrontItem(item)
	return item
}

func (l *list) pushFrontItem(item *ListItem) *ListItem {
	if l.len == 0 {
		l.front = item
		l.back = item
	} else {
		l.front.Prev = item
		item.Next = l.front
		l.front = item
	}
	l.len++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.len == 0 {
		l.front = item
		l.back = item
	} else {
		l.back.Next = item
		item.Prev = l.back
		l.back = item
	}
	l.len++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil && i.Next == nil {
		// элемент не принадлежит списку - ничего не делаем и выходим
		return
	}
	if i.Prev == nil {
		// если это первый элемент в списке - то заменяем его на следующий элемент
		l.front = i.Next
	} else {
		// если не первый - то переключаем выход предыдущего элемента на вход следующего
		i.Prev.Next = i.Next
	}
	if i.Next == nil {
		// если это последний элемент в списке - то заменяем его на предыдущий элемент
		l.back = i.Prev
	} else {
		// если не последний - то переключаем вход следующего элемента на выход предыдущего
		i.Next.Prev = i.Prev
	}
	// у самого элемента ссылки на предыдущий и следующий элемент тоже надо обнулить
	i.Next = nil
	i.Prev = nil
	// уменьшаем длину списка
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil && i.Next == nil {
		// элемент не принадлежит списку - ничего не делаем и выходим
		return
	}
	if l.front == i {
		return // уже все сделано)
	}
	l.Remove(i)
	l.pushFrontItem(i)
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func NewList() List {
	return new(list)
}
