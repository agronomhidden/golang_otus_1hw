package hw04_lru_cache //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var counter int64
var ErrOuterItem = errors.New("list is not contain received item")
var ErrNilledItem = errors.New("nilled item has been received")

func getID() uint {
	return uint(atomic.AddInt64(&counter, 1))
}

// хочется сделать защиту от передачи в список элемента,
// который не принадлежит данному списку
// для этого вводим некий идентификатор элемента
// и ошибки если передали чужеродный элемент
// + может быть передан нулевой указатель, как у нас в тестах: middle := l.Back().Next // 20
type List interface {
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(v interface{}) *listItem
	PushBack(v interface{}) *listItem
	Remove(i *listItem) error
	MoveToFront(i *listItem) error
}

// думаю есть необходимость защитить next и prev, чтобы извне нельзя было переопределить соседа
// на элемент из совсем другого списка.
type listItem struct {
	Value interface{}
	next  *listItem
	prev  *listItem
	id    uint
}

func (e listItem) Next() *listItem {
	return e.next
}

func (e listItem) Prev() *listItem {
	return e.prev
}

type list struct {
	first *listItem
	last  *listItem
	index map[uint]struct{}
	sync.Mutex
	len int
}

func (e *list) Len() int {
	return e.len
}

func (e *list) PushFront(v interface{}) *listItem {
	newItem := e.initNewItem(v)
	newItem.next = e.first
	if e.first == nil {
		e.last = newItem
	} else {
		newItem.next.prev = newItem
	}
	e.first = newItem

	return newItem
}

func (e *list) PushBack(v interface{}) *listItem {
	newItem := e.initNewItem(v)
	newItem.prev = e.last
	if e.last == nil {
		e.first = newItem
	} else {
		newItem.prev.next = newItem
	}
	e.last = newItem

	return newItem
}

func (e *list) Remove(item *listItem) error {
	if item == nil {
		return ErrNilledItem
	}
	if _, ok := e.index[item.id]; !ok {
		return nil
	}
	e.len--
	delete(e.index, item.id)
	item.id = 0

	switch {
	case item.prev == nil && item.next == nil:
		e.first = nil
		e.last = nil
	case item.prev == nil:
		e.first = item.next
		item.next = nil
	case item.next == nil:
		e.last = item.prev
		item.prev = nil
	default:
		item.prev.next = item.next
		item.next.prev = item.prev
		item.next = nil
		item.prev = nil
	}
	return nil
}

func (e *list) MoveToFront(item *listItem) error {
	if item == nil {
		return ErrNilledItem
	}
	if _, ok := e.index[item.id]; !ok {
		return ErrOuterItem
	}
	e.Remove(item)
	e.PushFront(item.Value)

	return nil
}

func (e *list) Front() *listItem {
	return e.first
}

func (e *list) Back() *listItem {
	return e.last
}

func (e *list) initNewItem(v interface{}) *listItem {

	newItem := &listItem{Value: v, id: getID()}
	e.index[newItem.id] = struct{}{}
	e.len++

	return newItem
}

func NewList() List {
	return &list{index: make(map[uint]struct{})}
}
