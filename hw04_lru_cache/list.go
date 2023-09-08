package hw04lrucache

import (
	"fmt"
	"strings"
)

type List interface {
	Len() int                          // длина списка
	Front() *ListItem                  // первый элемент списка
	Back() *ListItem                   // последний элемент списка
	PushFront(v interface{}) *ListItem // добавить значение в начало
	PushBack(v interface{}) *ListItem  // добавить значение в конец
	Remove(i *ListItem)                // удалить элемент
	MoveToFront(i *ListItem)           // переместить элемент в начало
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	first *ListItem
	last  *ListItem
	items map[*ListItem]*ListItem
}

func (l list) String() string {
	s := &strings.Builder{}
	cnt := 0
	i := l.first
	for i != nil {
		if i.Prev != nil {
			s.WriteString(fmt.Sprintf("(%v", i.Prev.Value))
		} else {
			s.WriteString("(nil")
		}
		s.WriteString(fmt.Sprintf("<-%v->", i.Value))
		if i.Next != nil {
			s.WriteString(fmt.Sprintf("%v) ", i.Next.Value))
		} else {
			s.WriteString("nil) ")
		}
		cnt++
		i = i.Next
	}
	return fmt.Sprintf("len: %v [ %v]", cnt, s.String())
}

func (l *list) Len() int {
	return len(l.items)
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.first != nil {
		l.first.Prev = item
		item.Next = l.first
	} else {
		l.last = item
	}
	l.first = item

	l.items[item] = item

	return l.first
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.last != nil {
		l.last.Next = item
		item.Prev = l.last
	} else {
		l.first = item
	}
	l.last = item

	l.items[item] = item

	return l.last
}

func (l *list) Remove(i *ListItem) {
	item, found := l.items[i]
	if !found {
		panic("item not found")
	}

	// может быть первым, последним или средним
	if item.Prev != nil {
		item.Prev.Next = item.Next
	} else {
		// если это первый элемент - голову на следующий
		l.first = item.Next
	}

	if item.Next != nil {
		item.Next.Prev = item.Prev
	} else {
		// если это последний элемент - хвост на предыдущий
		l.last = item.Prev
	}

	delete(l.items, i)
}

func (l *list) MoveToFront(i *ListItem) {
	item, found := l.items[i]
	if !found {
		panic("item not found")
	}

	if l.first == item {
		return
	}
	if item.Prev != nil {
		item.Prev.Next = item.Next
	}
	if item.Next != nil {
		item.Next.Prev = item.Prev
	}
	item.Next = l.first
	l.first = item
	if l.last == item {
		l.last = item.Prev
	}
	item.Prev = nil
}

func NewList() List {
	return &list{
		items: map[*ListItem]*ListItem{},
	}
}
