package hw04lrucache

const (
	back  = "back"
	front = "front"
)

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
	l.pushItem(item, front)
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	l.pushItem(item, back)
	return item
}

func (l *list) pushItem(item *ListItem, whereTo string) {
	switch {
	case l.len == 0:
		l.front = item
		l.back = item
	case whereTo == front:
		l.front.Prev = item
		item.Next = l.front
		l.front = item
	case whereTo == back:
		l.back.Next = item
		item.Prev = l.back
		l.back = item
	default:
		return
	}
	l.len++
}

func (l *list) Remove(i *ListItem) {
	// если удаляем первый элемент
	if l.Front() == i {
		l.front = i.Next
	}
	// если удаляем последний элемент
	if l.Back() == i {
		l.back = i.Prev
	}
	// соединяем предыдущий и следующий элементы вместо удаленного
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	i.Next = nil
	i.Prev = nil
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	// элемент уже первым является
	if l.front == i {
		return
	}
	l.Remove(i)
	l.pushItem(i, front)
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
