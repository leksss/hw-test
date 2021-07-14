package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	RemoveAll()
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	items map[*ListItem]*ListItem
	front *ListItem
	back  *ListItem
}

func NewList() List {
	l := new(list)
	l.items = make(map[*ListItem]*ListItem)
	return l
}

func (l *list) Len() int {
	return len(l.items)
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	curFirstItem := &ListItem{Value: v}

	if l.Len() > 0 {
		oldFirstItem := l.Front()
		curFirstItem.Next = oldFirstItem
		oldFirstItem.Prev = curFirstItem
	}

	l.items[curFirstItem] = curFirstItem
	l.front = curFirstItem
	if l.Len() == 1 {
		l.back = curFirstItem
	}
	return curFirstItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	curLastItem := &ListItem{Value: v}

	if l.Len() > 0 {
		oldLastItem := l.Back()
		curLastItem.Prev = oldLastItem
		oldLastItem.Next = curLastItem
	}

	l.items[curLastItem] = curLastItem
	l.back = curLastItem
	if l.Len() == 1 {
		l.front = curLastItem
	}
	return curLastItem
}

func (l *list) Remove(i *ListItem) {
	curItem, ok := l.items[i]
	if !ok {
		return
	}

	prevItem := curItem.Prev
	nextItem := curItem.Next

	if prevItem != nil {
		prevItem.Next = nextItem
	}
	if nextItem != nil {
		nextItem.Prev = prevItem
	}

	if i == l.front {
		l.front = i.Next
	}
	if i == l.back {
		l.back = i.Prev
	}

	delete(l.items, i)
}

func (l *list) RemoveAll() {
	l.items = make(map[*ListItem]*ListItem)
	l.front = nil
	l.back = nil
}

func (l *list) MoveToFront(i *ListItem) {
	if l.Front() == i {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}
